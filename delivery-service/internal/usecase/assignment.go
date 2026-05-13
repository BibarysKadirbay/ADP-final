package usecase

import (
	"context"
	"sort"

	"github.com/aitu/food-delivery/delivery-service/internal/domain/entities"
	"github.com/aitu/food-delivery/delivery-service/internal/domain/repositories"
	"github.com/aitu/food-delivery/delivery-service/internal/domain/services"
)

type AssignmentStrategy interface {
	SelectCourier(ctx context.Context, couriers []entities.Courier, deliveries repositories.DeliveryRepository) (*entities.Courier, error)
}

type RatingBalancedAssignmentStrategy struct {
	MaxActiveDeliveries int64
}

func (s RatingBalancedAssignmentStrategy) SelectCourier(ctx context.Context, couriers []entities.Courier, deliveries repositories.DeliveryRepository) (*entities.Courier, error) {
	if s.MaxActiveDeliveries <= 0 {
		s.MaxActiveDeliveries = 3
	}
	sort.SliceStable(couriers, func(i, j int) bool {
		if couriers[i].Rating == couriers[j].Rating {
			return couriers[i].TotalDeliveries < couriers[j].TotalDeliveries
		}
		return couriers[i].Rating > couriers[j].Rating
	})
	for i := range couriers {
		active, err := deliveries.ActiveCountByCourier(ctx, couriers[i].ID)
		if err != nil {
			return nil, err
		}
		if active < s.MaxActiveDeliveries {
			return &couriers[i], nil
		}
	}
	return nil, services.ErrNoCourierAvailable
}

type ETAEstimator interface {
	Calculate(routeDistanceKM float64, status entities.DeliveryStatus) int32
}

type SimpleETAEstimator struct {
	AverageSpeedKMPH float64
}

func (e SimpleETAEstimator) Calculate(routeDistanceKM float64, status entities.DeliveryStatus) int32 {
	if status == entities.StatusDelivered || status == entities.StatusCancelled {
		return 0
	}
	speed := e.AverageSpeedKMPH
	if speed <= 0 {
		speed = 28
	}
	if routeDistanceKM <= 0 {
		return 10
	}
	minutes := int32((routeDistanceKM / speed) * 60)
	if minutes < 5 {
		minutes = 5
	}
	if status == entities.StatusPickedUp || status == entities.StatusOnTheWay {
		minutes = int32(float64(minutes) * 0.65)
		if minutes < 3 {
			minutes = 3
		}
	}
	return minutes
}
