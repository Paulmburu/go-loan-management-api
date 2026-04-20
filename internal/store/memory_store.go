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
