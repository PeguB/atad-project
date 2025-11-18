package models

import "time"

type Budget struct {
	ID        int64     `json:"id"`
	Category  string    `json:"category"`
	Amount    float64   `json:"amount"`
	Period    string    `json:"period"`     // Always "custom"
	StartDate time.Time `json:"start_date"` // Required
	EndDate   time.Time `json:"end_date"`   // Required
}
