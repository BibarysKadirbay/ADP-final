package usecase

import (
	"fmt"

	"github.com/aitu/food-delivery/email-service/internal/infrastructure/smtp"
)

type PaymentCompletedEvent struct {
	PaymentID       string `json:"payment_id"`
	OrderID         string `json:"order_id"`
	UserID          string `json:"user_id"`
	Amount          int64  `json:"amount"`
	DeliveryAddress string `json:"delivery_address"`
	UserEmail       string `json:"user_email"`
}

type EmailUsecase struct {
	smtp *smtp.Client
}

func NewEmailUsecase(smtpClient *smtp.Client) *EmailUsecase {
	return &EmailUsecase{smtp: smtpClient}
}

func (u *EmailUsecase) HandlePaymentCompleted(event PaymentCompletedEvent) error {
	to := event.UserEmail
	if to == "" {
		to = fmt.Sprintf("user-%s@example.com", event.UserID)
	}
	subject := fmt.Sprintf("Payment confirmed for order %s", event.OrderID)
	body := fmt.Sprintf(
		"Hello,\n\nYour payment of %d was processed successfully.\nOrder ID: %s\nDelivery address: %s\n\nThank you for using Food Delivery!",
		event.Amount, event.OrderID, event.DeliveryAddress,
	)
	return u.smtp.Send(to, subject, body)
}
