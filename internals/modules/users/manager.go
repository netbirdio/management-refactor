package users

import (
	"context"

	"github.com/gorilla/mux"

	"github.com/netbirdio/management-refactor/internals/modules/users/types"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
	"github.com/netbirdio/management-refactor/pkg/logging"
)

var log = logging.LoggerForThisPackage()

type Manager struct {
	repo Repository
}

func NewManager(store *db.Store, router *mux.Router, permissionsManager permissions.Manager) *Manager {
	repo := newRepository(store)
	m := &Manager{repo: repo}
	api := newHandler(m, permissionsManager)
	api.RegisterEndpoints(router)
	return m
}

func (m *Manager) GetAllUsers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]types.User, error) {
	return m.repo.GetAllUsers(tx, strength, accountID)
}

func (m *Manager) GetUserByID(ctx context.Context, tx db.Transaction, strength db.LockingStrength, id string) (*types.User, error) {
	return m.repo.GetUserByID(tx, strength, id)
}
