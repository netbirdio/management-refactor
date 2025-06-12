package manager

//go:generate go run github.com/golang/mock/mockgen -package manager -destination=manager_mock.go -source=../interface.go -build_flags=-mod=mod

import (
	"context"

	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Manager struct {
	repo Repository
}

func NewManager(repo Repository) *Manager {
	return &Manager{repo: repo}
}

func (m *Manager) GetAllUsers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]users.User, error) {
	return m.repo.GetAllUsers(tx, strength, accountID)
}

func (m *Manager) GetUserByID(ctx context.Context, tx db.Transaction, strength db.LockingStrength, id string) (*users.User, error) {
	user, err := m.repo.GetUserByID(tx, strength, id)
	if err != nil {
		return nil, err
	}

	user.Issued = "****"

	return user, nil
}
