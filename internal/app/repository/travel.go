package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"lts/internal/app/ds"
)

type TravelRepositoryImpl struct {
	db *sqlx.DB
}

func NewTravelRepo(db *sqlx.DB) *TravelRepositoryImpl {
	return &TravelRepositoryImpl{
		db: db,
	}
}

func (t TravelRepositoryImpl) CreateTravel(ctx context.Context, travel ds.Travel) (ds.Travel, error) {
	travel.ID = uuid.New()

	_, err := t.db.ExecContext(ctx, "INSERT INTO travel VALUES ($1, $2, $3, $4, $5)", travel.ID, travel.Name, travel.Description, travel.DateStart.Time, travel.DateEnd.Time)
	if err != nil {
		return ds.Travel{}, fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return travel, nil
}

func (t TravelRepositoryImpl) SetTravelPreview(ctx context.Context, path string, uuid uuid.UUID) error {
	_, err := t.db.ExecContext(ctx, "UPDATE travel SET preview = $1 WHERE id = $2", path, uuid)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}
