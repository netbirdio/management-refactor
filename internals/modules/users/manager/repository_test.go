package manager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/netbirdio/management-refactor/internals/shared/db"
)

func TestGetUserByID(t *testing.T) {
	cfg := &db.Config{
		Engine:      "",
		PostgresDsn: "",
		DataDir:     "",
	}
	dbConn, err := db.NewDatabaseConn(context.Background(), cfg)
	assert.NoError(t, err)
	store := db.NewStore(context.Background(), dbConn)
	repo := NewRepository(store)

	result, err := repo.GetUserByID(nil, db.LockingStrengthShare, "5")
	assert.NoError(t, err)

	assert.Equal(t, "5", result.Id)
}
