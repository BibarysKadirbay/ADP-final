package grpc

import (
	"context"

	"order-service/internal/domain/entities"
	"order-service/internal/domain/services"
	orderpb "order-service/internal/infrastructure/grpc/orderpb"
	"order-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	uc *usecase.OrderUsecase
}

func NewOrderServer(uc *usecase.OrderUsecase) *OrderServer {
	return &OrderServer{uc: uc}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.OrderResponse, error) {
	order, err := s.uc.CreateOrder(ctx, req.UserId, req.RestaurantId, req.TotalPrice, req.Address, req.Comment, req.UserEmail)
	if err != nil {
		return nil, toStatus(err)
	}
	return mapOrder(order), nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.OrderResponse, error) {
	order, err := s.uc.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, toStatus(err)
	}
	return mapOrder(order), nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.OrderResponse, error) {
	if err := s.uc.UpdateOrderStatus(ctx, req.OrderId, req.Status); err != nil {
		return nil, toStatus(err)
	}
	order, err := s.uc.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, toStatus(err)
	}
	return mapOrder(order), nil
}

func (s *OrderServer) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.OrderResponse, error) {
	if err := s.uc.CancelOrder(ctx, req.OrderId); err != nil {
		return nil, toStatus(err)
	}
	order, err := s.uc.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, toStatus(err)
	}
	return mapOrder(order), nil
}

func (s *OrderServer) GetOrdersByUser(ctx context.Context, req *orderpb.GetOrdersByUserRequest) (*orderpb.OrdersResponse, error) {
	orders, err := s.uc.ListOrdersByUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	result := make([]*orderpb.OrderResponse, 0, len(orders))
	for _, o := range orders {
		order := o
		result = append(result, mapOrder(&order))
	}
	return &orderpb.OrdersResponse{Orders: result}, nil
}

func mapOrder(order *entities.Order) *orderpb.OrderResponse {
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

func toStatus(err error) error {
	if err == services.ErrOrderNotFound {
		return status.Error(codes.NotFound, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
