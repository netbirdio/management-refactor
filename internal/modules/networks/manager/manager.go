package manager

import (
	"context"
	"fmt"

	"github.com/netbirdio/netbird/management/server/account"
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

func NewManager(store *db.Store, permissionsManager permissions.Manager, resourceManager resources.Manager, routersManager routers.Manager) resources.Manager {
	repo := newRepository(store)
	m := &managerImpl{repo: repo}
	// api := newHandler(m, permissionsManager)
	// api.RegisterEndpoints(router)
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

	network, err := m.repo.GetNetworkByID(tx, db.LockingStrengthUpdate, accountID, networkID)
	if err != nil {
		return fmt.Errorf("failed to get network: %w", err)
	}

	resources, err := m.resourcesManager.GetNetworkResourcesByNetID(ctx, tx, db.LockingStrengthUpdate, accountID, userID, networkID)
	if err != nil {
		return fmt.Errorf("failed to get resources in network: %w", err)
	}

	for _, resource := range resources {
		err = m.resourcesManager.DeleteResource(ctx, tx, accountID, userID, networkID, resource.ID)
		if err != nil {
			return fmt.Errorf("failed to delete resource: %w", err)
		}
	}

	routers, err := m.routersManager.GetNetworkRoutersByNetID(ctx, tx, db.LockingStrengthUpdate, accountID, userID, networkID)
	if err != nil {
		return fmt.Errorf("failed to get routers in network: %w", err)
	}

	for _, router := range routers {
		err := m.routersManager.DeleteRouter(ctx, tx, accountID, userID, networkID, router.ID)
		if err != nil {
			return fmt.Errorf("failed to delete router: %w", err)
		}
	}

	err = m.repo.DeleteNetwork(tx, &networks.Network{ID: networkID})
	if err != nil {
		return fmt.Errorf("failed to delete network: %w", err)
	}

	// TODO: send delete event with network

	return nil
}
