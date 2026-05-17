package usecase

import (
	"context"
	"errors"
	"time"

	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"
	"user-service/internal/domain/services"
	infraauth "user-service/internal/infrastructure/auth"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo      repositories.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewUserUsecase(repo repositories.UserRepository, jwtSecret string) *UserUsecase {
	return &UserUsecase{
		repo:      repo,
		jwtSecret: jwtSecret,
		tokenTTL:  24 * time.Hour,
	}
}

func (u *UserUsecase) Register(ctx context.Context, name, email, password, phone, address string) (*entities.User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}
	user := &entities.User{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Phone:        phone,
		Address:      address,
		CreatedAt:    time.Now(),
	}
	if err := u.repo.Create(ctx, user); err != nil {
		return nil, "", err
	}
	token, err := infraauth.GenerateToken(user.ID, user.Email, u.jwtSecret, u.tokenTTL)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (u *UserUsecase) Login(ctx context.Context, email, password string) (*entities.User, string, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, "", services.ErrInvalidCredentials
		}
		return nil, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", services.ErrInvalidCredentials
	}
	token, err := infraauth.GenerateToken(user.ID, user.Email, u.jwtSecret, u.tokenTTL)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (u *UserUsecase) CreateUser(ctx context.Context, name, email, phone, address string) (*entities.User, error) {
	user := &entities.User{
		ID:        uuid.NewString(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Address:   address,
		CreatedAt: time.Now(),
	}
	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) GetUser(ctx context.Context, id string) (*entities.User, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UserUsecase) UpdateProfile(ctx context.Context, id, name, phone, address string) (*entities.User, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Name = name
	user.Phone = phone
	user.Address = address
	if err := u.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) ListUsers(ctx context.Context) ([]entities.User, error) {
	return u.repo.List(ctx)
}
