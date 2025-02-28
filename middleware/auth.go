package middleware

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/therealtoxicdev/chronoserve/utils"
)

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
	logger *utils.Logger
	config AuthConfig
)

// InitAuth initializes the authentication configuration
func InitAuth(cfg AuthConfig) {
	var err error
	appCfg := utils.GetConfig()
	logger, err = utils.NewLogger(utils.LoggerOptions{
		Level:      utils.INFO,
		MaxSize:    10,
		MaxBackups: 5,
		Directory:  appCfg.Logging.Directory,
		Filename:   "auth.log",
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	config = cfg
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
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.IssuedBy,
			Subject:   userID,
		},
		UserID: userID,
		Roles:  roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

func validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
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
