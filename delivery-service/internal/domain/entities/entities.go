package entities

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryStatus string

const (
	StatusPending   DeliveryStatus = "pending"
	StatusAssigned  DeliveryStatus = "assigned"
	StatusPickedUp  DeliveryStatus = "picked_up"
	StatusOnTheWay  DeliveryStatus = "on_the_way"
	StatusDelivered DeliveryStatus = "delivered"
	StatusCancelled DeliveryStatus = "cancelled"
)

type Courier struct {
	ID              uuid.UUID `db:"id"`
	UserID          uuid.UUID `db:"user_id"`
	FullName        string    `db:"full_name"`
	Phone           string    `db:"phone"`
	VehicleType     string    `db:"vehicle_type"`
	Rating          float64   `db:"rating"`
	TotalDeliveries int32     `db:"total_deliveries"`
	IsAvailable     bool      `db:"is_available"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type Delivery struct {
	ID                  uuid.UUID      `db:"id"`
	OrderID             uuid.UUID      `db:"order_id"`
	CourierID           uuid.UUID      `db:"courier_id"`
	RestaurantID        uuid.UUID      `db:"restaurant_id"`
	CustomerID          uuid.UUID      `db:"customer_id"`
	Status              DeliveryStatus `db:"status"`
	PickupAddress       string         `db:"pickup_address"`
	DeliveryAddress     string         `db:"delivery_address"`
	EstimatedETAMinutes int32          `db:"estimated_eta_minutes"`
	PickupTime          *time.Time     `db:"pickup_time"`
	DeliveredTime       *time.Time     `db:"delivered_time"`
	RouteDistanceKM     float64        `db:"route_distance_km"`
	CreatedAt           time.Time      `db:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at"`
}

type DeliveryStatusHistory struct {
	ID         uuid.UUID       `db:"id"`
	DeliveryID uuid.UUID       `db:"delivery_id"`
	OldStatus  *DeliveryStatus `db:"old_status"`
	NewStatus  DeliveryStatus  `db:"new_status"`
	ChangedAt  time.Time       `db:"changed_at"`
}

type CourierRating struct {
	ID         uuid.UUID `db:"id"`
	CourierID  uuid.UUID `db:"courier_id"`
	OrderID    uuid.UUID `db:"order_id"`
	CustomerID uuid.UUID `db:"customer_id"`
	Rating     int32     `db:"rating"`
	Comment    string    `db:"comment"`
	CreatedAt  time.Time `db:"created_at"`
}

type DeliveryFilter struct {
	Status        DeliveryStatus
	SortBy        string
	SortDirection string
	Limit         int
	Offset        int
}

type CourierFilter struct {
	VehicleType string
	Limit       int
	Offset      int
}

type PageMeta struct {
	Page     int
	PageSize int
	Total    int64
}

type DeliveryStats struct {
	TotalDeliveries     int64
	ActiveDeliveries    int64
	CompletedDeliveries int64
	CancelledDeliveries int64
	AverageETAMinutes   float64
	AverageDistanceKM   float64
}

type RestaurantSnapshot struct {
	ID      uuid.UUID
	Name    string
	Address string
}
