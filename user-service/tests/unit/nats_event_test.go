package unit

import (
	"encoding/json"
	"testing"
)

type UserCreatedEvent struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func TestNatsEvent(t *testing.T) {

	event := UserCreatedEvent{
		UserID: "user-1",
		Name:   "Bibarys",
		Email:  "bibarys@test.com",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("event payload is empty")
	}
}
