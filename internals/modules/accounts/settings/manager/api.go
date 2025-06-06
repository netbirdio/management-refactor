package manager

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	nbcontext "github.com/netbirdio/netbird/management/server/context"
	"github.com/netbirdio/netbird/management/server/http/util"

	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/internals/shared/errors"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/modules"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

type handler struct {
	manager            *Manager
	permissionsManager permissions.Manager
}

func newHandler(manager *Manager, permissionsManager permissions.Manager) *handler {
	return &handler{
		manager:            manager,
		permissionsManager: permissionsManager,
	}
}

func (h *handler) RegisterEndpoints(router *mux.Router) {
	router.HandleFunc("/account/{accountID}/settings", h.getSettings).Methods("GET", "OPTIONS")
	router.HandleFunc("/account/{accountID}/settings", h.updateSettings).Methods("PUT", "OPTIONS")
}

func (h *handler) getSettings(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}

	vars := mux.Vars(r)
	accountId := vars["accountID"]

	allowed, err := h.permissionsManager.ValidateUserPermissions(r.Context(), accountId, userAuth.UserId, modules.Settings, operations.Read)
	if err != nil {
		util.WriteError(r.Context(), errors.NewPermissionValidationError(err), w)
		return
	}
	if !allowed {
		util.WriteError(r.Context(), errors.NewPermissionDeniedError(), w)
	}

	users, err := h.manager.GetSettings(r.Context(), nil, db.LockingStrengthShare, accountId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(users)
}

func (h *handler) updateSettings(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}

	vars := mux.Vars(r)
	accountId := vars["accountID"]

	allowed, err := h.permissionsManager.ValidateUserPermissions(r.Context(), accountId, userAuth.UserId, modules.Settings, operations.Write)
	if err != nil {
		util.WriteError(r.Context(), errors.NewPermissionValidationError(err), w)
		return
	}
	if !allowed {
		util.WriteError(r.Context(), errors.NewPermissionDeniedError(), w)
	}

	settings, err := h.manager.UpdateSettings(r.Context(), nil, db.LockingStrengthShare, accountId)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(settings)
}
