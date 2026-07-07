FROM golang:alpine AS builder
WORKDIR /workspace

COPY go.work ./
COPY common/ ./common/
COPY services/ ./services/

WORKDIR /workspace/services/shipping
RUN go build -o /app/shipping-service ./cmd/api

FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/shipping-service .
COPY services/shipping/.env .
EXPOSE 8084
CMD ["./shipping-service"]
