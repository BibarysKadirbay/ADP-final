package repositories

import (
	"context"

	"github.com/aitu/food-delivery/payment-service/internal/domain/entities"
)

type PaymentRepository interface {
	CreateInTx(ctx context.Context, payment *entities.Payment) error
	GetByID(ctx context.Context, id string) (*entities.Payment, error)
	ListByOrder(ctx context.Context, orderID string) ([]entities.Payment, error)
}
