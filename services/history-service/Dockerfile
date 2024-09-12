FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o history-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/history-service .

EXPOSE 8081

CMD ["./history-service"]
