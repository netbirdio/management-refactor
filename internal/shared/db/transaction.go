package db

import "gorm.io/gorm"

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
}

type storeTx struct {
	db     *gorm.DB
	events []func()
}

func (t *storeTx) Commit() error {
	err := t.db.Commit().Error
	if err != nil {
		t.commitEvents()
	}
	return err
}

func (t *storeTx) Rollback() error {
	return t.db.Rollback().Error
}

func (t *storeTx) AddEvent(event func()) {
	if t.events == nil {
		t.events = make([]func(), 0)
	}
	t.events = append(t.events, event)
}

func (t *storeTx) commitEvents() {
	for _, event := range t.events {
		event()
	}
	t.events = nil
}
