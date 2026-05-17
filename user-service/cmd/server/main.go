package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"user-service/internal/config"
	"user-service/internal/infrastructure/database"
	"user-service/internal/infrastructure/redis"
	pb "user-service/internal/infrastructure/grpc/userpb"
	appmetrics "user-service/internal/infrastructure/metrics"
	"user-service/internal/middleware"
	grpcTransport "user-service/internal/transport/grpc"
	"user-service/internal/usecase"

	"github.com/nats-io/nats.go"
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

	var nc *nats.Conn
	if cfg.NATS.URL != "" {
		conn, err := nats.Connect(cfg.NATS.URL)
		if err != nil {
			log.Println("nats connect failed; events will not be published:", err)
		} else {
			nc = conn
			defer nc.Close()
		}
	}

	repo := database.NewUserRepository(db)
	uc := usecase.NewUserUsecase(repo, cfg.JWTSecret)
	m := appmetrics.New(cfg.ServiceName)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("user-service metrics started on :" + cfg.MetricsPort)
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

	pb.RegisterUserServiceServer(grpcServer, grpcTransport.NewUserServer(uc, nc))
	reflection.Register(grpcServer)

	log.Println("user-service gRPC started on :" + cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
