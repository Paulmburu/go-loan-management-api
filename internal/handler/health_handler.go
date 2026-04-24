package handler

import (
	"context"
	"database/sql"
	"go-loan-management-api/internal/response"
	"net/http"
	"time"
)

// HealthHandler is a simple HTTP handler for checking API status.
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler creates a new HealthHandler.
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// GetHealth checks whether the service and database are healthy.
func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	// Use a short timeout so health checks fail fast if the DB is hanging.
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Ping the database using the request context.
	if err := h.db.PingContext(ctx); err != nil {
		response.Error(w, http.StatusServiceUnavailable, "database is unreachable")
		return
	}

	response.Success(w, http.StatusOK, "Service is running", map[string]string{
		"status":   "ok",
		"database": "up",
	})

}
