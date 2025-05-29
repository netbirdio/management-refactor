package team

import (
	"management/internal/shared/db"

	"github.com/netbirdio/netbird/management/server/integration_reference"
)

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

func (Group) TableName() string {
	return "groups"
}

type Resource struct {
	ID   string
	Type string
}

// func (r *Resource) ToAPIResponse() *api.Resource {
// 	if r.ID == "" && r.Type == "" {
// 		return nil
// 	}

// 	return &api.Resource{
// 		Id:   r.ID,
// 		Type: api.ResourceType(r.Type),
// 	}
// }

// func (r *Resource) FromAPIRequest(req *api.Resource) {
// 	if req == nil {
// 		return
// 	}

// 	r.ID = req.Id
// 	r.Type = string(req.Type)
// }

type (
	GroupEvent      = db.ModelEvent[Group]
	GroupErrorEvent = db.ModelErrorEvent[Group]
)
