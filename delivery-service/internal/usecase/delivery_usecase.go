package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aitu/food-delivery/delivery-service/internal/domain/entities"
	"github.com/aitu/food-delivery/delivery-service/internal/domain/repositories"
	"github.com/aitu/food-delivery/delivery-service/internal/domain/services"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DeliveryUsecase struct {
	couriers    repositories.CourierRepository
	deliveries  repositories.DeliveryRepository
	ratings     repositories.RatingRepository
	cache       repositories.Cache
	publisher   repositories.EventPublisher
	restaurant  repositories.RestaurantClient
	assigner    AssignmentStrategy
	eta         ETAEstimator
	cacheTTL    time.Duration
	etaCacheTTL time.Duration
	log         *zap.Logger
}

func NewDeliveryUsecase(
	couriers repositories.CourierRepository,
	deliveries repositories.DeliveryRepository,
	ratings repositories.RatingRepository,
	cache repositories.Cache,
	publisher repositories.EventPublisher,
	restaurant repositories.RestaurantClient,
	assigner AssignmentStrategy,
	eta ETAEstimator,
	cacheTTL time.Duration,
	etaCacheTTL time.Duration,
	log *zap.Logger,
) *DeliveryUsecase {
	if assigner == nil {
		assigner = RatingBalancedAssignmentStrategy{MaxActiveDeliveries: 3}
	}
	if eta == nil {
		eta = SimpleETAEstimator{AverageSpeedKMPH: 28}
	}
	return &DeliveryUsecase{
		couriers: couriers, deliveries: deliveries, ratings: ratings, cache: cache, publisher: publisher,
		restaurant: restaurant, assigner: assigner, eta: eta, cacheTTL: cacheTTL, etaCacheTTL: etaCacheTTL, log: log,
	}
}

func (u *DeliveryUsecase) RegisterCourier(ctx context.Context, c *entities.Courier) (*entities.Courier, error) {
	if c.UserID == uuid.Nil || c.FullName == "" || c.Phone == "" || c.VehicleType == "" {
		return nil, services.ErrInvalidInput
	}
	if err := u.couriers.Create(ctx, c); err != nil {
		return nil, err
	}
	_ = u.cache.DeletePattern(ctx, "couriers:available:*")
	return c, nil
}

func (u *DeliveryUsecase) UpdateCourierAvailability(ctx context.Context, courierID uuid.UUID, available bool) (*entities.Courier, error) {
	c, err := u.couriers.UpdateAvailability(ctx, courierID, available)
	if err != nil {
		return nil, err
	}
	_ = u.cache.DeletePattern(ctx, "couriers:available:*")
	return c, nil
}

func (u *DeliveryUsecase) ListAvailableCouriers(ctx context.Context, filter entities.CourierFilter) ([]entities.Courier, entities.PageMeta, error) {
	normalizeCourierFilter(&filter)
	key := fmt.Sprintf("couriers:available:%s:%d:%d", filter.VehicleType, filter.Limit, filter.Offset)
	var cached struct {
		Items []entities.Courier
		Total int64
	}
	if ok, err := u.cache.Get(ctx, key, &cached); err == nil && ok {
		return cached.Items, metaFromFilter(filter.Limit, filter.Offset, cached.Total), nil
	}
	items, total, err := u.couriers.ListAvailable(ctx, filter)
	if err != nil {
		return nil, entities.PageMeta{}, err
	}
	_ = u.cache.Set(ctx, key, struct {
		Items []entities.Courier
		Total int64
	}{items, total}, u.cacheTTL)
	return items, metaFromFilter(filter.Limit, filter.Offset, total), nil
}

func (u *DeliveryUsecase) AssignDelivery(ctx context.Context, d *entities.Delivery) (*entities.Delivery, error) {
	if d.OrderID == uuid.Nil || d.RestaurantID == uuid.Nil || d.CustomerID == uuid.Nil || d.DeliveryAddress == "" {
		return nil, services.ErrInvalidInput
	}
	if d.PickupAddress == "" {
		if snapshot, err := u.restaurant.GetRestaurant(ctx, d.RestaurantID); err == nil && snapshot != nil {
			d.PickupAddress = snapshot.Address
		} else if u.log != nil {
			u.log.Warn("restaurant lookup degraded", zap.String("restaurant_id", d.RestaurantID.String()), zap.Error(err))
		}
	}
	if d.PickupAddress == "" {
		d.PickupAddress = "restaurant address unavailable"
	}
	couriers, _, err := u.couriers.ListAvailable(ctx, entities.CourierFilter{Limit: 50})
	if err != nil {
		return nil, err
	}
	courier, err := u.assigner.SelectCourier(ctx, couriers, u.deliveries)
	if err != nil {
		return nil, err
	}
	d.CourierID = courier.ID
	d.Status = entities.StatusAssigned
	d.EstimatedETAMinutes = u.eta.Calculate(d.RouteDistanceKM, d.Status)
	if err := u.deliveries.Assign(ctx, d); err != nil {
		return nil, err
	}
	_ = u.cache.DeletePattern(ctx, "couriers:available:*")
	_ = u.cache.Set(ctx, deliveryKey(d.ID), d, u.cacheTTL)
	_ = u.publisher.Publish(ctx, "delivery.assigned", deliveryEvent(d))
	return d, nil
}

