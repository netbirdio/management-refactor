package manager

import (
	"context"
	"management/internal/modules/team"
	"management/internal/shared/db"
)

// GetPATSByAccount implements team.Manager.
func (m *manager) GetPATSByAccount(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*team.PersonalAccessToken, error) {
	panic("unimplemented")
}

// GetPATById implements team.Manager.
func (m *manager) GetPATById(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string, patID string) (*team.PersonalAccessToken, error) {
	panic("unimplemented")
}

// CreatePAT implements team.Manager.
func (m *manager) CreatePAT(ctx context.Context, tx db.Transaction, group *team.PersonalAccessToken) (*team.PersonalAccessToken, error) {
	panic("unimplemented")
}

// DeletePAT implements team.Manager.
func (m *manager) DeletePAT(ctx context.Context, tx db.Transaction, group *team.PersonalAccessToken) error {
	panic("unimplemented")
}

// UpdatePAT implements team.Manager.
func (m *manager) UpdatePAT(ctx context.Context, tx db.Transaction, group *team.PersonalAccessToken) (*team.PersonalAccessToken, error) {
	panic("unimplemented")
}

// GetPATSByUser implements team.Manager.
func (m *manager) GetPATSByUser(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string, userID string) ([]*team.PersonalAccessToken, error) {
	panic("unimplemented")
}
