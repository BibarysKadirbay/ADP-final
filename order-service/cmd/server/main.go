package main

import (
	"context"
	"log"
	"net"

	"order-service/internal/config"
	"order-service/internal/infrastructure/database"
	pb "order-service/internal/infrastructure/grpc/orderpb"
	appmetrics "order-service/internal/infrastructure/metrics"
	"order-service/internal/middleware"
	grpcTransport "order-service/internal/transport/grpc"

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

	repo := database.NewOrderRepository(db)

	m := appmetrics.New(cfg.ServiceName)

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

	pb.RegisterOrderServiceServer(
		grpcServer,
		grpcTransport.NewOrderServer(repo, nil),
	)
	reflection.Register(grpcServer)
	log.Println("order-service gRPC started on :" + cfg.GRPCPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
