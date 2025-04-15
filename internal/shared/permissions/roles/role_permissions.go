package roles

import (
	"management/internal/modules/users/types"
	"management/internal/shared/permissions/modules"
	"management/internal/shared/permissions/operations"
)

type RolePermissions struct {
	Role         types.UserRole
	Permissions  Permissions
	AutoAllowNew map[operations.Operation]bool
}

type Permissions map[modules.Module]map[operations.Operation]bool

var RolesMap = map[types.UserRole]RolePermissions{
	types.UserRoleOwner: Owner,
	types.UserRoleAdmin: Admin,
	types.UserRoleUser:  User,
}
