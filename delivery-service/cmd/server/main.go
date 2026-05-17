package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aitu/food-delivery/delivery-service/internal/config"
	"github.com/aitu/food-delivery/delivery-service/internal/domain/repositories"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/database"
	infraGrpc "github.com/aitu/food-delivery/delivery-service/internal/infrastructure/grpc"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/grpc/deliverypb"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/logger"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/metrics"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/nats"
	cacheredis "github.com/aitu/food-delivery/delivery-service/internal/infrastructure/redis"
	"github.com/aitu/food-delivery/delivery-service/internal/infrastructure/telemetry"
	"github.com/aitu/food-delivery/delivery-service/internal/middleware"
	transportgrpc "github.com/aitu/food-delivery/delivery-service/internal/transport/grpc"
	"github.com/aitu/food-delivery/delivery-service/internal/usecase"
	otelgrpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log, err := logger.New(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	m := metrics.New()
	shutdownTracing, err := telemetry.InitTracing(ctx, cfg.ServiceName, cfg.OTELEnabled)
	if err != nil {
		log.Fatal("init tracing", zap.Error(err))
	}
	defer shutdownTracing(context.Background())

	db, err := database.NewPostgresPool(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatal("connect postgres", zap.Error(err))
	}
	defer db.Close()

	cache := cacheredis.New(cfg.Redis, m)
	defer cache.Close()
	if err := cache.Ping(ctx); err != nil {
		log.Warn("redis ping failed; cache fallback to postgres will be used", zap.Error(err))
	}

	publisher, err := nats.NewPublisher(cfg.NATS)
	if err != nil {
		log.Fatal("connect nats publisher", zap.Error(err))
	}
	defer publisher.Close()

	restaurantClient, err := infraGrpc.NewRestaurantClient(ctx, cfg.Restaurant)
	if err != nil {
		log.Warn("restaurant grpc unavailable at startup; graceful degradation enabled", zap.Error(err))
	}
	var restaurantDep repositories.RestaurantClient = infraGrpc.NoopRestaurantClient{}
	if restaurantClient != nil {
		defer restaurantClient.Close()
		restaurantDep = restaurantClient
	}

	courierRepo := database.NewCourierRepository(db, m)
	deliveryRepo := database.NewDeliveryRepository(db, m)
	ratingRepo := database.NewRatingRepository(db, m)
	uc := usecase.NewDeliveryUsecase(
		courierRepo, deliveryRepo, ratingRepo, cache, publisher,
		restaurantDep,
		usecase.RatingBalancedAssignmentStrategy{MaxActiveDeliveries: 3},
		usecase.SimpleETAEstimator{AverageSpeedKMPH: 28},
		cfg.CacheTTL,
		cfg.ETACacheTTL,
		log,
	)

	subscriber, err := nats.NewSubscriber(cfg.NATS, uc, log)
	if err != nil {
		log.Warn("nats subscriber unavailable; service will still accept grpc traffic", zap.Error(err))
	} else {
		defer subscriber.Close()
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDUnaryInterceptor(),
			middleware.JWTUnaryInterceptor(cfg.JWTSecret),
			middleware.MetricsUnaryInterceptor(m),
		),
	)
	deliverypb.RegisterDeliveryServiceServer(grpcServer, transportgrpc.NewServer(uc, cfg.ServiceName))
	reflection.Register(grpcServer)
	metricsServer := &http.Server{
		Addr:              ":" + cfg.MetricsPort,
		Handler:           m.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("metrics server failed", zap.Error(err))
		}
	}()

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal("listen grpc", zap.Error(err))
	}
	go func() {
		log.Info("delivery service started", zap.String("grpc_port", cfg.GRPCPort), zap.String("metrics_port", cfg.MetricsPort))
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("grpc server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("shutting down")
	grpcServer.GracefulStop()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = metricsServer.Shutdown(shutdownCtx)
}
