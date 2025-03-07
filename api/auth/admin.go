package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/toxic-development/sysmanix/middleware"
	"github.com/toxic-development/sysmanix/utils"
)

// @Summary      List all tokens (admin)
// @Description  Returns all valid tokens in the system (admin only)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  TokensResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auth/admin/tokens [get]
func AdminListAllTokensHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get all valid tokens
	tokens := middleware.ListValidTokens()

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
		Message: "All tokens retrieved successfully",
		Data:    tokenInfos,
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
}

// @Summary      List user tokens (admin)
// @Description  Returns all valid tokens for a specific user (admin only)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        userId query string true "User ID"
// @Success      200  {object}  TokensResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auth/admin/tokens/user [get]
func AdminListUserTokensHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from query parameters
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		utils.WriteErrorResponse(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	// Get all tokens for the specified user
	tokens := middleware.ListUserTokens(userID)

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

// @Summary      Revoke user tokens (admin)
// @Description  Revokes all tokens for a specific user (admin only)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body AdminRevokeRequest true "User whose tokens to revoke"
// @Success      200  {object}  TokenActionResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /auth/admin/tokens/revoke [post]
func AdminRevokeUserTokensHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req AdminRevokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		utils.WriteErrorResponse(w, "UserID is required", http.StatusBadRequest)
		return
	}

	// Revoke all tokens for the specified user
	count := middleware.RevokeUserTokens(req.UserID)

	response := TokenActionResponse{
		Status:  "success",
		Message: "All user tokens revoked successfully",
		Data: map[string]any{
			"count":  count,
			"userId": req.UserID,
		},
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
}
