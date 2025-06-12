package server

import (
	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/modules/networks/manager"
	"github.com/netbirdio/management-refactor/internals/modules/networks/resources"
	resourcesManager "github.com/netbirdio/management-refactor/internals/modules/networks/resources/manager"
	"github.com/netbirdio/management-refactor/internals/modules/networks/routers"
	routersManager "github.com/netbirdio/management-refactor/internals/modules/networks/routers/manager"
	"github.com/netbirdio/management-refactor/internals/modules/peers"
	peersManager "github.com/netbirdio/management-refactor/internals/modules/peers/manager"
	"github.com/netbirdio/management-refactor/internals/modules/users"
	usersManager "github.com/netbirdio/management-refactor/internals/modules/users/manager"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
)

func (s *BaseServer) NetworksManager() networks.Manager {
	return Create(s, func() networks.Manager {
		repo := manager.NewRepository(s.Store())
		return manager.NewManager(repo, s.ActivityManager(), s.ResourcesManager(), s.RoutersManager())
	})
}

func (s *BaseServer) ResourcesManager() resources.Manager {
	return Create(s, func() resources.Manager {
		repo := resourcesManager.NewRepository(s.Store())
		manager := resourcesManager.NewManager(repo, s.Router(), s.NetworksManager())
		return manager
	})
}

func (s *BaseServer) RoutersManager() routers.Manager {
	return Create(s, func() routers.Manager {
		repo := routersManager.NewRepository(s.Store())
		manager := routersManager.NewManager(repo)
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
		repo := peersManager.NewRepository(s.Store())
		manager := peersManager.NewManager(repo)
		s.AfterInit(func(s *BaseServer) {
			peersManager.RegisterEndpoints(s.Router(), s.PermissionsManager(), manager)
			manager.SetNetworkMapController(s.NetworkMapController())
		})
		return manager
	})
}

func (s *BaseServer) UsersManager() users.Manager {
	return Create(s, func() users.Manager {
		repo := usersManager.NewRepository(s.Store())
		manager := usersManager.NewManager(repo)
		s.AfterInit(func(s *BaseServer) {
			usersManager.RegisterEndpoints(s.Router(), s.PermissionsManager(), manager)
		})
		return manager
	})
}
