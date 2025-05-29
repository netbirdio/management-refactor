package manager

import (
	"context"
	"management/internal/modules/team"
	"management/internal/shared/db"
)

// GetUserById implements team.Manager.
func (m *manager) GetUserById(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string, userID string) (*team.User, error) {
	panic("unimplemented")
}

// GetUsersByAccount implements team.Manager.
func (m *manager) GetUsersByAccount(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*team.User, error) {
	panic("unimplemented")
}

// UpdateUser implements team.Manager.
func (m *manager) UpdateUser(ctx context.Context, tx db.Transaction, user *team.User) (*team.User, error) {
	panic("unimplemented")
}

// CreateUser implements team.Manager.
func (m *manager) CreateUser(ctx context.Context, tx db.Transaction, user *team.User) (*team.User, error) {
	panic("unimplemented")
}

// DeleteUser implements team.Manager.
func (m *manager) DeleteUser(ctx context.Context, tx db.Transaction, user *team.User) error {
	panic("unimplemented")
}
