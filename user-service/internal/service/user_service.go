package service

import (
	"context"
	"errors"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	userpb "food-delivery-system/proto/userpb"
)

var jwtSecret = []byte("secret-key")

type UserService struct {
	userpb.UnimplementedUserServiceServer
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, req *userpb.RegisterRequest) (*userpb.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}

func (s *UserService) LoginUser(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	return &userpb.LoginResponse{
		Token: tokenString,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	user, err := s.repo.GetByID(req.UserId)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	user := domain.User{
		ID:    req.UserId,
		Name:  req.Name,
		Email: req.Email,
	}

	err := s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return &userpb.UserResponse{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}, nil
}

func (s *UserService) ValidateToken(ctx context.Context, req *userpb.ValidateTokenRequest) (*userpb.ValidateTokenResponse, error) {
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return &userpb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	return &userpb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}
