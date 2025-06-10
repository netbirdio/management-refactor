package manager

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	nbcontext "github.com/netbirdio/management-refactor/management/server/context"
	"github.com/netbirdio/management-refactor/management/server/http/util"

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
	router.HandleFunc("/groups", h.getAllGroups).Methods("GET", "OPTIONS")
	router.HandleFunc("/groups", h.createGroup).Methods("POST", "OPTIONS")
	router.HandleFunc("/groups/{groupId}", h.updateGroup).Methods("PUT", "OPTIONS")
	router.HandleFunc("/groups/{groupId}", h.getGroup).Methods("GET", "OPTIONS")
	router.HandleFunc("/groups/{groupId}", h.deleteGroup).Methods("DELETE", "OPTIONS")
}

func (h *handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}

	allowed, err := h.permissionsManager.ValidateUserPermissions(r.Context(), userAuth.AccountId, userAuth.UserId, modules.Users, operations.Read)
	if err != nil {
		util.WriteError(r.Context(), errors.NewPermissionValidationError(err), w)
		return
	}
	if !allowed {
		util.WriteError(r.Context(), errors.NewPermissionDeniedError(), w)
	}

	users, err := h.manager.GetAllUsers(r.Context(), nil, db.LockingStrengthShare, userAuth.AccountId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(users)
}

func (h *handler) getUser(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}

	allowed, err := h.permissionsManager.ValidateUserPermissions(r.Context(), userAuth.AccountId, userAuth.UserId, modules.Users, operations.Read)
	if err != nil {
		util.WriteError(r.Context(), errors.NewPermissionValidationError(err), w)
		return
	}
	if !allowed {
		util.WriteError(r.Context(), errors.NewPermissionDeniedError(), w)
	}

	vars := mux.Vars(r)
	userId := vars["userId"]

	user, err := h.manager.GetUserByID(r.Context(), nil, db.LockingStrengthShare, userId)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(user)
}
