package grpc

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"user-service/internal/domain/entities"
	"user-service/internal/infrastructure/database"

	userpb "user-service/internal/infrastructure/grpc/userpb"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	repo *database.UserRepository
	nc   *nats.Conn
}

func NewUserServer(
	repo *database.UserRepository,
	nc *nats.Conn,
) *UserServer {

	return &UserServer{
		repo: repo,
		nc:   nc,
	}
}

func (s *UserServer) CreateUser(
	ctx context.Context,
	req *userpb.CreateUserRequest,
) (*userpb.UserResponse, error) {

	user := &entities.User{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	event := map[string]interface{}{
		"user_id": user.ID,
		"name":    user.Name,
		"email":   user.Email,
	}

	data, err := json.Marshal(event)
	if err == nil && s.nc != nil {

		if err := s.nc.Publish(
			"user.created",
			data,
		); err != nil {

			log.Println(
				"failed to publish user.created:",
				err,
			)
		}
	}

	return mapUserToResponse(user), nil
}

func (s *UserServer) GetUser(
	ctx context.Context,
	req *userpb.GetUserRequest,
) (*userpb.UserResponse, error) {

	user, err := s.repo.GetByID(
		ctx,
		req.UserId,
	)

	if err != nil {
		return nil, err
	}

	return mapUserToResponse(user), nil
}

func (s *UserServer) ListUsers(
	ctx context.Context,
	req *userpb.Empty,
) (*userpb.UsersResponse, error) {

	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []*userpb.UserResponse

	for _, user := range users {

		u := user

		result = append(
			result,
			mapUserToResponse(&u),
		)
	}

	return &userpb.UsersResponse{
		Users: result,
	}, nil
}

func mapUserToResponse(
	user *entities.User,
) *userpb.UserResponse {

	return &userpb.UserResponse{
		UserId:  user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Address: user.Address,
	}
}
