package unit

import (
	"context"
	"testing"
	"time"

	"github.com/aitu/food-delivery/delivery-service/internal/domain/entities"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/grpc/deliverypb"
	transportgrpc "github.com/aitu/food-delivery/delivery-service/internal/transport/grpc"
	"github.com/aitu/food-delivery/delivery-service/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	srv := transportgrpc.NewServer(nil, "delivery-service")
	resp, err := srv.HealthCheck(context.Background(), &deliverypb.HealthCheckRequest{})
	require.NoError(t, err)
	require.Equal(t, "SERVING", resp.GetStatus())
	require.Equal(t, "delivery-service", resp.GetService())
}

func TestRegisterCourierHandler(t *testing.T) {
	uc := usecase.NewDeliveryUsecase(&creatingCourierRepo{}, &fakeDeliveryRepo{}, &fakeRatingRepo{}, &fakeCache{}, &fakePublisher{}, fakeRestaurantClient{}, nil, nil, time.Minute, time.Minute, nil)
	srv := transportgrpc.NewServer(uc, "delivery-service")
	resp, err := srv.RegisterCourier(context.Background(), &deliverypb.RegisterCourierRequest{
		UserId: uuid.NewString(), FullName: "Ayan Courier", Phone: "+77010000000", VehicleType: "bike",
	})
	require.NoError(t, err)
	require.Equal(t, "Ayan Courier", resp.GetCourier().GetFullName())
	require.True(t, resp.GetCourier().GetIsAvailable())
}

type creatingCourierRepo struct {
	fakeCourierRepo
}

func (creatingCourierRepo) Create(_ context.Context, c *entities.Courier) error {
	c.ID = uuid.New()
	c.Rating = 5
	c.CreatedAt = time.Now()
	c.UpdatedAt = c.CreatedAt
	return nil
}
