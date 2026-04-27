package repository

import (
	"context"
	"testing"

	"go-loan-management-api/internal/model"
)

func TestPostgresCustomerRepository_CreateGetAndList(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	_, err := db.Exec(`TRUNCATE TABLE repayments, loans, customers RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("failed to truncate related tables: %v", err)
	}

	repo := NewPostgresCustomerRepository(db)

	createdCustomer, err := repo.Create(context.Background(), model.Customer{
		FullName: "Paul Mburu",
		Email:    "paul.customer@example.com",
		Phone:    "0712345678",
	})
	if err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}

	if createdCustomer.ID == 0 {
		t.Fatalf("expected customer to have ID")
	}

	fetchedCustomer, err := repo.GetByID(context.Background(), createdCustomer.ID)
	if err != nil {
		t.Fatalf("expected get by id to succeed, got %v", err)
	}

	if fetchedCustomer.Email != "paul.customer@example.com" {
		t.Fatalf("expected fetched email to match, got %s", fetchedCustomer.Email)
	}

	customers, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("expected list to succeed, got %v", err)
	}

	if len(customers) != 1 {
		t.Fatalf("expected 1 customer, got %d", len(customers))
	}
}
