package networks

import (
	"context"
	"management/internal/shared/db"
	"management/internal/shared/hook"
)

type NetworkEvent struct {
	hook.Event

	Context context.Context
	Tx      db.Transaction
	Network *Network
}
