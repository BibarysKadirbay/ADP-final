package entities

import "time"

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

type Payment struct {
	ID        string
	OrderID   string
	UserID    string
	Amount    int64
	Status    string
	Method    string
	CreatedAt time.Time
}
