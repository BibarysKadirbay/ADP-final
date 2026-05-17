package entities

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	Phone        string
	Address      string
	CreatedAt    time.Time
}
