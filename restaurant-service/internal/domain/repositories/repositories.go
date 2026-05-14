package repositories

import (
	"context"
	"time"

	"github.com/aitu/food-delivery/restaurant-service/internal/domain/entities"
	"github.com/google/uuid"
)

type RestaurantRepository interface {
	Create(ctx context.Context, r *entities.Restaurant) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Restaurant, error)
	Update(ctx context.Context, r *entities.Restaurant) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter entities.RestaurantFilter) ([]entities.Restaurant, int64, error)
	TopRated(ctx context.Context, city string, limit int) ([]entities.Restaurant, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, c *entities.Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Category, error)
	Update(ctx context.Context, c *entities.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]entities.Category, error)
}

type MenuRepository interface {
	CreateItem(ctx context.Context, item *entities.MenuItem) error
	GetItemByID(ctx context.Context, id uuid.UUID) (*entities.MenuItem, error)
	UpdateItem(ctx context.Context, item *entities.MenuItem) error
	DeleteItem(ctx context.Context, id uuid.UUID) error
	GetByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]entities.MenuCategory, error)
	SetAvailability(ctx context.Context, id uuid.UUID, available bool) (*entities.MenuItem, error)
}

type CacheRepository interface {
	Get(ctx context.Context, key string, dest any) (bool, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeletePattern(ctx context.Context, pattern string) error
}

type EventPublisher interface {
	Publish(ctx context.Context, subject string, payload any) error
	Close() error
}
