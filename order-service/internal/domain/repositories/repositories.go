package repositories

import (
	"context"

	"order-service/internal/domain/entities"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	GetByID(ctx context.Context, id string) (*entities.Order, error)
	List(ctx context.Context) ([]entities.Order, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdatePaymentStatus(ctx context.Context, id string, paymentStatus string) error
}
