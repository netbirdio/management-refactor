package peers

import (
	"context"

	"github.com/gorilla/mux"

	"management/internal/modules/peers/types"
	"management/internal/shared/db"
	"management/internal/shared/permissions"
	"management/pkg/logging"
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

func (m *Manager) GetPeer(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, peerID string) (*types.Peer, error) {
	return m.repo.GetPeerByID(tx, strength, accountID, peerID)
}

func (m *Manager) GetPeers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*types.Peer, error) {
	return m.repo.GetPeers(tx, strength, accountID)
}

func (m *Manager) GetFilteredPeers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, nameFilter, ipFilter string) ([]*types.Peer, error) {
	return m.repo.GetFilteredPeers(tx, strength, accountID, nameFilter, ipFilter)
}

import (
	"context"

	"github.com/gorilla/mux"

	"management/internal/modules/peers/types"
	"management/internal/shared/db"
	"management/internal/shared/permissions"
	"management/pkg/logging"
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

func (m *Manager) GetPeer(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, peerID string) (*types.Peer, error) {
	return m.repo.GetPeerByID(tx, strength, accountID, peerID)
}

func (m *Manager) GetPeers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*types.Peer, error) {
	return m.repo.GetPeers(tx, strength, accountID)
}

func (m *Manager) GetFilteredPeers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, nameFilter, ipFilter string) ([]*types.Peer, error) {
	return m.repo.GetFilteredPeers(tx, strength, accountID, nameFilter, ipFilter)
}

import (
	"context"

	"github.com/gorilla/mux"

	"management/internal/modules/peers/types"
	"management/internal/shared/db"
	"management/internal/shared/permissions"
	"management/pkg/logging"
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

func (m *Manager) GetPeer(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, peerID string) (*types.Peer, error) {
	return m.repo.GetPeerByID(tx, strength, accountID, peerID)
}

func (m *Manager) GetPeers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*types.Peer, error) {
	return m.repo.GetPeers(tx, strength, accountID)
}

func (m *Manager) GetFilteredPeers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, nameFilter, ipFilter string) ([]*types.Peer, error) {
	return m.repo.GetFilteredPeers(tx, strength, accountID, nameFilter, ipFilter)
}
