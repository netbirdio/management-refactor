package activity

import (
	"context"
	"time"

	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/pkg/configuration"
	"github.com/netbirdio/management-refactor/pkg/logging"
)

var log = logging.LoggerForThisPackage()

type Manager struct {
	cfg *config
	// eventStore is the event store
	eventStore Store
}

// NewManager creates a new activity manager
func NewManager(eventStore Store) *Manager {
	cfg, err := configuration.Parse[config]()
	if err != nil {
		log.Fatalf("failed to parse activity config: %v", err)
	}
	return &Manager{
		cfg:        cfg,
		eventStore: eventStore,
	}
}

func (m *Manager) StoreEvent(ctx context.Context, tx db.Transaction, initiatorID, targetID, accountID string, activityID ActivityDescriber, meta map[string]any) {
	if m.cfg.Enabled {
		eventFunc := func() {
			_, err := m.eventStore.Save(ctx, &Event{
				Timestamp:   time.Now().UTC(),
				Activity:    activityID,
				InitiatorID: initiatorID,
				TargetID:    targetID,
				AccountID:   accountID,
				Meta:        meta,
			})
			if err != nil {
				// todo add metric
				log.WithContext(ctx).Errorf("received an error while storing an activity event, error: %s", err)
			}
		}

		if tx != nil {
			tx.AddEvent(eventFunc)
			return
		}

		go eventFunc()
	}
}
