package server

// @note this file includes all the lower level dependencies, db, http and grpc server, metrics, logger, etc.

import (
	"context"
	"net/http"

	"management/internal/shared/activity"
	"management/internal/shared/activity/sqlite"
	"management/internal/shared/api/rest"
	"management/internal/shared/db"
)

func (s *Server) Store() *db.Store {
	return Create(s, func() *db.Store {
		ctx := context.Background()
		dbConn, err := db.NewDatabaseConn(ctx)
		if err != nil {
			log.Fatalf("error while creating database connection: %s", err)
		}

		return db.NewStore(ctx, dbConn)
	})
}

func (s *Server) HttpServer() *http.Server {
	return Create(s, func() *http.Server {
		router := rest.NewRouter()

		return &http.Server{
			Addr:    ":8080", // or from a config file
			Handler: router,
		}
	})
}

func (s *Server) EventStore() activity.Store {
	return Create(s, func() activity.Store {
		ctx := context.Background()
		store, err := sqlite.NewSQLiteStore(ctx, "dataDir", "encryptionKey")
		if err != nil {
			log.Fatalf("error while creating event store: %s", err)
		}
		return store
	})
}
