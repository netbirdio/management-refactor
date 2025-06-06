package peers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	nbcontext "github.com/netbirdio/netbird/management/server/context"
	"github.com/netbirdio/netbird/management/server/groups"
	"github.com/netbirdio/netbird/management/server/http/api"
	"github.com/netbirdio/netbird/management/server/http/util"

	"management/internal/shared/db"
	"management/internal/shared/permissions"
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

func (h *handler) getAllPeers(w http.ResponseWriter, r *http.Request) {
	userAuth, err := nbcontext.GetUserAuthFromContext(r.Context())
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}

	nameFilter := r.URL.Query().Get("name")
	ipFilter := r.URL.Query().Get("ip")

	peers, err := h.manager.GetFilteredPeers(r.Context(), nil, db.LockingStrengthShare, userAuth.AccountId, nameFilter, ipFilter)
	if err != nil {
		util.WriteError(r.Context(), err, w)
		return
	}

	dnsDomain := h.accountManager.GetDNSDomain()

	grps, _ := h.accountManager.GetAllGroups(r.Context(), accountID, userID)

	grpsInfoMap := groups.ToGroupsInfoMap(grps, len(peers))
	respBody := make([]*api.PeerBatch, 0, len(peers))
	for _, peer := range peers {
		peerToReturn, err := h.checkPeerStatus(peer)
		if err != nil {
			util.WriteError(r.Context(), err, w)
			return
		}

		respBody = append(respBody, toPeerListItemResponse(peerToReturn, grpsInfoMap[peer.ID], dnsDomain, 0))
	}

	validPeersMap, err := h.accountManager.GetValidatedPeers(r.Context(), accountID)
	if err != nil {
		log.WithContext(r.Context()).Errorf("failed to list appreoved peers: %v", err)
		util.WriteError(r.Context(), fmt.Errorf("internal error"), w)
		return
	}
	h.setApprovalRequiredFlag(respBody, validPeersMap)

	util.WriteJSONObject(r.Context(), w, respBody)
}

func (h *handler) getPeer(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) updatePeer(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) deletePeer(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) getAccessiblePeers(w http.ResponseWriter, r *http.Request) {

}
