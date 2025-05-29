package team

//go:generate go run github.com/golang/mock/mockgen -package team -destination=manager_mock.go -source=./manager.go -build_flags=-mod=mod

import (
	"context"
	"management/internal/shared/db"
	"management/internal/shared/hook"
)

// needed for go:generate's mock run
// https://github.com/golang/mock/issues/621
type (
	GroupHookEvent      = *hook.Hook[*GroupEvent]
	GroupHookErrorEvent = *hook.Hook[*GroupErrorEvent]
)

type Manager interface {
	// User methods
	GetUsersByAccount(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*User, error)
	GetUserById(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, userID string) (*User, error)
	CreateUser(ctx context.Context, tx db.Transaction, user *User) (*User, error)
	UpdateUser(ctx context.Context, tx db.Transaction, user *User) (*User, error)
	DeleteUser(ctx context.Context, tx db.Transaction, user *User) error

	// Group methods
	GetGroupsByAccount(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*Group, error)
	GetGroupsByUser(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, userID string) ([]*Group, error)
	GetGroupById(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, groupID string) (*Group, error)
	CreateGroup(ctx context.Context, tx db.Transaction, group *Group) (*Group, error)
	UpdateGroup(ctx context.Context, tx db.Transaction, group *Group) (*Group, error)
	DeleteGroup(ctx context.Context, tx db.Transaction, group *Group) error
	OnGroupCreate() GroupHookEvent
	OnGroupCreateError() GroupHookErrorEvent
	OnGroupUpdate() GroupHookEvent
	OnGroupUpdateError() GroupHookErrorEvent
	OnGroupDelete() GroupHookEvent
	OnGroupDeleteError() GroupHookErrorEvent

	// PAT methods
	GetPATSByAccount(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*PersonalAccessToken, error)
	GetPATSByUser(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, userID string) ([]*PersonalAccessToken, error)
	GetPATById(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID, patID string) (*PersonalAccessToken, error)
	CreatePAT(ctx context.Context, tx db.Transaction, group *PersonalAccessToken) (*PersonalAccessToken, error)
	UpdatePAT(ctx context.Context, tx db.Transaction, group *PersonalAccessToken) (*PersonalAccessToken, error)
	DeletePAT(ctx context.Context, tx db.Transaction, group *PersonalAccessToken) error
}
