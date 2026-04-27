package dto

// CreateLoanRequest represents the request body for creating a loan.
type CreateLoanRequest struct {
	CustomerID      int     `json:"customer_id"`
	PrincipalAmount float64 `json:"principal_amount"`
	InterestRate    float64 `json:"interest_rate"`
}
