package networks

import (
	"context"

	"management/internal/shared/db"
)

type Manager interface {
	Using(tx db.Transaction) Manager
	GetAllNetworks(ctx context.Context, strength db.LockingStrength, accountID, userID string) ([]*Network, error)
	GetNetwork(ctx context.Context, strength db.LockingStrength, accountID, userID, networkID string) (*Network, error)
	CreateNetwork(ctx context.Context, userID string, network *Network) (*Network, error)
	UpdateNetwork(ctx context.Context, userID string, network *Network) (*Network, error)
	DeleteNetwork(ctx context.Context, accountID, userID, networkID string) error
}
