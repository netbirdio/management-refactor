package repository

import (
	"management/internal/modules/team"
	"management/internal/shared/db"
)

var _ team.Repository = (*repository)(nil)

type repository struct {
	store *db.Store
}

func NewRepository(s *db.Store) *repository {
	return &repository{store: s}
}

func (r *repository) Store() *db.Store {
	return r.store
}
