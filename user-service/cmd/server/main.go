package main

import (
	"context"
	"log"
	"net"

	"user-service/internal/config"
	"user-service/internal/infrastructure/database"
	pb "user-service/internal/infrastructure/grpc/userpb"
	appmetrics "user-service/internal/infrastructure/metrics"
	"user-service/internal/middleware"
	grpcTransport "user-service/internal/transport/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresPool(
		ctx,
		cfg.PostgresDSN,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repo := database.NewUserRepository(db)

	m := appmetrics.New(cfg.ServiceName)

	lis, err := net.Listen(
		"tcp",
		":"+cfg.GRPCPort,
	)

	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDUnaryInterceptor(),
			middleware.MetricsUnaryInterceptor(m),
		),
	)

	pb.RegisterUserServiceServer(
		grpcServer,
		grpcTransport.NewUserServer(repo, nil),
	)

	reflection.Register(grpcServer)

	log.Println(
		"user-service gRPC started on :" + cfg.GRPCPort,
	)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
