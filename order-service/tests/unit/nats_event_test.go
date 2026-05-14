package unit

import (
	"encoding/json"
	"testing"
)

type OrderCreatedEvent struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Amount  int64  `json:"amount"`
}

func TestNatsEvent(t *testing.T) {

	event := OrderCreatedEvent{
		OrderID: "order-1",
		UserID:  "user-1",
		Amount:  5000,
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("event payload is empty")
	}
}
