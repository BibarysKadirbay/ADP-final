package services

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidUser        = errors.New("invalid user")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
