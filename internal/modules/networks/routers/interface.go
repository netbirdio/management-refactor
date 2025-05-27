package routers

import (
	"context"

	"management/internal/shared/db"
)

type Manager interface {
	Using(tx db.Transaction) Manager
	GetNetworkRoutersByNetID(ctx context.Context, tx db.Transaction, lockingStrength db.LockingStrength, accountID, userID, networkID string) ([]*NetworkRouter, error)
	DeleteRouter(ctx context.Context, tx db.Transaction, accountID, userID, networkID, routerID string) error
}
