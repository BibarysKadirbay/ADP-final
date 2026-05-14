package middleware

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func RequestIDUnaryInterceptor() grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		requestID := uuid.NewString()

		ctx = context.WithValue(
			ctx,
			"request_id",
			requestID,
		)

		return handler(ctx, req)
	}
}
