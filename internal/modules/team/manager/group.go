package manager

import (
	"context"
	"fmt"
	"management/internal/modules/team"
	"management/internal/shared/db"
	"management/internal/shared/hook"
)

// GetGroupById implements team.Manager.
func (m *manager) GetGroupById(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string, groupID string) (*team.Group, error) {
	panic("unimplemented")
}

// GetGroupsByAccount implements team.Manager.
func (m *manager) GetGroupsByAccount(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string) ([]*team.Group, error) {
	panic("unimplemented")
}

// GetGroupsByUser implements team.Manager.
func (m *manager) GetGroupsByUser(ctx context.Context, tx db.Transaction, strength db.LockingStrength, accountID string, userID string) ([]*team.Group, error) {
	panic("unimplemented")
}

// CreateGroup implements team.Manager.
func (m *manager) CreateGroup(ctx context.Context, tx db.Transaction, group *team.Group) (*team.Group, error) {
	panic("unimplemented")
}

// DeleteGroup implements team.Manager.
func (m *manager) DeleteGroup(ctx context.Context, tx db.Transaction, group *team.Group) error {
	panic("unimplemented")
}

// UpdateGroup implements team.Manager.
func (m *manager) UpdateGroup(ctx context.Context, tx db.Transaction, group *team.Group) (*team.Group, error) {
	err := db.WithTx(m.repository.Store(), tx, func(tx db.Transaction) error {
		ev := &team.GroupEvent{
			Tx:      tx,
			Context: ctx,
			Model:   group,
		}

		err := m.OnGroupUpdate().Trigger(ev, func(ge *team.GroupEvent) error {
			if err := m.repository.Store().Update(tx, ge.Model); err != nil {
				return fmt.Errorf("failed to update group: %w", err)
			}

			tx.AddEvent(func() {
				// addActivityEvent("Network deleted")
				// noop
			})
			return nil
		})

		return err
	})
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (m *manager) OnGroupUpdate() team.GroupHookEvent {
	return m.onGroupUpdate
}

func (m *manager) OnGroupUpdateError() team.GroupHookErrorEvent {
	return m.onGroupUpdateError
}

func (m *manager) OnGroupCreate() team.GroupHookEvent {
	panic("unimplemented")
}

func (m *manager) OnGroupCreateError() team.GroupHookErrorEvent {
	panic("unimplemented")
}

func (m *manager) OnGroupDelete() team.GroupHookEvent {
	panic("unimplemented")
}

func (m *manager) OnGroupDeleteError() team.GroupHookErrorEvent {
	panic("unimplemented")
}

func (m *manager) initGroupHooks() {
	m.onGroupCreate = &hook.Hook[*team.GroupEvent]{}
	m.onGroupCreateError = &hook.Hook[*team.GroupErrorEvent]{}
	m.onGroupUpdate = &hook.Hook[*team.GroupEvent]{}
	m.onGroupUpdateError = &hook.Hook[*team.GroupErrorEvent]{}
	m.onGroupDelete = &hook.Hook[*team.GroupEvent]{}
	m.onGroupDeleteError = &hook.Hook[*team.GroupErrorEvent]{}
}
