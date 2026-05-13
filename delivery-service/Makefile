APP=delivery-service
MIGRATE=migrate -path migrations -database "$$POSTGRES_DSN"

.PHONY: tidy test cover build run docker-up docker-down migrate-up migrate-down proto

tidy:
	go mod tidy

test:
	go test ./...

cover:
	go test ./... -coverpkg=./internal/...,./pkg/... -coverprofile=coverage.out
	go tool cover -func=coverage.out

build:
	go build -o bin/$(APP) ./cmd/server

run:
	go run ./cmd/server

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

proto:
	protoc -I proto --go_out=internal/infrastructure/grpc/deliverypb --go_opt=paths=source_relative --go-grpc_out=internal/infrastructure/grpc/deliverypb --go-grpc_opt=paths=source_relative proto/delivery.proto
	protoc -I proto --go_out=internal/infrastructure/grpc/restaurantpb --go_opt=paths=source_relative --go-grpc_out=internal/infrastructure/grpc/restaurantpb --go-grpc_opt=paths=source_relative proto/restaurant.proto
