package networks

import (
	"context"

	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Manager interface {
	GetAllNetworks(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, userID string) ([]*Network, error)
	GetNetwork(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, userID, networkID string) (*Network, error)
	CreateNetwork(ctx context.Context, tx db.Transaction, userID string, network *Network) (*Network, error)
	UpdateNetwork(ctx context.Context, tx db.Transaction, userID string, network *Network) (*Network, error)
	DeleteNetwork(ctx context.Context, tx db.Transaction, accountID, userID, networkID string) error
}
