package store

import "go-loan-management-api/internal/model"

// MemoryStore holds application data in memory.
// This is useful for Phase 1 because it avoids database complexity.
type MemoryStore struct {
	Customers  []model.Customer
	Loans      []model.Loan
	Repayments []model.Repayment
}

// NewMemoryStore creates and returns an empty in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Customers:  []model.Customer{},
		Loans:      []model.Loan{},
		Repayments: []model.Repayment{},
	}
}

// FindCustomerByID looks up a customer by ID.
// It returns the customer and true if found, otherwise a zero-value customer and false.
func (s *MemoryStore) FindCustomerByID(customerID int) (model.Customer, bool) {
	for _, customer := range s.Customers {
		if customer.ID == customerID {
			return customer, true
		}
	}

	return model.Customer{}, false
}

// FindLoanByID looks up a loan by ID.
// It returns the loan, its index in the slice, and true if found.
func (s *MemoryStore) FindLoanByID(loanID int) (model.Loan, int, bool) {
	for i, loan := range s.Loans {
		if loan.ID == loanID {
			return loan, i, true
		}
	}

	return model.Loan{}, -1, false
}

// UpdateLoan replaces the loan at the given index.
// It returns false if the index is invalid.
func (s *MemoryStore) UpdateLoan(index int, loan model.Loan) bool {
	if index < 0 || index >= len(s.Loans) {
		return false
	}

	s.Loans[index] = loan
	return true
}
