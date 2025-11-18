package models

import "time"

type Transaction struct {
	ID          int64     `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Type        string    `json:"type"` // "income" or "expense"
	CreatedAt   time.Time `json:"created_at"`
}
