package repository

import (
	"context"
	"go-loan-management-api/internal/model"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer model.Customer) (model.Customer, error)
	GetByID(ctx context.Context, id int) (model.Customer, error)
	List(ctx context.Context) ([]model.Customer, error)
}
