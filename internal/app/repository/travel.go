package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (t TravelRepositoryImpl) UpdateTravel(ctx context.Context, id uuid.UUID, travel ds.Travel) error {
	_, err := t.db.ExecContext(ctx, "UPDATE travel SET (name, description, date_start, date_end) = ($1, $2, $3, $4) WHERE id = $5", travel.Name, travel.Description, travel.DateStart.Time, travel.DateEnd.Time, id)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

func (t TravelRepositoryImpl) DeleteTravel(ctx context.Context, id uuid.UUID) error {
	_, err := t.db.ExecContext(ctx, "DELETE FROM travel WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

func (t TravelRepositoryImpl) SetTravelPreview(ctx context.Context, path string, uuid uuid.UUID) error {
	_, err := t.db.ExecContext(ctx, "UPDATE travel SET preview = $1 WHERE id = $2", path, uuid)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

func (t TravelRepositoryImpl) AddPlace(ctx context.Context, travelUUID, placeUUID uuid.UUID) error {
	_, err := t.db.ExecContext(ctx, "UPDATE travel SET places = array_append(places, $1) WHERE id = $2", placeUUID, travelUUID)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

func (t TravelRepositoryImpl) GetTravel(ctx context.Context, travelUUID uuid.UUID) (ds.Travel, error) {
	var travel ds.Travel
	var placesBytes []byte
	err := t.db.QueryRowContext(ctx, "SELECT id, name, description, date_start, date_end, places, preview FROM travel WHERE id = $1", travelUUID).Scan(
		&travel.ID, &travel.Name, &travel.Description, &travel.DateStart.Time, &travel.DateEnd.Time, &placesBytes, &travel.Preview,
	)
	if err != nil {
		return ds.Travel{}, fmt.Errorf("[db.ExecContext]: %w", err)
	}

	var places pq.StringArray
	err = places.Scan(placesBytes)
	if err != nil {
		return ds.Travel{}, fmt.Errorf("ошибка при сканировании places: %w", err)
	}

	travel.Places = make([]uuid.UUID, len(places))

	for i, placeID := range places {
		placeUUID, err := uuid.Parse(placeID)
		if err != nil {
			return ds.Travel{}, fmt.Errorf("неверный формат UUID места: %w", err)
		}
		travel.Places[i] = placeUUID
	}

	return travel, nil
}
