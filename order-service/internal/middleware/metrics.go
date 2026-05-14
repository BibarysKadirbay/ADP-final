package middleware

import (
	"context"
	"time"

	"order-service/internal/infrastructure/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MetricsUnaryInterceptor(m *metrics.Metrics) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		started := time.Now()

		resp, err := handler(ctx, req)

		code := status.Code(err).String()

		m.GRPCRequests.WithLabelValues(info.FullMethod, code).Inc()
		m.GRPCDuration.WithLabelValues(info.FullMethod).Observe(time.Since(started).Seconds())

		return resp, err
	}
}
