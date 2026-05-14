package grpc

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"order-service/internal/domain/entities"
	"order-service/internal/domain/services"
	"order-service/internal/infrastructure/database"

	orderpb "order-service/internal/infrastructure/grpc/orderpb"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	repo *database.OrderRepository
	nc   *nats.Conn
}

func NewOrderServer(
	repo *database.OrderRepository,
	nc *nats.Conn,
) *OrderServer {
	return &OrderServer{
		repo: repo,
		nc:   nc,
	}
}

func (s *OrderServer) CreateOrder(
	ctx context.Context,
	req *orderpb.CreateOrderRequest,
) (*orderpb.OrderResponse, error) {

	order := &entities.Order{
		ID:            uuid.NewString(),
		UserID:        req.UserId,
		RestaurantID:  req.RestaurantId,
		DeliveryID:    "",
		TotalPrice:    req.TotalPrice,
		Status:        services.OrderPending,
		PaymentStatus: services.PaymentPending,
		Address:       req.Address,
		Comment:       req.Comment,
		CreatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	event := map[string]interface{}{
		"order_id":      order.ID,
		"user_id":       order.UserID,
		"restaurant_id": order.RestaurantID,
		"amount":        order.TotalPrice,
		"status":        order.Status,
	}

	data, err := json.Marshal(event)
	if err == nil && s.nc != nil {
		if err := s.nc.Publish("order.created", data); err != nil {
			log.Println("failed to publish order.created:", err)
		}
	}

	return mapOrderToResponse(order), nil
}

func (s *OrderServer) GetOrder(
	ctx context.Context,
	req *orderpb.GetOrderRequest,
) (*orderpb.OrderResponse, error) {

	order, err := s.repo.GetByID(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return mapOrderToResponse(order), nil
}

func (s *OrderServer) UpdateOrderStatus(
	ctx context.Context,
	req *orderpb.UpdateOrderStatusRequest,
) (*orderpb.OrderResponse, error) {

	if err := s.repo.UpdateStatus(ctx, req.OrderId, req.Status); err != nil {
		return nil, err
	}

	order, err := s.repo.GetByID(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return mapOrderToResponse(order), nil
}

func (s *OrderServer) CancelOrder(
	ctx context.Context,
	req *orderpb.CancelOrderRequest,
) (*orderpb.OrderResponse, error) {

	if err := s.repo.UpdateStatus(ctx, req.OrderId, services.OrderCancelled); err != nil {
		return nil, err
	}

	order, err := s.repo.GetByID(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return mapOrderToResponse(order), nil
}

func (s *OrderServer) GetOrdersByUser(
	ctx context.Context,
	req *orderpb.GetOrdersByUserRequest,
) (*orderpb.OrdersResponse, error) {

	orders, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []*orderpb.OrderResponse

	for _, order := range orders {
		if order.UserID == req.UserId {
			o := order
			result = append(result, mapOrderToResponse(&o))
		}
	}

	return &orderpb.OrdersResponse{
		Orders: result,
	}, nil
}

func mapOrderToResponse(order *entities.Order) *orderpb.OrderResponse {
	return &orderpb.OrderResponse{
		OrderId:       order.ID,
		UserId:        order.UserID,
		RestaurantId:  order.RestaurantID,
		DeliveryId:    order.DeliveryID,
		TotalPrice:    order.TotalPrice,
		Status:        order.Status,
		PaymentStatus: order.PaymentStatus,
		Address:       order.Address,
		Comment:       order.Comment,
	}
}
