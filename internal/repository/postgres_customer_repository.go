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

// PostgresCustomerRepository handles database operations for the Customer model
type PostgresCustomerRepository struct {
	db *sql.DB // Reference to the database connection pool
}

// NewPostgresCustomerRepository initializes the repository with a database connection
func NewPostgresCustomerRepository(db *sql.DB) *PostgresCustomerRepository {
	return &PostgresCustomerRepository{db: db}
}

// Create adds a new customer and returns the saved record including generated fields
func (r *PostgresCustomerRepository) Create(ctx context.Context, customer model.Customer) (model.Customer, error) {
	// SQL query using $ placeholders to prevent SQL injection
	query := `INSERT INTO customers (full_name, email, phone)
              VALUES ($1, $2, $3)
              RETURNING id, full_name, email, phone, created_at
            `
	var savedCustomer model.Customer

	// Execute query and map returned values (id/created_at) into the struct
	err := r.db.QueryRowContext(
		ctx,
		query,
		customer.FullName,
		customer.Email,
		customer.Phone,
	).Scan(
		&savedCustomer.ID,
		&savedCustomer.FullName,
		&savedCustomer.Email,
		&savedCustomer.Phone,
		&savedCustomer.CreatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// 23505 = unique_violation
			if pqErr.Code == "23505" {
				return model.Customer{}, fmt.Errorf("%w: email already exists", apperror.ErrConflict)
			}
		}

		return model.Customer{}, fmt.Errorf("%w: failed to create customer", apperror.ErrInternal)
	}

	return savedCustomer, nil
}

// GetByID fetches a specific customer using their unique ID
func (r *PostgresCustomerRepository) GetByID(ctx context.Context, id int) (model.Customer, error) {
	query := `
        SELECT id, full_name, email, phone, created_at
        FROM customers
        WHERE id = $1
    `

	var customer model.Customer

	// QueryRowContext is used because we expect exactly one result
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID,
		&customer.FullName,
		&customer.Email,
		&customer.Phone,
		&customer.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Customer{}, fmt.Errorf("%w: customer not found", apperror.ErrNotFound)
		}
		return model.Customer{}, fmt.Errorf("%w: failed to get customer by id", apperror.ErrInternal)
	}

	return customer, nil
}

// List retrieves all customers from the database sorted by their ID
func (r *PostgresCustomerRepository) List(ctx context.Context) ([]model.Customer, error) {
	query := `
        SELECT id, full_name, email, phone, created_at
        FROM customers
        ORDER BY id ASC
    `

	// QueryContext is used for returning multiple rows
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	// Ensure the database result set is closed when the function exits
	defer rows.Close()

	customers := []model.Customer{}

	// Loop through each row in the result set
	for rows.Next() {
		var customer model.Customer

		// Copy columns from the current row into the customer struct
		err := rows.Scan(
			&customer.ID,
			&customer.FullName,
			&customer.Email,
			&customer.Phone,
			&customer.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer row: %w", err)
		}

		// Add the populated customer to our slice
		customers = append(customers, customer)
	}

	// Check for errors that occurred during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating customer rows: %w", err)
	}

	return customers, nil
}
