package httpx

import (
	"net/http"

	"go-loan-management-api/internal/response"
)

// MethodHandler maps HTTP methods to handlers.
type MethodHandler map[string]http.HandlerFunc

// HandleByMethod dispatches a request to the correct handler by HTTP method.
func HandleByMethod(handlers MethodHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := handlers[r.Method]; ok {
			handler(w, r)
			return
		}

		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
