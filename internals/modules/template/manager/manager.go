//go:build ignore
// +build ignore

package manager

import (
	"github.com/gorilla/mux"

	"github.com/netbirdio/management-refactor/internals/modules/template"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	appmetrics "github.com/netbirdio/management-refactor/internals/shared/metrics"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
)

type managerImpl struct {
	repo Repository
}

func NewManager(store *db.Store, router *mux.Router, metrics appmetrics.AppMetrics, permissionsManager permissions.Manager) template.Manager {
	repo := newRepository(store)
	m := &managerImpl{
		repo: repo,
	}

	api := newHandler(m, permissionsManager)
	api.RegisterEndpoints(router)
	return m
}
