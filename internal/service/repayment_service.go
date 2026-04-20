package service

import (
	"fmt"
	"go-loan-management-api/internal/model"
	"go-loan-management-api/internal/store"
)

// RepaymentService contains business logic related to repayments.
type RepaymentService struct {
	store *store.MemoryStore
}

// NewRepaymentService creates a new RepaymentService.
func NewRepaymentService(store *store.MemoryStore) *RepaymentService {
	return &RepaymentService{store: store}
}

// AddRepayment validates and records a repayment for a loan.
// It also updates the loan's outstanding balance and status.
func (s *RepaymentService) AddRepayment(loanID int, amount float64) (*model.Repayment, error) {
	// Validate the loan ID.
	if loanID <= 0 {
		return &model.Repayment{}, fmt.Errorf("loan_id must be a positive integer")
	}

	// Validate the repayment amount.
	if amount <= 0 {
		return &model.Repayment{}, fmt.Errorf("amount must be greater than zero")
	}

	// Find the loan in memory.
	loanIndex := -1 
	for i, loan := range s.store.Loans {
		if loan.ID == loanID {
			loanIndex = i
			break
		}
	}

	if loanIndex == -1 {
		return &model.Repayment{}, fmt.Errorf("loan with ID %d does not exist", loanID)
	}

	// Get the loan from the store.
	loan := s.store.Loans[loanIndex]

	// Only active loans can accept repayments.
	if loan.Status != "active" {
		return &model.Repayment{}, fmt.Errorf("only active loans can accept repayments")
	}

	// Prevent overpayment.
	if amount > loan.OutstandingAmount {
		return &model.Repayment{}, fmt.Errorf("repayment amount cannot exceed outstanding amount of %.2f", loan.OutstandingAmount)
	}

	// Create the repayment.
	repayment := model.Repayment{
		ID:     len(s.store.Repayments) + 1, // Auto-increment ID based on current count.
		LoanID: loanID,
		Amount: amount,
	}

	// Save the repayment.
	s.store.Repayments = append(s.store.Repayments, repayment)

	// Update the outstanding balance.
	loan.OutstandingAmount -= amount

	// If the loan is fully paid, update its status.
	if loan.OutstandingAmount == 0 {
		loan.Status = "paid"
	}

	// Write the updated loan back into the in-memory store.
	s.store.Loans[loanIndex] = loan

	return &repayment, nil
}