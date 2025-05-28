package db

import (
	"gorm.io/gorm"
)

type LockingStrength string

const (
	LockingStrengthUpdate      LockingStrength = "UPDATE"        // Strongest lock, preventing any changes by other transactions until your transaction completes.
	LockingStrengthShare       LockingStrength = "SHARE"         // Allows reading but prevents changes by other transactions.
	LockingStrengthNoKeyUpdate LockingStrength = "NO KEY UPDATE" // Similar to UPDATE but allows changes to related rows.
	LockingStrengthKeyShare    LockingStrength = "KEY SHARE"     // Protects against changes to primary/unique keys but allows other updates.
)

// Transaction interface (unchanged)
type Transaction interface {
	Commit() error
	Rollback() error
	AddEvent(event func())
	FlushEvents()
}

type TransactionalManager[T any] interface {
	UsingTx(tx Transaction) T
}

type storeTx struct {
	db     *gorm.DB
	events []func()
}

func (tx *storeTx) Commit() error {
	for _, e := range tx.events {
		e()
	}
	return tx.db.Commit().Error
}

func (tx *storeTx) Rollback() error {
	return tx.db.Rollback().Error
}

func (tx *storeTx) AddEvent(fn func()) {
	tx.events = append(tx.events, fn)
}

func (tx *storeTx) FlushEvents() {
	for _, e := range tx.events {
		e()
	}
	tx.events = nil
}

func WithTx(store *Store, parentTx Transaction, fn func(tx Transaction) error) error {
	if parentTx != nil {
		return fn(parentTx)
	}

	return store.RunInTx(func(tx Transaction) error {
		if err := fn(tx); err != nil {
			return err
		}

		tx.FlushEvents()
		return nil
	})
}
