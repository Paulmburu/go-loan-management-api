package middleware

import (
	"go-loan-management-api/internal/auth"
	"go-loan-management-api/internal/response"
	"net/http"
)

// RequireRoles allows access only to users with one of the allowed roles.
func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	allowed := map[string]struct{}{}

	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := auth.GetClaims(r.Context())
			if !ok {
				response.Error(w, http.StatusUnauthorized, "missing authentication claims")
				return
			}

			if _, exists := allowed[claims.Role]; !exists {
				response.Error(w, http.StatusForbidden, "forbidden: insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}

}
