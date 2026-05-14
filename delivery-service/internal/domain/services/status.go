package services

import "github.com/aitu/food-delivery/delivery-service/internal/domain/entities"

var allowedTransitions = map[entities.DeliveryStatus]map[entities.DeliveryStatus]bool{
	entities.StatusPending: {
		entities.StatusAssigned:  true,
		entities.StatusCancelled: true,
	},
	entities.StatusAssigned: {
		entities.StatusPickedUp:  true,
		entities.StatusCancelled: true,
	},
	entities.StatusPickedUp: {
		entities.StatusOnTheWay:  true,
		entities.StatusCancelled: true,
	},
	entities.StatusOnTheWay: {
		entities.StatusDelivered: true,
		entities.StatusCancelled: true,
	},
	entities.StatusDelivered: {},
	entities.StatusCancelled: {},
}

func ValidateStatusTransition(from, to entities.DeliveryStatus) error {
	if from == to {
		return nil
	}
	if allowedTransitions[from][to] {
		return nil
	}
	return ErrInvalidTransition
}

func IsTerminal(status entities.DeliveryStatus) bool {
	return status == entities.StatusDelivered || status == entities.StatusCancelled
}

func IsValidStatus(status entities.DeliveryStatus) bool {
	_, ok := allowedTransitions[status]
	return ok
}
