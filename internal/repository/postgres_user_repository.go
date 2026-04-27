package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-loan-management-api/internal/apperror"
	"go-loan-management-api/internal/model"

	"github.com/lib/pq"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	query := `
		INSERT INTO users (full_name, email, password_hash, role, customer_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, full_name, email, role, customer_id, created_at
	`

	var savedUser model.User

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CustomerID,
	).Scan(
		&savedUser.ID,
		&savedUser.FullName,
		&savedUser.Email,
		&savedUser.Role,
		&savedUser.CustomerID,
		&savedUser.CreatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return model.User{}, fmt.Errorf("%w: email already exists", apperror.ErrConflict)
			}
		}

		return model.User{}, fmt.Errorf("%w: failed to create user", apperror.ErrInternal)
	}
	return savedUser, nil
}

// GetByEmail fetches a user by email.
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, customer_id, created_at
		FROM users
		WHERE email = $1
	`

	var user model.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CustomerID,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("%w: user not found", apperror.ErrNotFound)
		}
		return model.User{}, fmt.Errorf("%w: failed to get user by email", apperror.ErrInternal)
	}

	return user, nil
}

// GetByID fetches a user by ID.
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (model.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, role, customer_id, created_at
		FROM users
		WHERE id = $1
	`

	var user model.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CustomerID,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("%w: user not found", apperror.ErrNotFound)
		}
		return model.User{}, fmt.Errorf("%w: failed to get user by id", apperror.ErrInternal)
	}

	return user, nil
}
