package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"lts/internal/app/ds"
)

type PlaceRepositoryImpl struct {
	db *sqlx.DB
}

func NewPlaceRepositoryImpl(db *sqlx.DB) *PlaceRepositoryImpl {
	return &PlaceRepositoryImpl{db: db}
}

func (p PlaceRepositoryImpl) CreatePlace(ctx context.Context, place ds.Place) (ds.Place, error) {
	place.ID = uuid.New()

	_, err := p.db.ExecContext(ctx, "INSERT INTO places VALUES ($1, $2, $3, $4)", place.ID, place.Name, place.Story, place.Date.Time)
	if err != nil {
		return ds.Place{}, fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return place, nil
}

func (p PlaceRepositoryImpl) SetExpenses(ctx context.Context, uuidExpense, uuidPlace uuid.UUID) error {
	_, err := p.db.ExecContext(ctx, "UPDATE places SET expenses = $1 WHERE id = $2", uuidExpense, uuidPlace)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

func (p PlaceRepositoryImpl) SetPreview(ctx context.Context, path string, uuid uuid.UUID) error {
	_, err := p.db.ExecContext(ctx, "UPDATE places SET preview = $1 WHERE id = $2", path, uuid)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

func (p PlaceRepositoryImpl) SetImages(ctx context.Context, paths []string, uuid uuid.UUID) error {
	_, err := p.db.ExecContext(ctx, "UPDATE places SET images = $1 WHERE id = $2", pq.StringArray(paths), uuid)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

func (p PlaceRepositoryImpl) DeletePlace(ctx context.Context, uuid uuid.UUID) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM places WHERE id = $1", uuid)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

func (p PlaceRepositoryImpl) UpdatePlace(ctx context.Context, id uuid.UUID, place ds.Place) error {
	_, err := p.db.ExecContext(ctx, "UPDATE places SET (name, story, date) = ($1, $2, $3) WHERE id = $4", place.Name, place.Story, place.Date.Time, id)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}

	return nil
}

func (p PlaceRepositoryImpl) GetPlace(ctx context.Context, id uuid.UUID) (ds.Place, error) {
	place := ds.Place{}

	var preview sql.NullString

	err := p.db.QueryRowContext(ctx, "SELECT id, name, story, date, images, expenses, preview FROM places WHERE id = $1", id).Scan(
		&place.ID, &place.Name, &place.Story, &place.Date.Time, &place.Images, &place.Expenses, &preview,
	)
	if err != nil {
		return ds.Place{}, fmt.Errorf("[db.ExecContext]: %w", err)
	}

	place.Preview = preview.String

	return place, nil

}
