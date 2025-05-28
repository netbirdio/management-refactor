package resources

import (
	"context"

	"management/internal/shared/db"
)

type Manager interface {
	// Create

	// Read
	GetNetworkResourcesByNetID(ctx context.Context, tx db.Transaction, lockingStrength db.LockingStrength, accountID, userID, networkID string) ([]*NetworkResource, error)

	// Update

	// Delete
	DeleteResource(ctx context.Context, tx db.Transaction, accountID, userID, networkID, resourceID string) error
	DeleteResourcesInNetwork(ctx context.Context, tx db.Transaction, accountID, userID, networkID string) error
}
