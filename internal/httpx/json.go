package httpx

import (
	"encoding/json"
	"net/http"

	"go-loan-management-api/internal/response"
)

// DecodeJSON decodes request JSON into dst and writes a standard 400 response on failure.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return false
	}

	return true
}
