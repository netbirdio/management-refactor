package settings

import (
	"management/internal/modules/accounts/settings/types"
	"management/internal/shared/db"
)

type Repository interface {
	RunInTx(fn func(tx db.Transaction) error) error
	GetAccountSettings(tx db.Transaction, strength db.LockingStrength, accountID string) (*types.Settings, error)
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

func (r *repository) GetAccountSettings(tx db.Transaction, strength db.LockingStrength, accountID string) (*types.Settings, error) {
	var settings types.Settings
	err := r.store.GetOne(tx, strength, &settings, "account_id = ?", accountID)
	return &settings, err
}
