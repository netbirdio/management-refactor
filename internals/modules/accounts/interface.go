package accounts

import (
	"context"

	"github.com/netbirdio/netbird/management/server/integrations/extra_settings"
	"github.com/netbirdio/netbird/management/server/types"

	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Manager interface {
	GetExtraSettingsManager() extra_settings.Manager
	GetSettings(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) (*types.Settings, error)
	GetExtraSettings(ctx context.Context, tx db.Transaction, accountID string) (*types.ExtraSettings, error)
	UpdateExtraSettings(ctx context.Context, accountID, userID string, extraSettings *types.ExtraSettings) (bool, error)
	UpdateSettings(ctx context.Context, tx db.Transaction, settings *types.Settings) (*types.Settings, error)
}
