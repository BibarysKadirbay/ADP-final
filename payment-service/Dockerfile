FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o payment-service ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/payment-service .
EXPOSE 50053 9103
CMD ["./payment-service"]
