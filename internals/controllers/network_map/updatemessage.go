package network_map

import (
	"github.com/netbirdio/netbird/management/proto"
	"github.com/netbirdio/netbird/management/server/types"
)

type UpdateMessage struct {
	Update     *proto.SyncResponse
	NetworkMap *types.NetworkMap
}
