package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/therealtoxicdev/chronoserve/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	Roles []string `json:"roles"`
}

// HandleLogin processes login requests and returns JWT tokens
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate credentials against config
	user, valid := validateCredentials(req.Username, req.Password)
	if !valid {
		utils.WriteErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token - Remove the ... since CreateToken now accepts []string
	token, err := CreateToken(req.Username, user.Roles)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// Return token to client
	response := LoginResponse{
		Token: token,
		Roles: user.Roles,
	}

	utils.WriteSuccessResponse(w, "Login successful", response)
}

// validateCredentials checks if the provided credentials are valid against the config
func validateCredentials(username, password string) (*utils.Credentials, bool) {
	config := utils.GetConfig()

	// Check if user exists
	user, exists := config.Auth.Users[username]
	if !exists {
		return nil, false
	}

	// Validate password
	if user.Password != password {
		return nil, false
	}

	return &user, true
}
