package resources

import (
	"context"

	"management/internal/modules/networks"
	"management/internal/shared/db"
)

type Manager interface {
	// Create

	// Read
	GetNetworkResourcesByNetID(ctx context.Context, tx db.Transaction, lockingStrength db.LockingStrength, network *networks.Network) ([]*NetworkResource, error)

	// Update

	// Delete
	DeleteResource(ctx context.Context, tx db.Transaction, resource *NetworkResource) error
	DeleteResourcesInNetwork(ctx context.Context, tx db.Transaction, network *networks.Network) error
}
