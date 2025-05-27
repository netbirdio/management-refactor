package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates and returns a mux.Router configured with default middleware
// and placeholder endpoints. You can add your own handlers here or in other files.
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Attach middlewares
	r.Use(loggingMiddleware)
	r.Use(recoveryMiddleware)

	// Example endpoint
	// r.HandleFunc("/health", healthCheckHandler).Methods("GET")

	return r
}

// loggingMiddleware is an example that logs each incoming request.
// Replace with your logger of choice.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("[%s] %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// recoveryMiddleware recovers from panics and returns a 500 Internal Server Error.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
