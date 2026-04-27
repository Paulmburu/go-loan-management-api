package dto

// CreateCustomerRequest represents the request body for creating a customer.
type CreateCustomerRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}