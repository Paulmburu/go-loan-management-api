package model

// Repayment represents a payment made against a loan.
type Repayment struct {
	ID     int     `json:"id"`
	LoanID int     `json:"loan_id"`
	Amount float64 `json:"amount"`
}
