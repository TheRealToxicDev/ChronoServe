package api

import (
	"net/http"

	"github.com/toxic-development/sysmanix/api/services"
	"github.com/toxic-development/sysmanix/utils"
)

// registerServiceRoutes registers service-related routes
func registerServiceRoutes(mux *http.ServeMux) {
	// Initialize appropriate service handler based on OS
	switch utils.GetOperatingSystem() {
	case "linux":
		registerLinuxServiceRoutes(mux)
	case "windows":
		registerWindowsServiceRoutes(mux)
	default:
		panic("Unsupported operating system")
	}
}

// Register Linux-specific service routes
func registerLinuxServiceRoutes(mux *http.ServeMux) {
	// Protected service endpoints
	registerRouteWithMiddleware(mux, "services", services.ListServices, true, []string{"admin", "viewer"})
	registerRouteWithMiddleware(mux, "services/start/", services.StartService, true, []string{"admin"})
}

// Register Windows-specific service routes
func registerWindowsServiceRoutes(mux *http.ServeMux) {
	// Protected service endpoints - same handlers but could be different if needed
	registerRouteWithMiddleware(mux, "services", services.ListServices, true, []string{"admin", "viewer"})
	registerRouteWithMiddleware(mux, "services/start/", services.StartService, true, []string{"admin"})
}
