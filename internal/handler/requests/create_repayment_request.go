package requests

// CreateRepaymentRequest represents the JSON body for creating a repayment.
type CreateRepaymentRequest struct {
	LoanID int     `json:"loan_id"`
	Amount float64 `json:"amount"`
}