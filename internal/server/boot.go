package server

// @note this file includes all the lower level dependencies, db, http and grpc server, metrics, logger, etc.

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"management/internal/controllers/network_map"
	"management/internal/modules/peers"
	"management/internal/shared/activity"
	"management/internal/shared/activity/sqlite"
	"management/internal/shared/api/rest"
	"management/internal/shared/db"
	"management/internal/shared/metrics"
	"management/internal/shared/permissions"
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

func (s *Server) Metrics() *metrics.AppMetrics {
	return Create(s, func() *metrics.AppMetrics {
		appMetrics, err := metrics.NewAppMetrics()
		if err != nil {
			log.Fatalf("error while creating app metrics: %s", err)
		}
		return appMetrics
	})
}

func (s *Server) Router() *mux.Router {
	return Create(s, func() *mux.Router {
		return mux.NewRouter()
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

func (s *Server) NetworkMapController() *network_map.Controller {
	return Create(s, func() *network_map.Controller {
		store := s.Store()
		metrics := s.Metrics()
		return network_map.NewController(store, metrics)
	})
}

func (s *Server) PermissionsManager() permissions.Manager {
	return Create(s, func() permissions.Manager {
		return permissions.NewManager()
	})
}

func (s *Server) PeersManager() *peers.Manager {
	return Create(s, func() *peers.Manager {
		store := s.Store()
		router := s.Router()
		permissionsManager := s.PermissionsManager()

		return peers.NewManager(store, router, permissionsManager)
	})
}
