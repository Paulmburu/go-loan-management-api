package service

import (
	"context"
	"errors"
	"testing"

	"go-loan-management-api/internal/apperror"
	"go-loan-management-api/internal/model"
)

type fakeRepaymentRepositoryForService struct {
	createFn       func(ctx context.Context, repayment model.Repayment) (model.Repayment, error)
	getByIDFn      func(ctx context.Context, id int) (model.Repayment, error)
	listByLoanIDFn func(ctx context.Context, loanID int) ([]model.Repayment, error)
}

func (f *fakeRepaymentRepositoryForService) Create(ctx context.Context, repayment model.Repayment) (model.Repayment, error) {
	return f.createFn(ctx, repayment)
}

func (f *fakeRepaymentRepositoryForService) GetByID(ctx context.Context, id int) (model.Repayment, error) {
	return f.getByIDFn(ctx, id)
}

func (f *fakeRepaymentRepositoryForService) ListByLoanID(ctx context.Context, loanID int) ([]model.Repayment, error) {
	return f.listByLoanIDFn(ctx, loanID)
}

func TestRepaymentService_ListRepaymentsByLoanID_Success(t *testing.T) {
	// This test proves the service:
	// - first verifies the loan exists
	// - then returns the repayments for that loan
	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			return model.Loan{ID: id, CustomerID: 1, Status: model.LoanStatusActive}, nil
		},
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	repaymentRepo := &fakeRepaymentRepositoryForService{
		createFn: func(ctx context.Context, repayment model.Repayment) (model.Repayment, error) {
			return model.Repayment{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Repayment, error) { return model.Repayment{}, nil },
		listByLoanIDFn: func(ctx context.Context, loanID int) ([]model.Repayment, error) {
			return []model.Repayment{
				{ID: 1, LoanID: loanID, Amount: 2000},
				{ID: 2, LoanID: loanID, Amount: 3000},
			}, nil
		},
	}

	service := NewRepaymentService(nil, loanRepo, repaymentRepo)

	repayments, err := service.ListRepaymentsByLoanID(context.Background(), 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(repayments) != 2 {
		t.Fatalf("expected 2 repayments, got %d", len(repayments))
	}
}

func TestRepaymentService_ListRepaymentsByLoanID_FailsWhenLoanIDInvalid(t *testing.T) {
	// This test proves the service validates the loan ID before any repository call.
	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			t.Fatalf("loan repository should not be called for invalid loan ID")
			return model.Loan{}, nil
		},
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	repaymentRepo := &fakeRepaymentRepositoryForService{
		createFn: func(ctx context.Context, repayment model.Repayment) (model.Repayment, error) {
			return model.Repayment{}, nil
		},
		getByIDFn:      func(ctx context.Context, id int) (model.Repayment, error) { return model.Repayment{}, nil },
		listByLoanIDFn: func(ctx context.Context, loanID int) ([]model.Repayment, error) { return nil, nil },
	}

	service := NewRepaymentService(nil, loanRepo, repaymentRepo)

	_, err := service.ListRepaymentsByLoanID(context.Background(), 0)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "loan_id must be greater than zero" {
		t.Fatalf("expected invalid loan_id error, got %q", err.Error())
	}
}

func TestRepaymentService_GetRepaymentByID_Success(t *testing.T) {
	// This test proves GetRepaymentByID returns a valid repayment when repository lookup succeeds.
	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn:          func(ctx context.Context, id int) (model.Loan, error) { return model.Loan{}, nil },
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	repaymentRepo := &fakeRepaymentRepositoryForService{
		createFn: func(ctx context.Context, repayment model.Repayment) (model.Repayment, error) {
			return model.Repayment{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Repayment, error) {
			return model.Repayment{ID: id, LoanID: 10, Amount: 2000}, nil
		},
		listByLoanIDFn: func(ctx context.Context, loanID int) ([]model.Repayment, error) { return nil, nil },
	}

	service := NewRepaymentService(nil, loanRepo, repaymentRepo)

	repayment, err := service.GetRepaymentByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repayment.ID != 5 {
		t.Fatalf("expected repayment ID 5, got %d", repayment.ID)
	}
}

func TestRepaymentService_GetRepaymentByID_FailsWhenInvalidID(t *testing.T) {
	// This test proves service validation happens before repository lookup.
	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn:          func(ctx context.Context, id int) (model.Loan, error) { return model.Loan{}, nil },
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	repaymentRepo := &fakeRepaymentRepositoryForService{
		createFn: func(ctx context.Context, repayment model.Repayment) (model.Repayment, error) {
			return model.Repayment{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Repayment, error) {
			t.Fatalf("repayment repository should not be called for invalid repayment ID")
			return model.Repayment{}, nil
		},
		listByLoanIDFn: func(ctx context.Context, loanID int) ([]model.Repayment, error) { return nil, nil },
	}

	service := NewRepaymentService(nil, loanRepo, repaymentRepo)

	_, err := service.GetRepaymentByID(context.Background(), 0)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err.Error() != "repayment_id must be greater than zero" {
		t.Fatalf("expected invalid repayment_id error, got %q", err.Error())
	}
}

func TestRepaymentService_GetRepaymentWithLoan_Success(t *testing.T) {
	// This test proves the service can:
	// - load a repayment
	// - then load the related loan
	// - and return both for ownership checks
	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn: func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn: func(ctx context.Context, id int) (model.Loan, error) {
			return model.Loan{ID: id, CustomerID: 7, Status: model.LoanStatusActive}, nil
		},
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	repaymentRepo := &fakeRepaymentRepositoryForService{
		createFn: func(ctx context.Context, repayment model.Repayment) (model.Repayment, error) {
			return model.Repayment{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Repayment, error) {
			return model.Repayment{ID: id, LoanID: 10, Amount: 1500}, nil
		},
		listByLoanIDFn: func(ctx context.Context, loanID int) ([]model.Repayment, error) { return nil, nil },
	}

	service := NewRepaymentService(nil, loanRepo, repaymentRepo)

	repayment, loan, err := service.GetRepaymentWithLoan(context.Background(), 5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repayment.ID != 5 {
		t.Fatalf("expected repayment ID 5, got %d", repayment.ID)
	}

	if loan.ID != 10 {
		t.Fatalf("expected related loan ID 10, got %d", loan.ID)
	}
}

func TestRepaymentService_GetRepaymentWithLoan_PropagatesRepaymentNotFound(t *testing.T) {
	// This test proves repayment repository not-found errors bubble up.
	loanRepo := &fakeLoanRepositoryForLoanService{
		createFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
		getByIDFn:          func(ctx context.Context, id int) (model.Loan, error) { return model.Loan{}, nil },
		listFn:             func(ctx context.Context) ([]model.Loan, error) { return nil, nil },
		listByCustomerIDFn: func(ctx context.Context, customerID int) ([]model.Loan, error) { return nil, nil },
		updateFn:           func(ctx context.Context, loan model.Loan) (model.Loan, error) { return model.Loan{}, nil },
	}

	repaymentRepo := &fakeRepaymentRepositoryForService{
		createFn: func(ctx context.Context, repayment model.Repayment) (model.Repayment, error) {
			return model.Repayment{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.Repayment, error) {
			return model.Repayment{}, apperror.ErrNotFound
		},
		listByLoanIDFn: func(ctx context.Context, loanID int) ([]model.Repayment, error) { return nil, nil },
	}

	service := NewRepaymentService(nil, loanRepo, repaymentRepo)

	_, _, err := service.GetRepaymentWithLoan(context.Background(), 123)
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Fatalf("expected apperror.ErrNotFound, got %v", err)
	}
}
