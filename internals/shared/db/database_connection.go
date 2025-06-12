package db

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/netbirdio/management-refactor/internals/shared/errors"
)

const (
	storeSqliteFileName = "licenses.db"
)

// DatabaseConn is a wrapper around the gorm database connection
type DatabaseConn struct {
	DB *gorm.DB
}

// NewDatabaseConn creates a new database connection based on the store engine
func NewDatabaseConn(ctx context.Context, cfg *Config) (*DatabaseConn, error) {
	log.WithContext(ctx).Infof("using %s store engine", cfg.Engine)

	var db *gorm.DB
	var err error
	switch Engine(cfg.Engine) {
	case SqliteStoreEngine:
		db, err = openSQLiteDB(cfg)
	case PostgresStoreEngine:
		db, err = openPostgresDB(cfg)
	case MemoryStoreEngine:
		db, err = openMemoryDB()
	default:
		err = errors.NewUnsupportedStoreEngineConfigError(cfg.Engine)
	}

	if err != nil || db == nil {
		return nil, fmt.Errorf("error while opening database: %w", err)
	}

	sql, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting sql db connection: %w", err)
	}

	conns := runtime.NumCPU()
	if Engine(cfg.Engine) == SqliteStoreEngine {
		conns = 1
	}
	sql.SetMaxOpenConns(conns)
	return &DatabaseConn{
		DB: db,
	}, nil
}

// openMemoryDB opens a new connection to an in-memory SQLite database for testing
func openMemoryDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// openSQLiteDB opens a new connection to a SQLite database
func openSQLiteDB(cfg *Config) (*gorm.DB, error) {
	storeStr := fmt.Sprintf("%s?cache=shared", storeSqliteFileName)
	if runtime.GOOS == "windows" {
		// To avoid `The process cannot access the file because it is being used by another process` on Windows
		storeStr = storeSqliteFileName
	}

	file := filepath.Join(cfg.DataDir, storeStr)
	db, err := gorm.Open(sqlite.Open(file), getGormConfig())
	if err != nil {
		return nil, err
	}

	return db, nil
}

// openPostgresDB opens a new connection to a Postgres database
func openPostgresDB(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.PostgresDsn), getGormConfig())
	if err != nil {
		return nil, err
	}
	return db, nil
}

// getGormConfig returns the gorm configuration
func getGormConfig() *gorm.Config {
	return &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Silent),
		CreateBatchSize: 400,
		PrepareStmt:     false,
	}
}
