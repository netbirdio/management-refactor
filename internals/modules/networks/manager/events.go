package manager

import (
	"management/internal/modules/networks"
	"management/internal/shared/hook"
)

func (m *managerImpl) OnNetworkDelete() *hook.Hook[*networks.NetworkEvent] {
	return m.onNetworkDelete
}
