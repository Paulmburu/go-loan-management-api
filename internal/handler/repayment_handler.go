package handler

import (
	"encoding/json"
	"go-loan-management-api/internal/handler/requests"
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req requests.CreateRepaymentRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	repayment, err := h.service.AddRepayment(req.LoanID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(repayment)
}
