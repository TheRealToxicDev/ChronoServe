package auth

import (
	"net/http"

	"github.com/toxic-development/sysmanix/middleware"
)

// @Summary      Authenticate user
// @Description  Authenticates a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body LoginRequest true "Login credentials"
// @Success      200 {object} LoginResponse
// @Router       /auth/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Call the middleware login handler
	middleware.HandleLogin(w, r)
}
