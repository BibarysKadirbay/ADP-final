package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/aitu/food-delivery/payment-service/internal/config"
	"github.com/aitu/food-delivery/payment-service/internal/infrastructure/database"
	"github.com/aitu/food-delivery/payment-service/internal/infrastructure/grpc/paymentpb"
	"github.com/aitu/food-delivery/payment-service/internal/infrastructure/nats"
	transportgrpc "github.com/aitu/food-delivery/payment-service/internal/transport/grpc"
	"github.com/aitu/food-delivery/payment-service/internal/usecase"
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

	publisher, err := nats.NewPublisher(cfg.NATSURL)
	if err != nil {
		log.Fatal("nats publisher:", err)
	}

	repo := database.NewPaymentRepository(db)
	uc := usecase.NewPaymentUsecase(repo, publisher)

	subscriber, err := nats.NewSubscriber(cfg.NATSURL, uc)
	if err != nil {
		log.Fatal("nats subscriber:", err)
	}
	defer subscriber.Close()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":"+cfg.MetricsPort, nil)
	}()

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	paymentpb.RegisterPaymentServiceServer(grpcServer, transportgrpc.NewServer(uc))
	reflection.Register(grpcServer)

	log.Println("payment-service started on", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
