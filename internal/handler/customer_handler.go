package handler

import (
	"encoding/json"
	"go-loan-management-api/internal/auth"
	"go-loan-management-api/internal/dto"
	"go-loan-management-api/internal/response"
	"go-loan-management-api/internal/service"
	"net/http"
	"strconv"
	"strings"
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

	// var req requests.CreateCustomerRequest

	// // Decode the incoming JSON body into the request struct.
	// err := json.NewDecoder(r.Body).Decode(&req)
	// if err != nil {
	// 	response.Error(w, http.StatusBadRequest, "invalid request body")
	// 	return
	// }

	var req dto.CreateCustomerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Call the service layer to apply business logic.
	customer, err := h.service.CreateCustomer(r.Context(), req.FullName, req.Email, req.Phone)

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

	customers, err := h.service.ListCustomers(r.Context())

	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(w, http.StatusOK, "customers fetched successfully", customers)
}

// GetCustomerByID handles GET /customers/{id}.
func (h *CustomerHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Expected path format: /customers/{id}
	if len(parts) != 2 || parts[0] != "customers" {
		response.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	customerID, err := strconv.Atoi(parts[1])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	// claims, ok := auth.GetClaims(r.Context())
	// if !ok {
	// 	response.Error(w, http.StatusUnauthorized, "unauthenticated")
	// 	return
	// }

	// if claims.Role == model.RoleCustomer && !auth.IsSameCustomer(claims, customerID) {
	// 	response.Error(w, http.StatusForbidden, "forbidden")
	// 	return
	// }

	authenticatedUser, ok := auth.GetAuthenticatedUser(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	if authenticatedUser.IsCustomer() {
		linkedCustomerID, exists := authenticatedUser.CustomerIDValue()
		if !exists || linkedCustomerID != customerID {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	customer, err := h.service.GetCustomerByID(r.Context(), customerID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "customer fetched successfully", customer)
}

// GetCustomerLoans handles GET /customers/{id}/loans.
func (h *CustomerHandler) GetCustomerLoans(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Expected path format: /customers/{id}/loans
	if len(parts) != 3 || parts[0] != "customers" || parts[2] != "loans" {
		response.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	customerID, err := strconv.Atoi(parts[1])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	loans, err := h.service.GetCustomerLoans(r.Context(), customerID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "customer loans fetched successfully", loans)
}
