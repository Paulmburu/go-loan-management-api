package handler

import (
	"encoding/json"
	"go-loan-management-api/internal/handler/requests"
	"go-loan-management-api/internal/response"
	"go-loan-management-api/internal/service"
	"net/http"
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

	repayment, err := h.service.AddRepayment(req.LoanID, req.Amount)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(w, http.StatusCreated, "repayment added successfully", repayment)
}
