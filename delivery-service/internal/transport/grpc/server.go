package grpc

import (
	"context"

	"github.com/aitu/food-delivery/delivery-service/internal/domain/entities"
	"github.com/aitu/food-delivery/delivery-service/internal/domain/services"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/grpc/deliverypb"
	"github.com/aitu/food-delivery/delivery-service/internal/usecase"
	"github.com/google/uuid"
)

type Server struct {
	deliverypb.UnimplementedDeliveryServiceServer
	uc          *usecase.DeliveryUsecase
	serviceName string
}

func NewServer(uc *usecase.DeliveryUsecase, serviceName string) *Server {
	return &Server{uc: uc, serviceName: serviceName}
}

func (s *Server) AssignDelivery(ctx context.Context, req *deliverypb.AssignDeliveryRequest) (*deliverypb.DeliveryResponse, error) {
	orderID, restaurantID, customerID, err := parseTriple(req.GetOrderId(), req.GetRestaurantId(), req.GetCustomerId())
	if err != nil {
		return nil, toStatus(err)
	}
	d, err := s.uc.AssignDelivery(ctx, &entities.Delivery{
		OrderID: orderID, RestaurantID: restaurantID, CustomerID: customerID, PickupAddress: req.GetPickupAddress(),
		DeliveryAddress: req.GetDeliveryAddress(), RouteDistanceKM: req.GetRouteDistanceKm(),
	})
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.DeliveryResponse{Delivery: toPBDelivery(d)}, nil
}

func (s *Server) GetDeliveryById(ctx context.Context, req *deliverypb.GetDeliveryByIdRequest) (*deliverypb.DeliveryResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	d, err := s.uc.GetDeliveryByID(ctx, id)
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.DeliveryResponse{Delivery: toPBDelivery(d)}, nil
}

func (s *Server) UpdateDeliveryStatus(ctx context.Context, req *deliverypb.UpdateDeliveryStatusRequest) (*deliverypb.DeliveryResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	d, err := s.uc.UpdateDeliveryStatus(ctx, id, entities.DeliveryStatus(req.GetStatus()))
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.DeliveryResponse{Delivery: toPBDelivery(d)}, nil
}

func (s *Server) GetDeliveriesByCourier(ctx context.Context, req *deliverypb.GetDeliveriesByCourierRequest) (*deliverypb.ListDeliveriesResponse, error) {
	courierID, err := uuid.Parse(req.GetCourierId())
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	filter := entities.DeliveryFilter{Limit: size(req.GetPagination()), Offset: (page(req.GetPagination()) - 1) * size(req.GetPagination())}
	if req.GetFilter() != nil {
		filter.Status = entities.DeliveryStatus(req.GetFilter().GetStatus())
		filter.SortBy = req.GetFilter().GetSortBy()
		filter.SortDirection = req.GetFilter().GetSortDirection()
	}
	items, meta, err := s.uc.GetDeliveriesByCourier(ctx, courierID, filter)
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.ListDeliveriesResponse{Deliveries: toPBDeliveries(items), Meta: toPBMeta(meta)}, nil
}

func (s *Server) GetDeliveriesByOrder(ctx context.Context, req *deliverypb.GetDeliveriesByOrderRequest) (*deliverypb.ListDeliveriesResponse, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	items, err := s.uc.GetDeliveriesByOrder(ctx, orderID)
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.ListDeliveriesResponse{Deliveries: toPBDeliveries(items), Meta: &deliverypb.PageMeta{Page: 1, PageSize: int32(len(items)), Total: int64(len(items))}}, nil
}

func (s *Server) ListAvailableCouriers(ctx context.Context, req *deliverypb.ListAvailableCouriersRequest) (*deliverypb.ListCouriersResponse, error) {
	filter := entities.CourierFilter{VehicleType: req.GetVehicleType(), Limit: size(req.GetPagination()), Offset: (page(req.GetPagination()) - 1) * size(req.GetPagination())}
	items, meta, err := s.uc.ListAvailableCouriers(ctx, filter)
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.ListCouriersResponse{Couriers: toPBCouriers(items), Meta: toPBMeta(meta)}, nil
}

