package roles

import (
	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/modules"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

var Admin = RolePermissions{
	Role: users.UserRoleAdmin,
	AutoAllowNew: map[operations.Operation]bool{
		operations.Read:  true,
		operations.Write: true,
	},
	Permissions: Permissions{
		modules.Accounts: {
			operations.Read:  true,
			operations.Write: false,
		},
	},
}
