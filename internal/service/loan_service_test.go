package service

import (
	"context"
	"errors"
	"testing"

	"go-loan-management-api/internal/apperror"
	"go-loan-management-api/internal/model"
)

type fakeCustomerRepositoryForLoanService struct {
	getByIDFn func(ctx context.Context, id int) (model.Customer, error)
}

func (f *fakeCustomerRepositoryForLoanService) Create(ctx context.Context, customer model.Customer) (model.Customer, error) {
	return model.Customer{}, nil
}

func (f *fakeCustomerRepositoryForLoanService) GetByID(ctx context.Context, id int) (model.Customer, error) {
	return f.getByIDFn(ctx, id)
}

func (f *fakeCustomerRepositoryForLoanService) List(ctx context.Context) ([]model.Customer, error) {
	return nil, nil
}

type fakeLoanRepositoryForLoanService struct {
	createFn           func(ctx context.Context, loan model.Loan) (model.Loan, error)
	getByIDFn          func(ctx context.Context, id int) (model.Loan, error)
	listFn             func(ctx context.Context) ([]model.Loan, error)
	listByCustomerIDFn func(ctx context.Context, customerID int) ([]model.Loan, error)
	updateFn           func(ctx context.Context, loan model.Loan) (model.Loan, error)
}

func (f *fakeLoanRepositoryForLoanService) Create(ctx context.Context, loan model.Loan) (model.Loan, error) {
	return f.createFn(ctx, loan)
}

func (f *fakeLoanRepositoryForLoanService) GetByID(ctx context.Context, id int) (model.Loan, error) {
	return f.getByIDFn(ctx, id)
}

func (f *fakeLoanRepositoryForLoanService) List(ctx context.Context) ([]model.Loan, error) {
	return f.listFn(ctx)
}

func (f *fakeLoanRepositoryForLoanService) ListByCustomerID(ctx context.Context, customerID int) ([]model.Loan, error) {
	return f.listByCustomerIDFn(ctx, customerID)
}

func (f *fakeLoanRepositoryForLoanService) Update(ctx context.Context, loan model.Loan) (model.Loan, error) {
	return f.updateFn(ctx, loan)
}

func TestLoanService_CreateLoan_Success(t *testing.T) {
	// This test proves that:
	// - an existing customer is required
	// - total amount is calculated correctly
	// - outstanding amount starts equal to total amount
	customerRepo := &fakeCustomerRepositoryForLoanService{
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{ID: 1, FullName: "Paul Mburu"}, nil
		},
	}

	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) {
			// Validate the calculated values before "saving"
			if loan.TotalAmount != 11200 {
				t.Fatalf("expected total amount 11200, got %f", loan.TotalAmount)
			}
			if loan.OutstandingAmount != 11200 {
				t.Fatalf("expected outstanding amount 11200, got %f", loan.OutstandingAmount)
			}
			if loan.Status != model.LoanStatusActive {
				t.Fatalf("expected loan status active, got %s", loan.Status)
			}

			loan.ID = 1
			return loan, nil
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

	service := NewLoanService(customerRepo, loanRepo)

	loan, err := service.CreateLoan(context.Background(), 1, 10000, 12)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if loan.ID != 1 {
		t.Fatalf("expected loan ID 1, got %d", loan.ID)
	}
}

func TestLoanService_CreateLoan_FailsWhenCustomerMissing(t *testing.T) {
	// This test proves that a loan cannot be created if the customer does not exist.
	customerRepo := &fakeCustomerRepositoryForLoanService{
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{}, errors.New("customer not found")
		},
	}

	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) {
			t.Fatalf("loan repository Create should not be called when customer is missing")
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

	service := NewLoanService(customerRepo, loanRepo)

	_, err := service.CreateLoan(context.Background(), 99, 10000, 12)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "customer not found" {
		t.Fatalf("expected customer not found, got %q", err.Error())
	}
}

