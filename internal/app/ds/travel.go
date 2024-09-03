package ds

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type FullTravel struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	DateStart   DateOnlyTime `json:"date_start"`
	DateEnd     DateOnlyTime `json:"date_end"`
	Places      []FullPlace  `json:"places"`
	Preview     string       `json:"preview"`
}

type Travel struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	DateStart   DateOnlyTime `json:"date_start"`
	DateEnd     DateOnlyTime `json:"date_end"`
	Places      []uuid.UUID  `json:"places"`
	Preview     string       `json:"preview"`
}

type DateOnlyTime struct {
	time.Time
}

func (t *DateOnlyTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(time.DateOnly, strings.Replace(string(b), "\"", "", -1))
	if err != nil {
		return err
	}
	t.Time = date
	return
}
