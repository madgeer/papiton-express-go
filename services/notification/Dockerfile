FROM golang:alpine AS builder
WORKDIR /workspace

COPY go.work ./
COPY common/ ./common/
COPY services/ ./services/

WORKDIR /workspace/services/notification
RUN go build -o /app/notification-service ./cmd/api

FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/notification-service .
COPY services/notification/.env .
EXPOSE 8087
CMD ["./notification-service"]
