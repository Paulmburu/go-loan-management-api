package service

import (
	"context"
	"database/sql"
	"fmt"
	"go-loan-management-api/internal/model"
	"go-loan-management-api/internal/repository"
)

// RepaymentService contains business logic related to repayments.
type RepaymentService struct {
	db                  *sql.DB
	loanRepository      repository.LoanRepository
	repaymentRepository repository.RepaymentRepository
}

// NewRepaymentService creates a new RepaymentService.
func NewRepaymentService(db *sql.DB, loanRepository repository.LoanRepository, repaymentRepository repository.RepaymentRepository) *RepaymentService {
	return &RepaymentService{
		db:                  db,
		loanRepository:      loanRepository,
		repaymentRepository: repaymentRepository,
	}
}

// AddRepayment validates and records a repayment for a loan.
// It also updates the loan's outstanding balance and status.
func (s *RepaymentService) AddRepayment(ctx context.Context, loanID int, amount float64) (model.Repayment, error) {
	// Validate the loan ID.
	if loanID <= 0 {
		return model.Repayment{}, fmt.Errorf("loan_id must be a positive integer")
	}

	// Validate the repayment amount.
	if amount <= 0 {
		return model.Repayment{}, fmt.Errorf("amount must be greater than zero")
	}

	// Start a database transaction. This ensures that the repayment creation and loan update happen atomically.
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Repayment{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var loan model.Loan

	// Step 1: fetch the loan inside the transaction.
	getLoanQuery := `
		SELECT id, customer_id, principal_amount, interest_rate, total_amount, outstanding_amount, status, created_at
		FROM loans
		WHERE id = $1
	`

	err = tx.QueryRowContext(ctx, getLoanQuery, loanID).Scan(
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
		if err == sql.ErrNoRows {
			return model.Repayment{}, fmt.Errorf("loan not found")
		}
		return model.Repayment{}, fmt.Errorf("failed to get loan: %w", err)
	}

	// Step 2: validate repayment rules.
	// if loan.Status != model.LoanStatusActive {
	// 	return model.Repayment{}, fmt.Errorf("only active loans can accept repayments")
	// }

	if !model.CanAcceptRepayment(loan) {
		return model.Repayment{}, fmt.Errorf("only active loans can accept repayments")
	}

	// Prevent overpayment.
	if amount > loan.OutstandingAmount {
		return model.Repayment{}, fmt.Errorf("repayment amount cannot exceed outstanding amount of %.2f", loan.OutstandingAmount)
	}

	// Step 3: insert repayment inside the transaction.
	var savedRepayment model.Repayment
	createRepaymentQuery := `
		INSERT INTO repayments (loan_id, amount)
		VALUES ($1, $2)
		RETURNING id, loan_id, amount, created_at
	`

	err = tx.QueryRowContext(ctx, createRepaymentQuery, loanID, amount).Scan(
		&savedRepayment.ID,
		&savedRepayment.LoanID,
		&savedRepayment.Amount,
		&savedRepayment.CreatedAt,
	)
	if err != nil {
		return model.Repayment{}, fmt.Errorf("failed to create repayment: %w", err)
	}

	// Step 4: update the in-memory copy of the loan state.
	loan.OutstandingAmount -= amount
	if loan.OutstandingAmount == 0 {
		loan.Status = model.LoanStatusPaid
	}

	// Step 5: persist the updated loan inside the same transaction.
	updateLoanQuery := `
		UPDATE loans
		SET outstanding_amount = $1, status = $2
		WHERE id = $3
	`

	_, err = tx.ExecContext(ctx, updateLoanQuery, loan.OutstandingAmount, loan.Status, loan.ID)
	if err != nil {
		return model.Repayment{}, fmt.Errorf("failed to update loan after repayment: %w", err)
	}

	// Step 6: commit the transaction.
	if err := tx.Commit(); err != nil {
		return model.Repayment{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return savedRepayment, nil
}

// ListRepaymentsByLoanID returns all repayments for a given loan.
// It first verifies that the loan exists.
func (s *RepaymentService) ListRepaymentsByLoanID(ctx context.Context, loanID int) ([]model.Repayment, error) {
	if loanID <= 0 {
		return nil, fmt.Errorf("loan_id must be greater than zero")
	}

	_, err := s.loanRepository.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	return s.repaymentRepository.ListByLoanID(ctx, loanID)
}

// GetRepaymentByID returns a repayment by ID.
func (s *RepaymentService) GetRepaymentByID(ctx context.Context, repaymentID int) (model.Repayment, error) {
	if repaymentID <= 0 {
		return model.Repayment{}, fmt.Errorf("repayment_id must be greater than zero")
	}

	return s.repaymentRepository.GetByID(ctx, repaymentID)
}

// GetRepaymentWithLoan returns a repayment and its related loan.
// This is useful for ownership checks.
func (s *RepaymentService) GetRepaymentWithLoan(ctx context.Context, repaymentID int) (model.Repayment, model.Loan, error) {
	if repaymentID <= 0 {
		return model.Repayment{}, model.Loan{}, fmt.Errorf("repayment_id must be greater than zero")
	}

	repayment, err := s.repaymentRepository.GetByID(ctx, repaymentID)
	if err != nil {
		return model.Repayment{}, model.Loan{}, err
	}

	loan, err := s.loanRepository.GetByID(ctx, repayment.LoanID)
	if err != nil {
		return model.Repayment{}, model.Loan{}, err
	}

	return repayment, loan, nil
}
