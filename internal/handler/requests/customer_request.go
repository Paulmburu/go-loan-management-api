package requests


// CreateCustomerRequest represents the JSON body for creating a customer.
type CreateCustomerRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
}