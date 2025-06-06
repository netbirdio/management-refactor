package manager

import (
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Repository interface {
	RunInTx(fn func(tx db.Transaction) error) error
}

type repository struct {
	store *db.Store
}

func newRepository(s *db.Store) Repository {
	return &repository{store: s}
}

func (r *repository) RunInTx(fn func(tx db.Transaction) error) error {
	return r.store.RunInTx(fn)
}
