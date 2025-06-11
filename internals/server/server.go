package server

import (
	"context"
	"time"

	"github.com/netbirdio/management-refactor/pkg/logging"
)

type Server interface {
	Start() error
	Stop() error
	GetContainer(key string) (any, bool)
	SetContainer(key string, container any)
}

// Server holds the HTTP BaseServer instance.
// Add any additional fields you need, such as database connections, config, etc.
type BaseServer struct {
	// container of dependencies, each dependency is identified by a unique string.
	container map[string]any
	// AfterInit is a function that will be called after the server is initialized
	afterInit []func(s *BaseServer)
}

var log = logging.LoggerForThisPackage()

// NewServer initializes and configures a new Server instance
func NewServer() *BaseServer {
	return &BaseServer{
		container: make(map[string]any),
	}
}

func (s *BaseServer) AfterInit(fn func(s *BaseServer)) {
	s.afterInit = append(s.afterInit, fn)
}

// Start begins listening for HTTP requests on the configured address
func (s *BaseServer) Start() error {
	s.PeersManager()

	for _, fn := range s.afterInit {
		if fn != nil {
			fn(s)
		}
	}

	// @todo instead of specifically starting httpserver
	// have a supervised start/stop of dependencies instead.
	// e.g. http, grpc, metrics, crons, etc
	return s.HttpServer().ListenAndServe()
}

// Stop attempts a graceful shutdown, waiting up to 5 seconds for active connections to finish
func (s *BaseServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.HttpServer().Shutdown(ctx)
}

// GetContainer retrieves a dependency from the BaseServer's container by its key
func (s *BaseServer) GetContainer(key string) (any, bool) {
	container, exists := s.container[key]
	return container, exists
}

// SetContainer stores a dependency in the BaseServer's container with the specified key
func (s *BaseServer) SetContainer(key string, container any) {
	if _, exists := s.container[key]; exists {
		log.Tracef("container with key %s already exists", key)
		return
	}
	s.container[key] = container
	log.Tracef("container with key %s set successfully", key)
}
