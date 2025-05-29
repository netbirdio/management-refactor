package manager

import (
	"context"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/rs/xid"

	"management/internal/modules/networks"
	"management/internal/shared/db"
	"management/internal/shared/hook"
	"management/internal/shared/permissions"
)

type managerImpl struct {
	repo Repository

	onNetworkDelete *hook.Hook[*networks.NetworkEvent]
}

func NewManager(store *db.Store, router *mux.Router, permissionsManager permissions.Manager) networks.Manager {
	repo := newRepository(store)
	m := &managerImpl{
		repo: repo,

		onNetworkDelete: &hook.Hook[*networks.NetworkEvent]{},
	}
	api := newHandler(m, permissionsManager)
	api.RegisterEndpoints(router)
	return m
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
		network := &networks.Network{ID: networkID}

		ev := &networks.NetworkEvent{
			Context: ctx,
			Tx:      tx,
			Network: network,
		}

		err := m.OnNetworkDelete().Trigger(ev, func(ne *networks.NetworkEvent) error {
			if err := m.repo.DeleteNetwork(ne.Tx, ne.Network); err != nil {
				return fmt.Errorf("failed to delete network: %w", err)
			}

			tx.AddEvent(func() {
				// addActivityEvent("Network deleted")
				// noop
			})
			return nil
		})

		return err
	})
}
