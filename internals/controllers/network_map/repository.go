package network_map

import (
	"time"

	nbpeer "github.com/netbirdio/netbird/management/server/peer"

	"github.com/netbirdio/management-refactor/internals/modules/accounts"
	"github.com/netbirdio/management-refactor/internals/modules/groups"
	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/modules/networks/resources"
	"github.com/netbirdio/management-refactor/internals/modules/networks/routers"
	"github.com/netbirdio/management-refactor/internals/modules/policies"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type NetworkMapData struct {
	// we have to name column to aid as it collides with Network.Id when work with associations
	Id string `gorm:"primaryKey"`

	Domain                 string `gorm:"index"`
	DomainCategory         string
	IsDomainPrimaryAccount bool
	Network                *accounts.Network  `gorm:"embedded;embeddedPrefix:network_"`
	Peers                  []nbpeer.Peer      `json:"-" gorm:"foreignKey:AccountID;references:id"`
	Groups                 []groups.Group     `json:"-" gorm:"foreignKey:AccountID;references:id"`
	Policies               []*policies.Policy `gorm:"foreignKey:AccountID;references:id"`

	Networks         []*networks.Network          `gorm:"foreignKey:AccountID;references:id"`
	NetworkRouters   []*routers.NetworkRouter     `gorm:"foreignKey:AccountID;references:id"`
	NetworkResources []*resources.NetworkResource `gorm:"foreignKey:AccountID;references:id"`
}

type Repository interface {
	GetNetworkMapData(accountID string) (*NetworkMapData, error)
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

func (r *repository) GetNetworkMapData(accountID string) (*NetworkMapData, error) {
	start := time.Now()
	var networkMapData NetworkMapData
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
