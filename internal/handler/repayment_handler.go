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

type RepaymentHandler struct {
	service *service.RepaymentService // You can add dependencies here, such as a service layer for handling business logic.
}

func NewRepaymentHandler(service *service.RepaymentService) *RepaymentHandler {
	return &RepaymentHandler{service: service}
}

// AddRepayment handles POST /repayments.
func (h *RepaymentHandler) AddRepayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req requests.CreateRepaymentRequest

	// Decode request body.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	repayment, err := h.service.AddRepayment(r.Context(), req.LoanID, req.Amount)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, "repayment added successfully", repayment)
}

// GetRepaymentByID handles GET /repayments/{id}.
func (h *RepaymentHandler) GetRepaymentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Expected path format: /repayments/{id}
	if len(parts) != 2 || parts[0] != "repayments" {
		response.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	repaymentID, err := strconv.Atoi(parts[1])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid repayment id")
		return
	}

	repayment, err := h.service.GetRepaymentByID(r.Context(), repaymentID)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "repayment fetched successfully", repayment)
}
