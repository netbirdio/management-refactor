package server

// @note this file includes all the lower level dependencies, db, http and grpc BaseServer, metrics, logger, etc.

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/netbirdio/management-refactor/internals/shared/activity"
	"github.com/netbirdio/management-refactor/internals/shared/activity/sqlite"
	"github.com/netbirdio/management-refactor/internals/shared/api/rest"
	"github.com/netbirdio/management-refactor/internals/shared/db"
	"github.com/netbirdio/management-refactor/internals/shared/metrics"
	"github.com/netbirdio/management-refactor/pkg/configuration"
)

func (s *BaseServer) Store() *db.Store {
	return Create(s, func() *db.Store {
		ctx := context.Background()

		cfg, err := configuration.Parse[db.Config]()
		if err != nil {
			log.Fatalf("failed to parse config: %v", err)
		}

		dbConn, err := db.NewDatabaseConn(ctx, cfg)
		if err != nil {
			log.Fatalf("error while creating database connection: %s", err)
		}

		return db.NewStore(ctx, dbConn)
	})
}

func (s *BaseServer) HttpServer() *http.Server {
	return Create(s, func() *http.Server {
		router := s.Router()

		return &http.Server{
			Addr:    ":8080", // or from a config file
			Handler: router,
		}
	})
}

func (s *BaseServer) Metrics() *metrics.AppMetrics {
	return Create(s, func() *metrics.AppMetrics {
		appMetrics, err := metrics.NewAppMetrics()
		if err != nil {
			log.Fatalf("error while creating app metrics: %s", err)
		}
		return appMetrics
	})
}

func (s *BaseServer) Router() *mux.Router {
	return Create(s, func() *mux.Router {
		return rest.NewRouter()
	})
}

func (s *BaseServer) ActivityManager() *activity.Manager {
	return Create(s, func() *activity.Manager {
		ctx := context.Background()
		store, err := sqlite.NewSQLiteStore(ctx, "dataDir", "encryptionKey")
		if err != nil {
			log.Fatalf("error while creating event store: %s", err)
		}

		return activity.NewManager(store)
	})
}
