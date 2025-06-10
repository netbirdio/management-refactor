package network_map

import (
	nbpeer "github.com/netbirdio/netbird/management/server/peer"

	"github.com/netbirdio/management-refactor/internals/modules/accounts"
	"github.com/netbirdio/management-refactor/internals/modules/groups"
	"github.com/netbirdio/management-refactor/internals/modules/networks"
	"github.com/netbirdio/management-refactor/internals/modules/networks/resources"
	"github.com/netbirdio/management-refactor/internals/modules/networks/routers"
	"github.com/netbirdio/management-refactor/internals/modules/policies"
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
