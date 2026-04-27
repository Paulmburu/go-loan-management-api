package dto

// RegisterRequest represents the request body for registration.
type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`

	// CustomerID uses a pointer (*int) to distinguish between:
	// 1. A value of 0 (an actual ID)
	// 2. A nil value (data was completely missing from the request)
	// omitempty ensures that if it's nil, it won't show up in the JSON output.
	CustomerID *int `json:"customer_id,omitempty"`
}

// LoginRequest represents the request body for login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
