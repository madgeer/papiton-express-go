FROM golang:alpine AS builder
WORKDIR /workspace

COPY go.work ./
COPY common/ ./common/
COPY services/ ./services/

WORKDIR /workspace/services/warehouse
RUN go build -o /app/warehouse-service ./cmd/api

FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/warehouse-service .
COPY services/warehouse/.env .
EXPOSE 8085
CMD ["./warehouse-service"]
