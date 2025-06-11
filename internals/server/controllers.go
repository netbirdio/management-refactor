package server

import (
	"github.com/netbirdio/management-refactor/internals/controllers/network_map"
	"github.com/netbirdio/management-refactor/internals/controllers/network_map/controller"
	"github.com/netbirdio/management-refactor/internals/controllers/network_map/update_channel"
)

func (s *BaseServer) NetworkMapController() network_map.Controller {
	return Create(s, func() network_map.Controller {
		return controller.NewController(s.Store(), s.Metrics(), s.NetworkMapUpdateChannel())
	})
}

func (s *BaseServer) NetworkMapUpdateChannel() network_map.UpdateChannel {
	return Create(s, func() network_map.UpdateChannel {
		return update_channel.NewUpdateChannel(s.Metrics())
	})
}
