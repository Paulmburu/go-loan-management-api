package model

import "time"

// Loan represents a loan issued to a customer.

// Loan represents a formal financial agreement between the lender and a customer.
//
// Invariants:
//   - TotalAmount = PrincipalAmount + CalculatedInterest
//   - OutstandingAmount <= TotalAmount
//   - OutstandingAmount = TotalAmount - TotalPaymentsReceived
//   - PrincipalAmount > 0
//   - InterestRate >= 0
type Loan struct {
	ID                int       `json:"id"`
	CustomerID        int       `json:"customer_id"`
	PrincipalAmount   float64   `json:"principal_amount"`
	InterestRate      float64   `json:"interest_rate"`
	TotalAmount       float64   `json:"total_amount"`
	OutstandingAmount float64   `json:"outstanding_amount"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}
