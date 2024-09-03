package ds

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type FullPlace struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Story    string         `json:"story"`
	Date     DateOnlyTime   `json:"date"`
	Images   pq.StringArray `json:"images"`
	Expenses Expense        `json:"expenses"`
	Preview  string         `json:"preview"`
}

type Place struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Story    string         `json:"story"`
	Date     DateOnlyTime   `json:"date"`
	Images   pq.StringArray `json:"images"`
	Expenses uuid.UUID      `json:"expenses"`
	Preview  string         `json:"preview"`
}
