package handler

import (
	"go-loan-management-api/internal/response"
	"net/http"
)

// HealthHandler is a simple HTTP handler for checking API status.
type HealthHandler struct{}

// GetHealth responds with a simple JSON status.
// This confirms that the server is running.
func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "Service is running", map[string]string{
		"status": "ok",
	})

}
