package peers

import (
	"context"

	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Manager interface {
	GetPeer(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, peerID string) (*Peer, error)
	GetPeers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*Peer, error)
	GetFilteredPeers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, nameFilter, ipFilter string) ([]*Peer, error)
	UpdatePeer(ctx context.Context, tx db.Transaction, peer *Peer) error
}
