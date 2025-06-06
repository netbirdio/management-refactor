package groups

import "github.com/netbirdio/netbird/management/server/integration_reference"

type Group struct {
	// ID of the group
	ID string `gorm:"primaryKey"`

	// AccountID is a reference to Account that this object belongs
	AccountID string `json:"-" gorm:"index"`

	// Name visible in the UI
	Name string

	// Issued defines how this group was created (enum of "api", "integration" or "jwt")
	Issued string

	// Peers list of the group
	Peers []string `gorm:"serializer:json"`

	// Resources contains a list of resources in that group
	Resources []Resource `gorm:"serializer:json"`

	IntegrationReference integration_reference.IntegrationReference `gorm:"embedded;embeddedPrefix:integration_ref_"`
}
