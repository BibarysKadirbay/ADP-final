package config

import "os"

type Config struct {
	HTTPPort           string
	UserGRPCAddr       string
	OrderGRPCAddr      string
	RestaurantGRPCAddr string
	PaymentGRPCAddr    string
	DeliveryGRPCAddr   string
	JWTSecret          string
}

func Load() Config {
	return Config{
		HTTPPort:           getenv("HTTP_PORT", "8080"),
		UserGRPCAddr:       getenv("USER_GRPC_ADDR", "user-service:50052"),
		OrderGRPCAddr:      getenv("ORDER_GRPC_ADDR", "order-service:50051"),
		RestaurantGRPCAddr: getenv("RESTAURANT_GRPC_ADDR", "restaurant-service:50055"),
		PaymentGRPCAddr:    getenv("PAYMENT_GRPC_ADDR", "payment-service:50053"),
		DeliveryGRPCAddr:   getenv("DELIVERY_GRPC_ADDR", "delivery-service:50056"),
		JWTSecret:          getenv("JWT_SECRET", "super-secret-key"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
