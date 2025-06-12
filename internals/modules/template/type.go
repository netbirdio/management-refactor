package template

import (
	"github.com/netbirdio/netbird/management/server/http/api"
)

type Template struct {
	Id string `gorm:"primary_key"`
}

// NewTemplate creates a new Template object
func NewTemplate() *Template {
	return &Template{}
}

// Copy the Template object
func (u *Template) Copy() *Template {
	return &Template{}
}

func (u *Template) ToApiResponse() *api.Template {

}

func (n *Template) FromAPIRequest(req *api.Template) {

}

func (n *Template) Validate() error {
	return nil
}
