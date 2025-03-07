package api

import (
	"net/http"

	"github.com/toxic-development/sysmanix/api/auth"
)

// registerAuthRoutes registers authentication-related routes
func registerAuthRoutes(mux *http.ServeMux) {
	// Public endpoints
	registerRouteWithMiddleware(mux, "auth/login", auth.LoginHandler, false, nil)

	// User token management (authenticated)
	registerRouteWithMiddleware(mux, "auth/tokens", auth.ListUserTokensHandler, true, []string{"admin", "viewer"})
	registerRouteWithMiddleware(mux, "auth/tokens/revoke", auth.RevokeTokenHandler, true, []string{"admin", "viewer"})
	registerRouteWithMiddleware(mux, "auth/tokens/revoke-all", auth.RevokeAllUserTokensHandler, true, []string{"admin", "viewer"})
	registerRouteWithMiddleware(mux, "auth/tokens/refresh", auth.RefreshTokenHandler, true, []string{"admin", "viewer"})

	// Admin token management (admin only)
	registerRouteWithMiddleware(mux, "auth/admin/tokens", auth.AdminListAllTokensHandler, true, []string{"admin"})
	registerRouteWithMiddleware(mux, "auth/admin/tokens/user", auth.AdminListUserTokensHandler, true, []string{"admin"})
	registerRouteWithMiddleware(mux, "auth/admin/tokens/revoke", auth.AdminRevokeUserTokensHandler, true, []string{"admin"})
}
