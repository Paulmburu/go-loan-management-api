package model

import "time"

// Customer represents a borrower in the system.
type Customer struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}
