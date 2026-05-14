package service

import (
	"context"
	"encoding/json"
	"log"

	"order-service/internal/domain"
	"order-service/internal/repository"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	orderpb "food-delivery-system/proto/orderpb"
)

type OrderService struct {
	orderpb.UnimplementedOrderServiceServer
	repo *repository.OrderRepository
	nc   *nats.Conn
}

func NewOrderService(
	repo *repository.OrderRepository,
	nc *nats.Conn,
) *OrderService {

	return &OrderService{
		repo: repo,
		nc:   nc,
	}
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	req *orderpb.CreateOrderRequest,
) (*orderpb.OrderResponse, error) {

	order := domain.Order{
		ID:           uuid.New().String(),
		UserID:       req.UserId,
		RestaurantID: req.RestaurantId,
		TotalPrice:   req.TotalPrice,
		Status:       "PENDING",
	}

	err := s.repo.Create(order)
	if err != nil {
		return nil, err
	}

	event := map[string]interface{}{
		"order_id": order.ID,
		"amount":   order.TotalPrice,
	}

	data, _ := json.Marshal(event)

	err = s.nc.Publish("order.created", data)
	if err != nil {
		log.Println(err)
	}

	return &orderpb.OrderResponse{
		OrderId:      order.ID,
		UserId:       order.UserID,
		RestaurantId: order.RestaurantID,
		TotalPrice:   order.TotalPrice,
		Status:       order.Status,
	}, nil
}

func (s *OrderService) GetOrder(
	ctx context.Context,
	req *orderpb.GetOrderRequest,
) (*orderpb.OrderResponse, error) {

	order, err := s.repo.GetByID(req.OrderId)
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{
		OrderId:      order.ID,
		UserId:       order.UserID,
		RestaurantId: order.RestaurantID,
		TotalPrice:   order.TotalPrice,
		Status:       order.Status,
	}, nil
}

func (s *OrderService) UpdateOrderStatus(
	ctx context.Context,
	req *orderpb.UpdateOrderStatusRequest,
) (*orderpb.OrderResponse, error) {

	err := s.repo.UpdateStatus(req.OrderId, req.Status)
	if err != nil {
		return nil, err
	}

	order, _ := s.repo.GetByID(req.OrderId)

	return &orderpb.OrderResponse{
		OrderId:      order.ID,
		UserId:       order.UserID,
		RestaurantId: order.RestaurantID,
		TotalPrice:   order.TotalPrice,
		Status:       order.Status,
	}, nil
}

func (s *OrderService) CancelOrder(
	ctx context.Context,
	req *orderpb.CancelOrderRequest,
) (*orderpb.OrderResponse, error) {

	err := s.repo.UpdateStatus(req.OrderId, "CANCELLED")
	if err != nil {
		return nil, err
	}

	order, _ := s.repo.GetByID(req.OrderId)

	return &orderpb.OrderResponse{
		OrderId:      order.ID,
		UserId:       order.UserID,
		RestaurantId: order.RestaurantID,
		TotalPrice:   order.TotalPrice,
		Status:       order.Status,
	}, nil
}

func (s *OrderService) GetOrdersByUser(
	ctx context.Context,
	req *orderpb.GetOrdersByUserRequest,
) (*orderpb.OrdersResponse, error) {

	return &orderpb.OrdersResponse{}, nil
}
