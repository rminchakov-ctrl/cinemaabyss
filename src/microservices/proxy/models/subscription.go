package models

import "time"

// Subscription ...
type Subscription struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	PlanType string    `json:"plan_type"`
	BeginAt  time.Time `json:"start_date"`
	EndAt    time.Time `json:"end_date"`
}
