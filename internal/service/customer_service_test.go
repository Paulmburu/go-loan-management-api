package service

import (
	"context"
	"errors"
	"testing"

	"go-loan-management-api/internal/apperror"
	"go-loan-management-api/internal/model"
)

type fakeCustomerRepositoryForCustomerService struct {
	createFn  func(ctx context.Context, customer model.Customer) (model.Customer, error)
	getByIDFn func(ctx context.Context, id int) (model.Customer, error)
	listFn    func(ctx context.Context) ([]model.Customer, error)
}

func (f *fakeCustomerRepositoryForCustomerService) Create(ctx context.Context, customer model.Customer) (model.Customer, error) {
	return f.createFn(ctx, customer)
}

func (f *fakeCustomerRepositoryForCustomerService) GetByID(ctx context.Context, id int) (model.Customer, error) {
	return f.getByIDFn(ctx, id)
}

func (f *fakeCustomerRepositoryForCustomerService) List(ctx context.Context) ([]model.Customer, error) {
	return f.listFn(ctx)
}

type fakeLoanRepositoryForCustomerService struct {
	listByCustomerIDFn func(ctx context.Context, customerID int) ([]model.Loan, error)
}

func (f *fakeLoanRepositoryForCustomerService) Create(ctx context.Context, loan model.Loan) (model.Loan, error) {
	return model.Loan{}, nil
}

func (f *fakeLoanRepositoryForCustomerService) GetByID(ctx context.Context, id int) (model.Loan, error) {
	return model.Loan{}, nil
}

func (f *fakeLoanRepositoryForCustomerService) List(ctx context.Context) ([]model.Loan, error) {
	return nil, nil
}

func (f *fakeLoanRepositoryForCustomerService) ListByCustomerID(ctx context.Context, customerID int) ([]model.Loan, error) {
	return f.listByCustomerIDFn(ctx, customerID)
}

func (f *fakeLoanRepositoryForCustomerService) Update(ctx context.Context, loan model.Loan) (model.Loan, error) {
	return model.Loan{}, nil
}

func TestCustomerService_CreateCustomer_Success(t *testing.T) {
	// This test proves customer creation:
	// - trims values
	// - validates required fields
	// - passes a clean customer object to the repository
	customerRepo := &fakeCustomerRepositoryForCustomerService{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			if customer.FullName != "Paul Mburu" {
				t.Fatalf("expected trimmed full name, got %q", customer.FullName)
			}
			if customer.Email != "paul@example.com" {
				t.Fatalf("expected trimmed email, got %q", customer.Email)
			}
			if customer.Phone != "0712345678" {
				t.Fatalf("expected trimmed phone, got %q", customer.Phone)
			}

			customer.ID = 1
			return customer, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{}, nil
		},
		listFn: func(ctx context.Context) ([]model.Customer, error) {
			return nil, nil
		},
	}

	loanRepo := &fakeLoanRepositoryForCustomerService{
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) {
			return nil, nil
		},
	}

	service := NewCustomerService(customerRepo, loanRepo)

	customer, err := service.CreateCustomer(context.Background(), " Paul Mburu ", " paul@example.com ", " 0712345678 ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if customer.ID != 1 {
		t.Fatalf("expected customer ID 1, got %d", customer.ID)
	}
}

func TestCustomerService_CreateCustomer_FailsWhenFullNameMissing(t *testing.T) {
	// This test proves the service rejects empty full_name before repository call.
	customerRepo := &fakeCustomerRepositoryForCustomerService{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			t.Fatalf("repository Create should not be called when full_name is missing")
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) { return model.Customer{}, nil },
		listFn:    func(ctx context.Context) ([]model.Customer, error) { return nil, nil },
	}

	loanRepo := &fakeLoanRepositoryForCustomerService{
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
	}

	service := NewCustomerService(customerRepo, loanRepo)

	_, err := service.CreateCustomer(context.Background(), "", "paul@example.com", "0712345678")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "full_name is required" {
		t.Fatalf("expected full_name is required, got %q", err.Error())
	}
}

