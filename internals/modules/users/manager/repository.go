package manager

import (
	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Repository interface {
	RunInTx(fn func(tx db.Transaction) error) error
	GetAllUsers(tx db.Transaction, strength db.LockingStrength, accountID string) ([]users.User, error)
	GetUserByID(tx db.Transaction, strength db.LockingStrength, id string) (*users.User, error)
	CreateUser(tx db.Transaction, u *users.User) error
}

type repository struct {
	store *db.Store
}

func NewRepository(s *db.Store) Repository {
	err := s.AutoMigrate(users.User{})
	if err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
	return &repository{store: s}
}

func (r *repository) RunInTx(fn func(tx db.Transaction) error) error {
	return r.store.RunInTx(fn)
}

func (r *repository) GetAllUsers(tx db.Transaction, strength db.LockingStrength, accountID string) ([]users.User, error) {
	var users []users.User
	err := r.store.GetMany(tx, strength, &users, "account_id = ?", accountID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetUserByID(tx db.Transaction, strength db.LockingStrength, id string) (*users.User, error) {
	var user users.User
	err := r.store.GetOne(tx, strength, &user, "id = ?", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) CreateUser(tx db.Transaction, u *users.User) error {
	return r.store.Create(tx, u)
}
