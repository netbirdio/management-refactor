package network_map

import (
	"github.com/netbirdio/netbird/management/server/types"
	"github.com/netbirdio/netbird/shared/management/proto"
)

type UpdateMessage struct {
	Update     *proto.SyncResponse
	NetworkMap *types.NetworkMap
}
