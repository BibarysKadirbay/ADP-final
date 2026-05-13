FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o user-service ./cmd/main.go

EXPOSE 50051

CMD ["./user-service"]