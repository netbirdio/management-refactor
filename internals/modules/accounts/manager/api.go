package manager

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	nbcontext "github.com/netbirdio/management-refactor/management/server/context"
	"github.com/netbirdio/management-refactor/management/server/http/util"

	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/internals/shared/errors"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/modules"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

type handler struct {
	manager            *Manager
	permissionsManager permissions.Manager
}

func newHandler(manager *Manager, permissionsManager permissions.Manager) *handler {
	return &handler{
		manager:            manager,
		permissionsManager: permissionsManager,
	}
}

func (h *handler) RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/accounts/{accountId}", h.updateAccount).Methods("PUT", "OPTIONS")
	router.HandleFunc("/accounts/{accountId}", h.deleteAccount).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/accounts", h.getAllAccounts).Methods("GET", "OPTIONS")
}

func (h *handler) updateAccount(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}

	allowed, err := h.permissionsManager.ValidateUserPermissions(r.Context(), userAuth.AccountId, userAuth.UserId, modules.Accounts, operations.Write)
	if err != nil {
		util.WriteError(r.Context(), errors.NewPermissionValidationError(err), w)
		return
	}
	if !allowed {
		util.WriteError(r.Context(), errors.NewPermissionDeniedError(), w)
	}

	vars := mux.Vars(r)
	accountId := vars["accountId"]

	users, err := h.manager.UpdateAccount(r.Context(), nil, accountId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(users)
}

func (h *handler) deleteAccount(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}

	allowed, err := h.permissionsManager.ValidateUserPermissions(r.Context(), userAuth.AccountId, userAuth.UserId, modules.Accounts, operations.Write)
	if err != nil {
		util.WriteError(r.Context(), errors.NewPermissionValidationError(err), w)
		return
	}
	if !allowed {
		util.WriteError(r.Context(), errors.NewPermissionDeniedError(), w)
	}

	vars := mux.Vars(r)
	accountId := vars["accountId"]

	user, err := h.manager.DeleteAccount(r.Context(), nil, db.LockingStrengthShare, accountId)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(user)
}

func (h *handler) getAllAccounts(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}

	allowed, err := h.permissionsManager.ValidateUserPermissions(r.Context(), userAuth.AccountId, userAuth.UserId, modules.Accounts, operations.Read)
	if err != nil {
		util.WriteError(r.Context(), errors.NewPermissionValidationError(err), w)
		return
	}
	if !allowed {
		util.WriteError(r.Context(), errors.NewPermissionDeniedError(), w)
	}

	accounts, err := h.manager.GetAllAccounts(r.Context(), nil, db.LockingStrengthShare)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(accounts)
}
