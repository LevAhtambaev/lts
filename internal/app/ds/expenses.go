package ds

import "github.com/google/uuid"

type Expense struct {
	ID            uuid.UUID `json:"id"`
	Road          int       `json:"road"`
	Residence     int       `json:"residence"`
	Food          int       `json:"food"`
	Entertainment int       `json:"entertainment"`
	Other         int       `json:"other"`
}