func (u *DeliveryUsecase) GetDeliveryByID(ctx context.Context, id uuid.UUID) (*entities.Delivery, error) {
	var cached entities.Delivery
	if ok, err := u.cache.Get(ctx, deliveryKey(id), &cached); err == nil && ok {
		return &cached, nil
	}
	d, err := u.deliveries.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = u.cache.Set(ctx, deliveryKey(id), d, u.cacheTTL)
	return d, nil
}

func (u *DeliveryUsecase) UpdateDeliveryStatus(ctx context.Context, id uuid.UUID, status entities.DeliveryStatus) (*entities.Delivery, error) {
	if !services.IsValidStatus(status) {
		return nil, services.ErrInvalidInput
	}
	current, err := u.deliveries.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := services.ValidateStatusTransition(current.Status, status); err != nil {
		return nil, err
	}
	updated, err := u.deliveries.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}
	_ = u.cache.Delete(ctx, deliveryKey(id), statsKey(updated.CourierID))
	_ = u.cache.DeletePattern(ctx, "couriers:available:*")
	switch status {
	case entities.StatusPickedUp, entities.StatusOnTheWay:
		_ = u.publisher.Publish(ctx, "delivery.started", deliveryEvent(updated))
	case entities.StatusDelivered:
		_ = u.publisher.Publish(ctx, "delivery.completed", deliveryEvent(updated))
	case entities.StatusCancelled:
		_ = u.publisher.Publish(ctx, "delivery.cancelled", deliveryEvent(updated))
	}
	return updated, nil
}

func (u *DeliveryUsecase) GetDeliveriesByCourier(ctx context.Context, courierID uuid.UUID, filter entities.DeliveryFilter) ([]entities.Delivery, entities.PageMeta, error) {
	normalizeDeliveryFilter(&filter)
	items, total, err := u.deliveries.ListByCourier(ctx, courierID, filter)
	if err != nil {
		return nil, entities.PageMeta{}, err
	}
	return items, metaFromFilter(filter.Limit, filter.Offset, total), nil
}

func (u *DeliveryUsecase) GetDeliveriesByOrder(ctx context.Context, orderID uuid.UUID) ([]entities.Delivery, error) {
	return u.deliveries.ListByOrder(ctx, orderID)
}

func (u *DeliveryUsecase) GetDeliveryHistory(ctx context.Context, deliveryID uuid.UUID) ([]entities.DeliveryStatusHistory, error) {
	return u.deliveries.History(ctx, deliveryID)
}

func (u *DeliveryUsecase) GetDeliveryStats(ctx context.Context, courierID uuid.UUID) (*entities.DeliveryStats, error) {
	var cached entities.DeliveryStats
	if ok, err := u.cache.Get(ctx, statsKey(courierID), &cached); err == nil && ok {
		return &cached, nil
	}
	stats, err := u.deliveries.Stats(ctx, courierID)
	if err != nil {
		return nil, err
	}
	_ = u.cache.Set(ctx, statsKey(courierID), stats, u.cacheTTL)
	return stats, nil
}

func (u *DeliveryUsecase) CalculateDeliveryETA(ctx context.Context, deliveryID uuid.UUID, distance float64, status entities.DeliveryStatus) (int32, float64, error) {
	if !services.IsValidStatus(status) {
		return 0, 0, services.ErrInvalidInput
	}
	key := fmt.Sprintf("eta:%s:%0.2f:%s", deliveryID, distance, status)
	var cached struct {
		ETA      int32
		Distance float64
	}
	if ok, err := u.cache.Get(ctx, key, &cached); err == nil && ok {
		return cached.ETA, cached.Distance, nil
	}
	if distance <= 0 && deliveryID != uuid.Nil {
		d, err := u.deliveries.GetByID(ctx, deliveryID)
		if err != nil {
			return 0, 0, err
		}
		distance = d.RouteDistanceKM
		if status == "" {
			status = d.Status
		}
	}
	eta := u.eta.Calculate(distance, status)
	_ = u.cache.Set(ctx, key, struct {
		ETA      int32
		Distance float64
	}{eta, distance}, u.etaCacheTTL)
	return eta, distance, nil
}

