package services

import "errors"

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrInvalidOrder     = errors.New("invalid order")
	ErrOrderAlreadyPaid = errors.New("order already paid")
	ErrPaymentFailed    = errors.New("payment failed")
)
