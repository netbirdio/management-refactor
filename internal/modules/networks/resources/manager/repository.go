package manager

import (
	"management/internal/modules/networks"
	"management/internal/shared/db"
)

type Repository interface {
	RunInTx(fn func(tx db.Transaction) error) error
	Using(tx db.Transaction) Repository
	CreateNetwork(tx db.Transaction, network *networks.Network) error
	UpdateNetwork(tx db.Transaction, network *networks.Network) error
	GetNetworkByID(tx db.Transaction, lockingStrength db.LockingStrength, accountID, networkID string) (*networks.Network, error)
	GetAccountNetworks(tx db.Transaction, lockingStrength db.LockingStrength, accountID string) ([]*networks.Network, error)
}

type repository struct {
	store *db.Store
}

func newRepository(s *db.Store) Repository {
	return &repository{store: s}
}

func (r *repository) RunInTx(fn func(tx db.Transaction) error) error {
	return r.store.RunInTx(fn)
}

func (r *repository) Using(tx db.Transaction) Repository {
	return &repository{store: r.store.Using(tx)}
}

func (r *repository) DeleteNetwork(tx db.Transaction, network *networks.Network) error {
	return r.store.Create(tx, network)
}

func (r *repository) UpdateNetwork(tx db.Transaction, network *networks.Network) error {
	return r.store.Update(tx, network)
}

func (r *repository) GetNetworkByID(tx db.Transaction, lockingStrength db.LockingStrength, accountID, networkID string) (*networks.Network, error) {
	var network networks.Network
	err := r.store.GetOne(tx, lockingStrength, &network, "account_id = ? AND network_id = ?", accountID, networkID)
	if err != nil {
		return nil, err
	}
	return &network, nil
}

func (r *repository) GetAccountNetworks(tx db.Transaction, lockingStrength db.LockingStrength, accountID string) ([]*networks.Network, error) {
	var networks []*networks.Network
	err := r.store.GetMany(tx, lockingStrength, &networks, "account_id = ?", accountID)
	if err != nil {
		return nil, err
	}
	return networks, nil
}

func (r *repository) CreateNetwork(tx db.Transaction, network *networks.Network) error {
	return r.store.Create(tx, network)
}
