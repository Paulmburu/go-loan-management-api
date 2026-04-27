package repository

import (
	"context"
	"go-loan-management-api/internal/model"
)

// UserRepository defines authenticated user data access behavior.
type UserRepository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByID(ctx context.Context, id int) (model.User, error)
}
