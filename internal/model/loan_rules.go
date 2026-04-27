package model

// CanAcceptRepayment returns true if the loan can accept a repayment.
func CanAcceptRepayment(loan Loan) bool {
	return loan.Status == LoanStatusActive && loan.OutstandingAmount > 0
}

