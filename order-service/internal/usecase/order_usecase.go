package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"order-service/internal/domain/entities"
	"order-service/internal/domain/repositories"
	"order-service/internal/domain/services"
)

type EventPublisher interface {
	Publish(subject string, payload interface{}) error
}

type OrderUsecase struct {
	repo      repositories.OrderRepository
	publisher EventPublisher
}

func NewOrderUsecase(repo repositories.OrderRepository, publisher EventPublisher) *OrderUsecase {
	return &OrderUsecase{repo: repo, publisher: publisher}
}

type OrderCreatedEvent struct {
	OrderID         string `json:"order_id"`
	UserID          string `json:"user_id"`
	RestaurantID    string `json:"restaurant_id"`
	Amount          int64  `json:"amount"`
	Status          string `json:"status"`
	DeliveryAddress string `json:"delivery_address"`
	UserEmail       string `json:"user_email"`
}

func (u *OrderUsecase) CreateOrder(
	ctx context.Context,
	userID, restaurantID string,
	totalPrice int64,
	address, comment, userEmail string,
) (*entities.Order, error) {
	order := &entities.Order{
		ID:            uuid.NewString(),
		UserID:        userID,
		RestaurantID:  restaurantID,
		TotalPrice:    totalPrice,
		Status:        services.OrderPending,
		PaymentStatus: services.PaymentPending,
		Address:       address,
		Comment:       comment,
		CreatedAt:     time.Now(),
	}

	if err := u.repo.CreateInTx(ctx, order); err != nil {
		return nil, err
	}

	if u.publisher != nil {
		_ = u.publisher.Publish("order.created", OrderCreatedEvent{
			OrderID:         order.ID,
			UserID:          order.UserID,
			RestaurantID:    order.RestaurantID,
			Amount:          order.TotalPrice,
			Status:          order.Status,
			DeliveryAddress: order.Address,
			UserEmail:       userEmail,
		})
	}
	return order, nil
}

func (u *OrderUsecase) GetOrder(ctx context.Context, id string) (*entities.Order, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *OrderUsecase) ListOrders(ctx context.Context) ([]entities.Order, error) {
	return u.repo.List(ctx)
}

func (u *OrderUsecase) ListOrdersByUser(ctx context.Context, userID string) ([]entities.Order, error) {
	return u.repo.ListByUser(ctx, userID)
}

func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, id, status string) error {
	return u.repo.UpdateStatus(ctx, id, status)
}

func (u *OrderUsecase) UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) error {
	return u.repo.UpdatePaymentStatus(ctx, id, paymentStatus)
}

func (u *OrderUsecase) CancelOrder(ctx context.Context, id string) error {
	if err := u.repo.UpdateStatus(ctx, id, services.OrderCancelled); err != nil {
		return err
	}
	if u.publisher != nil {
		_ = u.publisher.Publish("order.cancelled", map[string]string{"order_id": id})
	}
	return nil
}
