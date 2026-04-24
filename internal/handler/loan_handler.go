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

// LoanHandler handles loan-related HTTP requests.
type LoanHandler struct {
	loanService      *service.LoanService
	repaymentService *service.RepaymentService
}

// NewLoanHandler creates a new LoanHandler.
func NewLoanHandler(
	loanService *service.LoanService,
	repaymentService *service.RepaymentService,
) *LoanHandler {
	return &LoanHandler{
		loanService:      loanService,
		repaymentService: repaymentService,
	}
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
	loan, err := h.loanService.CreateLoan(r.Context(), req.CustomerID, req.PrincipalAmount, req.InterestRate)
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

	loans, err := h.loanService.ListLoans(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
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
	balance, err := h.loanService.GetLoanBalance(r.Context(), loanID)

	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	// Return the balance as JSON.
	response.Success(w, http.StatusOK, "loan balance fetched successfully", map[string]float64{
		"outstanding_balance": balance,
	})
}

// GetLoanByID handles GET /loans/{id}.
func (h *LoanHandler) GetLoanByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Expected path format: /loans/{id}
	if len(parts) != 2 || parts[0] != "loans" {
		response.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	loanID, err := strconv.Atoi(parts[1])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid loan id")
		return
	}

	loan, err := h.loanService.GetLoanByID(r.Context(), loanID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "loan fetched successfully", loan)
}

// GetLoanRepayments handles GET /loans/{id}/repayments.
func (h *LoanHandler) GetLoanRepayments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Expected path format: /loans/{id}/repayments
	if len(parts) != 3 || parts[0] != "loans" || parts[2] != "repayments" {
		response.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	loanID, err := strconv.Atoi(parts[1])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid loan id")
		return
	}

	repayments, err := h.repaymentService.ListRepaymentsByLoanID(r.Context(), loanID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}

	response.Success(w, http.StatusOK, "loan repayments fetched successfully", repayments)
}
