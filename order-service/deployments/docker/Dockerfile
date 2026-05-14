FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o order-service ./cmd/server/main.go

EXPOSE 50051

CMD ["./order-service"]