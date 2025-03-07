package auth

// AdminRevokeRequest represents the request for admin token revocation
type AdminRevokeRequest struct {
	UserID string `json:"userId"`
}

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

// TokenInfo represents token information in responses
type TokenInfo struct {
	TokenID   string   `json:"tokenId"`
	UserID    string   `json:"userId"`
	Roles     []string `json:"roles"`
	IssuedAt  string   `json:"issuedAt"`
	ExpiresAt string   `json:"expiresAt"`
}

// TokensResponse represents the response for listing tokens
type TokensResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    []TokenInfo `json:"data"`
}

// TokenActionResponse represents the response for token actions like revocation
type TokenActionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// RevokeTokenRequest represents the request for token revocation
type RevokeTokenRequest struct {
	TokenID string `json:"tokenId"`
}

// RefreshTokenResponse represents the response for token refresh
type RefreshTokenResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"An error occurred while processing your request"`
	Code    int    `json:"code,omitempty" example:"404"`
}

// SuccessResponse represents a success response structure
type SuccessResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Operation completed successfully"`
	Data    any    `json:"data,omitempty"`
}
