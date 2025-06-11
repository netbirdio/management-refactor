package roles

import (
	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/modules"
	"github.com/netbirdio/management-refactor/internals/shared/permissions/operations"
)

var NetworkAdmin = RolePermissions{
	Role: users.UserRoleNetworkAdmin,
	AutoAllowNew: map[operations.Operation]bool{
		operations.Read:  false,
		operations.Write: false,
	},
	Permissions: Permissions{
		modules.Networks: {
			operations.Read:  true,
			operations.Write: true,
		},
		modules.Groups: {
			operations.Read:  true,
			operations.Write: true,
		},
		modules.Settings: {
			operations.Read:  true,
			operations.Write: false,
		},
		modules.Accounts: {
			operations.Read:  true,
			operations.Write: false,
		},
		modules.Dns: {
			operations.Read:  true,
			operations.Write: true,
		},
		modules.Nameservers: {
			operations.Read:  true,
			operations.Write: true,
		},
		modules.Events: {
			operations.Read:  true,
			operations.Write: false,
		},
		modules.Policies: {
			operations.Read:  true,
			operations.Write: true,
		},
		modules.Routes: {
			operations.Read:  true,
			operations.Write: true,
		},
		modules.Users: {
			operations.Read:  true,
			operations.Write: false,
		},
		modules.SetupKeys: {
			operations.Read:  true,
			operations.Write: false,
		},
		modules.Pats: {
			operations.Read:  true,
			operations.Write: true,
		},
		modules.Peers: {
			operations.Read:  true,
			operations.Write: false,
		},
	},
}
