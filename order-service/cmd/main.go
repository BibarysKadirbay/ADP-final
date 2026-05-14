package main

import (
	"log"
	"net"

	"order-service/internal/database"
	"order-service/internal/repository"
	"order-service/internal/service"

	"github.com/nats-io/nats.go"

	orderpb "food-delivery-system/proto/orderpb"

	"google.golang.org/grpc"
)

func main() {

	db := database.Connect()

	repo := repository.NewOrderRepository(db)

	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}

	orderService := service.NewOrderService(repo, nc)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	orderpb.RegisterOrderServiceServer(
		grpcServer,
		orderService,
	)

	log.Println("Order Service running on 50052")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
