FROM golang:1.22-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/restaurant-service ./cmd/server

FROM alpine:3.20

WORKDIR /app
RUN adduser -D -H appuser && apk add --no-cache ca-certificates
COPY --from=builder /out/restaurant-service /app/restaurant-service
COPY .env /app/.env
USER appuser
EXPOSE 50055 9105
ENTRYPOINT ["/app/restaurant-service"]
