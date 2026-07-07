FROM golang:alpine AS builder
WORKDIR /workspace

COPY go.work ./
COPY common/ ./common/
COPY services/ ./services/

WORKDIR /workspace/services/tracking
RUN go build -o /app/tracking-service ./cmd/api

FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/tracking-service .
COPY services/tracking/.env .
EXPOSE 8086
CMD ["./tracking-service"]
