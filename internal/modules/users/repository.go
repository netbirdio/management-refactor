package users

import (
	"management/internal/modules/users/types"
	"management/internal/shared/db"
)

type Repository interface {
	RunInTx(fn func(tx db.Transaction) error) error
	GetAllUsers(tx db.Transaction, strength db.LockingStrength, accountID string) ([]types.User, error)
	GetUserByID(tx db.Transaction, strength db.LockingStrength, id string) (*types.User, error)
	CreateUser(tx db.Transaction, u *types.User) error
}

type repository struct {
	store *db.Store
}

func newRepository(s *db.Store) Repository {
	err := s.AutoMigrate(types.User{})
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
	return &repository{store: s}
}

func (r *repository) RunInTx(fn func(tx db.Transaction) error) error {
	return r.store.RunInTx(fn)
}

func (r *repository) GetAllUsers(tx db.Transaction, strength db.LockingStrength, accountID string) ([]types.User, error) {
	var users []types.User
	err := r.store.GetMany(tx, strength, &users, "account_id = ?", accountID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetUserByID(tx db.Transaction, strength db.LockingStrength, id string) (*types.User, error) {
	var user types.User
	err := r.store.GetOne(tx, strength, &user, "id = ?", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) CreateUser(tx db.Transaction, u *types.User) error {
	return r.store.Create(tx, u)
}
