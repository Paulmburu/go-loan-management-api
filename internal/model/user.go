package model

import "time"

// User represents an authenticated account in the system.
type User struct {
	ID           int       `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CustomerID   *int      `json:"customer_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
