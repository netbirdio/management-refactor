package manager

import (
	"management/internal/modules/team"
	"management/internal/shared/hook"
)

var _ team.Manager = (*manager)(nil)

type manager struct {
	repository team.Repository

	onGroupCreate      *hook.Hook[*team.GroupEvent]
	onGroupCreateError *hook.Hook[*team.GroupErrorEvent]
	onGroupUpdate      *hook.Hook[*team.GroupEvent]
	onGroupUpdateError *hook.Hook[*team.GroupErrorEvent]
	onGroupDelete      *hook.Hook[*team.GroupEvent]
	onGroupDeleteError *hook.Hook[*team.GroupErrorEvent]

	// @todo user and pat events
}

func NewManager(repository team.Repository) *manager {
	m := &manager{
		repository: repository,
	}
	m.initGroupHooks()
	return m
}
