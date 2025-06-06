package roles

import (
	"github.com/netbirdio/netbird/management/server/types"

	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

var User = RolePermissions{
	Role: types.UserRoleUser,
	AutoAllowNew: map[operations.Operation]bool{
		operations.Read:  false,
		operations.Write: false,
	},
}
