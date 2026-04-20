package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

// GetHealth responds with a simple JSON status.
// This confirms that the server is running.
func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"status": "ok"}

	// We use the blank identifier (_) here because if the JSON encoding fails
	// during a health check, the connection is likely already broken,
	// and there is no meaningful way to recover or log the error
	// without potentially causing a recursive failure.
	_ = json.NewEncoder(w).Encode(response)
}
