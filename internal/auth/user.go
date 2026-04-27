package auth

import "go-loan-management-api/internal/model"

// AuthenticatedUser represents the identity of the currently authenticated user.
type AuthenticatedUser struct {
	UserID     int
	Email      string
	Role       string
	CustomerID *int
}

// NewAuthenticatedUserFromClaims converts JWT claims into an authenticated user.
func NewAuthenticatedUserFromClaims(claims *Claims) *AuthenticatedUser {
	if claims == nil {
		return nil
	}

	return &AuthenticatedUser{
		UserID:     claims.UserID,
		Email:      claims.Email,
		Role:       claims.Role,
		CustomerID: claims.CustomerID,
	}
}

// IsAdmin returns true if the authenticated user is an admin.
func (u *AuthenticatedUser) IsAdmin() bool {
	return u != nil && u.Role == model.RoleAdmin
}

// IsLoanOfficer returns true if the authenticated user is a loan officer.
func (u *AuthenticatedUser) IsLoanOfficer() bool {
	return u != nil && u.Role == model.RoleLoanOfficer
}

// IsCustomer returns true if the authenticated user is a customer.
func (u *AuthenticatedUser) IsCustomer() bool {
	return u != nil && u.Role == model.RoleCustomer
}

// HasCustomerID returns true if the user is linked to a customer record.
func (u *AuthenticatedUser) HasCustomerID() bool {
	return u != nil && u.CustomerID != nil
}

// CustomerIDValue returns the linked customer ID and whether it exists.
func (u *AuthenticatedUser) CustomerIDValue() (int, bool) {
	if u == nil || u.CustomerID == nil {
		return 0, false
	}

	return *u.CustomerID, true
}
