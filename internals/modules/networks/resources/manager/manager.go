package manager

import (
	"context"
	"fmt"

	"github.com/gorilla/mux"

	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/modules/networks/resources"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type managerImpl struct {
	repo           Repository
	networkManager networks.Manager
}

func NewManager(repo Repository, router *mux.Router, networkManager networks.Manager) resources.Manager {
	return &managerImpl{
		repo:           repo,
		networkManager: networkManager,
	}
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
