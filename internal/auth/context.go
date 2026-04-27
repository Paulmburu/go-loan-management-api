package auth

import "context"

type contextKey string

const claimsContextKey = contextKey("claims")

// WithClaims stores JWT claims in request context.
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

// GetClaims retrieves JWT claims from request context.
func GetClaims(ctx context.Context) (*Claims, bool) {
	value := ctx.Value(claimsContextKey)
	claims, ok := value.(*Claims)
	return claims, ok
}
