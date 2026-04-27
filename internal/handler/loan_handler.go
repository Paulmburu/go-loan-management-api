package handler

import (
	"encoding/json"
	"go-loan-management-api/internal/auth"
	"go-loan-management-api/internal/dto"
	"go-loan-management-api/internal/model"
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

	var req dto.CreateLoanRequest

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

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "loans" || parts[2] != "balance" {
		response.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	loanID, err := strconv.Atoi(parts[1])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid loan id")
		return
	}

	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	loan, err := h.loanService.GetLoanByID(r.Context(), loanID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	if claims.Role == model.RoleCustomer && !auth.IsLoanOwnedByCustomer(claims, loan) {
		response.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	response.Success(w, http.StatusOK, "loan balance fetched successfully", map[string]float64{
		"outstanding_balance": loan.OutstandingAmount,
	})
}

func (h *LoanHandler) GetLoanByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "loans" {
		response.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	loanID, err := strconv.Atoi(parts[1])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid loan id")
		return
	}

	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	loan, err := h.loanService.GetLoanByID(r.Context(), loanID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	if claims.Role == model.RoleCustomer && !auth.IsLoanOwnedByCustomer(claims, loan) {
		response.Error(w, http.StatusForbidden, "forbidden")
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
	if len(parts) != 3 || parts[0] != "loans" || parts[2] != "repayments" {
		response.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	loanID, err := strconv.Atoi(parts[1])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid loan id")
		return
	}

	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	loan, err := h.loanService.GetLoanByID(r.Context(), loanID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	if claims.Role == model.RoleCustomer && !auth.IsLoanOwnedByCustomer(claims, loan) {
		response.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	repayments, err := h.repaymentService.ListRepaymentsByLoanID(r.Context(), loanID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "loan repayments fetched successfully", repayments)
}
