package middleware

import (
	"log"
	"net/http"

	"go-loan-management-api/internal/response"
)

// Recovery catches panics so the server can return a safe error response
// instead of crashing.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				// Read request ID from context so panic logs can be traced back to a request.
				requestID := GetRequestID(r)

				log.Printf(
					"request_id=%s method=%s path=%s panic=%v",
					requestID,
					r.Method,
					r.URL.Path,
					recovered,
				)

				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
