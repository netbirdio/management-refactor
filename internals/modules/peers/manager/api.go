package manager

import (
	"net/http"

	"github.com/gorilla/mux"
	nbcontext "github.com/netbirdio/netbird/management/server/context"
	"github.com/netbirdio/netbird/management/server/http/util"
	"github.com/netbirdio/netbird/management/server/status"

	"github.com/netbirdio/management-refactor/internals/modules/peers"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/modules"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

type handler struct {
	manager peers.Manager
}

func RegisterEndpoints(router *mux.Router, permissionsManager permissions.Manager, manager peers.Manager) {
	h := &handler{
		manager: manager,
	}

	router.HandleFunc("/peers", permissionsManager.WithPermission(modules.Peers, operations.Read, h.getAllPeers)).Methods("GET", "OPTIONS")
	router.HandleFunc("/peers/{peerId}", permissionsManager.WithPermission(modules.Peers, operations.Read, h.getPeer)).Methods("GET", "OPTIONS")
	router.HandleFunc("/peers/{peerId}", permissionsManager.WithPermission(modules.Peers, operations.Write, h.updatePeer)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/peers/{peerId}", permissionsManager.WithPermission(modules.Peers, operations.Write, h.deletePeer)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/peers/{peerId}/accessible-peers", permissionsManager.WithPermission(modules.Peers, operations.Read, h.getAccessiblePeers)).Methods("GET", "OPTIONS")
}

func (h *handler) getAllPeers(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {
	peers := []peers.Peer{{ID: "peer1"}, {ID: "peer2"}}
	util.WriteJSONObject(r.Context(), w, peers)
}

func (h *handler) getPeer(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {

}

func (h *handler) updatePeer(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {
	vars := mux.Vars(r)
	peerID, ok := vars["peerId"]
	if !ok {
		util.WriteError(r.Context(), status.Errorf(status.InvalidArgument, "peer ID field is missing"), w)
		return
	}
	if len(peerID) == 0 {
		util.WriteError(r.Context(), status.Errorf(status.InvalidArgument, "peer ID can't be empty"), w)
		return
	}

	err := h.manager.UpdatePeer(r.Context(), nil, &peers.Peer{ID: peerID, AccountID: peerID})
	if err != nil {
		util.WriteErrorResponse("Failed to update peer", http.StatusInternalServerError, w)
		return
	}

	util.WriteJSONObject(r.Context(), w, map[string]string{"status": "success"})
}

func (h *handler) deletePeer(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {

}

func (h *handler) getAccessiblePeers(w http.ResponseWriter, r *http.Request, userAuth *nbcontext.UserAuth) {

}
