package server

import (
	"context"
	"time"

	"management/pkg/logging"
)

// Server holds the HTTP server instance.
// Add any additional fields you need, such as database connections, config, etc.
type Server struct {
	// container of dependencies, each dependency is identified by a unique string.
	container map[string]any
}

var log = logging.LoggerForThisPackage()

// NewServer initializes and configures a new Server instance
func NewServer() *Server {
	return &Server{
		// @todo shared config
		container: make(map[string]any),
	}
}

// Start begins listening for HTTP requests on the configured address
func (s *Server) Start() error {
	// @todo instead of specifically starting httpserver
	// have a supervised start/stop of dependencies instead.
	// e.g. http, grpc, metrics, crons, etc
	return s.HttpServer().ListenAndServe()
}

// Stop attempts a graceful shutdown, waiting up to 5 seconds for active connections to finish
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.HttpServer().Shutdown(ctx)
}
