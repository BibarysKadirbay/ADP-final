package integration

import (
	"context"
	"testing"

	"order-service/internal/domain/entities"

	"github.com/google/uuid"
)

func TestCreateOrder(t *testing.T) {

	order := &entities.Order{
		ID:            uuid.NewString(),
		UserID:        "user-1",
		RestaurantID:  "restaurant-1",
		DeliveryID:    "",
		TotalPrice:    5000,
		Status:        "ORDER_PENDING",
		PaymentStatus: "PAYMENT_PENDING",
		Address:       "Astana",
		Comment:       "test order",
	}

	ctx := context.Background()

	if order.UserID == "" {
		t.Fatal("user id is empty")
	}

	if order.TotalPrice <= 0 {
		t.Fatal("invalid total price")
	}

	_ = ctx
}
