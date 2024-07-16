package repository

import (
	"context"
	"github.com/google/uuid"
	"lts/internal/app/ds"
)

type TravelRepository interface {
	CreateTravel(ctx context.Context, travel ds.Travel) (ds.Travel, error)
	SetTravelPreview(ctx context.Context, path string, uuid uuid.UUID) error
}

type ExpensesRepository interface {
	CreateExpense(ctx context.Context, expense ds.Expense) (ds.Expense, error)
	GetExpense(ctx context.Context, uuid uuid.UUID) (ds.Expense, error)
}
