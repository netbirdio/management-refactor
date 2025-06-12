package manager

import (
	"context"
	"fmt"

	"github.com/rs/xid"

	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/modules/networks/resources"
	"github.com/netbirdio/management-refactor/internals/modules/networks/routers"
	"github.com/netbirdio/management-refactor/internals/shared/activity"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type managerImpl struct {
	repo            Repository
	activityManager activity.Manager
	resourceManager resources.Manager
	routerManager   routers.Manager
}

func NewManager(repo Repository, activityManager activity.Manager, resourceManager resources.Manager, routerManager routers.Manager) networks.Manager {
	return &managerImpl{
		repo:            repo,
		activityManager: activityManager,
		resourceManager: resourceManager,
		routerManager:   routerManager,
	}
}

func (m *managerImpl) GetAllNetworks(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, userID string) ([]*networks.Network, error) {
	return m.repo.GetAccountNetworks(tx, strength, accountID)
}

func (m *managerImpl) CreateNetwork(ctx context.Context, tx db.Transaction, userID string, network *networks.Network) (*networks.Network, error) {
	network.ID = xid.New().String()

	err := m.repo.CreateNetwork(tx, network)
	if err != nil {
		return nil, fmt.Errorf("failed to save network: %w", err)
	}

	// m.accountManager.StoreEvent(ctx, userID, network.ID, network.AccountID, activity.NetworkCreated, network.EventMeta())

	return network, nil
}

func (m *managerImpl) GetNetwork(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, userID, networkID string) (*networks.Network, error) {
	return m.repo.GetNetworkByID(tx, strength, accountID, networkID)
}

func (m *managerImpl) UpdateNetwork(ctx context.Context, tx db.Transaction, userID string, network *networks.Network) (*networks.Network, error) {
	_, err := m.repo.GetNetworkByID(tx, db.LockingStrengthUpdate, network.AccountID, network.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get network: %w", err)
	}

	// m.accountManager.StoreEvent(ctx, userID, network.ID, network.AccountID, activity.NetworkUpdated, network.EventMeta())

	return network, m.repo.UpdateNetwork(tx, network)
}

func (m *managerImpl) DeleteNetwork(ctx context.Context, tx db.Transaction, accountID, userID, networkID string) error {
	return db.WithTx(m.repo.Store(), tx, func(tx db.Transaction) error {

		err := m.routerManager.DeleteRoutersInNetwork(ctx, tx, accountID, userID, networkID)
		if err != nil {
			return fmt.Errorf("failed to delete routers in network: %w", err)
		}

		err = m.resourceManager.DeleteResourcesInNetwork(ctx, tx, accountID, userID, networkID)
		if err != nil {
			return fmt.Errorf("failed to delete resources in network: %w", err)
		}

		network := &networks.Network{ID: networkID}
		err = m.repo.DeleteNetwork(tx, network)
		if err != nil {
			return err
		}

		m.activityManager.StoreEvent(ctx, tx, userID, networkID, accountID, activity.NetworkDeleted, network.EventMeta())

		return err
	})
}
