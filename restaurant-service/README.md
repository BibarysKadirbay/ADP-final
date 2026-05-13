# Restaurant Service

Production-ready Go microservice for restaurant and menu management in a Food Delivery Platform.

## Architecture

The service follows Clean Architecture:

- `domain`: entities, repository contracts, and domain errors
- `usecase`: business orchestration, validation, ownership checks, cache invalidation, event publishing
- `infrastructure`: PostgreSQL, Redis, NATS JetStream, metrics, tracing, logging
- `transport/grpc`: gRPC handlers and DTO mapping
- `middleware`: request ID, metrics, JWT/role extension points

Handlers contain no business logic. Dependencies are injected in `cmd/server/main.go`, keeping the service ready for API Gateway, frontend, mobile, and other microservice integrations.

## Features

- Restaurants: create, get, update, delete, list, search, top rated
- Menu categories: create, update, delete, list
- Menu items: create, update, delete, list menu by restaurant, availability changes
- Redis cache-aside for restaurant details, lists, menus, and top-rated restaurants
- NATS JetStream JSON events for restaurant and menu changes
- PostgreSQL migrations with indexes, foreign keys, and cascade rules
- Prometheus metrics and OpenTelemetry tracing
- Zap structured logging and graceful shutdown
- JWT and role middleware hooks for `customer`, `restaurant_owner`, and `admin`

## Environment

See `.env`.

```env
APP_ENV=development
SERVICE_NAME=restaurant-service
GRPC_PORT=50055
METRICS_PORT=9105
POSTGRES_DSN=postgres://restaurant:restaurant@postgres:5432/restaurant_service?sslmode=disable
REDIS_ADDR=redis:6379
NATS_URL=nats://nats:4222
CACHE_TTL=5m
JWT_SECRET=change-me
OTEL_ENABLED=true
```

## Run With Docker

```bash
docker compose up -d --build
```

gRPC listens on `localhost:50055`, metrics on `localhost:9105`.

## Migrations

Install `golang-migrate`, then:

```bash
export POSTGRES_DSN='postgres://restaurant:restaurant@localhost:5435/restaurant_service?sslmode=disable'
make migrate-up
make migrate-down
```

## Proto Generation

The repository includes a buildable generated-style package so development can start immediately. For real generated files:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
make proto
```

## Testing

```bash
go test ./...
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Integration tests are intentionally skipped unless run in a CI stage with Docker Compose dependencies.

## gRPC Examples

Create a restaurant:

```json
{
  "owner_id": "11111111-1111-1111-1111-111111111111",
  "name": "Green Bowl",
  "description": "Fresh bowls and salads",
  "cuisine_type": "healthy",
  "address": "10 Abay Ave",
  "city": "Almaty",
  "image_url": "https://cdn.example.com/green-bowl.jpg",
  "is_open": true
}
```

List/search request:

```json
{
  "pagination": { "page": 1, "page_size": 20 },
  "filter": {
    "query": "green",
    "cuisine_type": "healthy",
    "city": "Almaty",
    "open_only": true,
    "sort_by": "rating",
    "sort_direction": "desc"
  }
}
```

## Events

Published as JSON to JetStream subjects:

- `restaurant.created`
- `restaurant.updated`
- `restaurant.deleted`
- `menu.item.created`
- `menu.item.updated`
- `menu.item.deleted`
- `menu.item.availability_changed`

## Metrics

Prometheus endpoint: `http://localhost:9105/metrics`

- `restaurant_grpc_requests_total`
- `restaurant_grpc_request_duration_seconds`
- `restaurant_db_query_duration_seconds`
- `restaurant_cache_events_total`

## API Gateway And Frontend Readiness

The proto responses use frontend-friendly models with stable IDs, timestamps, pagination metadata, filtering, sorting, and grouped menu categories. The service can be exposed through grpc-gateway or an API Gateway without changing business logic.
