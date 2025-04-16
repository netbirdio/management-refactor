package settings

import (
	"context"
	"fmt"

	"github.com/netbirdio/netbird/management/server/integrations/extra_settings"

	"management/internal/modules/accounts/settings/types"
	"management/internal/modules/users"
	"management/internal/shared/db"
	"management/pkg/logging"
)

var log = logging.LoggerForThisPackage()

type Manager struct {
	repository           Repository
	extraSettingsManager extra_settings.Manager
	userManager          *users.Manager
}

func NewManager(store *db.Store, userManager *users.Manager, extraSettingsManager extra_settings.Manager) *Manager {
	return &Manager{
		repository:           newRepository(store),
		extraSettingsManager: extraSettingsManager,
		userManager:          userManager,
	}
}

func (m *Manager) GetExtraSettingsManager() extra_settings.Manager {
	return m.extraSettingsManager
}

func (m *Manager) GetSettings(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) (*types.Settings, error) {
	extraSettings, err := m.extraSettingsManager.GetExtraSettings(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get extra settings: %w", err)
	}

	settings, err := m.repository.GetAccountSettings(tx, strength, accountID)
	if err != nil {
		return nil, fmt.Errorf("get account settings: %w", err)
	}

	// Once we migrate the peer approval to settings manager this merging is obsolete
	if settings.Extra != nil {
		settings.Extra.FlowEnabled = extraSettings.FlowEnabled
		settings.Extra.FlowPacketCounterEnabled = extraSettings.FlowPacketCounterEnabled
		settings.Extra.FlowENCollectionEnabled = extraSettings.FlowENCollectionEnabled
		settings.Extra.FlowDnsCollectionEnabled = extraSettings.FlowDnsCollectionEnabled
	}

	return settings, nil
}

func (m *Manager) GetExtraSettings(ctx context.Context, tx db.Transaction, accountID string) (*types.ExtraSettings, error) {
	extraSettings, err := m.extraSettingsManager.GetExtraSettings(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get extra settings: %w", err)
	}

	settings, err := m.repository.GetAccountSettings(tx, db.LockingStrengthShare, accountID)
	if err != nil {
		return nil, fmt.Errorf("get account settings: %w", err)
	}

	// Once we migrate the peer approval to settings manager this merging is obsolete
	if settings.Extra == nil {
		settings.Extra = &types.ExtraSettings{}
	}

	settings.Extra.FlowEnabled = extraSettings.FlowEnabled

	return settings.Extra, nil
}

func (m *Manager) UpdateExtraSettings(ctx context.Context, accountID, userID string, extraSettings *types.ExtraSettings) (bool, error) {
	return m.extraSettingsManager.UpdateExtraSettings(ctx, accountID, userID, extraSettings)
}
