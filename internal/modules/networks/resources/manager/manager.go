package manager

import (
	"context"

	"github.com/gorilla/mux"

	"management/internal/modules/networks"
	"management/internal/modules/networks/resources"
	"management/internal/modules/networks/routers"
	"management/internal/shared/db"
	"management/internal/shared/errors"
	"management/internal/shared/permissions"
	"management/internal/shared/permissions/modules"
	"management/internal/shared/permissions/operations"
)

type managerImpl struct {
	repo               Repository
	permissionsManager permissions.Manager
	resourcesManager   resources.Manager
	routersManager     routers.Manager
}

func NewManager(store *db.Store, router *mux.Router, permissionsManager permissions.Manager, resourceManager resources.Manager, routersManager routers.Manager) resources.Manager {
	repo := newRepository(store)
	m := &managerImpl{
		repo:               repo,
		permissionsManager: permissionsManager,
		resourcesManager:   resourceManager,
		routersManager:     routersManager,
	}
	// api := newHandler(m, permissionsManager)
	// api.RegisterEndpoints(router)
	return m
}

func (m *managerImpl) DeleteResourcesInNetwork(ctx context.Context, tx db.Transaction, accountID, userID, networkID string) error {
	ok, err := m.permissionsManager.ValidateUserPermissions(ctx, accountID, userID, modules.Networks, operations.Read)
	if err != nil {
		return errors.NewPermissionValidationError(err)
	}
	if !ok {
		return errors.NewPermissionDeniedError()
	}

	resources, err := m.repo.GetResourcesByNetworkID(tx, db.LockingStrengthUpdate, networkID)
	if err != nil {
		return err
	}
	for _, resource := range resources {
		err = m.repo.DeleteResource(tx, resource.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
