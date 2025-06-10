package roles

import (
	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/modules"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

type RolePermissions struct {
	Role         users.UserRole
	Permissions  Permissions
	AutoAllowNew map[operations.Operation]bool
}

type Permissions map[modules.Module]map[operations.Operation]bool

var RolesMap = map[users.UserRole]RolePermissions{
	users.UserRoleOwner: Owner,
	users.UserRoleAdmin: Admin,
	users.UserRoleUser:  User,
}
