FROM golang:alpine AS builder
WORKDIR /workspace

COPY go.work ./
COPY common/ ./common/
COPY services/ ./services/

WORKDIR /workspace/services/payment
RUN go build -o /app/payment-service ./cmd/api

FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/payment-service .
COPY services/payment/.env .
EXPOSE 8083
CMD ["./payment-service"]
