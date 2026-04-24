package service

import (
	"context"
	"fmt"
	"go-loan-management-api/internal/model"
	"go-loan-management-api/internal/repository"
	"strings"
)

// CustomerService contains customer-related business logic.
type CustomerService struct {
	customerRepository repository.CustomerRepository
	loanRepository     repository.LoanRepository
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(
	customerRepository repository.CustomerRepository,
	loanRepository repository.LoanRepository,
) *CustomerService {
	return &CustomerService{
		customerRepository: customerRepository,
		loanRepository:     loanRepository,
	}
}

// CreateCustomer validates input, creates a customer, and stores it.
func (s *CustomerService) CreateCustomer(ctx context.Context, fullName, email, phone string) (model.Customer, error) {
	// Remove accidental spaces around input values.
	// Use '=' to update an existing variable's value.
	fullName = strings.TrimSpace(fullName)
	email = strings.TrimSpace(email)
	phone = strings.TrimSpace(phone)

	// Basic validation for Phase 1.
	if fullName == "" {
		return model.Customer{}, fmt.Errorf("full_name is required")
	}
	if email == "" {
		return model.Customer{}, fmt.Errorf("email is required")
	}
	if phone == "" {
		return model.Customer{}, fmt.Errorf("phone is required")
	}

	customer := model.Customer{
		FullName: fullName,
		Email:    email,
		Phone:    phone,
	}

	return s.customerRepository.Create(ctx, customer)
}

// ListCustomers returns all customers currently stored in memory.
func (s *CustomerService) ListCustomers(ctx context.Context) ([]model.Customer, error) {
	return s.customerRepository.List(ctx)
}

// GetCustomerByID returns a customer by ID.
func (s *CustomerService) GetCustomerByID(ctx context.Context, customerID int) (model.Customer, error) {
	if customerID <= 0 {
		return model.Customer{}, fmt.Errorf("customer_id must be greater than zero")
	}

	return s.customerRepository.GetByID(ctx, customerID)
}

// GetCustomerLoans returns all loans for a given customer.
// It first verifies that the customer exists.
func (s *CustomerService) GetCustomerLoans(ctx context.Context, customerID int) ([]model.Loan, error) {
	if customerID <= 0 {
		return nil, fmt.Errorf("customer_id must be greater than zero")
	}

	_, err := s.customerRepository.GetByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	return s.loanRepository.ListByCustomerID(ctx, customerID)
}
