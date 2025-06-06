package peers

import (
	"context"

	"github.com/gorilla/mux"

	"github.com/netbirdio/management-refactor/internals/modules/peers/types"
	"github.com/netbirdio/management-refactor/internals/shared/activity"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
	"github.com/netbirdio/management-refactor/pkg/logging"
)

var log = logging.LoggerForThisPackage()

type Manager struct {
	repo         Repository
	eventManager *activity.Manager
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

func (m *Manager) UpdatePeer(ctx context.Context, tx db.Transaction, peer *types.Peer) error {
	validateInput
	validatePermissions
	err := m.repo.RunInTx(func(tx db.Transaction) error {
		othermanager.UpdatePeers
		ourmanager.UpdateGroup
	})
	if err != nil {
		return err
	}

	err := sendPeerUpdateEvent(tx, peer) // -> goes to peerUpdtaeChannel
	if err != nil {
		log.Errorf("Failed to send peer update event: %v", err)
	}

	return m.repo.UpdatePeer(tx, peer)
}
