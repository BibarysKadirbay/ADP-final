package usecase

import (
	"context"
	"time"

	"github.com/aitu/food-delivery/payment-service/internal/domain/entities"
	"github.com/aitu/food-delivery/payment-service/internal/domain/repositories"
	"github.com/google/uuid"
)

type EventPublisher interface {
	Publish(subject string, payload interface{}) error
}

type OrderCreatedEvent struct {
	OrderID         string `json:"order_id"`
	UserID          string `json:"user_id"`
	RestaurantID    string `json:"restaurant_id"`
	Amount          int64  `json:"amount"`
	DeliveryAddress string `json:"delivery_address"`
	UserEmail       string `json:"user_email"`
}

type PaymentCompletedEvent struct {
	PaymentID       string `json:"payment_id"`
	OrderID         string `json:"order_id"`
	UserID          string `json:"user_id"`
	RestaurantID    string `json:"restaurant_id"`
	Amount          int64  `json:"amount"`
	DeliveryAddress string `json:"delivery_address"`
	UserEmail       string `json:"user_email"`
}

type PaymentUsecase struct {
	repo      repositories.PaymentRepository
	publisher EventPublisher
}

func NewPaymentUsecase(repo repositories.PaymentRepository, publisher EventPublisher) *PaymentUsecase {
	return &PaymentUsecase{repo: repo, publisher: publisher}
}

func (u *PaymentUsecase) ProcessOrderCreated(ctx context.Context, event OrderCreatedEvent) (*entities.Payment, error) {
	payment := &entities.Payment{
		ID:        uuid.NewString(),
		OrderID:   event.OrderID,
		UserID:    event.UserID,
		Amount:    event.Amount,
		Status:    entities.StatusCompleted,
		Method:    "card",
		CreatedAt: time.Now(),
	}
	if err := u.repo.CreateInTx(ctx, payment); err != nil {
		return nil, err
	}
	if u.publisher != nil {
		_ = u.publisher.Publish("payment.completed", PaymentCompletedEvent{
			PaymentID:       payment.ID,
			OrderID:         event.OrderID,
			UserID:          event.UserID,
			RestaurantID:    event.RestaurantID,
			Amount:          event.Amount,
			DeliveryAddress: event.DeliveryAddress,
			UserEmail:       event.UserEmail,
		})
	}
	return payment, nil
}

func (u *PaymentUsecase) GetPayment(ctx context.Context, id string) (*entities.Payment, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *PaymentUsecase) ListByOrder(ctx context.Context, orderID string) ([]entities.Payment, error) {
	return u.repo.ListByOrder(ctx, orderID)
}
