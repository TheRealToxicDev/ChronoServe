package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/toxic-development/sysmanix/middleware"
	"github.com/toxic-development/sysmanix/utils"
)

// @Summary      List user tokens
// @Description  Returns all valid tokens for the current user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  TokensResponse
// @Failure      401  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auth/tokens [get]
func ListUserTokensHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from JWT token
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		utils.WriteErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Get all tokens for the user
	tokens := middleware.ListUserTokens(claims.UserID)

	// Convert to response format
	tokenInfos := make([]TokenInfo, 0, len(tokens))
	for _, token := range tokens {
		tokenInfos = append(tokenInfos, TokenInfo{
			TokenID:   token.TokenID,
			UserID:    token.UserID,
			Roles:     token.Roles,
			IssuedAt:  token.IssuedAt.Format(time.RFC3339),
			ExpiresAt: token.ExpiresAt.Format(time.RFC3339),
		})
	}

	// Return the tokens
	response := TokensResponse{
		Status:  "success",
		Message: "User tokens retrieved successfully",
		Data:    tokenInfos,
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
}

// @Summary      Revoke token
// @Description  Revokes a specific token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RevokeTokenRequest true "Token to revoke"
// @Success      200  {object}  TokenActionResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auth/tokens/revoke [post]
func RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from JWT token
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		utils.WriteErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req RevokeTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TokenID == "" {
		utils.WriteErrorResponse(w, "TokenID is required", http.StatusBadRequest)
		return
	}

	// Check if the token belongs to the user unless admin
	isAdmin := false
	for _, role := range claims.Roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	// Only allow admins to revoke tokens that don't belong to them
	if !isAdmin {
		tokens := middleware.ListUserTokens(claims.UserID)
		isUserToken := false

		for _, token := range tokens {
			if token.TokenID == req.TokenID {
				isUserToken = true
				break
			}
		}

		if !isUserToken {
			utils.WriteErrorResponse(w, "Forbidden: You can only revoke your own tokens", http.StatusForbidden)
			return
		}
	}

	// Revoke the token
	success := middleware.RevokeToken(req.TokenID)

	response := TokenActionResponse{
		Status:  "success",
		Message: "Token revoked successfully",
	}

	if !success {
		response.Status = "warning"
		response.Message = "Token not found or already revoked"
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
}

// @Summary      Revoke all user tokens
// @Description  Revokes all tokens for the current user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  TokenActionResponse
// @Failure      401  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auth/tokens/revoke-all [post]
func RevokeAllUserTokensHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from JWT token
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		utils.WriteErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Revoke all tokens for the user
	count := middleware.RevokeUserTokens(claims.UserID)

	response := TokenActionResponse{
		Status:  "success",
		Message: "All tokens revoked successfully",
		Data: map[string]int{
			"count": count,
		},
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
}

// @Summary      Refresh token
// @Description  Generates a new token and invalidates the current one
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  RefreshTokenResponse
// @Failure      401  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auth/tokens/refresh [post]
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from JWT token
	claims := middleware.GetClaimsFromContext(r.Context())
	if claims == nil {
		utils.WriteErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Extract token ID from the current token
	tokenID := claims.ID
	if tokenID == "" {
		utils.WriteErrorResponse(w, "Invalid token: missing ID claim", http.StatusBadRequest)
		return
	}

	// Revoke current token
	middleware.RevokeToken(tokenID)

	// Generate new token
	newToken, err := middleware.CreateToken(claims.UserID, claims.Roles)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to create new token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := RefreshTokenResponse{
		Status:  "success",
		Message: "Token refreshed successfully",
		Data: struct {
			Token string `json:"token"`
		}{
			Token: newToken,
		},
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
}
