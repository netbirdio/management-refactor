package manager

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/netbirdio/management-refactor/internals/modules/templates"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/modules"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
	nbcontext "github.com/netbirdio/netbird/management/server/context"
)

type handler struct {
	manager templates.Manager
}

func RegisterEndpoints(router *mux.Router, permissionsManager permissions.Manager, manager templates.Manager) {
	h := &handler{
		manager: manager,
	}

	router.HandleFunc("/templates", permissionsManager.WithPermission(modules.Template, operations.Write, h.createTemplate)).Methods("POST", "OPTIONS")
	router.HandleFunc("/templates", permissionsManager.WithPermission(modules.Template, operations.Read, h.getAllTemplates)).Methods("GET", "OPTIONS")
	router.HandleFunc("/templates/{templateId}", permissionsManager.WithPermission(modules.Template, operations.Read, h.getTemplate)).Methods("GET", "OPTIONS")
	router.HandleFunc("/templates/{templateId}", permissionsManager.WithPermission(modules.Template, operations.Write, h.updateTemplate)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/templates/{templateId}", permissionsManager.WithPermission(modules.Template, operations.Write, h.deleteTemplate)).Methods("DELETE", "OPTIONS")
}

func (h *handler) getAllTemplates(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {
	users, err := h.manager.GetAllTemplates(r.Context(), nil, db.LockingStrengthShare, userAuth.AccountId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(users)
}

func (h *handler) getTemplate(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {
	vars := mux.Vars(r)
	templateId := vars["templateId"]

	user, err := h.manager.GetTemplateByID(r.Context(), nil, db.LockingStrengthShare, templateId)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(user)
}

func (h *handler) createTemplate(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {
	var templateReq api.Template
	if err := json.NewDecoder(r.Body).Decode(&templateReq); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	templateObj := &templates.Template{}
	templateObj.FromAPIRequest(templateReq)

	err := templateObj.Validate()
	if err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	templateResp, err := h.manager.CreateTemplate(r.Context(), nil, &templateReq)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(templateResp)
}

func (h *handler) updateTemplate(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {
	vars := mux.Vars(r)
	templateId := vars["templateId"]

	var templateReq api.Template
	if err := json.NewDecoder(r.Body).Decode(&templateReq); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	templateObj := &templates.Template{}
	templateObj.FromAPIRequest(templateReq)
	templateObj.Id = templateId

	err := templateObj.Validate()
	if err != nil {
		http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
		return
	}

	templateResp, err := h.manager.UpdateTemplate(r.Context(), nil, &templateReq)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(templateResp)
}

func (h *handler) deleteTemplate(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {
	vars := mux.Vars(r)
	templateId := vars["templateId"]

	err := h.manager.DeleteTemplate(r.Context(), nil, templateId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
