package users

import (
	"context"

	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Manager interface {
	GetAllUsers(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]User, error)
	GetUserByID(ctx context.Context, tx db.Transaction, strength db.LockingStrength, id string) (*User, error)
}
