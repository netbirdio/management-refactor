package network_map

import (
	nbpeer "github.com/netbirdio/netbird/management/server/peer"

	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/modules/networks/resources"
	"github.com/netbirdio/management-refactor/internals/modules/networks/routers"
)

type NetworkMapData struct {
	// we have to name column to aid as it collides with Network.Id when work with associations
	Id string `gorm:"primaryKey"`

	Domain                 string `gorm:"index"`
	DomainCategory         string
	IsDomainPrimaryAccount bool
	Peers                  []nbpeer.Peer `json:"-" gorm:"foreignKey:AccountID;references:id"`

	Networks         []*networks.Network          `gorm:"foreignKey:AccountID;references:id"`
	NetworkRouters   []*routers.NetworkRouter     `gorm:"foreignKey:AccountID;references:id"`
	NetworkResources []*resources.NetworkResource `gorm:"foreignKey:AccountID;references:id"`
}
