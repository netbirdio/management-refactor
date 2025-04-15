package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"management/internal/shared/errors"
	"management/pkg/configuration"
)

const (
	storeSqliteFileName = "licenses.db"
	storeDataDirEnv     = "NB_STORE_DATA_DIR"
	storeDefaultDataDir = "/var/lib/netbird"
)

// DatabaseConn is a wrapper around the gorm database connection
type DatabaseConn struct {
	DB *gorm.DB
}

// NewDatabaseConn creates a new database connection based on the store engine
func NewDatabaseConn(ctx context.Context) (*DatabaseConn, error) {
	cfg, err := configuration.Parse[config]()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	log.WithContext(ctx).Infof("using %s store engine", cfg.Engine)

	var db *gorm.DB
	switch Engine(cfg.Engine) {
	case SqliteStoreEngine:
		db, err = openSQLiteDB()
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
func openSQLiteDB() (*gorm.DB, error) {
	storeStr := fmt.Sprintf("%s?cache=shared", storeSqliteFileName)
	if runtime.GOOS == "windows" {
		// To avoid `The process cannot access the file because it is being used by another process` on Windows
		storeStr = storeSqliteFileName
	}

	dataDir, ok := os.LookupEnv(storeDataDirEnv)
	if !ok {
		dataDir = storeDefaultDataDir
	}

	file := filepath.Join(dataDir, storeStr)
	db, err := gorm.Open(sqlite.Open(file), getGormConfig())
	if err != nil {
		return nil, err
	}

	return db, nil
}

// openPostgresDB opens a new connection to a Postgres database
func openPostgresDB(cfg *config) (*gorm.DB, error) {
	dsn, ok := os.LookupEnv(cfg.PostgresDsnEnv)
	if !ok {
		return nil, fmt.Errorf("%s is not set", cfg.PostgresDsnEnv)
	}

	db, err := gorm.Open(postgres.Open(dsn), getGormConfig())
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
