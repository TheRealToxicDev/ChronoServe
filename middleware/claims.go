package middleware

import (
	"context"
)

type contextKey string

const claimsContextKey contextKey = "claims"

// AddClaimsToContext adds JWT claims to the request context
func AddClaimsToContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

// GetClaimsFromContext retrieves JWT claims from the context
func GetClaimsFromContext(ctx context.Context) *Claims {
	if claims, ok := ctx.Value(claimsContextKey).(*Claims); ok {
		return claims
	}
	return nil
}
