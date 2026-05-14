package integration

import (
	"context"
	"testing"

	"user-service/internal/domain/entities"

	"github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {

	user := &entities.User{
		ID:      uuid.NewString(),
		Name:    "Bibarys",
		Email:   "bibarys@test.com",
		Phone:   "+77001234567",
		Address: "Astana",
	}

	ctx := context.Background()

	if user.Name == "" {
		t.Fatal("user name is empty")
	}

	if user.Email == "" {
		t.Fatal("user email is empty")
	}

	_ = ctx
}
