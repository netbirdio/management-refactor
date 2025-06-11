package roles

import (
	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

var Auditor = RolePermissions{
	Role: users.UserRoleAuditor,
	AutoAllowNew: map[operations.Operation]bool{
		operations.Read:  true,
		operations.Write: false,
	},
}
