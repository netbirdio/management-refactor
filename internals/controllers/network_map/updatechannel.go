package network_map

import "context"

type UpdateChannel interface {
	SendUpdate(ctx context.Context, peerID string, update *UpdateMessage)
}
