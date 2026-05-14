package usecase

import (
	"context"
	"time"

	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"

	"github.com/google/uuid"
)

type UserUsecase struct {
	repo repositories.UserRepository
}

func NewUserUsecase(
	repo repositories.UserRepository,
) *UserUsecase {

	return &UserUsecase{
		repo: repo,
	}
}

func (u *UserUsecase) CreateUser(
	ctx context.Context,
	name string,
	email string,
	phone string,
	address string,
) (*entities.User, error) {

	user := &entities.User{
		ID:        uuid.NewString(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Address:   address,
		CreatedAt: time.Now(),
	}

	err := u.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) GetUser(
	ctx context.Context,
	id string,
) (*entities.User, error) {

	return u.repo.GetByID(ctx, id)
}

func (u *UserUsecase) ListUsers(
	ctx context.Context,
) ([]entities.User, error) {

	return u.repo.List(ctx)
}
