# Papiton Express - Microservices Backend (Go)

Papiton Express adalah sistem backend logistik ekspedisi berbasis mikroservis (microservices) yang dirancang menggunakan arsitektur asinkron *event-driven* (Kafka) dan pola arsitektur berlapis (*Clean/Layered Architecture*) di Go.

---

## 🗺️ Gambaran Arsitektur Proyek (Workspace Layout)

Proyek ini menggunakan fitur **Go Workspaces (`go.work`)** untuk mengelola multi-module dalam satu repositori lokal.

* **`common/`**: Modul bersama (*shared module*) yang diimpor oleh layanan lainnya untuk menghindari duplikasi kode (misal: koneksi database GORM).
* **`services/`**: Sub-direktori yang berisi layanan microservice independen dengan struktur folder modular standard Go:
  * **`services/order/`**: Mengelola pembuatan order pengiriman, kalkulasi tarif dinamis, alamat pengirim/penerima, dan transactional outbox event.
  * **`services/auth/`**: Mengelola registrasi user, login dengan BCrypt password hashing, JWT Access Token, dan Refresh Token.
  * **`services/payment/`**: Mengelola invoice pembayaran, status transaksi, dan simulasi callback webhook payment gateway.
  * **`services/shipping/`**: Mengelola registrasi kurir, penugasan kurir otomatis ke pesanan, dan pembaruan status logistik pengiriman paket.
  * **`services/warehouse/`**: Mengelola pendaftaran gudang (warehouses), pengaturan rute transit statis gudang, dan pencatatan riwayat paket masuk/keluar gudang (movements).

---

## 🚀 Port & Layanan (Local Ports Mapping)

| Layanan | Port API | Port Swagger UI | Database PostgreSQL |
| :--- | :--- | :--- | :--- |
| **Order Service** | `8081` | `http://localhost:8081/swagger/index.html` | `papiton_order` |
| **Auth Service** | `8082` | `http://localhost:8082/swagger/index.html` | `papiton_auth` |
| **Payment Service** | `8083` | `http://localhost:8083/swagger/index.html` | `papiton_payment` |
| **Shipping Service** | `8084` | `http://localhost:8084/swagger/index.html` | `papiton_shipping` |
| **Warehouse Service** | `8085` | `http://localhost:8085/swagger/index.html` | `papiton_warehouse` |

---

## 🛠️ Prasyarat Instalasi (Prerequisites)

Pastikan komputer lokal Anda sudah terpasang:
1. **Go 1.21+** atau lebih baru.
2. **Docker Desktop** (untuk PostgreSQL dan Kafka).
3. Tools Database Client (seperti **DBeaver** atau **pgAdmin**).

---

## 🏃‍♂️ Cara Setup & Menjalankan secara Lokal (Local Setup Guide)

### Langkah 1: Nyalakan Infrastruktur Docker
Jalankan kontainer database PostgreSQL 16, Kafka, dan Kafka UI menggunakan docker-compose di root folder:
```powershell
docker compose up -d
```
*Database akan otomatis terinisialisasi dengan membuat 7 database microservices terpisah sesuai skrip [init-db.sql](init-db.sql).*

### Langkah 2: Jalankan Migrasi Database
Sebelum menjalankan aplikasi, eksekusi skrip DDL SQL berikut pada database masing-masing menggunakan database client Anda:
1. **Auth DB (`papiton_auth`)**: Jalankan skrip [migrations/000001_init_auth_schema.up.sql](services/auth/migrations/000001_init_auth_schema.up.sql).
2. **Order DB (`papiton_order`)**: Jalankan skrip [migrations/000001_init_schema.up.sql](services/order/migrations/000001_init_schema.up.sql).
3. **Payment DB (`papiton_payment`)**: Jalankan skrip [migrations/000001_init_payment_schema.up.sql](services/payment/migrations/000001_init_payment_schema.up.sql).
4. **Shipping DB (`papiton_shipping`)**: Jalankan skrip [migrations/000001_init_shipping_schema.up.sql](services/shipping/migrations/000001_init_shipping_schema.up.sql).
5. **Warehouse DB (`papiton_warehouse`)**: Jalankan skrip [migrations/000001_init_warehouse_schema.up.sql](services/warehouse/migrations/000001_init_warehouse_schema.up.sql).

### Langkah 3: Jalankan Microservices
Buka terminal baru untuk masing-masing service dan jalankan perintah berikut:

* **Menjalankan Order Service**:
  ```powershell
  cd services/order
  go run ./cmd/api
  ```

* **Menjalankan Auth Service**:
  ```powershell
  cd services/auth
  go run ./cmd/api
  ```

* **Menjalankan Payment Service**:
  ```powershell
  cd services/payment
  go run ./cmd/api
  ```

* **Menjalankan Shipping Service**:
  ```powershell
  cd services/shipping
  go run ./cmd/api
  ```

* **Menjalankan Warehouse Service**:
  ```powershell
  cd services/warehouse
  go run ./cmd/api
  ```

---

## 🧪 Cara Menjalankan Pengujian (Testing Guide)

Proyek ini telah dilengkapi dengan **Unit Testing** (menggunakan mock repository) dan **Functional Integration Testing** (menguji HTTP Router ➔ Database riil).

* **Menjalankan seluruh test di workspace**:
  Dari folder root (`papiton/`), jalankan:
  ```powershell
  go test -v github.com/madgeer/papiton-express-go/services/order/... github.com/madgeer/papiton-express-go/services/payment/... github.com/madgeer/papiton-express-go/services/shipping/... github.com/madgeer/papiton-express-go/services/warehouse/...
  ```

* **Menjalankan test spesifik per layanan**:
  Masuk ke direktori layanan lalu jalankan test:
  ```powershell
  cd services/order
  go test -v ./...
  ```

---

## 📝 Pembaruan Dokumentasi API Swagger

Jika Anda melakukan penambahan/perubahan route handler di masa depan, perbarui dokumentasi Swagger dengan masuk ke direktori layanan dan jalankan:

* **Order Service**:
  ```powershell
  cd services/order
  go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g cmd/api/main.go
  ```
* **Auth Service**:
  ```powershell
  cd services/auth
  go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g cmd/api/main.go
  ```
* **Payment Service**:
  ```powershell
  cd services/payment
  go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g cmd/api/main.go
  ```
* **Shipping Service**:
  ```powershell
  cd services/shipping
  go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g cmd/api/main.go
  ```
* **Warehouse Service**:
  ```powershell
  cd services/warehouse
  go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g cmd/api/main.go
  ```
