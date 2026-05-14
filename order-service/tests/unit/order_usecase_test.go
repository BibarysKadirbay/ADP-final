package unit

import (
	"testing"

	"order-service/internal/domain/services"
)

func TestOrderStatus(t *testing.T) {

	status := services.OrderPending

	if status != "ORDER_PENDING" {
		t.Fatalf(
			"expected ORDER_PENDING, got %s",
			status,
		)
	}
}

func TestPaymentStatus(t *testing.T) {

	status := services.PaymentPending

	if status != "PAYMENT_PENDING" {
		t.Fatalf(
			"expected PAYMENT_PENDING, got %s",
			status,
		)
	}
}
