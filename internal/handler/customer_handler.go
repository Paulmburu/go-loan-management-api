package handler

import (
	"encoding/json"
	"go-loan-management-api/internal/handler/requests"
	"go-loan-management-api/internal/service"
	"net/http"
)

// CustomerHandler handles customer-related HTTP requests.
type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// CreateCustomer handles POST /customers.
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	// Ensure the endpoint only accepts POST requests.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var req requests.CreateCustomerRequest

	// Decode the incoming JSON body into the request struct.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	// Call the service layer to apply business logic.
	customer, err := h.service.CreateCustomer(req.FullName, req.Email, req.Phone)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Return the created customer as JSON.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(customer)
}
