package auth

import (
	"net/http"

	"github.com/toxic-development/sysmanix/config"
	"github.com/toxic-development/sysmanix/middleware"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string   `json:"token"`
	Roles []string `json:"roles"`
}

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

// validateCredentials checks if the provided credentials are valid against the config
func validateCredentials(username, password string) (*config.Credentials, bool) {
	config := config.GetConfig()

	// Check if user exists
	user, exists := config.Auth.Users[username]
	if !exists {
		return nil, false
	}

	// Validate password
	valid, err := user.VerifyPassword(password)
	if err != nil {
		return nil, false
	}

	if !valid {
		return nil, false
	}

	return &user, true
}
