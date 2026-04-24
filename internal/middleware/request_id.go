package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// contextKey is a private type used to avoid collisions in request context keys.
type contextKey string

const requestIDKey contextKey = "request_id"

// RequestID adds a request ID to both the response header and request context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a lightweight unique-ish request ID.
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())

		// Add it to the response header so the client can also see it.
		w.Header().Set("X-Request-ID", requestID)

		// Add it to the request context so downstream middleware/handlers can read it.
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// GetRequestID returns the request ID from context if present.
func GetRequestID(r *http.Request) string {
	value := r.Context().Value(requestIDKey)
	if requestID, ok := value.(string); ok {
		return requestID
	}

	return ""
}
