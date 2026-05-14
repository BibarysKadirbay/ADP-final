FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o order-service ./cmd/main.go

EXPOSE 50052

CMD ["./order-service"]