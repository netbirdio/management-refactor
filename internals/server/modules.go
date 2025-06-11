package server

import (
	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/modules/networks/manager"
	"github.com/netbirdio/management-refactor/internals/modules/networks/resources"
	resourcesManager "github.com/netbirdio/management-refactor/internals/modules/networks/resources/manager"
	"github.com/netbirdio/management-refactor/internals/modules/peers"
	peersManager "github.com/netbirdio/management-refactor/internals/modules/peers/manager"
	"github.com/netbirdio/management-refactor/internals/modules/users"
	usersManager "github.com/netbirdio/management-refactor/internals/modules/users/manager"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
)

func (s *BaseServer) NetworksManager() networks.Manager {
	return Create(s, func() networks.Manager {
		return manager.NewManager(s.Store(), s.Router(), s.PermissionsManager())
	})
}

func (s *BaseServer) ResourcesManager() resources.Manager {
	return Create(s, func() resources.Manager {
		manager := resourcesManager.NewManager(s.Store(), s.Router(), s.NetworksManager())
		return manager
	})
}

func (s *BaseServer) PermissionsManager() permissions.Manager {
	return Create(s, func() permissions.Manager {
		return permissions.NewManager(s.UsersManager())
	})
}

func (s *BaseServer) PeersManager() peers.Manager {
	return Create(s, func() peers.Manager {
		manager := peersManager.NewManager(s.Store())
		s.AfterInit(func(s *BaseServer) {
			peersManager.RegisterEndpoints(s.Router(), s.PermissionsManager(), manager)
		})
		return manager
	})
}

func (s *BaseServer) UsersManager() users.Manager {
	return Create(s, func() users.Manager {
		manager := usersManager.NewManager(s.Store())
		s.AfterInit(func(s *BaseServer) {
			usersManager.RegisterEndpoints(s.Router(), s.PermissionsManager(), manager)
		})
		return manager
	})
}
