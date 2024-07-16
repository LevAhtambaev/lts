package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"lts/internal/app/ds"
)

type ExpensesRepositoryImpl struct {
	db *sqlx.DB
}

func NewExpensesRepo(db *sqlx.DB) *ExpensesRepositoryImpl {
	return &ExpensesRepositoryImpl{
		db: db,
	}
}

func (e ExpensesRepositoryImpl) CreateExpense(ctx context.Context, expense ds.Expense) error {
	expense.ID = uuid.New()

	_, err := e.db.ExecContext(ctx, "INSERT INTO expenses VALUES ($1, $2, $3, $4, $5, $6)", expense.ID, expense.Road, expense.Residence, expense.Food, expense.Entertainment, expense.Other)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}

func (e ExpensesRepositoryImpl) GetExpense(ctx context.Context, uuid uuid.UUID) (ds.Expense, error) {
	var expense ds.Expense
	err := e.db.GetContext(ctx, &expense, "SELECT * FROM expenses WHERE id = $1", uuid)
	if err != nil {
		return ds.Expense{}, fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return expense, nil
}

func (e ExpensesRepositoryImpl) UpdateExpense(ctx context.Context, expense ds.Expense, uuid uuid.UUID) error {
	_, err := e.db.ExecContext(ctx, "UPDATE expenses SET (road, residence, food, entertainment, other) = ($1, $2, $3, $4, $5) WHERE id = $6", expense.Road, expense.Residence, expense.Food, expense.Entertainment, expense.Other, uuid)
	if err != nil {
		return fmt.Errorf("[db.ExecContext]: %w", err)
	}
	return nil
}
