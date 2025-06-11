package manager

import (
	"github.com/netbirdio/management-refactor/internals/modules/networks/resources"
	"github.com/netbirdio/management-refactor/internals/shared/db"
)

type Repository interface {
	Store() *db.Store
	GetResourcesByNetworkID(tx db.Transaction, strength db.LockingStrength, networkID string) ([]*resources.NetworkResource, error)
	DeleteResource(tx db.Transaction, resource *resources.NetworkResource) error
}

type repository struct {
	store *db.Store
}

func NewRepository(s *db.Store) Repository {
	return &repository{store: s}
}

func (r *repository) Store() *db.Store {
	return r.store
}

func (r *repository) GetResourcesByNetworkID(tx db.Transaction, strength db.LockingStrength, networkID string) ([]*resources.NetworkResource, error) {
	var resources []*resources.NetworkResource
	err := r.store.GetMany(tx, strength, &resources, "network_id = ?", networkID)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *repository) DeleteResource(tx db.Transaction, resource *resources.NetworkResource) error {
	return r.store.Delete(tx, resource)
}
