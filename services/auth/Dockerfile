FROM golang:alpine AS builder
WORKDIR /workspace

COPY go.work ./
COPY common/ ./common/
COPY services/ ./services/

WORKDIR /workspace/services/auth
RUN go build -o /app/auth-service ./cmd/api

FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/auth-service .
COPY services/auth/.env .
EXPOSE 8082
CMD ["./auth-service"]
