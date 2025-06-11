package manager

import (
	"context"

	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/pkg/logging"
)

var log = logging.LoggerForThisPackage()

type Manager struct {
	repo Repository
}

func NewManager(store *db.Store) *Manager {
	repo := newRepository(store)
	m := &Manager{repo: repo}
	return m
}

func (m *Manager) GetAllUsers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]users.User, error) {
	return m.repo.GetAllUsers(tx, strength, accountID)
}

func (m *Manager) GetUserByID(ctx context.Context, tx db.Transaction, strength db.LockingStrength, id string) (*users.User, error) {
	return m.repo.GetUserByID(tx, strength, id)
}
