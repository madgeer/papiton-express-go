# Entity Relationship Diagram (ERD): Papiton Express

Karena Papiton Express menggunakan arsitektur **Database-per-Service** (setiap microservice memiliki database PostgreSQL-nya sendiri), secara fisik **tidak ada Foreign Key (FK) lintas database**. 

Hubungan data antar-layanan terikat secara **logis** melalui *shared identifiers* (seperti `order_id` atau `user_id`) yang disinkronkan melalui event Kafka.

Berikut adalah visualisasi ERD lengkap untuk seluruh layanan Papiton Express:

```mermaid
erDiagram

%% ===========================
%% AUTH SERVICE
%% ===========================

users {
    uuid id PK
    varchar name
    varchar email UK
    varchar password
    varchar phone
    enum role
    timestamp created_at
    timestamp updated_at
}

refresh_tokens {
    uuid id PK
    uuid user_id FK
    varchar token
    timestamp expired_at
}

users ||--o{ refresh_tokens : owns

%% ===========================
%% ORDER SERVICE
%% ===========================

orders {
    uuid id PK
    uuid customer_id
    uuid service_type_id FK
    varchar tracking_number UK
    decimal weight
    decimal length
    decimal width
    decimal height
    decimal shipping_cost
    decimal insurance_fee
    decimal total_price
    enum status
    timestamp created_at
    timestamp updated_at
}

addresses {
    uuid id PK
    uuid order_id FK
    enum type
    varchar name
    varchar phone
    text address
    varchar city
    varchar province
    varchar postal_code
}

service_types {
    uuid id PK
    varchar name
    int estimated_day
}

tariffs {
    uuid id PK
    varchar origin_city
    varchar destination_city
    uuid service_type_id FK
    decimal price_per_kg
}

order_outbox {
    uuid id PK
    varchar aggregate_type
    uuid aggregate_id
    varchar event_type
    jsonb payload
    varchar status
    timestamp created_at
}

orders ||--o{ addresses : has
service_types ||--o{ orders : selected
service_types ||--o{ tariffs : owns
orders ||--o{ order_outbox : publish

%% ===========================
%% PAYMENT SERVICE
%% ===========================

payments {
    uuid id PK
    uuid order_id
    decimal amount
    varchar payment_method
    varchar provider_transaction_id
    enum status
    timestamp created_at
}

payment_logs {
    uuid id PK
    uuid payment_id FK
    text message
    timestamp created_at
}

payment_outbox {
    uuid id PK
    varchar aggregate_type
    uuid aggregate_id
    varchar event_type
    jsonb payload
    varchar status
    timestamp created_at
}

payments ||--o{ payment_logs : has
payments ||--o{ payment_outbox : publish

%% ===========================
%% SHIPPING SERVICE
%% ===========================

couriers {
    uuid id PK
    varchar name
    varchar phone
    enum status
    varchar vehicle_type
}

shipments {
    uuid id PK
    uuid order_id
    uuid courier_id FK
    enum status
    timestamp pickup_at
    timestamp delivered_at
    timestamp created_at
}

shipping_outbox {
    uuid id PK
    varchar aggregate_type
    uuid aggregate_id
    varchar event_type
    jsonb payload
    varchar status
    timestamp created_at
}

couriers ||--o{ shipments : assigned
shipments ||--o{ shipping_outbox : publish

%% ===========================
%% WAREHOUSE SERVICE
%% ===========================

warehouses {
    uuid id PK
    varchar name
    text location
}

warehouse_movements {
    uuid id PK
    uuid order_id
    uuid warehouse_id FK
    enum movement_type
    timestamp created_at
}

warehouse_outbox {
    uuid id PK
    varchar aggregate_type
    uuid aggregate_id
    varchar event_type
    jsonb payload
    varchar status
    timestamp created_at
}

warehouses ||--o{ warehouse_movements : contains
warehouse_movements ||--o{ warehouse_outbox : publish

%% ===========================
%% TRACKING SERVICE
%% ===========================

tracking_events {
    uuid id PK
    uuid order_id
    varchar event_type
    varchar status
    text location
    text notes
    timestamp created_at
}

%% ===========================
%% NOTIFICATION SERVICE
%% ===========================

notifications {
    uuid id PK
    uuid order_id
    varchar recipient
    varchar type
    enum status
    timestamp sent_at
    timestamp created_at
}

%% ===========================
%% LOGICAL RELATIONSHIP
%% ===========================

users ||--o{ orders : customer
orders ||--|| payments : paid
orders ||--|| shipments : shipped
orders ||--o{ warehouse_movements : warehouse
orders ||--o{ tracking_events : tracked
orders ||--o{ notifications : notify
```

---

## Deskripsi Relasi & Atribut Kunci

### 1. Pola Hubungan Outbox (`outbox_events` / `_outbox`)
Setiap DB operasional memiliki tabel outbox lokal. Struktur tabel ini seragam di seluruh service. Relasinya adalah **1-to-Many** antara entitas utama (seperti `orders` atau `payments`) dengan tabel outbox-nya. Ketika entitas utama dibuat atau statusnya berubah, baris baru ditambahkan ke tabel outbox dalam transaksi DB yang sama.

### 2. Logika Relasi Antar-Layanan
* **`orders.customer_id` ➔ `users.id`**: Relasi logis. Order Service tidak memvalidasi FK ke Auth DB secara langsung, tetapi API Gateway memastikan token JWT yang digunakan berisi `userID` yang valid milik customer.
* **`payments.order_id` ➔ `orders.id`**: Relasi logis. Saat Order dibuat, detail order dikirim melalui Kafka, lalu Payment Service merekam tagihan untuk `order_id` tersebut.
* **`shipments.courier_id` ➔ `couriers.id`**: Relasi fisik lokal (FK) di dalam Shipping Database.
* **`warehouse_movements.warehouse_id` ➔ `warehouses.id`**: Relasi fisik lokal (FK) di dalam Warehouse Database.
* **`tracking_events.order_id` ➔ `orders.id`**: Relasi logis. Tracking Service mengumpulkan seluruh event status dari berbagai service dan menampilkannya sebagai *historical timeline* untuk satu `order_id`.
