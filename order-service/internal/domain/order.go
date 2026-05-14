package domain

type Order struct {
	ID           string
	UserID       string
	RestaurantID string
	TotalPrice   int64
	Status       string
}
