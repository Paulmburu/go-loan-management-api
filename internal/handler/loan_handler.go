package handler

import (
	"encoding/json"
	"go-loan-management-api/internal/handler/requests"
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var req requests.CreateLoanRequest

	// Decode the request body.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service layer.
	loan, err := h.service.CreateLoan(req.CustomerID, req.PrincipalAmount, req.InterestRate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the created loan.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(loan)
}

// ListLoans handles GET /loans.
func (h *LoanHandler) ListLoans(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	loans := h.service.ListLoans()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(loans)

}

// GetLoanBalance handles GET /loans/{id}/balance.
func (h *LoanHandler) GetLoanBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Expected path format: /loans/{id}/balance
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// We expect exactly 3 parts: loans, {id}, balance
	if len(parts) != 3 || parts[0] != "loans" || parts[2] != "balance" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	// Convert the loan ID from string to int.
	loanID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	// Ask the service layer for the balance.
	balance, err := h.service.GetLoanBalance(loanID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return the balance as JSON.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]float64{"outstanding_balance": balance})
}
