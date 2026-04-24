package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-loan-management-api/internal/apperror"
	"go-loan-management-api/internal/model"
)

// PostgresRepaymentRepository implements RepaymentRepository using Postgres.
type PostgresRepaymentRepository struct {
	db *sql.DB
}

// NewPostgresRepaymentRepository creates a new Postgres-backed repayment repository.
func NewPostgresRepaymentRepository(db *sql.DB) *PostgresRepaymentRepository {
	return &PostgresRepaymentRepository{
		db: db,
	}
}

// Create inserts a new repayment into the database and returns the saved record.
func (r *PostgresRepaymentRepository) Create(ctx context.Context, repayment model.Repayment) (model.Repayment, error) {
	query := `
		INSERT INTO repayments (loan_id, amount)
		VALUES ($1, $2)
		RETURNING id, loan_id, amount, created_at
	`

	var savedRepayment model.Repayment

	err := r.db.QueryRowContext(
		ctx,
		query,
		repayment.LoanID,
		repayment.Amount,
	).Scan(
		&savedRepayment.ID,
		&savedRepayment.LoanID,
		&savedRepayment.Amount,
		&savedRepayment.CreatedAt,
	)
	if err != nil {
		return model.Repayment{}, fmt.Errorf("failed to create repayment: %w", err)
	}

	return savedRepayment, nil
}

// GetByID fetches a repayment by ID.
func (r *PostgresRepaymentRepository) GetByID(ctx context.Context, id int) (model.Repayment, error) {
	query := `
		SELECT id, loan_id, amount, created_at
		FROM repayments
		WHERE id = $1
	`

	var repayment model.Repayment

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&repayment.ID,
		&repayment.LoanID,
		&repayment.Amount,
		&repayment.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Repayment{}, fmt.Errorf("%w: repayment not found", apperror.ErrNotFound)
		}
		return model.Repayment{}, fmt.Errorf("%w: failed to get repayment by id", apperror.ErrInternal)
	}

	return repayment, nil
}

// ListByLoanID returns all repayments for a given loan ordered by ID.
func (r *PostgresRepaymentRepository) ListByLoanID(ctx context.Context, loanID int) ([]model.Repayment, error) {
	query := `
		SELECT id, loan_id, amount, created_at
		FROM repayments
		WHERE loan_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to list repayments by loan id: %w", err)
	}
	defer rows.Close()

	repayments := []model.Repayment{}

	for rows.Next() {
		var repayment model.Repayment

		err := rows.Scan(
			&repayment.ID,
			&repayment.LoanID,
			&repayment.Amount,
			&repayment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan repayment row: %w", err)
		}

		repayments = append(repayments, repayment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating repayment rows: %w", err)
	}

	return repayments, nil
}
