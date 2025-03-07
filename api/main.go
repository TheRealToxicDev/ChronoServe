package api

import (
	"net/http"

	"github.com/toxic-development/sysmanix/config"
	"github.com/toxic-development/sysmanix/middleware"
)

const (
	apiPrefix = "/"
)

// Route represents an API route with its handler and required role
type Route struct {
	Path        string
	Handler     http.HandlerFunc
	RequireAuth bool
	Roles       []string
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Auth-Token, X-Request-Id, X-Request-Start")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// SetupRoutes configures and returns the main HTTP router
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	cfg := config.GetConfig()

	// Register all route groups
	registerBaseRoutes(mux)
	registerAuthRoutes(mux)
	registerServiceRoutes(mux)

	// Only register Swagger if enabled in config
	if cfg.API.EnableSwagger {
		registerSwaggerRoutes(mux)
	}

	return corsMiddleware(mux)
}

// registerRouteWithMiddleware registers a route with appropriate middleware
func registerRouteWithMiddleware(mux *http.ServeMux, path string, handler http.HandlerFunc, requireAuth bool, roles []string) {
	if requireAuth {
		// Add authentication and role-based access
		chainedHandler := middleware.Chain(
			middleware.Recovery,
			middleware.Logger,
			middleware.AuthMiddleware,
			middleware.RequireAnyRole(roles...),
		)(http.HandlerFunc(handler))

		// Convert http.Handler back to http.HandlerFunc
		finalHandler := chainedHandler.ServeHTTP

		registerRouteWithCORS(mux, path, finalHandler)
	} else {
		// Only add basic middleware for public endpoints
		chainedHandler := middleware.Chain(
			middleware.Recovery,
			middleware.Logger,
		)(http.HandlerFunc(handler))

		// Convert http.Handler back to http.HandlerFunc
		finalHandler := chainedHandler.ServeHTTP

		registerRouteWithCORS(mux, path, finalHandler)
	}
}

// registerRouteWithCORS registers a route with CORS handling
func registerRouteWithCORS(mux *http.ServeMux, path string, handler http.HandlerFunc) {
	mux.HandleFunc(apiPrefix+path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Auth-Token, X-Request-Id, X-Request-Start")
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	})
}
