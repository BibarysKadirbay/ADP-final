package entities

import "time"

type Order struct {
	ID           string
	UserID       string
	RestaurantID string
	DeliveryID   string

	TotalPrice int64

	Status        string
	PaymentStatus string

	Address   string
	Comment   string
	CreatedAt time.Time
}
