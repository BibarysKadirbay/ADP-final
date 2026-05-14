package repositories

import (
	"context"

	"user-service/internal/domain/entities"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		user *entities.User,
	) error

	GetByID(
		ctx context.Context,
		id string,
	) (*entities.User, error)

	List(
		ctx context.Context,
	) ([]entities.User, error)
}
