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
}

type storeTx struct {
	db *gorm.DB
}

func (t *storeTx) Commit() error {
	return t.db.Commit().Error
}

func (t *storeTx) Rollback() error {
	return t.db.Rollback().Error
}
