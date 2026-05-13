package services

import "errors"

var (
	ErrInvalidInput          = errors.New("invalid input")
	ErrNotFound              = errors.New("not found")
	ErrNoCourierAvailable    = errors.New("no courier available")
	ErrInvalidTransition     = errors.New("invalid delivery status transition")
	ErrRestaurantUnavailable = errors.New("restaurant service unavailable")
	ErrConflict              = errors.New("conflict")
)
