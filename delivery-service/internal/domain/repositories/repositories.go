package repositories

import (
	"context"
	"time"

	"github.com/aitu/food-delivery/delivery-service/internal/domain/entities"
	"github.com/google/uuid"
)

type CourierRepository interface {
	Create(ctx context.Context, courier *entities.Courier) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Courier, error)
	UpdateAvailability(ctx context.Context, id uuid.UUID, available bool) (*entities.Courier, error)
	ListAvailable(ctx context.Context, filter entities.CourierFilter) ([]entities.Courier, int64, error)
	RecalculateRating(ctx context.Context, courierID uuid.UUID) error
}

type DeliveryRepository interface {
	Assign(ctx context.Context, delivery *entities.Delivery) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Delivery, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*entities.Delivery, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entities.DeliveryStatus) (*entities.Delivery, error)
	ListByCourier(ctx context.Context, courierID uuid.UUID, filter entities.DeliveryFilter) ([]entities.Delivery, int64, error)
	ListByOrder(ctx context.Context, orderID uuid.UUID) ([]entities.Delivery, error)
	History(ctx context.Context, deliveryID uuid.UUID) ([]entities.DeliveryStatusHistory, error)
	Stats(ctx context.Context, courierID uuid.UUID) (*entities.DeliveryStats, error)
	ActiveCountByCourier(ctx context.Context, courierID uuid.UUID) (int64, error)
	CancelByOrder(ctx context.Context, orderID uuid.UUID) (*entities.Delivery, error)
}

type RatingRepository interface {
	Create(ctx context.Context, rating *entities.CourierRating) error
}

type Cache interface {
	Get(ctx context.Context, key string, dest any) (bool, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeletePattern(ctx context.Context, pattern string) error
}

type EventPublisher interface {
	Publish(ctx context.Context, subject string, payload any) error
}

type RestaurantClient interface {
	GetRestaurant(ctx context.Context, id uuid.UUID) (*entities.RestaurantSnapshot, error)
}
