package manager

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/netbirdio/management-refactor/internals/controllers/network_map"
	"github.com/netbirdio/management-refactor/internals/modules/peers"
	"github.com/netbirdio/management-refactor/internals/shared/activity"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

var _ peers.Manager = (*Manager)(nil)

type Manager struct {
	repo                 Repository
	eventManager         *activity.Manager
	networkMapController network_map.Controller
}

func NewManager(repo Repository) *Manager {
	return &Manager{repo: repo}
}

func (m *Manager) SetNetworkMapController(networkMapController network_map.Controller) {
	log.Tracef("Setting network map controller for peers manager")
	m.networkMapController = networkMapController
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
	// err := m.repo.UpdatePeer(tx, peer)
	// if err != nil {
	// 	return fmt.Errorf("failed to update peer: %w", err)
	// }

	_ = m.networkMapController.UpdatePeers(ctx, peer.AccountID)

	return nil
}

func (m *Manager) GetAllEphemeralPeers(ctx context.Context, tx db.Transaction, strength db.LockingStrength) ([]*peers.Peer, error) {
	// TODO implement me
	panic("implement me")
}

func (m *Manager) DeletePeer(ctx context.Context, tx db.Transaction, accountID, peerID string) error {
	// TODO implement me
	panic("implement me")
}