func TestLoanService_CreateLoan_FailsWhenPrincipalInvalid(t *testing.T) {
	// This test proves principal_amount must be greater than zero.
	customerRepo := &fakeCustomerRepositoryForLoanService{
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{ID: 1}, nil
		},
	}

	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) {
			t.Fatalf("loan repository Create should not be called when principal is invalid")
			return model.Loan{}, nil
		},
		getByIDFn:          func(ctx context.Context, id int) (model.Loan, error) { return model.Loan{}, nil },
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	service := NewLoanService(customerRepo, loanRepo)

	_, err := service.CreateLoan(context.Background(), 1, 0, 12)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "principal_amount must be greater than zero" {
		t.Fatalf("expected principal validation error, got %q", err.Error())
	}
}

func TestLoanService_CreateLoan_FailsWhenInterestRateNegative(t *testing.T) {
	// This test proves interest_rate cannot be negative.
	customerRepo := &fakeCustomerRepositoryForLoanService{
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{ID: 1}, nil
		},
	}

	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) {
			t.Fatalf("loan repository Create should not be called when interest is invalid")
			return model.Loan{}, nil
		},
		getByIDFn:          func(ctx context.Context, id int) (model.Loan, error) { return model.Loan{}, nil },
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	service := NewLoanService(customerRepo, loanRepo)

	_, err := service.CreateLoan(context.Background(), 1, 10000, -1)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "interest_rate cannot be negative" {
		t.Fatalf("expected interest validation error, got %q", err.Error())
	}
}

func TestLoanService_GetLoanByID_Success(t *testing.T) {
	// This test proves a valid loan lookup returns the expected record.
	customerRepo := &fakeCustomerRepositoryForLoanService{
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{}, nil
		},
	}

	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			return model.Loan{ID: 5, CustomerID: 1, OutstandingAmount: 9200}, nil
		},
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	service := NewLoanService(customerRepo, loanRepo)

	loan, err := service.GetLoanByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if loan.ID != 5 {
		t.Fatalf("expected loan ID 5, got %d", loan.ID)
	}
}

func TestLoanService_ListLoans_Success(t *testing.T) {
	// This test proves list loans returns whatever the repository provides.
	customerRepo := &fakeCustomerRepositoryForLoanService{
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{}, nil
		},
	}

	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn:  func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) { return model.Loan{}, nil },
		listFn: func(ctx context.Context) ([]model.Loan, error) {
			return []model.Loan{
				{ID: 1, CustomerID: 1},
				{ID: 2, CustomerID: 2},
			}, nil
		},
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	service := NewLoanService(customerRepo, loanRepo)

	loans, err := service.ListLoans(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(loans) != 2 {
		t.Fatalf("expected 2 loans, got %d", len(loans))
	}
}

func TestLoanService_GetLoanBalance_Success(t *testing.T) {
	// This test proves loan balance is read from the loan record.
	customerRepo := &fakeCustomerRepositoryForLoanService{
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{}, nil
		},
	}

	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			return model.Loan{ID: 1, OutstandingAmount: 9200}, nil
		},
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	service := NewLoanService(customerRepo, loanRepo)

	balance, err := service.GetLoanBalance(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if balance != 9200 {
		t.Fatalf("expected balance 9200, got %f", balance)
	}
}

func TestLoanService_GetLoanByID_FailsWhenInvalidID(t *testing.T) {
	// This test proves service validation happens before repository lookup.
	customerRepo := &fakeCustomerRepositoryForLoanService{
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{}, nil
		},
	}

	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			t.Fatalf("repository GetByID should not be called for invalid ID")
			return model.Loan{}, nil
		},
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	service := NewLoanService(customerRepo, loanRepo)

	_, err := service.GetLoanByID(context.Background(), 0)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "loan_id must be greater than zero" {
		t.Fatalf("expected invalid ID error, got %q", err.Error())
	}
}

func TestLoanService_GetLoanByID_PropagatesRepositoryNotFound(t *testing.T) {
	// This test proves repository-level not found errors bubble up.
	customerRepo := &fakeCustomerRepositoryForLoanService{
		getByIDFn: func(ctx context.Context, id int) (model.Customer, error) {
			return model.Customer{}, nil
		},
	}

	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			return model.Loan{}, apperror.ErrNotFound
		},
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	service := NewLoanService(customerRepo, loanRepo)

	_, err := service.GetLoanByID(context.Background(), 123)
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Fatalf("expected apperror.ErrNotFound, got %v", err)
	}
}
