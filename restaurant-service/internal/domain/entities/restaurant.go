package entities

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID           uuid.UUID `json:"id"`
	OwnerID      uuid.UUID `json:"owner_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CuisineType  string    `json:"cuisine_type"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	Rating       float64   `json:"rating"`
	TotalReviews int32     `json:"total_reviews"`
	ImageURL     string    `json:"image_url"`
	IsOpen       bool      `json:"is_open"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Category struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
}

type MenuItem struct {
	ID           uuid.UUID `json:"id"`
	CategoryID   uuid.UUID `json:"category_id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	ImageURL     string    `json:"image_url"`
	IsAvailable  bool      `json:"is_available"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MenuCategory struct {
	Category Category   `json:"category"`
	Items    []MenuItem `json:"items"`
}

type RestaurantFilter struct {
	Query         string
	CuisineType   string
	City          string
	OpenOnly      *bool
	SortBy        string
	SortDirection string
	Limit         int
	Offset        int
}

type PageMeta struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}
