package manager

import (
	"context"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/rs/xid"

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

func NewManager(store *db.Store, router *mux.Router, permissionsManager permissions.Manager, resourceManager resources.Manager, routersManager routers.Manager) networks.Manager {
	repo := newRepository(store)
	m := &managerImpl{
		repo:               repo,
		permissionsManager: permissionsManager,
		resourcesManager:   resourceManager,
		routersManager:     routersManager,
	}
	api := newHandler(m, permissionsManager)
	api.RegisterEndpoints(router)
	return m
}

func (m *managerImpl) GetAllNetworks(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, userID string) ([]*networks.Network, error) {
	ok, err := m.permissionsManager.ValidateUserPermissions(ctx, accountID, userID, modules.Networks, operations.Read)
	if err != nil {
		return nil, errors.NewPermissionValidationError(err)
	}
	if !ok {
		return nil, errors.NewPermissionDeniedError()
	}

	return m.repo.GetAccountNetworks(tx, strength, accountID)
}

func (m *managerImpl) CreateNetwork(ctx context.Context, tx db.Transaction, userID string, network *networks.Network) (*networks.Network, error) {
	ok, err := m.permissionsManager.ValidateUserPermissions(ctx, network.AccountID, userID, modules.Networks, operations.Write)
	if err != nil {
		return nil, errors.NewPermissionValidationError(err)
	}
	if !ok {
		return nil, errors.NewPermissionDeniedError()
	}

	network.ID = xid.New().String()

	err = m.repo.CreateNetwork(tx, network)
	if err != nil {
		return nil, fmt.Errorf("failed to save network: %w", err)
	}

	// m.accountManager.StoreEvent(ctx, userID, network.ID, network.AccountID, activity.NetworkCreated, network.EventMeta())

	return network, nil
}

func (m *managerImpl) GetNetwork(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, userID, networkID string) (*networks.Network, error) {
	ok, err := m.permissionsManager.ValidateUserPermissions(ctx, accountID, userID, modules.Networks, operations.Read)
	if err != nil {
		return nil, errors.NewPermissionValidationError(err)
	}
	if !ok {
		return nil, errors.NewPermissionDeniedError()
	}

	return m.repo.GetNetworkByID(tx, strength, accountID, networkID)
}

func (m *managerImpl) UpdateNetwork(ctx context.Context, tx db.Transaction, userID string, network *networks.Network) (*networks.Network, error) {
	ok, err := m.permissionsManager.ValidateUserPermissions(ctx, network.AccountID, userID, modules.Networks, operations.Write)
	if err != nil {
		return nil, errors.NewPermissionValidationError(err)
	}
	if !ok {
		return nil, errors.NewPermissionDeniedError()
	}

	_, err = m.repo.GetNetworkByID(tx, db.LockingStrengthUpdate, network.AccountID, network.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network: %w", err)
	}

	// m.accountManager.StoreEvent(ctx, userID, network.ID, network.AccountID, activity.NetworkUpdated, network.EventMeta())

	return network, m.repo.UpdateNetwork(tx, network)
}

func (m *managerImpl) DeleteNetwork(ctx context.Context, tx db.Transaction, accountID, userID, networkID string) error {
	ok, err := m.permissionsManager.ValidateUserPermissions(ctx, accountID, userID, modules.Networks, operations.Write)
	if err != nil {
		return errors.NewPermissionValidationError(err)
	}
	if !ok {
		return errors.NewPermissionDeniedError()
	}

	return db.WithTx(m.repo.Store(), tx, func(tx db.Transaction) error {

		sendEvent(tx, "network_deleted")

		if err := m.repo.DeleteNetwork(tx, &networks.Network{ID: networkID}); err != nil {
			return fmt.Errorf("failed to delete network: %w", err)
		}

		tx.AddEvent(func() {
			addActivityEvent("Network deleted")
			// noop
		})

		return nil
	})
}
