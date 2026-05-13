APP=restaurant-service
MIGRATE=migrate -path migrations -database "$$POSTGRES_DSN"

.PHONY: tidy test cover build run docker-up docker-down migrate-up migrate-down proto

tidy:
	go mod tidy

test:
	go test ./...

cover:
	go test ./... -coverprofile=coverage.out
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
	protoc -I proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/restaurant.proto
