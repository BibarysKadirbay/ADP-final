FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o user-service ./cmd/server/main.go

EXPOSE 50052

CMD ["./user-service"]