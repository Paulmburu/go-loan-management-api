package service

import (
	"fmt"
	"go-loan-management-api/internal/model"
	"go-loan-management-api/internal/store"
)

// LoanService contains business logic related to loans.
type LoanService struct {
	store *store.MemoryStore
}

// NewLoanService creates a new LoanService with access to the shared memory store.
func NewLoanService(store *store.MemoryStore) *LoanService {
	return &LoanService{store: store}
}

// CreateLoan validates loan input, checks that the customer exists,
// calculates the total and outstanding amounts, then stores the loan.
func (s *LoanService) CreateLoan(customerID int, principalAmount, interestRate float64) (*model.Loan, error) {
	// Validate the customer ID.
	if customerID <= 0 {
		return &model.Loan{}, fmt.Errorf("customer_id must be a positive integer")
	}

	// Validate the principal amount.
	if principalAmount <= 0 {
		return &model.Loan{}, fmt.Errorf("principal_amount must be greater than zero")
	}

	// Validate the interest rate.
	if interestRate < 0 {
		return &model.Loan{}, fmt.Errorf("interest_rate cannot be negative")
	}

	customerExists := false

	for _, customer := range s.store.Customers {
		if customer.ID == customerID {
			customerExists = true
			break
		}
	}

	if !customerExists {
		return &model.Loan{}, fmt.Errorf("customer with ID %d does not exist", customerID)
	}

	// Calculate simple interest amount.
	interestAmount := principalAmount * (interestRate / 100)

	// Calculate the total loan amount.
	totalAmount := principalAmount + interestAmount

	loan := model.Loan{
		ID:                len(s.store.Loans) + 1, // Auto-increment ID based on current count.
		CustomerID:        customerID,
		PrincipalAmount:   principalAmount,
		InterestRate:      interestRate,
		TotalAmount:       totalAmount,
		OutstandingAmount: totalAmount, // Initially, the outstanding amount is the total amount.
		Status:            "active",    // New loans start with an "active" status.
	}

	// Save the loan in the in-memory store.
	s.store.Loans = append(s.store.Loans, loan)

	return &loan, nil
}

func (s *LoanService) ListLoans() []model.Loan {
	return s.store.Loans
}

// GetLoanBalance returns the outstanding balance for a given loan ID.
func (s *LoanService) GetLoanBalance(loanID int) (float64, error) {

	for _, loan := range s.store.Loans {
		if loan.ID == loanID {
			return loan.OutstandingAmount, nil
		}
	}

	return 0, fmt.Errorf("loan with ID %d does not exist", loanID)
}
