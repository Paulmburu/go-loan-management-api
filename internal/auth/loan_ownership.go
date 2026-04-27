package auth

import "go-loan-management-api/internal/model"

// IsLoanOwnedByCustomer returns true if the given loan belongs to the authenticated customer.
func IsLoanOwnedByCustomer(claims *Claims, loan model.Loan) bool {
	if claims == nil || claims.CustomerID == nil {
		return false
	}

	return loan.CustomerID == *claims.CustomerID
}
