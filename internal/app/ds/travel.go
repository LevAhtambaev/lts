package ds

import (
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

type Travel struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	DateStart   DateOnlyTime `json:"date_start"`
	DateEnd     DateOnlyTime `json:"date_end"`
	Preview     string       `json:"preview"`
}

type DateOnlyTime struct {
	time.Time
}

func (t *DateOnlyTime) UnmarshalJSON(b []byte) (err error) {
	log.Print(string(b))
	date, err := time.Parse(time.DateOnly, strings.Replace(string(b), "\"", "", -1))
	if err != nil {
		return err
	}
	t.Time = date
	return
}
