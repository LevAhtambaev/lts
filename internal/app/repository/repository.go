package repository

import (
	"context"
	"github.com/google/uuid"
	"lts/internal/app/ds"
)

type TravelRepository interface {
	CreateTravel(ctx context.Context, travel ds.Travel) (ds.Travel, error)
	SetTravelPreview(ctx context.Context, path string, uuid uuid.UUID) error
	AddPlace(ctx context.Context, travelUUID, placeUUID uuid.UUID) error
	GetTravel(ctx context.Context, travelUUID uuid.UUID) (ds.Travel, error)
	UpdateTravel(ctx context.Context, id uuid.UUID, travel ds.Travel) error
}

type PlaceRepository interface {
	CreatePlace(ctx context.Context, place ds.Place) (ds.Place, error)
	SetExpenses(ctx context.Context, uuidExpense, uuidPlace uuid.UUID) error
	SetPreview(ctx context.Context, path string, uuid uuid.UUID) error
	SetImages(ctx context.Context, paths []string, uuid uuid.UUID) error
	DeletePlace(ctx context.Context, uuid uuid.UUID) error
	UpdatePlace(ctx context.Context, id uuid.UUID, place ds.Place) error
	GetPlace(ctx context.Context, id uuid.UUID) (ds.Place, error)
}

type ExpensesRepository interface {
	CreateExpense(ctx context.Context, expense ds.Expense) (ds.Expense, error)
	GetExpense(ctx context.Context, uuid uuid.UUID) (ds.Expense, error)
	UpdateExpense(ctx context.Context, expense ds.Expense, uuid uuid.UUID) (ds.Expense, error)
	DeleteExpense(ctx context.Context, uuid uuid.UUID) error
}
