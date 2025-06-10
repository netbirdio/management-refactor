package server

import (
	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/modules/networks/manager"
	"github.com/netbirdio/management-refactor/internals/modules/networks/resources"
	resourcesManager "github.com/netbirdio/management-refactor/internals/modules/networks/resources/manager"
	"github.com/netbirdio/management-refactor/internals/modules/peers"
	peersManager "github.com/netbirdio/management-refactor/internals/modules/peers/manager"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
)

func (s *BaseServer) NetworksManager() networks.Manager {
	return Create(s, func() networks.Manager {
		return manager.NewManager(s.Store(), s.Router(), s.PermissionsManager())
	})
}

func (s *BaseServer) ResourcesManager() resources.Manager {
	return Create(s, func() resources.Manager {
		return resourcesManager.NewManager(s.Store(), s.Router(), s.NetworksManager())
	})
}

func (s *BaseServer) PermissionsManager() permissions.Manager {
	return Create(s, func() permissions.Manager {
		return permissions.NewManager()
	})
}

func (s *BaseServer) PeersManager() peers.Manager {
	return Create(s, func() peers.Manager {
		store := s.Store()
		router := s.Router()
		permissionsManager := s.PermissionsManager()

		return peersManager.NewManager(store, router, permissionsManager)
	})
}
