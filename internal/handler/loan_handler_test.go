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

type fakeRepaymentRepository struct {
	createFn       func(ctx context.Context, repayment model.Repayment) (model.Repayment, error)
	getByIDFn      func(ctx context.Context, id int) (model.Repayment, error)
	listByLoanIDFn func(ctx context.Context, loanID int) ([]model.Repayment, error)
}

func (f *fakeRepaymentRepository) Create(ctx context.Context, repayment model.Repayment) (model.Repayment, error) {
	return f.createFn(ctx, repayment)
}

func (f *fakeRepaymentRepository) GetByID(ctx context.Context, id int) (model.Repayment, error) {
	return f.getByIDFn(ctx, id)
}

func (f *fakeRepaymentRepository) ListByLoanID(ctx context.Context, loanID int) ([]model.Repayment, error) {
	return f.listByLoanIDFn(ctx, loanID)
}

func TestLoanHandler_GetLoanByID_AllowsOwnedLoan(t *testing.T) {
	customerRepo := &fakeCustomerRepository{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) { return model.Customer{}, nil },
		listFn:    func(ctx context.Context) ([]model.Customer, error) { return nil, nil },
	}

	loanRepo := &fakeLoanRepository{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			return model.Loan{
				ID:         10,
				CustomerID: 7,
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
		getByIDFn: func(ctx context.Context, id int) (model.Repayment, error) { return model.Repayment{}, nil },
		listByLoanIDFn: func(ctx context.Context, loanID int) ([]model.Repayment, error) {
			return nil, nil
		},
	}

	loanService := service.NewLoanService(customerRepo, loanRepo)
	repaymentService := service.NewRepaymentService(nil, loanRepo, repaymentRepo)
	handler := NewLoanHandler(loanService, repaymentService)

	customerID := 7
	claims := &auth.Claims{
		UserID:     1,
		Email:      "customer@example.com",
		Role:       model.RoleCustomer,
		CustomerID: &customerID,
	}

	req := httptest.NewRequest(http.MethodGet, "/loans/10", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()

	handler.GetLoanByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestLoanHandler_GetLoanByID_RejectsOtherCustomersLoan(t *testing.T) {
	customerRepo := &fakeCustomerRepository{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) { return model.Customer{}, nil },
		listFn:    func(ctx context.Context) ([]model.Customer, error) { return nil, nil },
	}

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
		getByIDFn: func(ctx context.Context, id int) (model.Repayment, error) { return model.Repayment{}, nil },
		listByLoanIDFn: func(ctx context.Context, loanID int) ([]model.Repayment, error) {
			return nil, nil
		},
	}

	loanService := service.NewLoanService(customerRepo, loanRepo)
	repaymentService := service.NewRepaymentService(nil, loanRepo, repaymentRepo)
	handler := NewLoanHandler(loanService, repaymentService)

	customerID := 7
	claims := &auth.Claims{
		UserID:     1,
		Email:      "customer@example.com",
		Role:       model.RoleCustomer,
		CustomerID: &customerID,
	}

	req := httptest.NewRequest(http.MethodGet, "/loans/10", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()

	handler.GetLoanByID(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "forbidden") {
		t.Fatalf("expected forbidden response, got %q", rec.Body.String())
	}
}

func TestLoanHandler_GetLoanRepayments_AllowsOwnedLoan(t *testing.T) {
	customerRepo := &fakeCustomerRepository{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) { return model.Customer{}, nil },
		listFn:    func(ctx context.Context) ([]model.Customer, error) { return nil, nil },
	}

	loanRepo := &fakeLoanRepository{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			return model.Loan{
				ID:         10,
				CustomerID: 7,
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
		getByIDFn: func(ctx context.Context, id int) (model.Repayment, error) { return model.Repayment{}, nil },
		listByLoanIDFn: func(ctx context.Context, loanID int) ([]model.Repayment, error) {
			return []model.Repayment{
				{ID: 1, LoanID: loanID, Amount: 2000},
			}, nil
		},
	}

	loanService := service.NewLoanService(customerRepo, loanRepo)
	repaymentService := service.NewRepaymentService(nil, loanRepo, repaymentRepo)
	handler := NewLoanHandler(loanService, repaymentService)

	customerID := 7
	claims := &auth.Claims{
		UserID:     1,
		Email:      "customer@example.com",
		Role:       model.RoleCustomer,
		CustomerID: &customerID,
	}

	req := httptest.NewRequest(http.MethodGet, "/loans/10/repayments", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()

	handler.GetLoanRepayments(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "loan repayments fetched successfully") {
		t.Fatalf("expected success message, got %q", rec.Body.String())
	}
}
