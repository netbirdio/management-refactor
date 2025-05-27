package resources

import (
	"context"

	"management/internal/shared/db"
)

type Manager interface {
	DeleteResource(ctx context.Context, tx db.Transaction, accountID, userID, networkID, resourceID string) error
	GetNetworkResourcesByNetID(ctx context.Context, tx db.Transaction, lockingStrength db.LockingStrength, accountID, userID, networkID string) ([]*NetworkResource, error)
}
