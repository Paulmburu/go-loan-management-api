package response

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the standard response format for this API.
// It can be used for both success and error responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSON writes a JSON response with the given HTTP status code.
func WriteJSON(w http.ResponseWriter, statusCode int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(payload)
}

// Success writes a standard success response.
// data interface{} allows this function to accept ANY Go type.
func Success(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	// The JSON encoder will use "Reflection" to figure out how to
	// format the specific type passed into 'data'.
	WriteJSON(w, statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error writes a standard error response.
func Error(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, APIResponse{
		Success: false,
		Message: message,
	})
}
