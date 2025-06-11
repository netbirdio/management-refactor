package manager

import (
	"context"

	"github.com/netbirdio/management-refactor/internals/modules/peers"
	"github.com/netbirdio/management-refactor/internals/shared/activity"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/pkg/logging"
)

var log = logging.LoggerForThisPackage()

var _ peers.Manager = (*Manager)(nil)

type Manager struct {
	repo         Repository
	eventManager *activity.Manager
}

func NewManager(store *db.Store) *Manager {
	return &Manager{repo: newRepository(store)}

}

func (m *Manager) GetPeer(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, peerID string) (*peers.Peer, error) {
	return m.repo.GetPeerByID(tx, strength, accountID, peerID)
}

func (m *Manager) GetPeers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*peers.Peer, error) {
	return m.repo.GetPeers(tx, strength, accountID)
}

func (m *Manager) GetFilteredPeers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, nameFilter, ipFilter string) ([]*peers.Peer, error) {
	return m.repo.GetFilteredPeers(tx, strength, accountID, nameFilter, ipFilter)
}

func (m *Manager) UpdatePeer(ctx context.Context, tx db.Transaction, peer *peers.Peer) error {
	return m.repo.UpdatePeer(tx, peer)
}
