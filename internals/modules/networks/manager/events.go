package manager

import (
	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/shared/hook"
)

func (m *managerImpl) OnNetworkDelete() *hook.Hook[*networks.NetworkEvent] {
	return m.onNetworkDelete
}
