package grpc

import (
	"context"
	"encoding/json"
	"log"

	"user-service/internal/domain/entities"
	"user-service/internal/domain/services"
	userpb "user-service/internal/infrastructure/grpc/userpb"
	"user-service/internal/usecase"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	uc *usecase.UserUsecase
	nc *nats.Conn
}

func NewUserServer(uc *usecase.UserUsecase, nc *nats.Conn) *UserServer {
	return &UserServer{uc: uc, nc: nc}
}

func (s *UserServer) RegisterUser(ctx context.Context, req *userpb.RegisterUserRequest) (*userpb.AuthResponse, error) {
	user, token, err := s.uc.Register(ctx, req.Name, req.Email, req.Password, req.Phone, req.Address)
	if err != nil {
		return nil, toStatus(err)
	}
	s.publishCreated(user)
	return &userpb.AuthResponse{Token: token, User: mapUser(user)}, nil
}

func (s *UserServer) LoginUser(ctx context.Context, req *userpb.LoginUserRequest) (*userpb.AuthResponse, error) {
	user, token, err := s.uc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, toStatus(err)
	}
	return &userpb.AuthResponse{Token: token, User: mapUser(user)}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	user, err := s.uc.CreateUser(ctx, req.Name, req.Email, req.Phone, req.Address)
	if err != nil {
		return nil, toStatus(err)
	}
	s.publishCreated(user)
	return mapUser(user), nil
}

func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	user, err := s.uc.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, toStatus(err)
	}
	return mapUser(user), nil
}

func (s *UserServer) UpdateUserProfile(ctx context.Context, req *userpb.UpdateUserProfileRequest) (*userpb.UserResponse, error) {
	user, err := s.uc.UpdateProfile(ctx, req.UserId, req.Name, req.Phone, req.Address)
	if err != nil {
		return nil, toStatus(err)
	}
	return mapUser(user), nil
}

func (s *UserServer) ListUsers(ctx context.Context, _ *userpb.Empty) (*userpb.UsersResponse, error) {
	users, err := s.uc.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*userpb.UserResponse, 0, len(users))
	for _, u := range users {
		user := u
		result = append(result, mapUser(&user))
	}
	return &userpb.UsersResponse{Users: result}, nil
}

func (s *UserServer) publishCreated(user *entities.User) {
	if s.nc == nil || user == nil {
		return
	}
	event := map[string]string{"user_id": user.ID, "name": user.Name, "email": user.Email}
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	if err := s.nc.Publish("user.created", data); err != nil {
		log.Println("failed to publish user.created:", err)
	}
}

func mapUser(user *entities.User) *userpb.UserResponse {
	return &userpb.UserResponse{
		UserId:  user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Address: user.Address,
	}
}

func toStatus(err error) error {
	switch err {
	case services.ErrInvalidCredentials:
		return status.Error(codes.Unauthenticated, err.Error())
	case services.ErrEmailAlreadyExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case services.ErrUserNotFound:
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
