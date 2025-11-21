package models

import "time"

// Payment ...
type Payment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}
