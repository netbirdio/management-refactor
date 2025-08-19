package permissions

//go:generate go run github.com/golang/mock/mockgen -package permissions -destination=manager_mock.go -source=./manager.go -build_flags=-mod=mod

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/activity"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/modules"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/roles"
	nbcontext "github.com/netbirdio/netbird/management/server/context"
	"github.com/netbirdio/netbird/shared/management/http/util"
	"github.com/netbirdio/netbird/shared/management/status"
)

type Manager interface {
	WithPermission(module modules.Module, operation operations.Operation, handlerFunc func(w http.ResponseWriter, r *http.Request, auth *nbcontext.UserAuth)) http.HandlerFunc
	ValidateUserPermissions(ctx context.Context, accountID, userID string, module modules.Module, operation operations.Operation) (bool, error)
	ValidateRoleModuleAccess(ctx context.Context, accountID string, role roles.RolePermissions, module modules.Module, operation operations.Operation) bool
	ValidateAccountAccess(ctx context.Context, accountID string, user *users.User, allowOwnerAndAdmin bool) error
}

type managerImpl struct {
	userManager users.Manager
}

func NewManager(userManager users.Manager) Manager {
	return &managerImpl{
		userManager: userManager,
	}
}

func (m *managerImpl) ValidateUserPermissions(
	ctx context.Context,
	accountID string,
	userID string,
	module modules.Module,
	operation operations.Operation,
) (bool, error) {
	if userID == activity.SystemInitiator || userID == "allowedUser" {
		return true, nil
	}

	user, err := m.userManager.GetUserByID(ctx, nil, db.LockingStrengthShare, userID)
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, status.NewUserNotFoundError(userID)
	}

	if user.IsBlocked() {
		return false, status.NewUserBlockedError()
	}

	if err := m.ValidateAccountAccess(ctx, accountID, user, false); err != nil {
		return false, err
	}

	if operation == operations.Read && user.IsServiceUser {
		return true, nil // this should be replaced by proper granular access role
	}

	role, ok := roles.RolesMap[user.Role]
	if !ok {
		return false, status.NewUserRoleNotFoundError(string(user.Role))
	}

	return m.ValidateRoleModuleAccess(ctx, accountID, role, module, operation), nil
}

func (m *managerImpl) ValidateRoleModuleAccess(
	ctx context.Context,
	accountID string,
	role roles.RolePermissions,
	module modules.Module,
	operation operations.Operation,
) bool {
	if permissions, ok := role.Permissions[module]; ok {
		if allowed, exists := permissions[operation]; exists {
			return allowed
		}
		log.WithContext(ctx).Tracef("operation %s not found on module %s for role %s", operation, module, role.Role)
		return false
	}

	return role.AutoAllowNew[operation]
}

func (m *managerImpl) ValidateAccountAccess(ctx context.Context, accountID string, user *users.User, allowOwnerAndAdmin bool) error {
	if user.AccountID != accountID {
		return status.NewUserNotPartOfAccountError()
	}
	return nil
}

type UserAuthExtractor func(ctx context.Context) (*nbcontext.UserAuth, error)

type PermissionValidator interface {
	ValidateUserPermissions(ctx context.Context, accountID, userID string, module modules.Module, operation operations.Operation) (bool, error)
}

// WithPermission wraps an HTTP handler with permission checking logic.
func (m *managerImpl) WithPermission(
	module modules.Module,
	operation operations.Operation,
	handlerFunc func(w http.ResponseWriter, r *http.Request, auth *nbcontext.UserAuth),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
		if err != nil {
			log.WithContext(r.Context()).Errorf("failed to get user auth from context: %v", err)
			util.WriteError(r.Context(), err, w)
			return
		}

		allowed, err := m.ValidateUserPermissions(r.Context(), userAuth.AccountId, userAuth.UserId, module, operation)
		if err != nil {
			log.WithContext(r.Context()).Errorf("failed to validate permissions for user %s on account %s: %v", userAuth.UserId, userAuth.AccountId, err)
			util.WriteError(r.Context(), status.NewPermissionValidationError(err), w)
			return
		}

		if !allowed {
			log.WithContext(r.Context()).Tracef("user %s on account %s is not allowed to %s in %s", userAuth.UserId, userAuth.AccountId, operation, module)
			util.WriteError(r.Context(), status.NewPermissionDeniedError(), w)
			return
		}

		handlerFunc(w, r, &userAuth)
	}
}
