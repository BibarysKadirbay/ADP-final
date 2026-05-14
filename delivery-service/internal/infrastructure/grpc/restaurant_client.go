package grpc

import (
	"context"
	"time"

	"github.com/aitu/food-delivery/delivery-service/internal/config"
	"github.com/aitu/food-delivery/delivery-service/internal/domain/entities"
	"github.com/aitu/food-delivery/delivery-service/internal/domain/services"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/grpc/restaurantpb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RestaurantClient struct {
	conn    *grpc.ClientConn
	client  restaurantpb.RestaurantServiceClient
	timeout time.Duration
	retries int
}

func NewRestaurantClient(ctx context.Context, cfg config.RestaurantGRPCConfig) (*RestaurantClient, error) {
	conn, err := grpc.DialContext(ctx, cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 2 * time.Second
	}
	if cfg.Retries < 0 {
		cfg.Retries = 0
	}
	return &RestaurantClient{conn: conn, client: restaurantpb.NewRestaurantServiceClient(conn), timeout: cfg.Timeout, retries: cfg.Retries}, nil
}

func (c *RestaurantClient) GetRestaurant(ctx context.Context, id uuid.UUID) (*entities.RestaurantSnapshot, error) {
	if c == nil || c.client == nil {
		return nil, services.ErrRestaurantUnavailable
	}
	var last error
	for attempt := 0; attempt <= c.retries; attempt++ {
		callCtx, cancel := context.WithTimeout(ctx, c.timeout)
		resp, err := c.client.GetRestaurantById(callCtx, &restaurantpb.GetRestaurantByIdRequest{Id: id.String()})
		cancel()
		if err == nil {
			if resp.GetRestaurant() == nil {
				return nil, services.ErrNotFound
			}
			return &entities.RestaurantSnapshot{ID: id, Name: resp.GetRestaurant().GetName(), Address: resp.GetRestaurant().GetAddress()}, nil
		}
		last = err
	}
	if last != nil {
		return nil, last
	}
	return nil, services.ErrRestaurantUnavailable
}

func (c *RestaurantClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

type NoopRestaurantClient struct{}

func (NoopRestaurantClient) GetRestaurant(context.Context, uuid.UUID) (*entities.RestaurantSnapshot, error) {
	return nil, services.ErrRestaurantUnavailable
}
