package service

import (
	"fmt"
	"go-loan-management-api/internal/model"
	"go-loan-management-api/internal/store"
	"strings"
)

// CustomerService contains customer-related business logic.
type CustomerService struct {
	store *store.MemoryStore
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(store *store.MemoryStore) *CustomerService {
	return &CustomerService{store: store}
}

// CreateCustomer validates input, creates a customer, and stores it.
func (s *CustomerService) CreateCustomer(fullName, email, phone string) (*model.Customer, error) {
	// Remove accidental spaces around input values.
	// Use '=' to update an existing variable's value.
	fullName = strings.TrimSpace(fullName)
	email = strings.TrimSpace(email)
	phone = strings.TrimSpace(phone)

	// Basic validation for Phase 1.
	if fullName == "" {
		return &model.Customer{}, fmt.Errorf("full_name is required")
	}
	if email == "" {
		return &model.Customer{}, fmt.Errorf("email is required")
	}
	if phone == "" {
		return &model.Customer{}, fmt.Errorf("phone is required")
	}

	// Create the customer using the next in-memory ID.
	// Use ':=' to declare a new variable and infer its type.
	customer := model.Customer{
		ID:       len(s.store.Customers) + 1, // Returns the number of elements currently in the slice.
		FullName: fullName,
		Email:    email,
		Phone:    phone,
	}

	// Store the customer in memory.
	s.store.Customers = append(s.store.Customers, customer)

	return &customer, nil
}
