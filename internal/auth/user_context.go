package auth

import "context"

type authenticatedUserContextKey string

const userContextKey authenticatedUserContextKey = "authenticated_user"

// WithAuthenticatedUser stores the authenticated user in request context.
func WithAuthenticatedUser(ctx context.Context, user *AuthenticatedUser) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetAuthenticatedUser retrieves the authenticated user from request context.
func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	value := ctx.Value(userContextKey)
	user, ok := value.(*AuthenticatedUser)
	return user, ok
}
