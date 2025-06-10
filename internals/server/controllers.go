package server

import (
	"github.com/netbirdio/management-refactor/internals/controllers/network_map"
	"github.com/netbirdio/management-refactor/internals/controllers/network_map/controller"
)

func (s *BaseServer) NetworkMapController() network_map.Controller {
	return Create(s, func() network_map.Controller {
		store := s.Store()
		metrics := s.Metrics()
		return controller.NewController(store, metrics)
	})
}
