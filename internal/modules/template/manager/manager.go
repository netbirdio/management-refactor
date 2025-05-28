//go:build ignore
// +build ignore

package manager

import (
	"github.com/gorilla/mux"

	"management/internal/modules/template"
	"management/internal/shared/db"
	appmetrics "management/internal/shared/metrics"
	"management/internal/shared/permissions"
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
