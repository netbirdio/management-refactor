package server

import (
	"context"
	"net/http"
	"time"

	"github.com/netbirdio/management-integrations/integrations"

	"management/internal/modules/settings"
	"management/internal/modules/users"
	"management/internal/shared/api"
	"management/internal/shared/api/rest"
	"management/internal/shared/db"
	"management/internal/shared/permissions"
	"management/pkg/logging"
)

// Server holds the HTTP server instance.
// Add any additional fields you need, such as database connections, config, etc.
type Server struct {
	httpServer *http.Server
}

var log = logging.LoggerForThisPackage()

// NewServer initializes and configures a new Server instance
func NewServer() *Server {
	ctx := context.Background()

	dbConn, err := db.NewDatabaseConn(ctx)
	if err != nil {
		log.Fatalf("error while creating database connection: %s", err)
	}

	store := db.NewStore(ctx, dbConn)

	router := rest.NewRouter()

	extraSettingsManager := integrations.NewManager()
	permissionsManager := permissions.NewManager(store)
	userManager := users.NewManager(store, permissions.NewManager(store))
	settingsManager := settings.NewManager(store, userManager, extraSettingsManager)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":8080", // or from a config file
			Handler: router,
		},
	}
}

// Start begins listening for HTTP requests on the configured address
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Stop attempts a graceful shutdown, waiting up to 5 seconds for active connections to finish
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
