package manager

import (
	"context"

	"github.com/netbirdio/management-refactor/internals/modules/networks/routers"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type managerImpl struct {
	repo Repository
}

func NewManager(repo Repository) routers.Manager {
	return &managerImpl{
		repo: repo,
	}
}

func (m managerImpl) GetNetworkRoutersByNetID(ctx context.Context, tx db.Transaction, lockingStrength db.LockingStrength, accountID, userID, networkID string) ([]*routers.NetworkRouter, error) {
	// TODO implement me
	panic("implement me")
}

func (m managerImpl) DeleteRouter(ctx context.Context, tx db.Transaction, accountID, userID, networkID, routerID string) error {
	// TODO implement me
	panic("implement me")
}

func (m managerImpl) DeleteRoutersInNetwork(ctx context.Context, tx db.Transaction, accountID, userID, networkID string) error {
	// TODO implement me
	panic("implement me")
}
