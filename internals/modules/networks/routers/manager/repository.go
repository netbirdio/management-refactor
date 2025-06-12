package manager

import (
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Repository interface {
	Store() *db.Store
}

type repository struct {
	store *db.Store
}

func NewRepository(s *db.Store) Repository {
	return &repository{store: s}
}

func (r *repository) Store() *db.Store {
	return r.store
}
