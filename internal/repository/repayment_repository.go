package repository

import (
	"context"
	"go-loan-management-api/internal/model"
)

// RepaymentRepository defines repayment data access behavior.
type RepaymentRepository interface {
	Create(ctx context.Context, repayment model.Repayment) (model.Repayment, error)
	GetByID(ctx context.Context, id int) (model.Repayment, error)
	ListByLoanID(ctx context.Context, loanID int) ([]model.Repayment, error)
}
