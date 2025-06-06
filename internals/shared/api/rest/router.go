package rest

import (
	"github.com/gorilla/mux"

	"management/internal/shared/api/rest/middleware"
)

// NewRouter creates and returns a mux.Router configured with default middleware
// and placeholder endpoints. You can add your own handlers here or in other files.
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	return r
}
