package network_map

import "context"

type UpdateChannel interface {
	UpdatePeers(accountID string) error
	SendUpdate(ctx context.Context, peerID string, update *UpdateMessage)
}
