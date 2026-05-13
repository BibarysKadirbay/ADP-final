# Delivery Service

Production-ready Go microservice for courier assignment and delivery lifecycle management in a Food Delivery Platform.

## Architecture

The service follows Clean Architecture:

- `domain`: entities, repository contracts, status rules, domain errors
- `usecase`: courier assignment, delivery lifecycle, ETA, cache invalidation, event orchestration
- `infrastructure`: PostgreSQL, Redis, NATS JetStream, Restaurant Service gRPC client, metrics, tracing, logging
- `transport/grpc`: gRPC handlers and DTO mapping only
- `middleware`: request IDs, metrics, JWT and role hooks

Delivery Service owns its PostgreSQL database and never reads another service database. Restaurant data is accessed only through the Restaurant Service gRPC contract, and order events arrive through NATS JetStream.

## Features

- Courier registration and availability management
- Highest-rated balanced courier assignment with overload protection
- Delivery status flow: `pending`, `assigned`, `picked_up`, `on_the_way`, `delivered`, `cancelled`
- Delivery history and courier rating
- ETA estimation based on route distance, average courier speed, and delivery status
- Redis cache for active deliveries, available couriers, stats, and ETA
- NATS JetStream subscribers for `order.confirmed` and `order.cancelled`
- NATS publishers for `delivery.assigned`, `delivery.started`, `delivery.completed`, `delivery.cancelled`
- Prometheus metrics and OpenTelemetry tracing
- Zap structured logging and graceful shutdown

## Environment

See `.env`.

```env
APP_ENV=development
SERVICE_NAME=delivery-service
GRPC_PORT=50056
METRICS_PORT=9106
POSTGRES_DSN=postgres://delivery:delivery@postgres:5432/delivery_service?sslmode=disable
REDIS_ADDR=redis:6379
NATS_URL=nats://nats:4222
RESTAURANT_GRPC_ADDR=restaurant-service:50055
CACHE_TTL=5m
ETA_CACHE_TTL=2m
JWT_SECRET=change-me
OTEL_ENABLED=true
```

## Run With Docker

```bash
docker compose up -d --build
```

gRPC listens on `localhost:50056`, metrics on `localhost:9106`, PostgreSQL on `localhost:5436`, Redis on `localhost:6386`, and NATS on `localhost:4226`.

## Migrations

Install `golang-migrate`, then:

```bash
export POSTGRES_DSN='postgres://delivery:delivery@localhost:5436/delivery_service?sslmode=disable'
make migrate-up
make migrate-down
```

The first migration creates:

- `couriers`
- `deliveries`
- `delivery_status_history`
- `courier_ratings`

with indexes, constraints, foreign keys, and status validation.

## gRPC API

Implemented endpoints:

- `AssignDelivery`
- `GetDeliveryById`
- `UpdateDeliveryStatus`
- `GetDeliveriesByCourier`
- `GetDeliveriesByOrder`
- `ListAvailableCouriers`
- `GetDeliveryStats`
- `CalculateDeliveryETA`
- `GetDeliveryHistory`
- `RateCourier`
- `RegisterCourier`
- `UpdateCourierAvailability`
- `HealthCheck`

Register a courier:

```json
{
  "user_id": "11111111-1111-1111-1111-111111111111",
  "full_name": "Ayan Courier",
  "phone": "+77010000000",
  "vehicle_type": "bike"
}
```

Assign a delivery:

```json
{
  "order_id": "22222222-2222-2222-2222-222222222222",
  "restaurant_id": "33333333-3333-3333-3333-333333333333",
  "customer_id": "44444444-4444-4444-4444-444444444444",
  "pickup_address": "10 Abay Ave",
  "delivery_address": "20 Dostyk Ave",
  "route_distance_km": 4.5
}
```

## Microservice Communication

Delivery Service subscribes to:

- `order.confirmed`
- `order.cancelled`

Delivery Service publishes:

- `delivery.assigned`
- `delivery.started`
- `delivery.completed`
- `delivery.cancelled`

Restaurant Service integration is gRPC-only. If Restaurant Service is unavailable, Delivery Service logs the failure and continues with graceful degradation.

## Redis Usage

Redis stores short-lived cache entries for:

- active delivery lookups: `delivery:{id}`
- available courier lists: `couriers:available:*`
- courier stats: `delivery:stats:{courier_id}`
- ETA calculations: `eta:*`

All cache paths are best-effort. PostgreSQL remains the source of truth.

## Observability

Prometheus endpoint:

```text
http://localhost:9106/metrics
```

Metrics include:

- `delivery_grpc_requests_total`
- `delivery_grpc_request_duration_seconds`
- `delivery_db_query_duration_seconds`
- `delivery_cache_events_total`

OpenTelemetry tracing is initialized at startup and can be wired to an exporter later without changing business logic.

## Testing

```bash
go test ./...
go test ./... -coverpkg=./internal/...,./pkg/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Integration tests are intentionally skipped unless run in CI with Docker Compose dependencies and migrations applied.

## Proto Generation

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
make proto
```
