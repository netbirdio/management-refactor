package manager

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/netbirdio/management-refactor/internals/shared/permissions"
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
	router.HandleFunc("/peers", h.getAllPeers).Methods("GET", "OPTIONS")
	router.HandleFunc("/peers/{peerId}", h.getPeer).Methods("GET", "OPTIONS")
	router.HandleFunc("/peers/{peerId}", h.updatePeer).Methods("PUT", "OPTIONS")
	router.HandleFunc("/peers/{peerId}", h.deletePeer).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/peers/{peerId}/accessible-peers", h.getAccessiblePeers).Methods("GET", "OPTIONS")
}

func (h *handler) getAllPeers(w http.ResponseWriter, r *http.Request) {}

func (h *handler) getPeer(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) updatePeer(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) deletePeer(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) getAccessiblePeers(w http.ResponseWriter, r *http.Request) {

}
