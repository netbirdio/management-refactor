package manager

import (
	"github.com/netbirdio/management-refactor/internals/modules/peers"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Repository interface {
	RunInTx(fn func(tx db.Transaction) error) error
	GetPeerByID(tx db.Transaction, strength db.LockingStrength, accountID, peerId string) (*peers.Peer, error)
	GetPeers(tx db.Transaction, strength db.LockingStrength, accountID string) ([]*peers.Peer, error)
	GetFilteredPeers(tx db.Transaction, strength db.LockingStrength, accountID string, nameFilter, ipFilter string) ([]*peers.Peer, error)
	UpdatePeer(tx db.Transaction, peer *peers.Peer) error
}

type repository struct {
	store *db.Store
}

func newRepository(s *db.Store) Repository {
	return &repository{store: s}
}

func (r *repository) RunInTx(fn func(tx db.Transaction) error) error {
	return r.store.RunInTx(fn)
}

func (r *repository) GetPeerByID(tx db.Transaction, strength db.LockingStrength, accountID, peerId string) (*peers.Peer, error) {
	var peer peers.Peer
	err := r.store.GetOne(tx, strength, &peer, "account_id = ? AND id = ?", accountID, peerId)
	if err != nil {
		return nil, err
	}
	return &peer, nil
}

func (r *repository) GetPeers(tx db.Transaction, strength db.LockingStrength, accountID string) ([]*peers.Peer, error) {
	var peers []*peers.Peer
	err := r.store.GetMany(tx, strength, &peers, "account_id = ?", accountID)
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func (r *repository) GetFilteredPeers(tx db.Transaction, strength db.LockingStrength, accountID string, nameFilter, ipFilter string) ([]*peers.Peer, error) {
	query := "account_id = ?"
	args := []interface{}{accountID}

	if nameFilter != "" {
		query += " AND name LIKE ?"
		args = append(args, nameFilter)
	}

	if ipFilter != "" {
		query += " AND ip LIKE ?"
		args = append(args, ipFilter)
	}

	var peers []*peers.Peer
	err := r.store.GetMany(tx, strength, &peers, query, args)
	if err != nil {
		return nil, err
	}
	return peers, nil
}

func (r *repository) UpdatePeer(tx db.Transaction, peer *peers.Peer) error {
	return r.store.Update(tx, peer)
}
