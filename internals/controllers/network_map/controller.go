package network_map

import "context"

type Controller interface {
	UpdatePeers(ctx context.Context, accountID string) error
}
