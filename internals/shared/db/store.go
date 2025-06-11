package db

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/netbirdio/management-refactor/pkg/logging"
)

var log = logging.LoggerForThisPackage

type Store struct {
	db *gorm.DB
}

func NewStore(ctx context.Context, dbConn *DatabaseConn) *Store {
	return &Store{db: dbConn.DB}
}

func (s *Store) Begin() (Transaction, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &storeTx{db: tx}, nil
}

// AutoMigrate automatically migrates your schema, to keep your database up to date.
func (s *Store) AutoMigrate(value interface{}) error {
	return s.db.AutoMigrate(value)
}

// Using picks the underlying db
func (s *Store) Using(tx Transaction) *gorm.DB {
	if tx == nil {
		return s.db
	}
	if st, ok := tx.(*storeTx); ok {
		return st.db
	}
	return s.db
}

// RunInTx is the new helper that starts a transaction, calls fn, and commits/rolls back automatically
func (s *Store) RunInTx(fn func(tx Transaction) error) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Store) Create(tx Transaction, value interface{}) error {
	return s.Using(tx).Create(value).Error
}

func (s *Store) GetOne(tx Transaction, strength LockingStrength, dest interface{}, query string, args ...interface{}) error {
	db := s.Using(tx).Clauses(clause.Locking{Strength: string(strength)})

	if query != "" && len(args) > 0 {
		db.Where(query, args...)
	}

	return db.First(dest).Error
}

func (s *Store) GetMany(tx Transaction, strength LockingStrength, dest interface{}, query string, args ...interface{}) error {
	db := s.Using(tx).Clauses(clause.Locking{Strength: string(strength)})

	if query != "" && len(args) > 0 {
		db.Where(query, args...)
	}

	return db.Find(dest).Error
}

func (s *Store) Delete(tx Transaction, value interface{}) error {
	return s.Using(tx).Delete(value).Error
}

func (s *Store) Update(tx Transaction, value interface{}) error {
	return s.Using(tx).Save(value).Error
}
