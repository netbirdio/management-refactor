package manager

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	nbcontext "github.com/netbirdio/netbird/management/server/context"

	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/modules"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

type handler struct {
	manager users.Manager
}

func RegisterEndpoints(router *mux.Router, permissionsManager permissions.Manager, manager users.Manager) {
	h := &handler{
		manager: manager,
	}

	router.HandleFunc("/users", permissionsManager.WithPermission(modules.Users, operations.Read, h.getAllUsers)).Methods("GET", "OPTIONS")
	router.HandleFunc("/users/{userId}", permissionsManager.WithPermission(modules.Users, operations.Read, h.getUser)).Methods("GET", "OPTIONS")
}

func (h *handler) getAllUsers(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {
	users, err := h.manager.GetAllUsers(r.Context(), nil, db.LockingStrengthShare, userAuth.AccountId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(users)
}

func (h *handler) getUser(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	user, err := h.manager.GetUserByID(r.Context(), nil, db.LockingStrengthShare, userId)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(user)
}
