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

// Fake repositories used to build a real CustomerService for handler tests.
type fakeCustomerRepository struct {
	createFn  func(ctx context.Context, customer model.Customer) (model.Customer, error)
	getByIDFn func(ctx context.Context, id int) (model.Customer, error)
	listFn    func(ctx context.Context) ([]model.Customer, error)
}

func (f *fakeCustomerRepository) Create(ctx context.Context, customer model.Customer) (model.Customer, error) {
	return f.createFn(ctx, customer)
}

func (f *fakeCustomerRepository) GetByID(ctx context.Context, id int) (model.Customer, error) {
	return f.getByIDFn(ctx, id)
}

func (f *fakeCustomerRepository) List(ctx context.Context) ([]model.Customer, error) {
	return f.listFn(ctx)
}

type fakeLoanRepository struct {
	createFn           func(ctx context.Context, loan model.Loan) (model.Loan, error)
	getByIDFn          func(ctx context.Context, id int) (model.Loan, error)
	listFn             func(ctx context.Context) ([]model.Loan, error)
	listByCustomerIDFn func(ctx context.Context, customerID int) ([]model.Loan, error)
	updateFn           func(ctx context.Context, loan model.Loan) (model.Loan, error)
}

func (f *fakeLoanRepository) Create(ctx context.Context, loan model.Loan) (model.Loan, error) {
	return f.createFn(ctx, loan)
}

func (f *fakeLoanRepository) GetByID(ctx context.Context, id int) (model.Loan, error) {
	return f.getByIDFn(ctx, id)
}

func (f *fakeLoanRepository) List(ctx context.Context) ([]model.Loan, error) {
	return f.listFn(ctx)
}

func (f *fakeLoanRepository) ListByCustomerID(ctx context.Context, customerID int) ([]model.Loan, error) {
	return f.listByCustomerIDFn(ctx, customerID)
}

func (f *fakeLoanRepository) Update(ctx context.Context, loan model.Loan) (model.Loan, error) {
	return f.updateFn(ctx, loan)
}

func TestCustomerHandler_GetCustomerByID_AllowsOwnCustomerRecord(t *testing.T) {
	customerRepo := &fakeCustomerRepository{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{
				ID:       7,
				FullName: "Paul Mburu",
				Email:    "paul@example.com",
				Phone:    "0712345678",
			}, nil
		},
		listFn: func(ctx context.Context) ([]model.Customer, error) {
			return nil, nil
		},
	}

	loanRepo := &fakeLoanRepository{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) {
			return model.Loan{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			return model.Loan{}, nil
		},
		listFn: func(ctx context.Context) ([]model.Loan, error) {
			return nil, nil
		},
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) {
			return nil, nil
		},
		updateFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) {
			return model.Loan{}, nil
		},
	}

	customerService := service.NewCustomerService(customerRepo, loanRepo)
	handler := NewCustomerHandler(customerService)

	customerID := 7
	claims := &auth.Claims{
		UserID:     1,
		Email:      "paul@example.com",
		Role:       model.RoleCustomer,
		CustomerID: &customerID,
	}

	req := httptest.NewRequest(http.MethodGet, "/customers/7", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()

	handler.GetCustomerByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "customer fetched successfully") {
		t.Fatalf("expected success message, got %q", rec.Body.String())
	}
}

func TestCustomerHandler_GetCustomerByID_RejectsOtherCustomerRecord(t *testing.T) {
	customerRepo := &fakeCustomerRepository{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{
				ID:       99,
				FullName: "Other Customer",
			}, nil
		},
		listFn: func(ctx context.Context) ([]model.Customer, error) {
			return nil, nil
		},
	}

	loanRepo := &fakeLoanRepository{
		createFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn:          func(ctx context.Context, id int) (model.Loan, error) { return model.Loan{}, nil },
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	customerService := service.NewCustomerService(customerRepo, loanRepo)
	handler := NewCustomerHandler(customerService)

	customerID := 7
	claims := &auth.Claims{
		UserID:     1,
		Email:      "paul@example.com",
		Role:       model.RoleCustomer,
		CustomerID: &customerID,
	}

	req := httptest.NewRequest(http.MethodGet, "/customers/99", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()

	handler.GetCustomerByID(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "forbidden") {
		t.Fatalf("expected forbidden response, got %q", rec.Body.String())
	}
}

func TestCustomerHandler_GetCustomerLoans_AllowsAdmin(t *testing.T) {
	customerRepo := &fakeCustomerRepository{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{ID: 7, FullName: "Paul Mburu"}, nil
		},
		listFn: func(ctx context.Context) ([]model.Customer, error) { return nil, nil },
	}

	loanRepo := &fakeLoanRepository{
		createFn:  func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) { return model.Loan{}, nil },
		listFn:    func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) {
			return []model.Loan{
				{ID: 1, CustomerID: 7, PrincipalAmount: 10000},
			}, nil
		},
		updateFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	customerService := service.NewCustomerService(customerRepo, loanRepo)
	handler := NewCustomerHandler(customerService)

	claims := &auth.Claims{
		UserID: 1,
		Email:  "admin@example.com",
		Role:   model.RoleAdmin,
	}

	req := httptest.NewRequest(http.MethodGet, "/customers/7/loans", nil)
	req = req.WithContext(auth.WithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()

	handler.GetCustomerLoans(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "customer loans fetched successfully") {
		t.Fatalf("expected success message, got %q", rec.Body.String())
	}
}
