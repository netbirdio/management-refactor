package manager

import (
	"context"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/netbirdio/management-refactor/management/server/http/api"
	"github.com/netbirdio/management-refactor/management/server/store"
	"github.com/netbirdio/management-refactor/management/server/types"

	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/internals/shared/permissions"
)

type Manager struct {
	repo Repository
}

func NewManager(store *db.Store, router *mux.Router, permissionsManager permissions.Manager) *Manager {
	repo := newRepository(store)
	m := &Manager{repo: repo}
	api := newHandler(m, permissionsManager)
	api.RegisterEndpoints(router)
	return m
}

func (m *Manager) GetAllGroups(ctx context.Context, accountID, userID string) ([]*types.Group, error) {
	groups, err := m.repo.GetAccountGroups(ctx, store.LockingStrengthShare, accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account groups: %w", err)
	}

	return groups, nil
}

func (m *Manager) GetAllGroupsMap(ctx context.Context, accountID, userID string) (map[string]*types.Group, error) {
	groups, err := m.GetAllGroups(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}

	groupsMap := make(map[string]*types.Group)
	for _, group := range groups {
		groupsMap[group.ID] = group
	}

	return groupsMap, nil
}

func (m *Manager) AddResourceToGroup(ctx context.Context, accountID, userID, groupID string, resource *types.Resource) error {
	event, err := m.AddResourceToGroupInTransaction(ctx, m.store, accountID, userID, groupID, resource)
	if err != nil {
		return fmt.Errorf("error adding resource to group: %w", err)
	}

	event()

	return nil
}

func (m *Manager) AddResourceToGroupInTransaction(ctx context.Context, tx db.Transaction, accountID, userID, groupID string, resource *types.Resource) (func(), error) {
	err := transaction.AddResourceToGroup(ctx, accountID, groupID, resource)
	if err != nil {
		return nil, fmt.Errorf("error adding resource to group: %w", err)
	}

	group, err := transaction.GetGroupByID(ctx, store.LockingStrengthShare, accountID, groupID)
	if err != nil {
		return nil, fmt.Errorf("error getting group: %w", err)
	}

	// TODO: at some point, this will need to become a switch statement
	networkResource, err := transaction.GetNetworkResourceByID(ctx, store.LockingStrengthShare, accountID, resource.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting network resource: %w", err)
	}

	event := func() {
		m.accountManager.StoreEvent(ctx, userID, groupID, accountID, activity.ResourceAddedToGroup, group.EventMetaResource(networkResource))
	}

	return event, nil
}

func (m *Manager) RemoveResourceFromGroupInTransaction(ctx context.Context, transaction store.Store, accountID, userID, groupID, resourceID string) (func(), error) {
	err := transaction.RemoveResourceFromGroup(ctx, accountID, groupID, resourceID)
	if err != nil {
		return nil, fmt.Errorf("error removing resource from group: %w", err)
	}

	group, err := transaction.GetGroupByID(ctx, store.LockingStrengthShare, accountID, groupID)
	if err != nil {
		return nil, fmt.Errorf("error getting group: %w", err)
	}

	// TODO: at some point, this will need to become a switch statement
	networkResource, err := transaction.GetNetworkResourceByID(ctx, store.LockingStrengthShare, accountID, resourceID)
	if err != nil {
		return nil, fmt.Errorf("error getting network resource: %w", err)
	}

	event := func() {
		m.accountManager.StoreEvent(ctx, userID, groupID, accountID, activity.ResourceRemovedFromGroup, group.EventMetaResource(networkResource))
	}

	return event, nil
}

func (m *Manager) GetResourceGroupsInTransaction(ctx context.Context, transaction store.Store, lockingStrength store.LockingStrength, accountID, resourceID string) ([]*types.Group, error) {
	return transaction.GetResourceGroups(ctx, lockingStrength, accountID, resourceID)
}

func ToGroupsInfoMap(groups []*types.Group, idCount int) map[string][]api.GroupMinimum {
	groupsInfoMap := make(map[string][]api.GroupMinimum, idCount)
	groupsChecked := make(map[string]struct{}, len(groups)) // not sure why this is needed (left over from old implementation)
	for _, group := range groups {
		_, ok := groupsChecked[group.ID]
		if ok {
			continue
		}

		groupsChecked[group.ID] = struct{}{}
		for _, pk := range group.Peers {
			info := api.GroupMinimum{
				Id:             group.ID,
				Name:           group.Name,
				PeersCount:     len(group.Peers),
				ResourcesCount: len(group.Resources),
			}
			groupsInfoMap[pk] = append(groupsInfoMap[pk], info)
		}
		for _, rk := range group.Resources {
			info := api.GroupMinimum{
				Id:             group.ID,
				Name:           group.Name,
				PeersCount:     len(group.Peers),
				ResourcesCount: len(group.Resources),
			}
			groupsInfoMap[rk.ID] = append(groupsInfoMap[rk.ID], info)
		}
	}
	return groupsInfoMap
}
