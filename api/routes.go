package api

import (
	"net/http"

	"github.com/therealtoxicdev/chronoserve/middleware"
	"github.com/therealtoxicdev/chronoserve/services"
	"github.com/therealtoxicdev/chronoserve/utils"
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

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Initialize service handler based on OS
	var serviceHandler services.ServiceHandler
	switch utils.GetOperatingSystem() {
	case "linux":
		serviceHandler = services.NewSystemdService()
	case "windows":
		serviceHandler = services.NewWindowsService()
	default:
		panic("Unsupported operating system")
	}

	// Define routes
	routes := []Route{
		// Public endpoints
		{Path: "health", Handler: utils.HealthCheck, RequireAuth: false},
		{Path: "auth/login", Handler: middleware.HandleLogin, RequireAuth: false},

		// Protected service endpoints
		{Path: "services", Handler: serviceHandler.ListServices, RequireAuth: true, Roles: []string{"admin", "viewer"}},
		{Path: "services/start/", Handler: serviceHandler.StartService, RequireAuth: true, Roles: []string{"admin"}},
		{Path: "services/stop/", Handler: serviceHandler.StopService, RequireAuth: true, Roles: []string{"admin"}},
		{Path: "services/logs/", Handler: serviceHandler.ViewServiceLogs, RequireAuth: true, Roles: []string{"admin", "viewer"}},
		{Path: "services/status/", Handler: serviceHandler.GetServiceStatus, RequireAuth: true, Roles: []string{"admin", "viewer"}},
	}

	// Register routes
	for _, route := range routes {
		handler := route.Handler

		if route.RequireAuth {
			// Add authentication and role-based access
			chainedHandler := middleware.Chain(
				middleware.Recovery,
				middleware.Logger,
				middleware.AuthMiddleware,
				middleware.RequireAnyRole(route.Roles...), // Updated to use RequireAnyRole
			)(http.HandlerFunc(handler))

			// Convert http.Handler back to http.HandlerFunc
			handler = chainedHandler.ServeHTTP
		} else {
			// Only add basic middleware for public endpoints
			chainedHandler := middleware.Chain(
				middleware.Recovery,
				middleware.Logger,
			)(http.HandlerFunc(handler))

			// Convert http.Handler back to http.HandlerFunc
			handler = chainedHandler.ServeHTTP
		}

		// Preflight handling
		mux.HandleFunc(apiPrefix+route.Path, func(w http.ResponseWriter, r *http.Request) {
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

	return corsMiddleware(mux)
}

// healthHandler returns service health information
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	utils.HealthCheck(w, r)
}
