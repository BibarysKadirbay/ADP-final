package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"order-service/internal/domain/entities"
	"order-service/internal/domain/repositories"
	"order-service/internal/domain/services"
)

type OrderUsecase struct {
	repo repositories.OrderRepository
}

func NewOrderUsecase(repo repositories.OrderRepository) *OrderUsecase {
	return &OrderUsecase{
		repo: repo,
	}
}

func (u *OrderUsecase) CreateOrder(
	ctx context.Context,
	userID string,
	restaurantID string,
	deliveryID string,
	totalPrice int64,
	address string,
	comment string,
) (*entities.Order, error) {

	order := &entities.Order{
		ID:            uuid.NewString(),
		UserID:        userID,
		RestaurantID:  restaurantID,
		DeliveryID:    deliveryID,
		TotalPrice:    totalPrice,
		Status:        services.OrderPending,
		PaymentStatus: services.PaymentPending,
		Address:       address,
		Comment:       comment,
		CreatedAt:     time.Now(),
	}

	err := u.repo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *OrderUsecase) GetOrder(
	ctx context.Context,
	id string,
) (*entities.Order, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *OrderUsecase) ListOrders(
	ctx context.Context,
) ([]entities.Order, error) {
	return u.repo.List(ctx)
}

func (u *OrderUsecase) UpdateOrderStatus(
	ctx context.Context,
	id string,
	status string,
) error {
	return u.repo.UpdateStatus(ctx, id, status)
}

func (u *OrderUsecase) UpdatePaymentStatus(
	ctx context.Context,
	id string,
	paymentStatus string,
) error {
	return u.repo.UpdatePaymentStatus(ctx, id, paymentStatus)
}

func (u *OrderUsecase) CancelOrder(
	ctx context.Context,
	id string,
) error {
	return u.repo.UpdateStatus(ctx, id, services.OrderCancelled)
}
