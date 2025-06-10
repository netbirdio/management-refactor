package manager

import (
	"context"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/netbirdio/management-refactor/management/server/integrations/extra_settings"
	types2 "github.com/netbirdio/management-refactor/management/server/types"

	"github.com/netbirdio/management-refactor/internals/modules/users"
	"github.com/netbirdio/management-refactor/internals/shared/activity"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
	"github.com/netbirdio/management-refactor/pkg/logging"
)

var log = logging.LoggerForThisPackage()

type Manager struct {
	repository           Repository
	extraSettingsManager extra_settings.Manager
	userManager          *users.Manager
	eventManager         *activity.Manager
}

func NewManager(store *db.Store, router *mux.Router, eventManager *activity.Manager, permissionsManager permissions.Manager, userManager *users.Manager, extraSettingsManager extra_settings.Manager) *Manager {
	repo := newRepository(store)
	m := &Manager{
		repository:           repo,
		extraSettingsManager: extraSettingsManager,
		userManager:          userManager,
		eventManager:         eventManager,
	}
	api := newHandler(m, permissionsManager)
	api.RegisterEndpoints(router)
	return m
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
	return m.extraSettingsManager.UpdateExtraSettings(ctx, accountID, userID, (*types2.ExtraSettings)(extraSettings))
}

func (m *Manager) UpdateSettings(ctx context.Context, tx db.Transaction, settings *types.Settings) (*types.Settings, error) {
	return m.repository.UpdateSettings(tx, settings)
}
