package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Role string

const (
	RoleCustomer        Role = "customer"
	RoleRestaurantOwner Role = "restaurant_owner"
	RoleAdmin           Role = "admin"
)

type Claims struct {
	UserID string
	Role   Role
}

type claimsKey struct{}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(Claims)
	return claims, ok
}

func JWTUnaryInterceptor(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		auth := first(md.Get("authorization"))
		if auth == "" {
			return handler(ctx, req)
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		parts := strings.Split(token, ":")
		if len(parts) != 2 {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization token")
		}
		ctx = context.WithValue(ctx, claimsKey{}, Claims{UserID: parts[0], Role: Role(parts[1])})
		return handler(ctx, req)
	}
}

func first(v []string) string {
	if len(v) == 0 {
		return ""
	}
	return v[0]
}
