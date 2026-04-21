package handler

import (
	"encoding/json"
	"go-loan-management-api/internal/handler/requests"
	"go-loan-management-api/internal/response"
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
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req requests.CreateCustomerRequest

	// Decode the incoming JSON body into the request struct.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Call the service layer to apply business logic.
	customer, err := h.service.CreateCustomer(req.FullName, req.Email, req.Phone)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return the created customer as JSON.
	// Return the created customer as JSON.
	response.Success(w, http.StatusCreated, "customer created successfully", customer)
}

// ListCustomers handles GET /customers.
func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests.
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	customers := h.service.ListCustomers()
	response.Success(w, http.StatusOK, "customers fetched successfully", customers)
}
