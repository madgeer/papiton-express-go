FROM golang:alpine AS builder
WORKDIR /workspace

COPY go.work ./
COPY common/ ./common/
COPY services/ ./services/

WORKDIR /workspace/services/order
RUN go build -o /app/order-service ./cmd/api

FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/order-service .
COPY services/order/.env .
EXPOSE 8081
CMD ["./order-service"]
