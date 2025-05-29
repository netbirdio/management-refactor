package manager

import (
	"context"
	"fmt"

	"github.com/gorilla/mux"

	"management/internal/modules/networks"
	"management/internal/modules/networks/resources"
	"management/internal/shared/db"
)

type managerImpl struct {
	repo           Repository
	networkManager networks.Manager
}

func NewManager(store *db.Store, router *mux.Router, networkManager networks.Manager) resources.Manager {
	repo := newRepository(store)
	m := &managerImpl{
		repo:           repo,
		networkManager: networkManager,
	}

	networkManager.OnNetworkDelete().BindFunc(func(e *networks.NetworkEvent) error {
		if err := m.DeleteResourcesInNetwork(e.Context, e.Tx, e.Network); err != nil {
			return fmt.Errorf("failed to delete resources in network: %w", err)
		}

		return e.Next()
	})

	// api := newHandler(m, permissionsManager)
	// api.RegisterEndpoints(router)
	return m
}

func (m *managerImpl) GetNetworkResourcesByNetID(ctx context.Context, tx db.Transaction, lockingStrength db.LockingStrength, network *networks.Network) ([]*resources.NetworkResource, error) {
	return m.repo.GetResourcesByNetworkID(tx, lockingStrength, network.ID)
}

func (m *managerImpl) DeleteResourcesInNetwork(ctx context.Context, tx db.Transaction, network *networks.Network) error {
	resources, err := m.GetNetworkResourcesByNetID(ctx, tx, db.LockingStrengthUpdate, network)
	if err != nil {
		return err
	}

	for _, resource := range resources {
		err = m.DeleteResource(ctx, tx, resource)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *managerImpl) DeleteResource(ctx context.Context, tx db.Transaction, resource *resources.NetworkResource) error {
	if err := m.repo.DeleteResource(tx, resource); err != nil {
		return fmt.Errorf("failed to delete network: %w", err)
	}

	tx.AddEvent(func() {
		// addActivityEvent("resource deleted")
		// noop
	})
	return nil
}
