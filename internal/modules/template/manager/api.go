//go:build ignore
// +build ignore

package manager

type handler struct {
	manager            template.Manager
	permissionsManager permissions.Manager
}

func newHandler(manager template.Manager, permissionsManager permissions.Manager) *handler {
	return &handler{
		manager:            manager,
		permissionsManager: permissionsManager,
	}
}

func (h *handler) RegisterEndpoints(router *mux.Router) {
	// Register the API endpoints for the module
}
