package grpc

import (
	"github.com/aitu/food-delivery/delivery-service/internal/domain/entities"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/grpc/deliverypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toPBCourier(c *entities.Courier) *deliverypb.Courier {
	if c == nil {
		return nil
	}
	return &deliverypb.Courier{
		Id: c.ID.String(), UserId: c.UserID.String(), FullName: c.FullName, Phone: c.Phone, VehicleType: c.VehicleType,
		Rating: c.Rating, TotalDeliveries: c.TotalDeliveries, IsAvailable: c.IsAvailable,
		CreatedAt: timestamppb.New(c.CreatedAt), UpdatedAt: timestamppb.New(c.UpdatedAt),
	}
}

func toPBCouriers(items []entities.Courier) []*deliverypb.Courier {
	out := make([]*deliverypb.Courier, 0, len(items))
	for i := range items {
		out = append(out, toPBCourier(&items[i]))
	}
	return out
}

func toPBDelivery(d *entities.Delivery) *deliverypb.Delivery {
	if d == nil {
		return nil
	}
	pb := &deliverypb.Delivery{
		Id: d.ID.String(), OrderId: d.OrderID.String(), CourierId: d.CourierID.String(), RestaurantId: d.RestaurantID.String(), CustomerId: d.CustomerID.String(),
		Status: string(d.Status), PickupAddress: d.PickupAddress, DeliveryAddress: d.DeliveryAddress, EstimatedEtaMinutes: d.EstimatedETAMinutes,
		RouteDistanceKm: d.RouteDistanceKM, CreatedAt: timestamppb.New(d.CreatedAt), UpdatedAt: timestamppb.New(d.UpdatedAt),
	}
	if d.PickupTime != nil {
		pb.PickupTime = timestamppb.New(*d.PickupTime)
	}
	if d.DeliveredTime != nil {
		pb.DeliveredTime = timestamppb.New(*d.DeliveredTime)
	}
	return pb
}

func toPBDeliveries(items []entities.Delivery) []*deliverypb.Delivery {
	out := make([]*deliverypb.Delivery, 0, len(items))
	for i := range items {
		out = append(out, toPBDelivery(&items[i]))
	}
	return out
}

func toPBHistory(items []entities.DeliveryStatusHistory) []*deliverypb.DeliveryStatusHistory {
	out := make([]*deliverypb.DeliveryStatusHistory, 0, len(items))
	for i := range items {
		oldStatus := ""
		if items[i].OldStatus != nil {
			oldStatus = string(*items[i].OldStatus)
		}
		out = append(out, &deliverypb.DeliveryStatusHistory{
			Id: items[i].ID.String(), DeliveryId: items[i].DeliveryID.String(), OldStatus: oldStatus,
			NewStatus: string(items[i].NewStatus), ChangedAt: timestamppb.New(items[i].ChangedAt),
		})
	}
	return out
}

func toPBRating(r *entities.CourierRating) *deliverypb.CourierRating {
	if r == nil {
		return nil
	}
	return &deliverypb.CourierRating{
		Id: r.ID.String(), CourierId: r.CourierID.String(), OrderId: r.OrderID.String(), CustomerId: r.CustomerID.String(),
		Rating: r.Rating, Comment: r.Comment, CreatedAt: timestamppb.New(r.CreatedAt),
	}
}

func toPBStats(s *entities.DeliveryStats) *deliverypb.DeliveryStats {
	if s == nil {
		return nil
	}
	return &deliverypb.DeliveryStats{
		TotalDeliveries: s.TotalDeliveries, ActiveDeliveries: s.ActiveDeliveries, CompletedDeliveries: s.CompletedDeliveries,
		CancelledDeliveries: s.CancelledDeliveries, AverageEtaMinutes: s.AverageETAMinutes, AverageDistanceKm: s.AverageDistanceKM,
	}
}

func toPBMeta(m entities.PageMeta) *deliverypb.PageMeta {
	return &deliverypb.PageMeta{Page: int32(m.Page), PageSize: int32(m.PageSize), Total: m.Total}
}

func page(p *deliverypb.Pagination) int {
	if p == nil || p.GetPage() <= 0 {
		return 1
	}
	return int(p.GetPage())
}

func size(p *deliverypb.Pagination) int {
	if p == nil || p.GetPageSize() <= 0 {
		return 20
	}
	return int(p.GetPageSize())
}
