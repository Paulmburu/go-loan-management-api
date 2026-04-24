package service

import (
	"context"
	"fmt"
	"go-loan-management-api/internal/model"
	"go-loan-management-api/internal/repository"
)

// LoanService contains business logic related to loans.
type LoanService struct {
	customerRepository repository.CustomerRepository
	loanRepository     repository.LoanRepository
}

// NewLoanService creates a new LoanService with access to the shared memory store.
func NewLoanService(customerRepository repository.CustomerRepository, loanRepository repository.LoanRepository) *LoanService {
	return &LoanService{
		customerRepository: customerRepository,
		loanRepository:     loanRepository,
	}
}

// CreateLoan validates loan input, checks that the customer exists,
// calculates the total and outstanding amounts, then stores the loan.
func (s *LoanService) CreateLoan(ctx context.Context, customerID int, principalAmount, interestRate float64) (model.Loan, error) {
	// Validate the customer ID.
	if customerID <= 0 {
		return model.Loan{}, fmt.Errorf("customer_id must be a positive integer")
	}

	// Validate the principal amount.
	if principalAmount <= 0 {
		return model.Loan{}, fmt.Errorf("principal_amount must be greater than zero")
	}

	// Validate the interest rate.
	if interestRate < 0 {
		return model.Loan{}, fmt.Errorf("interest_rate cannot be negative")
	}

	// Check whether the customer exists.
	_, err := s.customerRepository.GetByID(ctx, customerID)
	if err != nil {
		return model.Loan{}, fmt.Errorf("customer not found")
	}
	// Calculate simple interest amount.
	interestAmount := principalAmount * (interestRate / 100)

	// Calculate the total loan amount.
	totalAmount := principalAmount + interestAmount

	loan := model.Loan{
		CustomerID:        customerID,
		PrincipalAmount:   principalAmount,
		InterestRate:      interestRate,
		TotalAmount:       totalAmount,
		OutstandingAmount: totalAmount, // Initially, the outstanding amount is the total amount.
		Status:            "active",    // New loans start with an "active" status.
	}

	return s.loanRepository.Create(ctx, loan)
}

func (s *LoanService) ListLoans(ctx context.Context) ([]model.Loan, error) {
	return s.loanRepository.List(ctx)
}

// GetLoanBalance returns the outstanding balance for a given loan ID.
func (s *LoanService) GetLoanBalance(ctx context.Context, loanID int) (float64, error) {
	if loanID <= 0 {
		return 0, fmt.Errorf("loan_id must be greater than zero")
	}

	loan, err := s.loanRepository.GetByID(ctx, loanID)
	if err != nil {
		return 0, err
	}

	return loan.OutstandingAmount, nil
}

func (s *LoanService) GetLoanByID(ctx context.Context, loanID int) (model.Loan, error) {
	if loanID <= 0 {
		return model.Loan{}, fmt.Errorf("loan_id must be greater than zero")
	}

	return s.loanRepository.GetByID(ctx, loanID)
}
