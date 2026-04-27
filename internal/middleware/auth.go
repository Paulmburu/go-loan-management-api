package middleware

import (
	"go-loan-management-api/internal/auth"
	"go-loan-management-api/internal/response"
	"net/http"
	"strings"
)

// AuthMiddleware validates JWT bearer tokens and stores claims in request context.
type AuthMiddleware struct {
	jwtSecret string
}

// NewAuthMiddleware creates a new AuthMiddleware.
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

// RequireAuth protects a route by requiring a valid Bearer token.
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			response.Error(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		if tokenString == "" {
			response.Error(w, http.StatusUnauthorized, "missing bearer token")
			return
		}

		// claims, err := auth.ParseToken(tokenString, m.jwtSecret)
		// if err != nil {
		// 	response.Error(w, http.StatusUnauthorized, "invalid or expired token")
		// 	return
		// }

		// ctx := auth.WithClaims(r.Context(), claims)
		// next.ServeHTTP(w, r.WithContext(ctx))

		claims, err := auth.ParseToken(tokenString, m.jwtSecret)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		authenticatedUser := auth.NewAuthenticatedUserFromClaims(claims)

		ctx := auth.WithClaims(r.Context(), claims)
		ctx = auth.WithAuthenticatedUser(ctx, authenticatedUser)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
