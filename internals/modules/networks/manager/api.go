package manager

import (
	"net/http"

	"github.com/gorilla/mux"
	nbcontext "github.com/netbirdio/netbird/management/server/context"

	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/modules"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

type handler struct {
	manager networks.Manager
}

func RegisterEndpoints(router *mux.Router, permissionsManager permissions.Manager, manager networks.Manager) {
	h := &handler{
		manager: manager,
	}

	router.HandleFunc("/networks/{id}", permissionsManager.WithPermission(modules.Networks, operations.Write, h.deleteNetwork)).Methods("DELETE", "OPTIONS")
}

func (h *handler) deleteNetwork(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {
	err := h.manager.DeleteNetwork(r.Context(), nil, userAuth.AccountId, userAuth.UserId, mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// _ = json.NewEncoder(w).Encode(users)
}