func (s *Server) GetDeliveryStats(ctx context.Context, req *deliverypb.GetDeliveryStatsRequest) (*deliverypb.DeliveryStatsResponse, error) {
	courierID, err := uuid.Parse(req.GetCourierId())
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	stats, err := s.uc.GetDeliveryStats(ctx, courierID)
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.DeliveryStatsResponse{Stats: toPBStats(stats)}, nil
}

func (s *Server) CalculateDeliveryETA(ctx context.Context, req *deliverypb.CalculateDeliveryETARequest) (*deliverypb.DeliveryETAResponse, error) {
	var deliveryID uuid.UUID
	var err error
	if req.GetDeliveryId() != "" {
		deliveryID, err = uuid.Parse(req.GetDeliveryId())
		if err != nil {
			return nil, toStatus(services.ErrInvalidInput)
		}
	}
	eta, distance, err := s.uc.CalculateDeliveryETA(ctx, deliveryID, req.GetRouteDistanceKm(), entities.DeliveryStatus(req.GetStatus()))
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.DeliveryETAResponse{EtaMinutes: eta, RouteDistanceKm: distance}, nil
}

func (s *Server) GetDeliveryHistory(ctx context.Context, req *deliverypb.GetDeliveryHistoryRequest) (*deliverypb.DeliveryHistoryResponse, error) {
	id, err := uuid.Parse(req.GetDeliveryId())
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	items, err := s.uc.GetDeliveryHistory(ctx, id)
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.DeliveryHistoryResponse{History: toPBHistory(items)}, nil
}

func (s *Server) RateCourier(ctx context.Context, req *deliverypb.RateCourierRequest) (*deliverypb.CourierRatingResponse, error) {
	courierID, orderID, customerID, err := parseTriple(req.GetCourierId(), req.GetOrderId(), req.GetCustomerId())
	if err != nil {
		return nil, toStatus(err)
	}
	rating, err := s.uc.RateCourier(ctx, &entities.CourierRating{CourierID: courierID, OrderID: orderID, CustomerID: customerID, Rating: req.GetRating(), Comment: req.GetComment()})
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.CourierRatingResponse{Rating: toPBRating(rating)}, nil
}

func (s *Server) RegisterCourier(ctx context.Context, req *deliverypb.RegisterCourierRequest) (*deliverypb.CourierResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	c, err := s.uc.RegisterCourier(ctx, &entities.Courier{UserID: userID, FullName: req.GetFullName(), Phone: req.GetPhone(), VehicleType: req.GetVehicleType(), IsAvailable: true})
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.CourierResponse{Courier: toPBCourier(c)}, nil
}

func (s *Server) UpdateCourierAvailability(ctx context.Context, req *deliverypb.UpdateCourierAvailabilityRequest) (*deliverypb.CourierResponse, error) {
	id, err := uuid.Parse(req.GetCourierId())
	if err != nil {
		return nil, toStatus(services.ErrInvalidInput)
	}
	c, err := s.uc.UpdateCourierAvailability(ctx, id, req.GetIsAvailable())
	if err != nil {
		return nil, toStatus(err)
	}
	return &deliverypb.CourierResponse{Courier: toPBCourier(c)}, nil
}

func (s *Server) HealthCheck(context.Context, *deliverypb.HealthCheckRequest) (*deliverypb.HealthCheckResponse, error) {
	return &deliverypb.HealthCheckResponse{Status: "SERVING", Service: s.serviceName}, nil
}

func parseTriple(a, b, c string) (uuid.UUID, uuid.UUID, uuid.UUID, error) {
	left, err := uuid.Parse(a)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, services.ErrInvalidInput
	}
	mid, err := uuid.Parse(b)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, services.ErrInvalidInput
	}
	right, err := uuid.Parse(c)
	if err != nil {
		return uuid.Nil, uuid.Nil, uuid.Nil, services.ErrInvalidInput
	}
	return left, mid, right, nil
}
