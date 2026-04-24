package middleware

import (
	"log"
	"net/http"
	"time"
)

// statusRecorder wraps http.ResponseWriter so we can capture the status code.
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before passing it through.
func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Logging logs the request ID, method, path, status code, and duration.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture start time before the request continues.
		start := time.Now()

		// Wrap the response writer so we can capture the final status code.
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default if WriteHeader is never called explicitly
		}

		// Continue the request chain using the wrapped writer.
		next.ServeHTTP(recorder, r)

		// Read request ID from context.
		requestID := GetRequestID(r)

		// Log after the handler finishes.
		log.Printf(
			"request_id=%s method=%s path=%s status=%d duration=%s",
			requestID,
			r.Method,
			r.URL.Path,
			recorder.statusCode,
			time.Since(start),
		)
	})
}
