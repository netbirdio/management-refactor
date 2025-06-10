package network_map

import (
	"github.com/netbirdio/netbird/management/proto"
	"github.com/netbirdio/netbird/management/server/groups"
	"github.com/netbirdio/netbird/management/server/types"

	"github.com/netbirdio/management-refactor/internals/modules/peers"
	"github.com/netbirdio/management-refactor/internals/modules/policies"
)

type UpdateMessage struct {
	Update        *proto.SyncResponse
	NetworkMap    *types.NetworkMap
	PeerManager   *peers.Manager
	PolicyManager *policies.Manager
	GroupManager  *groups.Manager
}
