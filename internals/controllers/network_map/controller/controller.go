package controller

import (
	"context"
	"fmt"
	"sync"

	nbpeer "github.com/netbirdio/netbird/management/server/peer"
	log "github.com/sirupsen/logrus"

	"github.com/netbirdio/management-refactor/internals/controllers/network_map"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	appmetrics "github.com/netbirdio/management-refactor/internals/shared/metrics"
)

type Controller struct {
	repo          Repository
	metrics       *metrics
	updateChannel *UpdateChannel
}

func NewController(store *db.Store, metrics *appmetrics.AppMetrics) *Controller {
	cMetrics, err := appmetrics.RegisterMetrics(metrics, newMetrics)
	if err != nil {
		log.Fatalf("Failed to register app metrics: %v", err)
	}
	return &Controller{
		repo:          newRepository(store, cMetrics),
		metrics:       cMetrics,
		updateChannel: NewUpdateChannel(cMetrics),
	}
}

func (c *Controller) CalculateNetworkMap(data *network_map.NetworkMapData) (*network_map.NetworkMap, error) {
	// Do calc on data

	return &network_map.NetworkMap{}, nil
}

func (c *Controller) UpdatePeers(ctx context.Context, accountID string) error {
	data, err := c.repo.GetNetworkMapData(accountID)
	if err != nil {
		return fmt.Errorf("get network map data: %w", err)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)

	for _, peer := range data.Peers {
		if !c.updateChannel.HasChannel(peer.ID) {
			log.WithContext(ctx).Tracef("peer %s doesn't have a channel, skipping network map update", peer.ID)
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}
		go func(p *nbpeer.Peer) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// TODO: posture checks

			_, err := c.CalculateNetworkMap(data)
			if err != nil {
				log.WithContext(ctx).Errorf("failed to calculate network map for peer %s: %v", p.ID, err)
				return
			}

			c.updateChannel.SendUpdate(ctx, p.ID, &network_map.UpdateMessage{})
		}(&peer)
	}

	return nil
}
