package network_map

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/netbirdio/management-refactor/internals/shared/db"
	appmetrics "github.com/netbirdio/management-refactor/internals/shared/metrics"
)

type Controller struct {
	repo    Repository
	metrics *metrics
}

func NewController(store *db.Store, metrics *appmetrics.AppMetrics) *Controller {
	cMetrics, err := appmetrics.RegisterMetrics(metrics, newMetrics)
	if err != nil {
		log.Fatalf("Failed to register app metrics: %v", err)
	}
	return &Controller{
		repo:    newRepository(store, cMetrics),
		metrics: cMetrics,
	}
}

func (c *Controller) UpdatePeers(accountID string) error {
	_, err := c.repo.GetNetworkMapData(accountID)
	if err != nil {
		return fmt.Errorf("get network map data: %w", err)
	}

	return nil
}
