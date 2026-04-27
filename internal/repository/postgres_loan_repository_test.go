package repository

import (
	"context"
	"testing"

	"go-loan-management-api/internal/model"
)

func TestPostgresLoanRepository_CreateGetListAndUpdate(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	_, err := db.Exec(`TRUNCATE TABLE repayments, loans, customers RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("failed to truncate related tables: %v", err)
	}

	customerRepo := NewPostgresCustomerRepository(db)
	loanRepo := NewPostgresLoanRepository(db)

	customer, err := customerRepo.Create(context.Background(), model.Customer{
		FullName: "Loan Customer",
		Email:    "loan.customer@example.com",
		Phone:    "0700000000",
	})
	if err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}

	createdLoan, err := loanRepo.Create(context.Background(), model.Loan{
		CustomerID:        customer.ID,
		PrincipalAmount:   10000,
		InterestRate:      12,
		TotalAmount:       11200,
		OutstandingAmount: 11200,
		Status:            model.LoanStatusActive,
	})
	if err != nil {
		t.Fatalf("expected create loan to succeed, got %v", err)
	}

	fetchedLoan, err := loanRepo.GetByID(context.Background(), createdLoan.ID)
	if err != nil {
		t.Fatalf("expected get by id to succeed, got %v", err)
	}

	if fetchedLoan.CustomerID != customer.ID {
		t.Fatalf("expected customer ID %d, got %d", customer.ID, fetchedLoan.CustomerID)
	}

	loans, err := loanRepo.List(context.Background())
	if err != nil {
		t.Fatalf("expected list to succeed, got %v", err)
	}

	if len(loans) != 1 {
		t.Fatalf("expected 1 loan, got %d", len(loans))
	}

	customerLoans, err := loanRepo.ListByCustomerID(context.Background(), customer.ID)
	if err != nil {
		t.Fatalf("expected list by customer ID to succeed, got %v", err)
	}

	if len(customerLoans) != 1 {
		t.Fatalf("expected 1 customer loan, got %d", len(customerLoans))
	}

	createdLoan.OutstandingAmount = 0
	createdLoan.Status = model.LoanStatusPaid

	updatedLoan, err := loanRepo.Update(context.Background(), createdLoan)
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}

	if updatedLoan.Status != model.LoanStatusPaid {
		t.Fatalf("expected updated status to be paid, got %s", updatedLoan.Status)
	}
}
