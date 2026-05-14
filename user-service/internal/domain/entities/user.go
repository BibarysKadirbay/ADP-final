package entities

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	Phone     string
	Address   string
	CreatedAt time.Time
}
