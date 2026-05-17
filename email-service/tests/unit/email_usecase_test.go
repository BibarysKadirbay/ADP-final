package unit

import (
	"testing"

	"github.com/aitu/food-delivery/email-service/internal/infrastructure/smtp"
	"github.com/aitu/food-delivery/email-service/internal/usecase"
)

func TestHandlePaymentCompletedWithoutSMTP(t *testing.T) {
	uc := usecase.NewEmailUsecase(smtp.New("smtp.gmail.com", "587", "", ""))
	err := uc.HandlePaymentCompleted(usecase.PaymentCompletedEvent{
		OrderID: "o1", UserID: "u1", Amount: 500, UserEmail: "test@example.com",
	})
	if err == nil {
		t.Fatal("expected smtp credentials error")
	}
}
