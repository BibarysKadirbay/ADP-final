package grpc

import (
	"errors"

	"github.com/aitu/food-delivery/restaurant-service/internal/domain/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toStatus(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, services.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, services.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, services.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, services.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
