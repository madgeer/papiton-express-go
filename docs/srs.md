# Software Requirements Specification (SRS)

## Daftar Isi

- [Cover](#cover)
- [Revision History](#revision-history)
- [BAB 1: Pendahuluan](#bab-1-pendahuluan)
- [BAB 2: Gambaran Umum](#bab-2-gambaran-umum)
- [BAB 3: Kebutuhan Sistem](#bab-3-kebutuhan-sistem)
- [BAB 4: Use Case](#bab-4-use-case)
- [BAB 5: Data Requirement](#bab-5-data-requirement)
- [BAB 6: Dokumentasi Teknis Pendukung](#bab-6-dokumentasi-teknis-pendukung)

---

## Cover

Software Requirements Specification untuk sistem ekspedisi Papiton Express.

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.1 | 2026-07-05 | Team | Revisi berdasarkan masukan arsitektur dan batasan domain |

---

## BAB 1: Pendahuluan

### 1.1 Latar Belakang
Papiton Express adalah sistem informasi ekspedisi yang dirancang untuk mendukung proses bisnis pengiriman paket dari awal hingga selesai. Sistem ini mencakup pembuatan order, pembayaran, penugasan kurir, pengelolaan gudang, pelacakan status, dan pemberian notifikasi.

Dikembangkan dengan pendekatan microservice, sistem ini menekankan pemisahan tanggung jawab tiap layanan sehingga perubahan pada satu domain tidak mengganggu domain lain. Pendekatan event-driven juga digunakan untuk menjaga integrasi antar layanan secara asynchronous.

### 1.2 Tujuan
Dokumen SRS ini disusun untuk menetapkan kebutuhan fungsional dan non-fungsional sistem Papiton Express agar dapat digunakan sebagai acuan analisis, desain, implementasi, pengujian, dan pemeliharaan.

Tujuan utama sistem adalah:
- menyediakan proses pengiriman paket yang terintegrasi;
- memisahkan domain bisnis ke dalam layanan yang jelas;
- mendukung transaksi dan status pengiriman yang dapat dilacak;
- memudahkan pengembangan fitur baru di masa mendatang.

### 1.3 Ruang Lingkup
Ruang lingkup sistem mencakup layanan berikut:
- Auth Service untuk autentikasi dan otorisasi.
- Gateway sebagai entry point untuk client.
- Order Service untuk lifecycle order.
- Payment Service untuk create payment dan verify payment.
- Shipping Service untuk penugasan kurir dan perjalanan pengiriman.
- Warehouse Service untuk proses masuk, sortir, dan keluar paket.
- Tracking Service untuk menyimpan dan menampilkan riwayat perjalanan paket.
- Notification Service untuk menyampaikan notifikasi kepada pelanggan.

### 1.4 Definisi Istilah
| Istilah | Definisi |
|---------|----------|
| Order | Permintaan pengiriman paket yang dibuat oleh customer. |
| Shipment | Proses pengiriman paket dari pengirim ke penerima. |
| Tracking Event | Catatan perubahan status atau lokasi paket. |
| Payment | Proses verifikasi pembayaran untuk suatu order. |
| Gateway | Entry point untuk client ke layanan backend. |
| Event | Informasi yang diproduksi oleh satu layanan untuk diproses layanan lain. |

### 1.5 Referensi
Dokumen teknis pendukung tersedia pada:
- [api-specification.md](api-specification.md)
- [event-flow.md](event-flow.md)
- [architecture.md](architecture.md)
- [deployment.md](deployment.md)

---

## BAB 2: Gambaran Umum

### 2.1 Deskripsi Sistem
Papiton Express adalah sistem ekspedisi digital yang menghubungkan customer, kurir, petugas gudang, dan admin melalui serangkaian layanan backend. Sistem ini mengatur alur dari pembuatan order hingga paket diterima dan status pengiriman selesai.

### 2.2 Tujuan Sistem
Tujuan sistem adalah memberikan pengalaman pengiriman yang terstruktur, aman, dan mudah dilacak. Sistem juga harus memungkinkan tim engineering untuk mengembangkan layanan secara mandiri.

### 2.3 Stakeholder
- Customer: membuat order, melihat status, dan menerima notifikasi.
- Admin: memantau operasi dan mengelola data operasional.
- Courier: mengambil dan mengantarkan paket.
- Warehouse Staff: menangani paket masuk, sortir, dan keluar gudang.
- System Administrator: mengelola infrastruktur dan deployment.

### 2.4 Alur Bisnis Utama
1. Customer mendaftar dan login melalui Auth Service.
2. Customer membuat order melalui Order Service.
3. Order menunggu pembayaran.
4. Payment Service menerima pembayaran dan melakukan verifikasi.
5. Shipping Service menugaskan kurir.
6. Kurir mengambil paket dan status berpindah sesuai tahapan.
7. Warehouse Service menangani proses sortir dan perpindahan paket.
8. Tracking Service mencatat peristiwa perjalanan paket.
9. Notification Service mengirimkan notifikasi kepada customer.

### 2.5 Domain dan Layanan
| Domain | Layanan | Fokus Utama |
|--------|---------|-------------|
| Auth | Auth Service | Register, login, JWT, refresh token |
| Order | Order Service | CRUD order, resi, status order |
| Payment | Payment Service | Create payment, verify payment |
| Shipping | Shipping Service | Penugasan kurir, pickup, delivery |
| Warehouse | Warehouse Service | Paket masuk, sorting, paket keluar |
| Tracking | Tracking Service | Timeline tracking dan riwayat status |
| Notification | Notification Service | Notifikasi pelanggan |

### 2.6 Arsitektur Tingkat Tinggi
Alur utama sistem adalah:
Gateway → Auth → Order → Payment → Shipping → Warehouse → Tracking → Notification.
Komunikasi antar layanan dilakukan menggunakan REST API dan event melalui Kafka. Tracking dan Notification bertindak sebagai consumer, bukan publisher.

---

## BAB 3: Kebutuhan Sistem

### 3.1 Functional Requirement

#### Auth Service
- FR-1.1: Sistem harus menyediakan registrasi pengguna.
- FR-1.2: Sistem harus menyediakan login pengguna.
- FR-1.3: Sistem harus menghasilkan access token dan refresh token.
- FR-1.4: Sistem harus mendukung refresh token.

#### Order Service
- FR-2.1: Sistem harus memungkinkan customer membuat order.
- FR-2.2: Sistem harus menyimpan alamat pengirim dan penerima.
- FR-2.3: Sistem harus menghasilkan nomor resi unik.
- FR-2.4: Sistem harus menyimpan status order sesuai alur bisnis.
- FR-2.5: Sistem harus mendukung perubahan status order dari waiting payment sampai completed.

#### Payment Service
- FR-3.1: Sistem harus menyediakan create payment.
- FR-3.2: Sistem harus menyediakan verify payment.
- FR-3.3: Sistem harus mengubah status order menjadi paid setelah pembayaran terverifikasi.

#### Shipping Service
- FR-4.1: Sistem harus menugaskan kurir berdasarkan ketersediaan.
- FR-4.2: Sistem harus mengubah status pengiriman menjadi picked up, in delivery, dan delivered.
- FR-4.3: Sistem harus mendukung penanganan pengiriman gagal.

#### Warehouse Service
- FR-5.1: Sistem harus mencatat paket masuk ke gudang.
- FR-5.2: Sistem harus melakukan proses sorting.
- FR-5.3: Sistem harus mencatat paket keluar dari gudang.

#### Tracking Service
- FR-6.1: Sistem harus mengonsumsi event penting dan menyimpan timeline pengiriman.
- FR-6.2: Sistem harus menyediakan endpoint untuk melihat status paket.

#### Notification Service
- FR-7.1: Sistem harus mengonsumsi event penting dan mengirim notifikasi.
- FR-7.2: Sistem harus mendukung notifikasi email dan push notification.

### 3.2 Non-Functional Requirement
- NFR-1: Response time API utama maksimal 2 detik pada beban normal.
- NFR-2: Sistem harus tersedia 99.5% dalam operasional produksi.
- NFR-3: Sistem harus dapat diskalakan untuk beban trafik meningkat.
- NFR-4: Semua komunikasi harus aman melalui HTTPS dan JWT.
- NFR-5: Sistem harus memiliki logging, monitoring, dan tracing yang memadai.
- NFR-6: Setiap layanan harus dapat dikembangkan dan diuji secara independen.

### 3.3 Business Rules
- BR-1: Order hanya dapat diproses jika data pengirim dan penerima lengkap.
- BR-2: Pembayaran harus terverifikasi sebelum kurir dapat ditugaskan.
- BR-3: Status order harus mengikuti urutan yang telah ditentukan.
- BR-4: Paket hanya dapat dipindahkan ke tahap berikutnya setelah status valid.
- BR-5: Tracking dan Notification hanya membaca event, tidak mempublikasikan event baru.

### 3.4 Constraint
- CON-1: Arsitektur harus berbasis microservice.
- CON-2: Setiap layanan memiliki database sendiri.
- CON-3: Implementasi menggunakan Go, Gin, PostgreSQL, Kafka, Docker, Kubernetes, dan JWT.
- CON-4: Gateway menjadi titik masuk utama untuk seluruh client.

### 3.5 Assumption
- ASM-1: Infrastruktur Kafka dan PostgreSQL tersedia untuk lingkungan development dan production.
- ASM-2: Tim operasional dapat menyediakan akun kurir dan petugas gudang.
- ASM-3: Auth Service sudah tersedia sebelum client dapat mengakses layanan lain.

---

## BAB 4: Use Case

### 4.1 Actor
- Customer
- Admin
- Courier
- Warehouse Staff
- System Administrator

### 4.2 Main Use Cases
1. Register dan login user.
2. Membuat order pengiriman.
3. Membuat dan memverifikasi pembayaran.
4. Menugaskan kurir dan mengubah status pengiriman.
5. Menangani paket masuk, sorting, dan paket keluar gudang.
6. Melihat tracking dan menerima notifikasi.

### 4.3 Catatan Desain
Diagram use case, ERD, dan sequence diagram disajikan dalam dokumen arsitektur pendukung agar SRS tetap fokus pada kebutuhan bisnis.

---

## BAB 5: Data Requirement

### 5.1 Domain Data
- Auth: users, refresh_tokens
- Order: orders, shipping_addresses, service_types, tariffs
- Payment: payments, payment_methods, payment_logs
- Shipping: shipments, couriers, courier_locations
- Warehouse: warehouses, warehouse_movements, warehouse_locations
- Tracking: tracking_events
- Notification: notification_logs, notification_templates

### 5.2 Data Requirement Utama
- Setiap order wajib memiliki data pengirim, penerima, dan status.
- Setiap pembayaran harus terhubung dengan order yang valid.
- Setiap shipment harus memiliki kurir dan lokasi saat ini.
- Setiap warehouse movement harus mencatat paket masuk dan keluar.
- Setiap perubahan status harus tercatat sebagai tracking event.

---

## BAB 6: Dokumentasi Teknis Pendukung

Dokumen teknis yang lebih rinci disimpan terpisah agar SRS tetap ringkas dan mudah dibaca:
- [api-specification.md](api-specification.md)
- [event-flow.md](event-flow.md)
- [architecture.md](architecture.md)
- [deployment.md](deployment.md)

Dokumen-dokumen tersebut memuat endpoint API, alur event Kafka, diagram arsitektur, ERD, sequence diagram, serta rencana deployment.
