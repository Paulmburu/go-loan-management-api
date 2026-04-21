package handler

import (
	"encoding/json"
	"go-loan-management-api/internal/handler/requests"
	"go-loan-management-api/internal/response"
	"go-loan-management-api/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type LoanHandler struct {
	service *service.LoanService
}

func NewLoanHandler(service *service.LoanService) *LoanHandler {
	return &LoanHandler{service: service}
}

// CreateLoan handles POST /loans.
func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	// Ensure the endpoint only accepts POST requests.
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req requests.CreateLoanRequest

	// Decode the request body.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Call the service layer.
	loan, err := h.service.CreateLoan(req.CustomerID, req.PrincipalAmount, req.InterestRate)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Return the created loan.
	response.Success(w, http.StatusCreated, "loan created successfully", loan)
}

// ListLoans handles GET /loans.
func (h *LoanHandler) ListLoans(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests.
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	loans := h.service.ListLoans()
	response.Success(w, http.StatusOK, "loans fetched successfully", loans)
}

// GetLoanBalance handles GET /loans/{id}/balance.
func (h *LoanHandler) GetLoanBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Expected path format: /loans/{id}/balance
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// We expect exactly 3 parts: loans, {id}, balance
	if len(parts) != 3 || parts[0] != "loans" || parts[2] != "balance" {
		response.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	// Convert the loan ID from string to int.
	loanID, err := strconv.Atoi(parts[1])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid loan id")
		return
	}

	// Ask the service layer for the balance.
	balance, err := h.service.GetLoanBalance(loanID)

	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	// Return the balance as JSON.
	response.Success(w, http.StatusOK, "loan balance fetched successfully", map[string]float64{
		"outstanding_balance": balance,
	})
}