func TestCustomerService_CreateCustomer_FailsWhenEmailMissing(t *testing.T) {
	// This test proves the service rejects empty email before repository call.
	customerRepo := &fakeCustomerRepositoryForCustomerService{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			t.Fatalf("repository Create should not be called when email is missing")
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) { return model.Customer{}, nil },
		listFn:    func(ctx context.Context) ([]model.Customer, error) { return nil, nil },
	}

	loanRepo := &fakeLoanRepositoryForCustomerService{
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
	}

	service := NewCustomerService(customerRepo, loanRepo)

	_, err := service.CreateCustomer(context.Background(), "Paul Mburu", "", "0712345678")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "email is required" {
		t.Fatalf("expected email is required, got %q", err.Error())
	}
}

func TestCustomerService_ListCustomers_Success(t *testing.T) {
	// This test proves list customers returns repository results unchanged.
	customerRepo := &fakeCustomerRepositoryForCustomerService{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{}, nil
		},
		listFn: func(ctx context.Context) ([]model.Customer, error) {
			return []model.Customer{
				{ID: 1, FullName: "A"},
				{ID: 2, FullName: "B"},
			}, nil
		},
	}

	loanRepo := &fakeLoanRepositoryForCustomerService{
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
	}

	service := NewCustomerService(customerRepo, loanRepo)

	customers, err := service.ListCustomers(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(customers) != 2 {
		t.Fatalf("expected 2 customers, got %d", len(customers))
	}
}

func TestCustomerService_GetCustomerByID_Success(t *testing.T) {
	// This test proves a valid customer lookup returns the record from repository.
	customerRepo := &fakeCustomerRepositoryForCustomerService{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{ID: id, FullName: "Paul Mburu"}, nil
		},
		listFn: func(ctx context.Context) ([]model.Customer, error) { return nil, nil },
	}

	loanRepo := &fakeLoanRepositoryForCustomerService{
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
	}

	service := NewCustomerService(customerRepo, loanRepo)

	customer, err := service.GetCustomerByID(context.Background(), 7)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if customer.ID != 7 {
		t.Fatalf("expected customer ID 7, got %d", customer.ID)
	}
}

func TestCustomerService_GetCustomerByID_FailsWhenInvalidID(t *testing.T) {
	// This test proves invalid IDs are rejected before repository access.
	customerRepo := &fakeCustomerRepositoryForCustomerService{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			t.Fatalf("repository GetByID should not be called for invalid customer ID")
			return model.Customer{}, nil
		},
		listFn: func(ctx context.Context) ([]model.Customer, error) { return nil, nil },
	}

	loanRepo := &fakeLoanRepositoryForCustomerService{
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
	}

	service := NewCustomerService(customerRepo, loanRepo)

	_, err := service.GetCustomerByID(context.Background(), 0)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "customer_id must be greater than zero" {
		t.Fatalf("expected invalid customer_id error, got %q", err.Error())
	}
}

func TestCustomerService_GetCustomerLoans_Success(t *testing.T) {
	// This test proves customer loans are returned only after customer existence is verified.
	customerRepo := &fakeCustomerRepositoryForCustomerService{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{ID: id, FullName: "Paul Mburu"}, nil
		},
		listFn: func(ctx context.Context) ([]model.Customer, error) { return nil, nil },
	}

	loanRepo := &fakeLoanRepositoryForCustomerService{
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) {
			return []model.Loan{
				{ID: 1, CustomerID: customerID},
				{ID: 2, CustomerID: customerID},
			}, nil
		},
	}

	service := NewCustomerService(customerRepo, loanRepo)

	loans, err := service.GetCustomerLoans(context.Background(), 7)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(loans) != 2 {
		t.Fatalf("expected 2 loans, got %d", len(loans))
	}
}

func TestCustomerService_GetCustomerLoans_FailsWhenCustomerMissing(t *testing.T) {
	// This test proves loan lookup is not attempted if the customer does not exist.
	customerRepo := &fakeCustomerRepositoryForCustomerService{
		createFn: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
			return model.Customer{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{}, apperror.ErrNotFound
		},
		listFn: func(ctx context.Context) ([]model.Customer, error) { return nil, nil },
	}

	loanRepo := &fakeLoanRepositoryForCustomerService{
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) {
			t.Fatalf("loan repository should not be called when customer is missing")
			return nil, nil
		},
	}

	service := NewCustomerService(customerRepo, loanRepo)

	_, err := service.GetCustomerLoans(context.Background(), 99)
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Fatalf("expected apperror.ErrNotFound, got %v", err)
	}
}
