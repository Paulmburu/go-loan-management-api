package repository

import (
	"context"
	"go-loan-management-api/internal/model"
)

type LoanRepository interface {
	Create(ctx context.Context, loan model.Loan) (model.Loan, error)
	GetByID(ctx context.Context, id int) (model.Loan, error)
	List(ctx context.Context) ([]model.Loan, error)
	ListByCustomerID(ctx context.Context, customerID int) ([]model.Loan, error)
	Update(ctx context.Context, loan model.Loan) (model.Loan, error)
}
