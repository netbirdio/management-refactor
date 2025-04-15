package db

import "slices"

// Engine represents the db engine to use.
type Engine string

const (
	SqliteStoreEngine   Engine = "sqlite"
	PostgresStoreEngine Engine = "postgres"
	MemoryStoreEngine   Engine = "memory"
	MysqlStoreEngine    Engine = "mysql"
)

var supportedEngines = []Engine{SqliteStoreEngine, PostgresStoreEngine, MysqlStoreEngine}

func IsSupportedEngine(engine Engine) bool {
	return slices.Contains(supportedEngines, engine)
}
