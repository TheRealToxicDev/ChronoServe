package api

import (
	"net/http"

	"github.com/toxic-development/sysmanix/api/auth"
	"github.com/toxic-development/sysmanix/api/base"
)

// registerBaseRoutes registers the basic routes
func registerBaseRoutes(mux *http.ServeMux) {
	// Public endpoints
	registerRouteWithMiddleware(mux, "health", base.HealthHandler, false, nil)
	registerRouteWithMiddleware(mux, "auth/login", auth.LoginHandler, false, nil)
}
