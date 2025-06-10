package roles

import (
	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

var Owner = RolePermissions{
	Role: users.UserRoleOwner,
	AutoAllowNew: map[operations.Operation]bool{
		operations.Read:  true,
		operations.Write: true,
	},
}
