package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-loan-management-api/internal/apperror"
	"go-loan-management-api/internal/model"
)

// PostgresLoanRepository implements LoanRepository using Postgres.
type PostgresLoanRepository struct {
	db *sql.DB
}

// NewPostgresLoanRepository creates a new Postgres-backed loan repository.
func NewPostgresLoanRepository(db *sql.DB) *PostgresLoanRepository {
	return &PostgresLoanRepository{
		db: db,
	}
}

func (r *PostgresLoanRepository) Create(ctx context.Context, loan model.Loan) (model.Loan, error) {
	query := `
		INSERT INTO loans (
		customer_id,
		principal_amount,
		interest_rate,
		total_amount,
		outstanding_amount,
		status
		) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, customer_id, principal_amount, interest_rate, total_amount, outstanding_amount, status, created_at
	`

	var savedLoan model.Loan
	err := r.db.QueryRowContext(
		ctx,
		query,
		loan.CustomerID,
		loan.PrincipalAmount,
		loan.InterestRate,
		loan.TotalAmount,
		loan.OutstandingAmount,
		loan.Status,
	).Scan(
		&savedLoan.ID,
		&savedLoan.CustomerID,
		&savedLoan.PrincipalAmount,
		&savedLoan.InterestRate,
		&savedLoan.TotalAmount,
		&savedLoan.OutstandingAmount,
		&savedLoan.Status,
		&savedLoan.CreatedAt,
	)

	if err != nil {
		return model.Loan{}, fmt.Errorf("failed to create loan: %w", err)
	}

	return savedLoan, nil
}

// GetByID fetches a loan by ID.
func (r *PostgresLoanRepository) GetByID(ctx context.Context, id int) (model.Loan, error) {
	query := `
		SELECT id, customer_id, principal_amount, interest_rate, total_amount, outstanding_amount, status, created_at
		FROM loans
		WHERE id = $1
	`

	var loan model.Loan

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&loan.ID,
		&loan.CustomerID,
		&loan.PrincipalAmount,
		&loan.InterestRate,
		&loan.TotalAmount,
		&loan.OutstandingAmount,
		&loan.Status,
		&loan.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Loan{}, fmt.Errorf("%w: loan not found", apperror.ErrNotFound)
		}
		return model.Loan{}, fmt.Errorf("%w: failed to get loan by id", apperror.ErrInternal)
	}

	return loan, nil
}

// List returns all loans ordered by ID.
func (r *PostgresLoanRepository) List(ctx context.Context) ([]model.Loan, error) {
	query := `
		SELECT id, customer_id, principal_amount, interest_rate, total_amount, outstanding_amount, status, created_at
		FROM loans
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list loans: %w", err)
	}
	defer rows.Close()

	loans := []model.Loan{}

	for rows.Next() {
		var loan model.Loan

		err := rows.Scan(
			&loan.ID,
			&loan.CustomerID,
			&loan.PrincipalAmount,
			&loan.InterestRate,
			&loan.TotalAmount,
			&loan.OutstandingAmount,
			&loan.Status,
			&loan.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan loan row: %w", err)
		}

		loans = append(loans, loan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating loan rows: %w", err)
	}

	return loans, nil
}

// ListByCustomerID returns all loans for a given customer.
func (r *PostgresLoanRepository) ListByCustomerID(ctx context.Context, customerID int) ([]model.Loan, error) {
	query := `
		SELECT id, customer_id, principal_amount, interest_rate, total_amount, outstanding_amount, status, created_at
		FROM loans
		WHERE customer_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list loans by customer id: %w", err)
	}
	defer rows.Close()

	loans := []model.Loan{}

	for rows.Next() {
		var loan model.Loan

		err := rows.Scan(
			&loan.ID,
			&loan.CustomerID,
			&loan.PrincipalAmount,
			&loan.InterestRate,
			&loan.TotalAmount,
			&loan.OutstandingAmount,
			&loan.Status,
			&loan.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan loan row: %w", err)
		}

		loans = append(loans, loan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating customer loan rows: %w", err)
	}

	return loans, nil
}

// Update updates the mutable fields of a loan and returns the updated record.
func (r *PostgresLoanRepository) Update(ctx context.Context, loan model.Loan) (model.Loan, error) {
	query := `
		UPDATE loans
		SET
			customer_id = $1,
			principal_amount = $2,
			interest_rate = $3,
			total_amount = $4,
			outstanding_amount = $5,
			status = $6
		WHERE id = $7
		RETURNING id, customer_id, principal_amount, interest_rate, total_amount, outstanding_amount, status, created_at
	`

	var updatedLoan model.Loan

	err := r.db.QueryRowContext(
		ctx,
		query,
		loan.CustomerID,
		loan.PrincipalAmount,
		loan.InterestRate,
		loan.TotalAmount,
		loan.OutstandingAmount,
		loan.Status,
		loan.ID,
	).Scan(
		&updatedLoan.ID,
		&updatedLoan.CustomerID,
		&updatedLoan.PrincipalAmount,
		&updatedLoan.InterestRate,
		&updatedLoan.TotalAmount,
		&updatedLoan.OutstandingAmount,
		&updatedLoan.Status,
		&updatedLoan.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Loan{}, fmt.Errorf("%w: loan not found", apperror.ErrNotFound)
		}
		return model.Loan{}, fmt.Errorf("%w: failed to update loan", apperror.ErrInternal)
	}

	return updatedLoan, nil
}
