package unit

import (
	"context"
	"testing"

	"github.com/aitu/food-delivery/payment-service/internal/domain/entities"
	"github.com/aitu/food-delivery/payment-service/internal/usecase"
)

type mockRepo struct{}

func (m *mockRepo) CreateInTx(ctx context.Context, payment *entities.Payment) error {
	return nil
}
func (m *mockRepo) GetByID(ctx context.Context, id string) (*entities.Payment, error) {
	return &entities.Payment{ID: id, Status: entities.StatusCompleted}, nil
}
func (m *mockRepo) ListByOrder(ctx context.Context, orderID string) ([]entities.Payment, error) {
	return nil, nil
}

type mockPub struct{}

func (m *mockPub) Publish(subject string, payload interface{}) error { return nil }

func TestProcessOrderCreated(t *testing.T) {
	uc := usecase.NewPaymentUsecase(&mockRepo{}, &mockPub{})
	p, err := uc.ProcessOrderCreated(context.Background(), usecase.OrderCreatedEvent{
		OrderID: "o1", UserID: "u1", RestaurantID: "r1", Amount: 1000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.Status != entities.StatusCompleted {
		t.Fatalf("expected completed, got %s", p.Status)
	}
}
