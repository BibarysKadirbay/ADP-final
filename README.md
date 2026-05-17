# Food Delivery Microservices (ADP Final)

Microservices food-delivery platform in Go with API Gateway, NATS, PostgreSQL, Redis, SMTP email, React frontend, and Prometheus/Grafana observability.

## Architecture

```
Frontend (React) → API Gateway (REST)
    ├── User Service (gRPC) — auth & profiles
    ├── Restaurant Service (gRPC) — restaurants & menus
    ├── Order Service (gRPC) — orders + transactions
    ├── Payment Service (gRPC + NATS worker)
    ├── Delivery Service (gRPC + NATS consumer)
    └── Email Service (NATS + SMTP)

NATS flow: order.created → payment.completed → delivery assign + email
```

## Quick start

```bash
# Optional: Gmail SMTP for email notifications
export SMTP_EMAIL=your@gmail.com
export SMTP_PASSWORD=your-app-password

docker compose up --build -d
```

| Service    | URL                    |
|-----------|------------------------|
| Frontend  | http://localhost:3000  |
| API       | http://localhost:8080  |
| Grafana   | http://localhost:3001  |
| Prometheus| http://localhost:9090  |

## API Gateway routes (by team member)

**Team Member 1** — `api-gateway/internal/handlers/auth_handler.go`
- `POST /register`, `POST /login`, `GET /users/:id`

**Team Member 2** — `api-gateway/internal/handlers/restaurant_handler.go`
- `POST /restaurants`, `GET /restaurants`, `GET /restaurants/:id`
- `POST /restaurants/:id/menu`, `GET /restaurants/:id/menu`
- `PATCH /restaurants/:id/menu/:menuItemId/availability`

**Team Member 3** — `api-gateway/internal/handlers/order_handler.go`
- `POST /orders`, `GET /orders/:id`, `GET /users/:id/orders`
- `PATCH /orders/:id/status`, `PATCH /orders/:id/cancel`

## Requirements checklist

| Requirement | Weight | Status | Where implemented |
|-------------|--------|--------|-------------------|
| **Clean Architecture** (domain, repository, usecase, transport, config) | 20% | ✅ | Each service under `internal/{domain,infrastructure,usecase,transport,config}` |
| **≥12 gRPC endpoints** | 20% | ✅ (42 total) | See gRPC table below |
| **NATS message queue** | 20% | ✅ | `order.created`, `payment.completed`, `order.cancelled`, delivery events |
| **PostgreSQL + migrations + transactions + Redis** | 20% | ✅ | `*/migrations/`, `CreateInTx` in order & payment repos, Redis in user/order/restaurant/delivery |
| **SMTP email** (Google/Microsoft) | 10% | ✅ | `email-service/internal/infrastructure/smtp/client.go` |
| **Unit + integration tests** | 10% | ✅ | `*/tests/unit/`, `*/tests/integration/` |
| **Bonus: Frontend** | +10% | ✅ | `frontend/` React + Vite |
| **Bonus: Grafana + Prometheus** | +10% | ✅ | `docker-compose.yml`, `observability/prometheus/` |

### gRPC endpoints (42)

| Service | RPCs | File |
|---------|------|------|
| User (6) | RegisterUser, LoginUser, CreateUser, GetUser, UpdateUserProfile, ListUsers | `user-service/proto/user.proto` |
| Order (5) | CreateOrder, GetOrder, UpdateOrderStatus, CancelOrder, GetOrdersByUser | `order-service/proto/order.proto` |
| Restaurant (17) | CRUD restaurants/categories/menu + HealthCheck | `restaurant-service/proto/restaurant.proto` |
| Delivery (12) | Assign, track, couriers, ETA, ratings + HealthCheck | `delivery-service/proto/delivery.proto` |
| Payment (2) | GetPayment, ListPaymentsByOrder | `payment-service/proto/payment.proto` |

### NATS subjects

| Subject | Publisher | Consumer |
|---------|-----------|----------|
| `user.created` | user-service | — |
| `order.created` | order-service | payment-service |
| `payment.completed` | payment-service | delivery-service, email-service |
| `order.cancelled` | order-service | delivery-service |
| `delivery.*` | delivery-service | — |

### Business features

| Feature | Service | Implementation |
|---------|---------|----------------|
| Register / Login / JWT | user-service | `usecase/user_usecase.go`, `infrastructure/auth/jwt.go` |
| Get / update profile | user-service | `UpdateUserProfile` gRPC |
| Restaurant CRUD & menu | restaurant-service | `usecase/restaurant_usecase.go` |
| Create order (transaction) | order-service | `order_repository.CreateInTx`, `usecase/order_usecase.go` |
| Payment on order.created | payment-service | `infrastructure/nats/nats.go` subscriber |
| Delivery on payment.completed | delivery-service | `HandlePaymentCompleted` |
| Email on payment.completed | email-service | `usecase/email_usecase.go` + SMTP |

## Project structure

```
├── api-gateway/          # REST → gRPC (3 handler files = 3 team members)
├── user-service/
├── order-service/
├── restaurant-service/
├── delivery-service/
├── payment-service/      # NEW
├── email-service/        # NEW
├── frontend/             # React SPA
├── observability/prometheus/
└── docker-compose.yml    # Full stack
```

## Development

```bash
# Per-service
cd user-service && go test ./...
cd order-service && go test ./...
cd frontend && npm run dev   # proxies /api → localhost:8080
```

## SMTP configuration

Set Gmail app password (or Microsoft SMTP) via environment:

```bash
SMTP_EMAIL=you@gmail.com
SMTP_PASSWORD=xxxx xxxx xxxx xxxx
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

Without credentials, email-service logs errors but the rest of the stack runs.
