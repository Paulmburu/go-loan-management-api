package repository

import (
	"context"
	"testing"

	"go-loan-management-api/internal/model"
)

func TestPostgresRepaymentRepository_CreateGetAndListByLoanID(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	_, err := db.Exec(`TRUNCATE TABLE repayments, loans, customers RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("failed to truncate related tables: %v", err)
	}

	customerRepo := NewPostgresCustomerRepository(db)
	loanRepo := NewPostgresLoanRepository(db)
	repaymentRepo := NewPostgresRepaymentRepository(db)

	customer, err := customerRepo.Create(context.Background(), model.Customer{
		FullName: "Repayment Customer",
		Email:    "repayment.customer@example.com",
		Phone:    "0700000001",
	})
	if err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}

	loan, err := loanRepo.Create(context.Background(), model.Loan{
		CustomerID:        customer.ID,
		PrincipalAmount:   10000,
		InterestRate:      12,
		TotalAmount:       11200,
		OutstandingAmount: 11200,
		Status:            model.LoanStatusActive,
	})
	if err != nil {
		t.Fatalf("failed to create loan: %v", err)
	}

	repayment, err := repaymentRepo.Create(context.Background(), model.Repayment{
		LoanID: loan.ID,
		Amount: 2000,
	})
	if err != nil {
		t.Fatalf("expected create repayment to succeed, got %v", err)
	}

	if repayment.ID == 0 {
		t.Fatalf("expected repayment to have ID")
	}

	fetchedRepayment, err := repaymentRepo.GetByID(context.Background(), repayment.ID)
	if err != nil {
		t.Fatalf("expected get by id to succeed, got %v", err)
	}

	if fetchedRepayment.Amount != 2000 {
		t.Fatalf("expected repayment amount 2000, got %f", fetchedRepayment.Amount)
	}

	repayments, err := repaymentRepo.ListByLoanID(context.Background(), loan.ID)
	if err != nil {
		t.Fatalf("expected list by loan ID to succeed, got %v", err)
	}

	if len(repayments) != 1 {
		t.Fatalf("expected 1 repayment, got %d", len(repayments))
	}
}
