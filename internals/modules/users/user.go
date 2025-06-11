package users

import (
	"strings"
	"time"

	"github.com/netbirdio/netbird/management/server/integration_reference"

	"github.com/netbirdio/management-refactor/internals/modules/users/pats"
)

const (
	UserRoleOwner        UserRole = "owner"
	UserRoleAdmin        UserRole = "admin"
	UserRoleUser         UserRole = "user"
	UserRoleUnknown      UserRole = "unknown"
	UserRoleBillingAdmin UserRole = "billing_admin"
	UserRoleAuditor      UserRole = "auditor"
	UserRoleNetworkAdmin UserRole = "network_admin"

	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
	UserStatusInvited  UserStatus = "invited"

	UserIssuedAPI         = "api"
	UserIssuedIntegration = "integration"
)

// StrRoleToUserRole returns UserRole for a given strRole or UserRoleUnknown if the specified role is unknown
func StrRoleToUserRole(strRole string) UserRole {
	switch strings.ToLower(strRole) {
	case "owner":
		return UserRoleOwner
	case "admin":
		return UserRoleAdmin
	case "user":
		return UserRoleUser
	case "billing_admin":
		return UserRoleBillingAdmin
	default:
		return UserRoleUnknown
	}
}

// UserStatus is the status of a User
type UserStatus string

// UserRole is the role of a User
type UserRole string

type UserInfo struct {
	ID                   string                                     `json:"id"`
	Email                string                                     `json:"email"`
	Name                 string                                     `json:"name"`
	Role                 string                                     `json:"role"`
	AutoGroups           []string                                   `json:"auto_groups"`
	Status               string                                     `json:"-"`
	IsServiceUser        bool                                       `json:"is_service_user"`
	IsBlocked            bool                                       `json:"is_blocked"`
	NonDeletable         bool                                       `json:"non_deletable"`
	LastLogin            time.Time                                  `json:"last_login"`
	Issued               string                                     `json:"issued"`
	IntegrationReference integration_reference.IntegrationReference `json:"-"`
	Permissions          UserPermissions                            `json:"permissions"`
}

type UserPermissions struct {
	DashboardView string `json:"dashboard_view"`
}

// User represents a user of the system
type User struct {
	Id string `gorm:"primaryKey"`
	// AccountID is a reference to Account that this object belongs
	AccountID     string `json:"-" gorm:"index"`
	Role          UserRole
	IsServiceUser bool
	// NonDeletable indicates whether the service user can be deleted
	NonDeletable bool
	// ServiceUserName is only set if IsServiceUser is true
	ServiceUserName string
	// AutoGroups is a list of Group IDs to auto-assign to peers registered by this user
	AutoGroups []string                             `gorm:"serializer:json"`
	PATs       map[string]*pats.PersonalAccessToken `gorm:"-"`
	PATsG      []pats.PersonalAccessToken           `json:"-" gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE;"`
	// Blocked indicates whether the user is blocked. Blocked users can't use the system.
	Blocked bool
	// LastLogin is the last time the user logged in to IdP
	LastLogin *time.Time
	// CreatedAt records the time the user was created
	CreatedAt time.Time

	// Issued of the user
	Issued string `gorm:"default:api"`

	IntegrationReference integration_reference.IntegrationReference `gorm:"embedded;embeddedPrefix:integration_ref_"`
}

// IsBlocked returns true if the user is blocked, false otherwise
func (u *User) IsBlocked() bool {
	return u.Blocked
}

func (u *User) LastDashboardLoginChanged(lastLogin time.Time) bool {
	return lastLogin.After(u.GetLastLogin()) && !u.GetLastLogin().IsZero()
}

// GetLastLogin returns the last login time of the user.
func (u *User) GetLastLogin() time.Time {
	if u.LastLogin != nil {
		return *u.LastLogin
	}
	return time.Time{}
}

// HasAdminPower returns true if the user has admin or owner roles, false otherwise
func (u *User) HasAdminPower() bool {
	return u.Role == UserRoleAdmin || u.Role == UserRoleOwner
}

// IsAdminOrServiceUser checks if the user has admin power or is a service user.
func (u *User) IsAdminOrServiceUser() bool {
	return u.HasAdminPower() || u.IsServiceUser
}

// IsRegularUser checks if the user is a regular user.
func (u *User) IsRegularUser() bool {
	return !u.HasAdminPower() && !u.IsServiceUser
}

// Copy the user
func (u *User) Copy() *User {
	autoGroups := make([]string, len(u.AutoGroups))
	copy(autoGroups, u.AutoGroups)
	pats := make(map[string]*pats.PersonalAccessToken, len(u.PATs))
	for k, v := range u.PATs {
		pats[k] = v.Copy()
	}
	return &User{
		Id:                   u.Id,
		AccountID:            u.AccountID,
		Role:                 u.Role,
		AutoGroups:           autoGroups,
		IsServiceUser:        u.IsServiceUser,
		NonDeletable:         u.NonDeletable,
		ServiceUserName:      u.ServiceUserName,
		PATs:                 pats,
		Blocked:              u.Blocked,
		LastLogin:            u.LastLogin,
		CreatedAt:            u.CreatedAt,
		Issued:               u.Issued,
		IntegrationReference: u.IntegrationReference,
	}
}

// NewUser creates a new user
func NewUser(id string, role UserRole, isServiceUser bool, nonDeletable bool, serviceUserName string, autoGroups []string, issued string) *User {
	return &User{
		Id:              id,
		Role:            role,
		IsServiceUser:   isServiceUser,
		NonDeletable:    nonDeletable,
		ServiceUserName: serviceUserName,
		AutoGroups:      autoGroups,
		Issued:          issued,
		CreatedAt:       time.Now().UTC(),
	}
}

// NewRegularUser creates a new user with role UserRoleUser
func NewRegularUser(id string) *User {
	return NewUser(id, UserRoleUser, false, false, "", []string{}, UserIssuedAPI)
}

// NewAdminUser creates a new user with role UserRoleAdmin
func NewAdminUser(id string) *User {
	return NewUser(id, UserRoleAdmin, false, false, "", []string{}, UserIssuedAPI)
}

// NewOwnerUser creates a new user with role UserRoleOwner
func NewOwnerUser(id string) *User {
	return NewUser(id, UserRoleOwner, false, false, "", []string{}, UserIssuedAPI)
}
