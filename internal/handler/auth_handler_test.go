package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-loan-management-api/internal/auth"
	"go-loan-management-api/internal/model"
	"go-loan-management-api/internal/service"
)

// fakeUserRepository is used to test AuthHandler through AuthService
// without depending on a real database.
type fakeUserRepository struct {
	createFn     func(ctx context.Context, user model.User) (model.User, error)
	getByEmailFn func(ctx context.Context, email string) (model.User, error)
	getByIDFn    func(ctx context.Context, id int) (model.User, error)
}

func (f *fakeUserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	return f.createFn(ctx, user)
}

func (f *fakeUserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	return f.getByEmailFn(ctx, email)
}

func (f *fakeUserRepository) GetByID(ctx context.Context, id int) (model.User, error) {
	return f.getByIDFn(ctx, id)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	repo := &fakeUserRepository{
		createFn: func(ctx context.Context, user model.User) (model.User, error) {
			user.ID = 1
			return user, nil
		},
		getByEmailFn: func(ctx context.Context, email string) (model.User, error) {
			return model.User{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.User, error) {
			return model.User{}, nil
		},
	}

	authService := service.NewAuthService(repo, "test-secret", 24)
	handler := NewAuthHandler(authService)

	body := `{
		"full_name":"Paul Mburu",
		"email":"paul@example.com",
		"password":"secret123",
		"role":"admin"
	}`

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	responseBody := rec.Body.String()
	if !strings.Contains(responseBody, "user registered successfully") {
		t.Fatalf("expected success message, got %q", responseBody)
	}
	if !strings.Contains(responseBody, "token") {
		t.Fatalf("expected token in response, got %q", responseBody)
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	hashedPassword, err := auth.HashPassword("secret123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	repo := &fakeUserRepository{
		createFn: func(ctx context.Context, user model.User) (model.User, error) {
			return model.User{}, nil
		},
		getByEmailFn: func(ctx context.Context, email string) (model.User, error) {
			return model.User{
				ID:           1,
				FullName:     "Paul Mburu",
				Email:        "paul@example.com",
				PasswordHash: hashedPassword,
				Role:         model.RoleAdmin,
			}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.User, error) {
			return model.User{}, nil
		},
	}

	authService := service.NewAuthService(repo, "test-secret", 24)
	handler := NewAuthHandler(authService)

	body := `{
		"email":"paul@example.com",
		"password":"secret123"
	}`

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	responseBody := rec.Body.String()
	if !strings.Contains(responseBody, "login successful") {
		t.Fatalf("expected login success message, got %q", responseBody)
	}
	if !strings.Contains(responseBody, "token") {
		t.Fatalf("expected token in response, got %q", responseBody)
	}
}

func TestAuthHandler_Login_FailsWithInvalidJSON(t *testing.T) {
	repo := &fakeUserRepository{
		createFn: func(ctx context.Context, user model.User) (model.User, error) {
			return model.User{}, nil
		},
		getByEmailFn: func(ctx context.Context, email string) (model.User, error) {
			return model.User{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.User, error) {
			return model.User{}, nil
		},
	}

	authService := service.NewAuthService(repo, "test-secret", 24)
	handler := NewAuthHandler(authService)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{invalid-json`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "invalid request body") {
		t.Fatalf("expected invalid request body error, got %q", rec.Body.String())
	}
}
