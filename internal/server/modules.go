package server

import (
	"management/internal/modules/networks"
	"management/internal/modules/networks/manager"
	"management/internal/modules/networks/resources"
	resourcesManager "management/internal/modules/networks/resources/manager"
)

func (s *Server) NetworksManager() networks.Manager {
	return Create(s, func() networks.Manager {
		return manager.NewManager(s.Store(), s.Router(), s.PermissionsManager())
	})
}

func (s *Server) ResourcesManager() resources.Manager {
	return Create(s, func() resources.Manager {
		return resourcesManager.NewManager(s.Store(), s.Router(), s.NetworksManager())
	})
}
