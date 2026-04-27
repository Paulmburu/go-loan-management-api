package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"go-loan-management-api/internal/model"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set; skipping integration test")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	return db
}

func TestPostgresUserRepository_CreateAndGetByEmail(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	// Clean table so the test starts with predictable state.
	_, err := db.Exec(`TRUNCATE TABLE users RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("failed to truncate users table: %v", err)
	}

	repo := NewPostgresUserRepository(db)

	createdUser, err := repo.Create(context.Background(), model.User{
		FullName:     "Admin User",
		Email:        "admin.integration@example.com",
		PasswordHash: "hashed-password",
		Role:         model.RoleAdmin,
	})
	if err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}

	if createdUser.ID == 0 {
		t.Fatalf("expected created user to have ID")
	}

	fetchedUser, err := repo.GetByEmail(context.Background(), "admin.integration@example.com")
	if err != nil {
		t.Fatalf("expected get by email to succeed, got %v", err)
	}

	if fetchedUser.Email != "admin.integration@example.com" {
		t.Fatalf("expected fetched email to match, got %s", fetchedUser.Email)
	}

	if fetchedUser.Role != model.RoleAdmin {
		t.Fatalf("expected fetched role to be admin, got %s", fetchedUser.Role)
	}
}

func TestPostgresUserRepository_GetByID(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	_, err := db.Exec(`TRUNCATE TABLE users RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("failed to truncate users table: %v", err)
	}

	repo := NewPostgresUserRepository(db)

	createdUser, err := repo.Create(context.Background(), model.User{
		FullName:     "Loan Officer",
		Email:        "officer.integration@example.com",
		PasswordHash: "hashed-password",
		Role:         model.RoleLoanOfficer,
	})
	if err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}

	fetchedUser, err := repo.GetByID(context.Background(), createdUser.ID)
	if err != nil {
		t.Fatalf("expected get by id to succeed, got %v", err)
	}

	if fetchedUser.ID != createdUser.ID {
		t.Fatalf("expected ID %d, got %d", createdUser.ID, fetchedUser.ID)
	}
}
