# Papiton Express Development Planning

## Project Goal

Membangun sistem ekspedisi berbasis Microservice Architecture menggunakan Go, PostgreSQL, Apache Kafka, Docker, dan Kubernetes dengan penerapan Clean Architecture serta Event Driven Architecture.

---

# Phase 1 - Analysis & Design

Target:
Menyelesaikan seluruh dokumentasi dan desain sistem sebelum implementasi.

Checklist

* [x] Software Requirement Specification (SRS)
* [x] Software Architecture
* [x] Event Flow
* [x] API Specification
* [x] ERD
* [ ] Sequence Diagram
* [ ] Activity Diagram
* [ ] Deployment Diagram

Output

```
docs/
```

---

# Phase 2 - Development Environment

Target:
Mempersiapkan lingkungan pengembangan.

Checklist

* [x] Git Repository
* [x] Docker Compose
* [x] PostgreSQL
* [x] Apache Kafka
* [x] Go Workspace
* [x] Common Package
* [x] Environment Configuration

Output

```
docker-compose.yml
go.work
init-db.sql
common/
```

---

# Phase 3 - Shared Components

Target:
Membuat komponen yang digunakan seluruh service.

Checklist

* [x] Configuration Loader (godotenv + .env)
* [x] Database Package (GORM postgres connection helper)
* [ ] Logger (Structured slog/zap)
* [ ] Response Helper
* [ ] Error Handler
* [ ] JWT Helper
* [ ] Password Hash Helper
* [ ] Kafka Helper
* [ ] Middleware
* [ ] Validation

Output

```
common/
```

---

# Phase 4 - Auth Service

Status

🟢 Completed (Milestone 1 - Standalone REST API)

Target

Authentication dan Authorization.

Checklist

## Database

* [x] Migration
* [x] Models

## Backend

* [x] Repository
* [x] Service
* [x] Handler
* [x] Routes
* [x] JWT
* [x] Refresh Token
* [x] Password Hash (BCrypt)
* [x] Swagger
* [ ] Unit Test
* [ ] Dockerfile

---

# Phase 5 - Order Service

Status

🟢 Completed (Milestone 1 - Standalone REST API)

Target

Order Management.

Checklist

## Database

* [x] Migration
* [x] Models
* [x] Seeder

## Backend

* [x] Repository
* [x] Business Service (Tariff & Price Calculation)
* [x] CRUD API
* [x] Address API
* [x] Tariff API
* [x] Validation
* [x] Error Handling
* [x] Swagger
* [x] Unit Test (services_test.go)
* [x] Integration Test (api_test.go)
* [ ] Dockerfile

---

# Phase 6 - Payment Service

Status

🟢 Completed (Milestone 1 - Standalone REST API)

Target

Payment Management.

Checklist

## Database

* [x] Migration
* [x] Models

## Backend

* [x] Repository
* [x] Service (Invoice & Webhook)
* [x] Payment API
* [x] Payment Verification
* [x] Swagger
* [x] Unit Test (services_test.go)
* [x] Integration Test (api_test.go)
* [ ] Dockerfile

---

# Phase 7 - Shipping Service

Status

🟢 Completed (Milestone 1 - Standalone REST API)

Checklist

* [x] Migration
* [x] Models
* [x] Courier CRUD
* [x] Shipment CRUD
* [x] Assignment Logic
* [x] Swagger
* [x] Unit Test (shipping_service_test.go & api_test.go)
* [ ] Dockerfile

---

# Phase 8 - Warehouse Service

Status

🟢 Completed (Milestone 1 - Standalone REST API)

Checklist

* [x] Migration
* [x] Models
* [x] Warehouse CRUD
* [x] Package Movement (Transit Routing Logic)
* [x] Swagger
* [x] Unit Test (warehouse_service_test.go & api_test.go)
* [ ] Dockerfile

---

# Phase 9 - Tracking Service

Status

🟢 Completed (Milestone 1 - Standalone REST API)

