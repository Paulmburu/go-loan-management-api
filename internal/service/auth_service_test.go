package service

import (
	"context"
	"testing"

	"go-loan-management-api/internal/auth"
	"go-loan-management-api/internal/model"
)

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

func TestAuthService_Register_Success(t *testing.T) {
	repo := &fakeUserRepository{
		createFn: func(ctx context.Context, user model.User) (model.User, error) {
			if user.FullName != "Paul Mburu" {
				t.Fatalf("expected full name to be saved")
			}
			if user.Email != "paul@example.com" {
				t.Fatalf("expected email to be normalized and saved")
			}
			if user.PasswordHash == "" {
				t.Fatalf("expected password hash to be set")
			}
			if user.PasswordHash == "secret123" {
				t.Fatalf("expected password to be hashed, not stored in plain text")
			}
			if user.Role != model.RoleAdmin {
				t.Fatalf("expected role to be admin")
			}

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

	service := NewAuthService(repo, "test-secret", 24)

	result, err := service.Register(context.Background(), RegisterInput{
		FullName: "Paul Mburu",
		Email:    "paul@example.com",
		Password: "secret123",
		Role:     model.RoleAdmin,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.User.ID != 1 {
		t.Fatalf("expected user ID to be 1, got %d", result.User.ID)
	}

	if result.User.Email != "paul@example.com" {
		t.Fatalf("expected returned email to match")
	}

	if result.Token == "" {
		t.Fatalf("expected token to be generated")
	}
}

func TestAuthService_Register_FailsWhenFullNameMissing(t *testing.T) {
	repo := &fakeUserRepository{
		createFn: func(ctx context.Context, user model.User) (model.User, error) {
			t.Fatalf("repository Create should not be called when validation fails")
			return model.User{}, nil
		},
		getByEmailFn: func(ctx context.Context, email string) (model.User, error) {
			return model.User{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.User, error) {
			return model.User{}, nil
		},
	}

	service := NewAuthService(repo, "test-secret", 24)

	_, err := service.Register(context.Background(), RegisterInput{
		FullName: "",
		Email:    "paul@example.com",
		Password: "secret123",
		Role:     model.RoleAdmin,
	})
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	if err.Error() != "full_name is required" {
		t.Fatalf("expected full_name is required, got %q", err.Error())
	}
}

func TestAuthService_Register_FailsWhenPasswordTooShort(t *testing.T) {
	repo := &fakeUserRepository{
		createFn: func(ctx context.Context, user model.User) (model.User, error) {
			t.Fatalf("repository Create should not be called when validation fails")
			return model.User{}, nil
		},
		getByEmailFn: func(ctx context.Context, email string) (model.User, error) {
			return model.User{}, nil
		},
		getByIDFn: func(ctx context.Context, id int) (model.User, error) {
			return model.User{}, nil
		},
	}

	service := NewAuthService(repo, "test-secret", 24)

	_, err := service.Register(context.Background(), RegisterInput{
		FullName: "Paul Mburu",
		Email:    "paul@example.com",
		Password: "123",
		Role:     model.RoleAdmin,
	})
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	if err.Error() != "password must be at least 6 characters" {
		t.Fatalf("expected short password error, got %q", err.Error())
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	hashedPassword, err := auth.HashPassword("secret123")
	if err != nil {
		t.Fatalf("failed to hash password for test: %v", err)
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

	service := NewAuthService(repo, "test-secret", 24)

	result, err := service.Login(context.Background(), LoginInput{
		Email:    "paul@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.User.ID != 1 {
		t.Fatalf("expected user ID to be 1, got %d", result.User.ID)
	}

	if result.Token == "" {
		t.Fatalf("expected token to be generated")
	}
}

func TestAuthService_Login_FailsWhenPasswordIncorrect(t *testing.T) {
	hashedPassword, err := auth.HashPassword("secret123")
	if err != nil {
		t.Fatalf("failed to hash password for test: %v", err)
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

	service := NewAuthService(repo, "test-secret", 24)

	_, err = service.Login(context.Background(), LoginInput{
		Email:    "paul@example.com",
		Password: "wrong-password",
	})
	if err == nil {
		t.Fatalf("expected login error, got nil")
	}

	if err.Error() != "invalid email or password" {
		t.Fatalf("expected invalid email or password, got %q", err.Error())
	}
}
