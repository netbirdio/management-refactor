package manager

import (
	"net/http"

	"github.com/gorilla/mux"
	nbcontext "github.com/netbirdio/netbird/management/server/context"
	"github.com/netbirdio/netbird/management/server/http/util"

	"management/internal/modules/networks"
	"management/internal/shared/errors"
	"management/internal/shared/permissions"
	"management/internal/shared/permissions/modules"
	"management/internal/shared/permissions/operations"
)

type handler struct {
	manager            networks.Manager
	permissionsManager permissions.Manager
}

func newHandler(manager networks.Manager, permissionsManager permissions.Manager) *handler {
	return &handler{
		manager:            manager,
		permissionsManager: permissionsManager,
	}
}

func (h *handler) RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/networks/{id}", h.deleteNetwork).Methods("DELETE", "OPTIONS")
}

func (h *handler) deleteNetwork(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}

	allowed, err := h.permissionsManager.ValidateUserPermissions(r.Context(), userAuth.AccountId, userAuth.UserId, modules.Networks, operations.Write)
	if err != nil {
		util.WriteError(r.Context(), errors.NewPermissionValidationError(err), w)
		return
	}
	if !allowed {
		util.WriteError(r.Context(), errors.NewPermissionDeniedError(), w)
	}

	err = h.manager.DeleteNetwork(r.Context(), nil, userAuth.AccountId, userAuth.UserId, mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// _ = json.NewEncoder(w).Encode(users)
}
