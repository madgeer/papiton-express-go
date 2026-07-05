# Deployment

## 1. Technology Stack
- Go
- Gin
- PostgreSQL
- Kafka
- Docker
- Kubernetes
- JWT

## 2. Deployment Plan
1. Setiap service dikemas sebagai container Docker.
2. Service di-deploy ke Kubernetes dengan deployment terpisah.
3. PostgreSQL dipasang per layanan atau per domain sesuai kebutuhan.
4. Kafka digunakan untuk event bus antar layanan.
5. Gateway menjadi entry point publik.

## 3. Observability
- Prometheus + Grafana untuk metrik.
- ELK Stack untuk log.
- Jaeger untuk tracing.

## 4. Reliability
- health check per service;
- autoscaling berdasarkan CPU dan request rate;
- backup database dilakukan berkala;
- rollback menggunakan deployment strategy yang aman.
