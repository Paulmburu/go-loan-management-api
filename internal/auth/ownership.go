package auth

// IsSameCustomer returns true if the authenticated claims belong to the target customer.
func IsSameCustomer(claims *Claims, customerID int) bool {
	if claims == nil || claims.CustomerID == nil {
		return false
	}

	return *claims.CustomerID == customerID
}
