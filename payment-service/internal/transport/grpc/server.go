package grpc

import (
	"context"

	"github.com/aitu/food-delivery/payment-service/internal/domain/entities"
	"github.com/aitu/food-delivery/payment-service/internal/infrastructure/grpc/paymentpb"
	"github.com/aitu/food-delivery/payment-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	paymentpb.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUsecase
}

func NewServer(uc *usecase.PaymentUsecase) *Server {
	return &Server{uc: uc}
}

func (s *Server) GetPayment(ctx context.Context, req *paymentpb.GetPaymentRequest) (*paymentpb.PaymentResponse, error) {
	p, err := s.uc.GetPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &paymentpb.PaymentResponse{Payment: toPB(p)}, nil
}

func (s *Server) ListPaymentsByOrder(ctx context.Context, req *paymentpb.ListPaymentsByOrderRequest) (*paymentpb.ListPaymentsResponse, error) {
	items, err := s.uc.ListByOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}
	resp := &paymentpb.ListPaymentsResponse{}
	for _, p := range items {
		item := p
		resp.Payments = append(resp.Payments, toPB(&item))
	}
	return resp, nil
}

func toPB(p *entities.Payment) *paymentpb.Payment {
	return &paymentpb.Payment{
		Id: p.ID, OrderId: p.OrderID, UserId: p.UserID,
		Amount: p.Amount, Status: p.Status, Method: p.Method,
	}
}
