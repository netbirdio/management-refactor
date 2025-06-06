package server

import (
	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/modules/networks/manager"
	"github.com/netbirdio/management-refactor/internals/modules/networks/resources"
	resourcesManager "github.com/netbirdio/management-refactor/internals/modules/networks/resources/manager"
	"github.com/netbirdio/management-refactor/internals/modules/peers"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
)

func (s *server) NetworksManager() networks.Manager {
	return Create(s, func() networks.Manager {
		return manager.NewManager(s.Store(), s.Router(), s.PermissionsManager())
	})
}

func (s *server) ResourcesManager() resources.Manager {
	return Create(s, func() resources.Manager {
		return resourcesManager.NewManager(s.Store(), s.Router(), s.NetworksManager())
	})
}

func (s *server) PermissionsManager() permissions.Manager {
	return Create(s, func() permissions.Manager {
		return permissions.NewManager()
	})
}

func (s *server) PeersManager() *peers.Manager {
	return Create(s, func() *peers.Manager {
		store := s.Store()
		router := s.Router()
		permissionsManager := s.PermissionsManager()

		return peers.NewManager(store, router, permissionsManager)
	})
}