Checklist

* [x] Migration
* [x] Models
* [x] Tracking CRUD
* [x] Tracking Timeline (Timeline Events Preloaded)
* [x] Swagger
* [x] Unit Test (tracking_service_test.go & api_test.go)
* [ ] Dockerfile

---

# Phase 10 - Notification Service

Status

⚪ Not Started

Checklist

* [ ] Migration
* [ ] Models
* [ ] Notification CRUD
* [ ] Email Sender
* [ ] Swagger
* [ ] Unit Test
* [ ] Dockerfile

---

# Phase 11 - Event Driven Architecture (Milestone 2)

Target

Menghubungkan seluruh service menggunakan Apache Kafka secara asinkron.

Checklist

* [ ] Kafka Topic Design
* [ ] Outbox Pattern (Daemon Worker)
* [ ] Kafka Producer
* [ ] Kafka Consumer
* [ ] Retry Mechanism
* [ ] Dead Letter Queue (DLQ)

---

# Phase 12 - API Gateway (Milestone 3)

Checklist

* [ ] Gateway Service
* [ ] JWT Middleware Validation
* [ ] Rate Limiter
* [ ] Reverse Proxy
* [ ] CORS

---

# Phase 13 - E2E Testing (Milestone 3)

Checklist

* [ ] End-to-End Test (E2E)
* [ ] Load Test

---

# Phase 14 - DevOps & Deployment (Milestone 3)

Checklist

* [ ] Dockerfile Semua Service
* [ ] Docker Compose Production
* [ ] CI/CD Pipeline (Jenkins/GitHub Actions)
* [ ] Kubernetes Manifests (AKS Deployment)

---

# Final Deliverables

## Documentation

* [x] SRS
* [x] Architecture
* [x] Event Flow
* [x] API Specification
* [x] ERD
* [x] README

## Source Code

* [x] Auth Service (REST API)
* [x] Order Service (REST API)
* [x] Payment Service (REST API)
* [x] Shipping Service (REST API)
* [x] Warehouse Service (REST API)
* [x] Tracking Service (REST API)
* [ ] Notification Service

## Infrastructure

* [x] Docker Compose (Postgres, Kafka, Kafka UI)
* [ ] Kubernetes
* [ ] CI/CD Pipeline

## Testing

* [x] Unit Test (Order, Payment)
* [x] Integration Test (Order, Payment)
* [ ] End-to-End Test

---

# Milestones

| Milestone | Target | Status |
| --------- | ------ | ------ |
| **Milestone 1** | **Standalone REST APIs (No Kafka)** | 🟡 In Progress |
| | - Analisis & Dokumentasi | ✓ Complete |
| | - Lingkungan Pengembangan & Go Workspace | ✓ Complete |
| | - Shared database connection (Modul `common`) | ✓ Complete |
| | - Auth Service (REST API CRUD) | ✓ Complete |
| | - Order Service (REST API CRUD + Test) | ✓ Complete |
| | - Payment Service (REST API CRUD + Test) | ✓ Complete |
| | - Shipping Service (REST API CRUD + Test) | ✓ Complete |
| | - Warehouse Service (REST API CRUD + Test) | ✓ Complete |
| | - Tracking Service (REST API CRUD + Test) | ✓ Complete |
| | - Notification Service (REST API CRUD) | ⚪ Not Started |
| **Milestone 2** | **Event-Driven Integration (Kafka)** | ⚪ Not Started |
| | - Outbox Worker & Kafka Producer | ⚪ Not Started |
| | - Kafka Consumers (Order, Tracking, Notification) | ⚪ Not Started |
| **Milestone 3** | **Gateway, Testing & DevOps Deployment** | ⚪ Not Started |
| | - API Gateway & JWT Security | ⚪ Not Started |
| | - E2E Integration Testing | ⚪ Not Started |
| | - Dockerization, CI/CD, Kubernetes (AKS) | ⚪ Not Started |
}
