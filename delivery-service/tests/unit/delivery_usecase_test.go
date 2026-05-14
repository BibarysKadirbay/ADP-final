package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aitu/food-delivery/delivery-service/internal/domain/entities"
	"github.com/aitu/food-delivery/delivery-service/internal/domain/services"
	"github.com/aitu/food-delivery/delivery-service/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStatusTransitions(t *testing.T) {
	require.NoError(t, services.ValidateStatusTransition(entities.StatusAssigned, entities.StatusPickedUp))
	require.NoError(t, services.ValidateStatusTransition(entities.StatusOnTheWay, entities.StatusDelivered))
	require.ErrorIs(t, services.ValidateStatusTransition(entities.StatusDelivered, entities.StatusAssigned), services.ErrInvalidTransition)
}

func TestETAEstimator(t *testing.T) {
	eta := usecase.SimpleETAEstimator{AverageSpeedKMPH: 30}
	require.Equal(t, int32(20), eta.Calculate(10, entities.StatusAssigned))
	require.Equal(t, int32(0), eta.Calculate(10, entities.StatusDelivered))
}

func TestAssignDeliverySelectsHighestRatedAvailableCourier(t *testing.T) {
	ctx := context.Background()
	highRatedID := uuid.New()
	courierRepo := &fakeCourierRepo{items: []entities.Courier{
		{ID: uuid.New(), UserID: uuid.New(), FullName: "Low", Phone: "1", VehicleType: "bike", Rating: 3.8, IsAvailable: true},
		{ID: highRatedID, UserID: uuid.New(), FullName: "High", Phone: "2", VehicleType: "bike", Rating: 4.9, IsAvailable: true},
	}}
	deliveryRepo := &fakeDeliveryRepo{}
	uc := usecase.NewDeliveryUsecase(
		courierRepo, deliveryRepo, &fakeRatingRepo{}, &fakeCache{}, &fakePublisher{}, fakeRestaurantClient{},
		usecase.RatingBalancedAssignmentStrategy{MaxActiveDeliveries: 2},
		usecase.SimpleETAEstimator{AverageSpeedKMPH: 30},
		time.Minute, time.Minute, nil,
	)
	delivery, err := uc.AssignDelivery(ctx, &entities.Delivery{
		OrderID: uuid.New(), RestaurantID: uuid.New(), CustomerID: uuid.New(),
		PickupAddress: "Cafe", DeliveryAddress: "Home", RouteDistanceKM: 6,
	})
	require.NoError(t, err)
	require.Equal(t, highRatedID, delivery.CourierID)
	require.Equal(t, entities.StatusAssigned, delivery.Status)
	require.NotEqual(t, uuid.Nil, delivery.ID)
	require.Len(t, deliveryRepo.items, 1)
}

func TestAssignDeliveryReturnsNoCourierAvailable(t *testing.T) {
	uc := usecase.NewDeliveryUsecase(
		&fakeCourierRepo{}, &fakeDeliveryRepo{}, &fakeRatingRepo{}, &fakeCache{}, &fakePublisher{}, fakeRestaurantClient{},
		nil, nil, time.Minute, time.Minute, nil,
	)
	_, err := uc.AssignDelivery(context.Background(), &entities.Delivery{
		OrderID: uuid.New(), RestaurantID: uuid.New(), CustomerID: uuid.New(),
		PickupAddress: "Cafe", DeliveryAddress: "Home", RouteDistanceKM: 3,
	})
	require.ErrorIs(t, err, services.ErrNoCourierAvailable)
}

type fakeCourierRepo struct {
	items []entities.Courier
}

func (f *fakeCourierRepo) Create(context.Context, *entities.Courier) error { return nil }
func (f *fakeCourierRepo) GetByID(context.Context, uuid.UUID) (*entities.Courier, error) {
	return nil, services.ErrNotFound
}
func (f *fakeCourierRepo) UpdateAvailability(context.Context, uuid.UUID, bool) (*entities.Courier, error) {
	return nil, services.ErrNotFound
}
func (f *fakeCourierRepo) ListAvailable(context.Context, entities.CourierFilter) ([]entities.Courier, int64, error) {
	return f.items, int64(len(f.items)), nil
}
func (f *fakeCourierRepo) RecalculateRating(context.Context, uuid.UUID) error { return nil }

type fakeDeliveryRepo struct {
	items []entities.Delivery
}

func (f *fakeDeliveryRepo) Assign(_ context.Context, d *entities.Delivery) error {
	d.ID = uuid.New()
	d.CreatedAt = time.Now()
	d.UpdatedAt = d.CreatedAt
	f.items = append(f.items, *d)
	return nil
}
func (f *fakeDeliveryRepo) GetByID(_ context.Context, id uuid.UUID) (*entities.Delivery, error) {
	for i := range f.items {
		if f.items[i].ID == id {
			return &f.items[i], nil
		}
	}
	return nil, services.ErrNotFound
}
func (f *fakeDeliveryRepo) GetByOrderID(context.Context, uuid.UUID) (*entities.Delivery, error) {
	return nil, services.ErrNotFound
}
func (f *fakeDeliveryRepo) UpdateStatus(context.Context, uuid.UUID, entities.DeliveryStatus) (*entities.Delivery, error) {
	return nil, services.ErrNotFound
}
func (f *fakeDeliveryRepo) ListByCourier(context.Context, uuid.UUID, entities.DeliveryFilter) ([]entities.Delivery, int64, error) {
	return f.items, int64(len(f.items)), nil
}
func (f *fakeDeliveryRepo) ListByOrder(context.Context, uuid.UUID) ([]entities.Delivery, error) {
	return f.items, nil
}
func (f *fakeDeliveryRepo) History(context.Context, uuid.UUID) ([]entities.DeliveryStatusHistory, error) {
	return nil, nil
}
func (f *fakeDeliveryRepo) Stats(context.Context, uuid.UUID) (*entities.DeliveryStats, error) {
	return &entities.DeliveryStats{}, nil
}
func (f *fakeDeliveryRepo) ActiveCountByCourier(context.Context, uuid.UUID) (int64, error) {
	return 0, nil
}
func (f *fakeDeliveryRepo) CancelByOrder(context.Context, uuid.UUID) (*entities.Delivery, error) {
	return nil, services.ErrNotFound
}

type fakeRatingRepo struct{}

func (fakeRatingRepo) Create(context.Context, *entities.CourierRating) error { return nil }

type fakeCache struct{}

func (fakeCache) Get(context.Context, string, any) (bool, error)        { return false, nil }
func (fakeCache) Set(context.Context, string, any, time.Duration) error { return nil }
func (fakeCache) Delete(context.Context, ...string) error               { return nil }
func (fakeCache) DeletePattern(context.Context, string) error           { return nil }

type fakePublisher struct {
	err error
}

func (f fakePublisher) Publish(context.Context, string, any) error { return f.err }

type fakeRestaurantClient struct{}

func (fakeRestaurantClient) GetRestaurant(context.Context, uuid.UUID) (*entities.RestaurantSnapshot, error) {
	return nil, errors.New("unavailable")
}
