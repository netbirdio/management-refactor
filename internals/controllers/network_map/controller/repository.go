package controller

import (
	"time"

	"github.com/netbirdio/management-refactor/internals/controllers/network_map"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Repository interface {
	GetNetworkMapData(accountID string) (*network_map.NetworkMapData, error)
}

type repository struct {
	store   *db.Store
	metrics *metrics
}

func newRepository(s *db.Store, metrics *metrics) Repository {
	return &repository{
		store:   s,
		metrics: metrics,
	}
}

func (r *repository) GetNetworkMapData(accountID string) (*network_map.NetworkMapData, error) {
	start := time.Now()
	var networkMapData network_map.NetworkMapData
	err := r.store.GetOne(nil, db.LockingStrengthShare, &networkMapData, "id = ?", accountID)
	if err != nil {
		return nil, err
	}

	// if err := r.store.Load(&networkMapData, "Peers", "Groups", "Policies", "Networks", "NetworkRouters", "NetworkResources"); err != nil {
	// 	return nil, err
	// }

	r.metrics.RecordDBAccessDuration(time.Since(start))

	return &networkMapData, nil
}
