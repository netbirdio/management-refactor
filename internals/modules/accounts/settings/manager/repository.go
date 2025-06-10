package manager

import (
	"github.com/netbirdio/management-refactor/internals/modules/accounts/settings"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Repository interface {
	RunInTx(fn func(tx db.Transaction) error) error
	GetAccountSettings(tx db.Transaction, strength db.LockingStrength, accountID string) (*settings.Settings, error)
	UpdateSettings(tx db.Transaction, settings *settings.Settings) (*settings.Settings, error)
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

func (r *repository) GetAccountSettings(tx db.Transaction, strength db.LockingStrength, accountID string) (*settings.Settings, error) {
	var settings settings.Settings
	err := r.store.GetOne(tx, strength, &settings, "account_id = ?", accountID)
	return &settings, err
}

func (r *repository) UpdateSettings(tx db.Transaction, settings *settings.Settings) (*settings.Settings, error) {
	err := r.store.Update(tx, settings)
	if err != nil {
		return nil, err
	}
	return settings, nil
}
