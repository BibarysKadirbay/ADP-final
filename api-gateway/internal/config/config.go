package config

import "os"

type Config struct {
	HTTPPort            string
	UserGRPCAddr        string
	OrderGRPCAddr       string
	RestaurantGRPCAddr  string
	JWTSecret           string
}

func Load() Config {
	return Config{
		HTTPPort:           getenv("HTTP_PORT", "8080"),
		UserGRPCAddr:       getenv("USER_GRPC_ADDR", "user-service:50052"),
		OrderGRPCAddr:      getenv("ORDER_GRPC_ADDR", "order-service:50051"),
		RestaurantGRPCAddr: getenv("RESTAURANT_GRPC_ADDR", "restaurant-service:50055"),
		JWTSecret:          getenv("JWT_SECRET", "super-secret-key"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
