package networks

import (
	"context"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/internals/shared/hook"
)

type NetworkEvent struct {
	hook.Event

	Context context.Context
	Tx      db.Transaction
	Network *Network
}
