package resources

import (
	"context"

	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Manager interface {
	// Create

	// Read
	GetNetworkResourcesByNetID(ctx context.Context, tx db.Transaction, lockingStrength db.LockingStrength, network *networks.Network) ([]*NetworkResource, error)

	// Update

	// Delete
	DeleteResource(ctx context.Context, tx db.Transaction, resource *NetworkResource) error
	DeleteResourcesInNetwork(ctx context.Context, tx db.Transaction, accountID, userID, networkID string) error
}
