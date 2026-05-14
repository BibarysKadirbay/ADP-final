package middleware

import (
	"context"

	"google.golang.org/grpc"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleCourier  Role = "courier"
	RoleCustomer Role = "customer"
)

func JWTUnaryInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Hook for gateway-issued JWT validation and role authorization.
		return handler(ctx, req)
	}
}
