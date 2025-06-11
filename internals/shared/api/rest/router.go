package rest

import (
	"github.com/gorilla/mux"

	"github.com/netbirdio/management-refactor/internals/shared/api/rest/middleware"
)

// NewRouter creates and returns a mux.Router configured with default middleware
// and placeholder endpoints. You can add your own handlers here or in other files.
func NewRouter() *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	authMiddleware := middleware.NewAuthMiddleware()

	r.Use(authMiddleware.Handler)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	return r
}
