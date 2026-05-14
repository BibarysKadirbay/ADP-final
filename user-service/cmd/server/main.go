package main

import (
	"log"
	"net"

	"user-service/internal/database"
	"user-service/internal/repository"
	"user-service/internal/service"

	userpb "food-delivery-system/proto/userpb"

	"google.golang.org/grpc"
)

func main() {
	db := database.Connect()
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, userService)

	log.Println("User Service is running on port 50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
