package unit

import (
	"context"
	"testing"
	"time"

	"github.com/aitu/food-delivery/restaurant-service/internal/infrastructure/grpc/restaurantpb"
	restaurantgrpc "github.com/aitu/food-delivery/restaurant-service/internal/transport/grpc"
	"github.com/aitu/food-delivery/restaurant-service/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	repos := newFakes()
	uc := usecase.NewRestaurantUsecase(repos.restaurants, repos.categories, repos.menu, repos.cache, repos.events, time.Minute)
	server := restaurantgrpc.NewServer(uc, "restaurant-service")

	resp, err := server.HealthCheck(context.Background(), &restaurantpb.HealthCheckRequest{})

	require.NoError(t, err)
	require.Equal(t, "SERVING", resp.Status)
}

func TestCreateRestaurantInvalidOwnerID(t *testing.T) {
	repos := newFakes()
	uc := usecase.NewRestaurantUsecase(repos.restaurants, repos.categories, repos.menu, repos.cache, repos.events, time.Minute)
	server := restaurantgrpc.NewServer(uc, "restaurant-service")

	_, err := server.CreateRestaurant(context.Background(), &restaurantpb.CreateRestaurantRequest{
		OwnerId: "bad", Name: "A", CuisineType: "thai", City: "Almaty",
	})

	require.Error(t, err)
}

func TestCreateRestaurantHandler(t *testing.T) {
	repos := newFakes()
	uc := usecase.NewRestaurantUsecase(repos.restaurants, repos.categories, repos.menu, repos.cache, repos.events, time.Minute)
	server := restaurantgrpc.NewServer(uc, "restaurant-service")

	resp, err := server.CreateRestaurant(context.Background(), &restaurantpb.CreateRestaurantRequest{
		OwnerId: uuid.NewString(), Name: "Green Bowl", CuisineType: "healthy", City: "Astana", IsOpen: true,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp.Restaurant.Id)
}
