package users

import (
	"context"

	"management/internal/modules/users/types"
	"management/internal/shared/db"
	"management/internal/shared/permissions"
	"management/pkg/logging"
)

var log = logging.LoggerForThisPackage()

type Manager struct {
	repo    Repository
	handler *handler
}

func NewManager(store *db.Store, permissionsManager permissions.Manager) *Manager {
	repo := newRepository(store)
	m := &Manager{repo: repo}
	m.handler = newHandler(m, permissionsManager)
	return m
}

func (m *Manager) GetAllUsers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]types.User, error) {
	return m.repo.GetAllUsers(tx, strength, accountID)
}

func (m *Manager) GetUserByID(ctx context.Context, tx db.Transaction, strength db.LockingStrength, id string) (*types.User, error) {
	return m.repo.GetUserByID(tx, strength, id)
}