func (u *DeliveryUsecase) RateCourier(ctx context.Context, rating *entities.CourierRating) (*entities.CourierRating, error) {
	if rating.CourierID == uuid.Nil || rating.OrderID == uuid.Nil || rating.CustomerID == uuid.Nil || rating.Rating < 1 || rating.Rating > 5 {
		return nil, services.ErrInvalidInput
	}
	if err := u.ratings.Create(ctx, rating); err != nil {
		return nil, err
	}
	if err := u.couriers.RecalculateRating(ctx, rating.CourierID); err != nil && !errors.Is(err, services.ErrNotFound) {
		return nil, err
	}
	_ = u.cache.Delete(ctx, statsKey(rating.CourierID))
	return rating, nil
}

func (u *DeliveryUsecase) HandleOrderConfirmed(ctx context.Context, event OrderConfirmedEvent) error {
	_, err := u.AssignDelivery(ctx, &entities.Delivery{
		OrderID: event.OrderID, RestaurantID: event.RestaurantID, CustomerID: event.CustomerID,
		PickupAddress: event.PickupAddress, DeliveryAddress: event.DeliveryAddress, RouteDistanceKM: event.RouteDistanceKM,
	})
	return err
}

type PaymentCompletedEvent struct {
	OrderID         string `json:"order_id"`
	RestaurantID    string `json:"restaurant_id"`
	UserID          string `json:"user_id"`
	DeliveryAddress string `json:"delivery_address"`
	Amount          int64  `json:"amount"`
}

func (u *DeliveryUsecase) HandlePaymentCompleted(ctx context.Context, event PaymentCompletedEvent) error {
	orderID, err := uuid.Parse(event.OrderID)
	if err != nil {
		return err
	}
	restaurantID, err := uuid.Parse(event.RestaurantID)
	if err != nil {
		restaurantID = uuid.Nil
	}
	customerID, err := uuid.Parse(event.UserID)
	if err != nil {
		return err
	}
	return u.HandleOrderConfirmed(ctx, OrderConfirmedEvent{
		OrderID: orderID, RestaurantID: restaurantID, CustomerID: customerID,
		PickupAddress: "Restaurant pickup", DeliveryAddress: event.DeliveryAddress, RouteDistanceKM: 3.5,
	})
}

func (u *DeliveryUsecase) HandleOrderCancelled(ctx context.Context, orderID uuid.UUID) error {
	d, err := u.deliveries.CancelByOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			return nil
		}
		return err
	}
	_ = u.cache.Delete(ctx, deliveryKey(d.ID), statsKey(d.CourierID))
	return u.publisher.Publish(ctx, "delivery.cancelled", deliveryEvent(d))
}

type OrderConfirmedEvent struct {
	OrderID         uuid.UUID `json:"order_id"`
	RestaurantID    uuid.UUID `json:"restaurant_id"`
	CustomerID      uuid.UUID `json:"customer_id"`
	PickupAddress   string    `json:"pickup_address"`
	DeliveryAddress string    `json:"delivery_address"`
	RouteDistanceKM float64   `json:"route_distance_km"`
}

func normalizeCourierFilter(f *entities.CourierFilter) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
}

func normalizeDeliveryFilter(f *entities.DeliveryFilter) {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	if f.SortBy == "" {
		f.SortBy = "created_at"
	}
	if f.SortDirection == "" {
		f.SortDirection = "desc"
	}
}

func metaFromFilter(limit, offset int, total int64) entities.PageMeta {
	page := 1
	if limit > 0 {
		page = offset/limit + 1
	}
	return entities.PageMeta{Page: page, PageSize: limit, Total: total}
}

func deliveryKey(id uuid.UUID) string {
	return "delivery:" + id.String()
}

func statsKey(courierID uuid.UUID) string {
	return "delivery:stats:" + courierID.String()
}

func deliveryEvent(d *entities.Delivery) map[string]any {
	return map[string]any{
		"id": d.ID, "order_id": d.OrderID, "courier_id": d.CourierID, "restaurant_id": d.RestaurantID,
		"customer_id": d.CustomerID, "status": d.Status, "eta_minutes": d.EstimatedETAMinutes,
	}
}
