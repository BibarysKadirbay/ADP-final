package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"order-service/internal/config"
	"order-service/internal/infrastructure/database"
	pb "order-service/internal/infrastructure/grpc/orderpb"
	appmetrics "order-service/internal/infrastructure/metrics"
	ordernats "order-service/internal/infrastructure/nats"
	"order-service/internal/infrastructure/redis"
	"order-service/internal/middleware"
	grpcTransport "order-service/internal/transport/grpc"
	"order-service/internal/usecase"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresPool(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if cache, err := redis.New(cfg.Redis.Addr); err != nil {
		log.Println("redis unavailable, continuing without cache:", err)
	} else {
		defer cache.Close()
	}

	var publisher usecase.EventPublisher
	if cfg.NATS.URL != "" {
		pub, err := ordernats.NewPublisher(cfg.NATS.URL)
		if err != nil {
			log.Println("nats publisher failed:", err)
		} else {
			publisher = pub
			defer pub.Close()
		}
	}

	repo := database.NewOrderRepository(db)
	uc := usecase.NewOrderUsecase(repo, publisher)
	m := appmetrics.New(cfg.ServiceName)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("order-service metrics started on :" + cfg.MetricsPort)
		if err := http.ListenAndServe(":"+cfg.MetricsPort, nil); err != nil {
			log.Println("metrics server error:", err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDUnaryInterceptor(),
			middleware.MetricsUnaryInterceptor(m),
		),
	)

	pb.RegisterOrderServiceServer(grpcServer, grpcTransport.NewOrderServer(uc))
	reflection.Register(grpcServer)

	log.Println("order-service gRPC started on :" + cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
