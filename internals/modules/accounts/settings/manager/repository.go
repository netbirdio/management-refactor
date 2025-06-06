package manager

import (
	"management/internal/modules/accounts/settings/types"
	"management/internal/shared/db"
)

type Repository interface {
	RunInTx(fn func(tx db.Transaction) error) error
	GetAccountSettings(tx db.Transaction, strength db.LockingStrength, accountID string) (*types.Settings, error)
	UpdateSettings(tx db.Transaction, settings *types.Settings) (*types.Settings, error)
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

func (r *repository) UpdateSettings(tx db.Transaction, settings *types.Settings) (*types.Settings, error) {
	err := r.store.Update(tx, settings)
	if err != nil {
		return nil, err
	}
	return settings, nil
}
