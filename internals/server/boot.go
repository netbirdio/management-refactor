package server

// @note this file includes all the lower level dependencies, db, http and grpc server, metrics, logger, etc.

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"management/internal/shared/activity"
	"management/internal/shared/activity/sqlite"
	"management/internal/shared/api/rest"
	"management/internal/shared/db"
	"management/internal/shared/metrics"
)

func (s *server) Store() *db.Store {
	return Create(s, func() *db.Store {
		ctx := context.Background()
		dbConn, err := db.NewDatabaseConn(ctx)
		if err != nil {
			log.Fatalf("error while creating database connection: %s", err)
		}

		return db.NewStore(ctx, dbConn)
	})
}

func (s *server) HttpServer() *http.Server {
	return Create(s, func() *http.Server {
		router := s.Router()

		return &http.Server{
			Addr:    ":8080", // or from a config file
			Handler: router,
		}
	})
}

func (s *server) Metrics() *metrics.AppMetrics {
	return Create(s, func() *metrics.AppMetrics {
		appMetrics, err := metrics.NewAppMetrics()
		if err != nil {
			log.Fatalf("error while creating app metrics: %s", err)
		}
		return appMetrics
	})
}

func (s *server) Router() *mux.Router {
	return Create(s, func() *mux.Router {
		return rest.NewRouter()
	})
}

func (s *server) EventStore() activity.Store {
	return Create(s, func() activity.Store {
		ctx := context.Background()
		store, err := sqlite.NewSQLiteStore(ctx, "dataDir", "encryptionKey")
		if err != nil {
			log.Fatalf("error while creating event store: %s", err)
		}
		return store
	})
}
