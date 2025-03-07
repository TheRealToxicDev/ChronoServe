package middleware

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/toxic-development/sysmanix/config"
	"github.com/toxic-development/sysmanix/utils"
)

// TokenInfo stores information about active tokens
type TokenInfo struct {
	TokenID   string    // Unique identifier for the token (e.g. jti claim)
	UserID    string    // User ID associated with the token
	Roles     []string  // User roles
	IssuedAt  time.Time // When the token was issued
	ExpiresAt time.Time // When the token expires
}

// TokenStore manages active tokens
type TokenStore struct {
	tokens map[string]TokenInfo
	mu     sync.RWMutex
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string   `json:"uid"`
	Roles  []string `json:"roles"`
}

type AuthConfig struct {
	SecretKey     string        `yaml:"secretKey"`
	TokenDuration time.Duration `yaml:"tokenDuration"`
	IssuedBy      string        `yaml:"issuedBy"`
}

var (
	logger     *utils.Logger
	authConfig AuthConfig
	tokenStore = &TokenStore{
		tokens: make(map[string]TokenInfo),
	}
)

const tokenFileName = "tokens.json"

// getTokenFilePath returns the full path to the tokens file in the "assets" directory
func getTokenFilePath() string {
	return filepath.Join("assets", tokenFileName)
}

// loadTokens loads tokens from a file
func loadTokens() error {
	filePath := getTokenFilePath()
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No tokens file exists yet
		}
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&tokenStore.tokens)
}

// saveTokens saves tokens to a file
func saveTokens() error {
	filePath := getTokenFilePath()
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(tokenStore.tokens)
}

// InitAuth initializes the authentication configuration
func InitAuth(cfg AuthConfig) {
	var err error
	appCfg := config.GetConfig()
	logger, err = utils.NewLogger(utils.LoggerOptions{
		Level:      utils.GetLogLevel(appCfg.Logging.Level),
		MaxSize:    appCfg.Logging.MaxSize,
		MaxBackups: appCfg.Logging.MaxBackups,
		Directory:  appCfg.Logging.Directory,
		Filename:   "auth.log",
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	authConfig = cfg

	// Load tokens from file
	if err := loadTokens(); err != nil {
		panic(fmt.Sprintf("Failed to load tokens: %v", err))
	}

	// Start token cleanup routine
	go cleanupExpiredTokens()
}

// AuthMiddleware provides JWT authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractToken(r)
		if err != nil {
			logger.Error("Auth failed: %v", err)
			utils.WriteErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := validateToken(token)
		if err != nil {
			logger.Error("Token validation failed: %v", err)
			utils.WriteErrorResponse(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		ctx := AddClaimsToContext(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateToken generates a new JWT token
func CreateToken(userID string, roles []string) (string, error) {
	// Generate a unique token ID
	tokenID, err := utils.GenerateUUID()
	if err != nil {
		return "", fmt.Errorf("failed to generate token ID: %w", err)
	}

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(authConfig.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    authConfig.IssuedBy,
			Subject:   userID,
			ID:        tokenID, // Add JTI claim for token identification
		},
		UserID: userID,
		Roles:  roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(authConfig.SecretKey))

	if err == nil {
		// Store token information
		storeToken(TokenInfo{
			TokenID:   tokenID,
			UserID:    userID,
			Roles:     roles,
			IssuedAt:  time.Now(),
			ExpiresAt: time.Now().Add(authConfig.TokenDuration),
		})
	}

	return signedToken, err
}

func validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(authConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	// Use constant-time comparison
	if subtle.ConstantTimeCompare([]byte(parts[0]), []byte("Bearer")) != 1 {
		return "", fmt.Errorf("invalid authorization type")
	}

	return parts[1], nil
}

// RequireRole middleware checks if the user has the required role
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r.Context())
			if claims == nil {
				logger.Warn("No claims found in context")
				utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Check if user has the required role
			hasRole := false
			for _, userRole := range claims.Roles {
				if userRole == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				logger.Warn("Unauthorized role access attempt. Required: %s, Has: %v", role, claims.Roles)
				utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole helper function to check if the user has any of the required roles
func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r.Context())
			if claims == nil {
				logger.Warn("No claims found in context")
				utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Check if user has any of the required roles
			hasRole := false
			for _, requiredRole := range roles {
				for _, userRole := range claims.Roles {
					if userRole == requiredRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				logger.Warn("Unauthorized role access attempt. Required any of: %v, Has: %v", roles, claims.Roles)
				utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ListValidTokens returns a list of all valid tokens
func ListValidTokens() []TokenInfo {
	tokenStore.mu.RLock()
	defer tokenStore.mu.RUnlock()

	validTokens := make([]TokenInfo, 0, len(tokenStore.tokens))
	now := time.Now()

	for _, token := range tokenStore.tokens {
		if token.ExpiresAt.After(now) {
			validTokens = append(validTokens, token)
		}
	}

	return validTokens
}

// ListUserTokens returns all valid tokens for a specific user
func ListUserTokens(userID string) []TokenInfo {
	tokenStore.mu.RLock()
	defer tokenStore.mu.RUnlock()

	userTokens := make([]TokenInfo, 0)
	now := time.Now()

	for _, token := range tokenStore.tokens {
		if token.UserID == userID && token.ExpiresAt.After(now) {
			userTokens = append(userTokens, token)
		}
	}

	return userTokens
}

// storeToken adds a token to the token store
func storeToken(info TokenInfo) {
	tokenStore.mu.Lock()
	defer tokenStore.mu.Unlock()

	tokenStore.tokens[info.TokenID] = info
	logger.Debug("Token stored for user %s, expires at %v", info.UserID, info.ExpiresAt)

	// Save tokens to file
	if err := saveTokens(); err != nil {
		logger.Error("Failed to save tokens: %v", err)
	}
}

// cleanupExpiredTokens periodically removes expired tokens from the store
func cleanupExpiredTokens() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		removeExpired()
	}
}

// removeExpired removes all expired tokens from the token store
func removeExpired() {
	tokenStore.mu.Lock()
	defer tokenStore.mu.Unlock()

	now := time.Now()
	removed := 0

	for id, token := range tokenStore.tokens {
		if token.ExpiresAt.Before(now) {
			delete(tokenStore.tokens, id)
			removed++
		}
	}

	if removed > 0 {
		logger.Debug("Removed %d expired tokens", removed)
		// Save tokens to file
		if err := saveTokens(); err != nil {
			logger.Error("Failed to save tokens: %v", err)
		}
	}
}

// RevokeToken removes a specific token from the token store
func RevokeToken(tokenID string) bool {
	tokenStore.mu.Lock()
	defer tokenStore.mu.Unlock()

	if _, exists := tokenStore.tokens[tokenID]; exists {
		delete(tokenStore.tokens, tokenID)
		logger.Debug("Token %s revoked", tokenID)

		// Save tokens to file
		if err := saveTokens(); err != nil {
			logger.Error("Failed to save tokens: %v", err)
		}
		return true
	}

	return false
}

// RevokeUserTokens revokes all tokens for a specific user
func RevokeUserTokens(userID string) int {
	tokenStore.mu.Lock()
	defer tokenStore.mu.Unlock()

	count := 0
	for id, token := range tokenStore.tokens {
		if token.UserID == userID {
			delete(tokenStore.tokens, id)
			count++
		}
	}

	if count > 0 {
		logger.Debug("Revoked %d tokens for user %s", count, userID)
		// Save tokens to file
		if err := saveTokens(); err != nil {
			logger.Error("Failed to save tokens: %v", err)
		}
	}

	return count
}
