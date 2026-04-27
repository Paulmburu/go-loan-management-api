package service

import (
	"context"
	"fmt"
	"go-loan-management-api/internal/auth"
	"go-loan-management-api/internal/model"
	"go-loan-management-api/internal/repository"
	"strings"
)

// AuthService contains authentication-related business logic.
type AuthService struct {
	userRepository repository.UserRepository
	jwtSecret      string
	jwtExpiryHours int
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepository repository.UserRepository, jwtSecret string,
	jwtExpiryHours int) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

// RegisterInput holds data for user registration.
type RegisterInput struct {
	FullName   string
	Email      string
	Password   string
	Role       string
	CustomerID *int
}

// LoginInput holds data for user login.
type LoginInput struct {
	Email    string
	Password string
}

// AuthResult is returned after a successful auth operation.
type AuthResult struct {
	User  model.User `json:"user"`
	Token string     `json:"token"`
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (AuthResult, error) {
	fullName := strings.TrimSpace(input.FullName)
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)
	role := strings.TrimSpace(input.Role)

	if fullName == "" {
		return AuthResult{}, fmt.Errorf("full_name is required")
	}
	if email == "" {
		return AuthResult{}, fmt.Errorf("email is required")
	}
	if password == "" {
		return AuthResult{}, fmt.Errorf("password is required")
	}
	if len(password) < 6 {
		return AuthResult{}, fmt.Errorf("password must be at least 6 characters")
	}
	if role == "" {
		return AuthResult{}, fmt.Errorf("role is required")
	}

	if role != model.RoleAdmin && role != model.RoleLoanOfficer && role != model.RoleCustomer {
		return AuthResult{}, fmt.Errorf("invalid role")
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return AuthResult{}, fmt.Errorf("failed to hash password")
	}

	user := model.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CustomerID:   input.CustomerID,
	}

	savedUser, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return AuthResult{}, err
	}

	token, err := auth.GenerateToken(savedUser, s.jwtSecret, s.jwtExpiryHours)
	if err != nil {
		return AuthResult{}, fmt.Errorf("failed to generate token")
	}

	return AuthResult{
		User:  savedUser,
		Token: token,
	}, nil
}

// Login verifies a user's credentials.
func (s *AuthService) Login(ctx context.Context, input LoginInput) (AuthResult, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)

	if email == "" {
		return AuthResult{}, fmt.Errorf("email is required")
	}
	if password == "" {
		return AuthResult{}, fmt.Errorf("password is required")
	}

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return AuthResult{}, fmt.Errorf("invalid email or password")
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		return AuthResult{}, fmt.Errorf("invalid email or password")
	}

	token, err := auth.GenerateToken(user, s.jwtSecret, s.jwtExpiryHours)
	if err != nil {
		return AuthResult{}, fmt.Errorf("failed to generate token")
	}

	return AuthResult{
		User:  user,
		Token: token,
	}, nil

}

// GetMe fetches the currently authenticated user by user ID.
func (s *AuthService) GetMe(ctx context.Context, userID int) (model.User, error) {
	if userID <= 0 {
		return model.User{}, fmt.Errorf("invalid user id")
	}

	return s.userRepository.GetByID(ctx, userID)
}
