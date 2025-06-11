package controller

import (
	"context"

	"github.com/netbirdio/management-refactor/internals/controllers/network_map"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	appmetrics "github.com/netbirdio/management-refactor/internals/shared/metrics"
	"github.com/netbirdio/management-refactor/pkg/logging"
)

var log = logging.LoggerForThisPackage

type Controller struct {
	repo          Repository
	metrics       *metrics
	UpdateChannel network_map.UpdateChannel
}

func NewController(store *db.Store, metrics *appmetrics.AppMetrics, updateChannel network_map.UpdateChannel) *Controller {
	cMetrics, err := appmetrics.RegisterMetrics(metrics, newMetrics)
	if err != nil {
		log().Fatalf("Failed to register app metrics: %v", err)
	}
	return &Controller{
		repo:          newRepository(store, cMetrics),
		metrics:       cMetrics,
		UpdateChannel: updateChannel,
	}
}

func (c *Controller) CalculateNetworkMap(accountID string) (*network_map.NetworkMap, error) {
	_, err := c.repo.GetNetworkMapData(accountID)
	if err != nil {
		// usually return error
	}

	// Do calc on data

	log().Tracef("Calculating network map for account on public")

	return &network_map.NetworkMap{}, nil
}

func (c *Controller) UpdatePeers(ctx context.Context, accountID string) error {
	_, err := c.CalculateNetworkMap(accountID)
	if err != nil {
		log().Errorf("Failed to calculate network map for account %s: %v", accountID, err)
		return err
	}
	c.UpdateChannel.SendUpdate(ctx, accountID, &network_map.UpdateMessage{})

	return nil
}
