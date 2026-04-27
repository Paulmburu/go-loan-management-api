package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-loan-management-api/internal/auth"
	"go-loan-management-api/internal/model"
	"go-loan-management-api/internal/service"
)

func TestRepaymentHandler_GetRepaymentByID_AllowsOwnedRepayment(t *testing.T) {
	loanRepo := &fakeLoanRepository{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			return model.Loan{
				ID:         10,
				CustomerID: 7,
				Status:    model.LoanStatusActive,
			}, nil
		},
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	repaymentRepo := &fakeRepaymentRepository{
		createFn: func(ctx context.Context, repayment model.Repayment) (model.Repayment, error) {
			return model.Repayment{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Repayment, error) {
			return model.Repayment{
				ID:     5,
				LoanID: 10,
				Amount: 2000,
			}, nil
		},
		listByLoanIDFn: func(ctx context.Context, loanID int) ([]model.Repayment, error) { return nil, nil },
	}

	repaymentService := service.NewRepaymentService(nil, loanRepo, repaymentRepo)
	handler := NewRepaymentHandler(repaymentService)

	customerID := 7
	claims := &auth.Claims{
		UserID:     1,
		Email:      "customer@example.com",
		Role:       model.RoleCustomer,
		CustomerID: &customerID,
	}

	req := httptest.NewRequest(http.MethodGet, "/repayments/5", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()

	handler.GetRepaymentByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "repayment fetched successfully") {
		t.Fatalf("expected success message, got %q", rec.Body.String())
	}
}

func TestRepaymentHandler_GetRepaymentByID_RejectsOtherCustomersRepayment(t *testing.T) {
	loanRepo := &fakeLoanRepository{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			return model.Loan{
				ID:         10,
				CustomerID: 99,
				Status:     model.LoanStatusActive,
			}, nil
		},
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	repaymentRepo := &fakeRepaymentRepository{
		createFn: func(ctx context.Context, repayment model.Repayment) (model.Repayment, error) {
			return model.Repayment{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Repayment, error) {
			return model.Repayment{
				ID:     5,
				LoanID: 10,
				Amount: 2000,
			}, nil
		},
		listByLoanIDFn: func(ctx context.Context, loanID int) ([]model.Repayment, error) { return nil, nil },
	}

	repaymentService := service.NewRepaymentService(nil, loanRepo, repaymentRepo)
	handler := NewRepaymentHandler(repaymentService)

	customerID := 7
	claims := &auth.Claims{
		UserID:     1,
		Email:      "customer@example.com",
		Role:       model.RoleCustomer,
		CustomerID: &customerID,
	}

	req := httptest.NewRequest(http.MethodGet, "/repayments/5", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()

	handler.GetRepaymentByID(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "forbidden") {
		t.Fatalf("expected forbidden response, got %q", rec.Body.String())
	}
}
