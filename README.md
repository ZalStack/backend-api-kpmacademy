# 🎓 KPM Academy — Backend API

> **Production-grade RESTful API** untuk platform belajar UTBK/SNBT  
> Dibangun dengan **Go · Fiber v2 · GORM · MySQL**

---

## 📋 Daftar Isi

- [Tech Stack](#-tech-stack)
- [Struktur Proyek](#-struktur-proyek)
- [Cara Menjalankan](#-cara-menjalankan)
- [Environment Variables](#-environment-variables)
- [Autentikasi](#-autentikasi)
- [Format Respons](#-format-respons)
- [Rate Limiting](#-rate-limiting)
- [Daftar Endpoint](#-daftar-endpoint)
  - [Health & Info](#1-health--info)
  - [Auth](#2-auth)
  - [Profile](#3-profile)
  - [Dashboard](#4-dashboard)
  - [Notifikasi](#5-notifikasi)
  - [Paket](#6-paket)
  - [Order](#7-order)
  - [Pembayaran (Webhook)](#8-pembayaran-webhook)
  - [Latihan Soal](#9-latihan-soal)
  - [Video](#10-video)
  - [Testimoni](#11-testimoni)
  - [Kontak Form](#12-kontak-form)
  - [Support Ticket](#13-support-ticket)
  - [Chat](#14-chat)
  - [Admin — Dashboard](#15-admin--dashboard)
  - [Admin — User](#16-admin--user)
  - [Admin — Paket](#17-admin--paket)
  - [Admin — Order](#18-admin--order)
  - [Admin — Transaksi](#19-admin--transaksi)
  - [Admin — Enroll Keys](#20-admin--enroll-keys)
  - [Admin — Statistik Latihan](#21-admin--statistik-latihan)
  - [Admin — Laporan](#22-admin--laporan)
  - [Admin — Testimoni](#23-admin--testimoni)
  - [Admin — Support](#24-admin--support)
  - [Admin — Kontak Form](#25-admin--kontak-form)
  - [Admin — Video](#26-admin--video)
  - [Admin — Video Orders](#27-admin--video-orders)
  - [Admin — Notifikasi](#28-admin--notifikasi)
- [Error Codes](#-error-codes)
- [Dokumentasi Interaktif](#-dokumentasi-interaktif)
- [Testing Endpoint](#-testing-endpoint)

---

## 🛠 Tech Stack

| Layer | Teknologi |
|---|---|
| **Runtime** | Go 1.22+ |
| **HTTP Framework** | Fiber v2 |
| **ORM** | GORM |
| **Database** | MySQL 8.x |
| **Auth** | JWT (Access + Refresh Token Rotation) |
| **Payment** | Midtrans |
| **Rate Limiting** | In-memory sliding window |
| **Dokumentasi** | OpenAPI 3.0 (Scalar / Swagger / Redoc) |

---

## 📁 Struktur Proyek

```
backend-api-kpmacademy/
├── config/             # Konfigurasi aplikasi (env loader)
├── database/           # Koneksi database GORM
├── docs/               # Dokumentasi API (OpenAPI, HTML UI)
│   ├── openapi.json    # Spesifikasi OpenAPI 3.0 (JSON)
│   ├── swagger.yaml    # Spesifikasi OpenAPI 3.0 (YAML)
│   ├── landing.html    # Landing page docs
│   ├── scalar.html     # Scalar API Reference UI
│   ├── swagger.html    # Swagger UI
│   └── redoc.html      # Redoc UI
├── handlers/           # Handler per fitur (HTTP logic)
├── middleware/         # Auth, CORS, Rate Limiter, dll.
├── migrations/         # AutoMigrate & seed admin
├── models/             # Struct GORM model
├── routes/             # Registrasi semua route
├── utils/              # Helper: JWT, response, validator, dll.
├── uploads/            # File upload (avatar, pdf, dll.)
├── .env                # Konfigurasi lokal (JANGAN commit)
├── .env.example        # Template env
├── main.go             # Entry point
└── test_all_endpoints.sh  # Script test semua endpoint (118 cases)
```

---

## 🚀 Cara Menjalankan

### Prasyarat
- Go 1.22+
- MySQL 8.x (running di port 3306)
- File `.env` sudah dikonfigurasi

### 1. Clone & Install Dependencies

```bash
git clone <repo-url>
cd backend-api-kpmacademy
go mod tidy
```

### 2. Konfigurasi Environment

```bash
cp .env.example .env
# Edit .env sesuai konfigurasi lokal kamu
```

### 3. Jalankan Server

```bash
go run .
# Server berjalan di http://localhost:3000
```

### 4. Jalankan dengan Hot Reload (Development)

```bash
# Install air jika belum ada
go install github.com/air-verse/air@latest

air
```

### 5. Jalankan Tests

```bash
# Unit test
go test ./...

# End-to-end test semua endpoint (memerlukan MySQL & server)
bash test_all_endpoints.sh
```

---

## ⚙️ Environment Variables

| Variabel | Wajib | Default | Keterangan |
|---|---|---|---|
| `APP_NAME` | ✗ | `KPM Academy` | Nama aplikasi |
| `APP_PORT` | ✗ | `3000` | Port server |
| `APP_ENV` | ✗ | `development` | `development` / `production` |
| `APP_DEBUG` | ✗ | `false` | Aktifkan request logger |
| `DB_HOST` | ✓ | `127.0.0.1` | Host MySQL |
| `DB_PORT` | ✗ | `3306` | Port MySQL |
| `DB_USER` | ✓ | — | Username MySQL |
| `DB_PASSWORD` | ✓ | — | Password MySQL |
| `DB_NAME` | ✓ | — | Nama database |
| `JWT_SECRET` | ✓ | — | Secret untuk access token |
| `JWT_REFRESH_SECRET` | ✓ | — | Secret untuk refresh token |
| `JWT_EXPIRATION` | ✗ | `3600` | Durasi access token (detik) |
| `JWT_REFRESH_EXPIRATION` | ✗ | `604800` | Durasi refresh token (detik) |
| `MIDTRANS_SERVER_KEY` | ✓ | — | Server key Midtrans |
| `MIDTRANS_CLIENT_KEY` | ✗ | — | Client key Midtrans |
| `MIDTRANS_IS_PRODUCTION` | ✗ | `false` | Mode production Midtrans |
| `UPLOAD_DIR` | ✗ | `./uploads` | Direktori upload file |
| `MAX_UPLOAD_SIZE` | ✗ | `10485760` | Maks ukuran upload (bytes) |
| `RATE_LIMIT` | ✗ | `100` | Maks request per IP / menit |
| `CORS_ALLOWED_ORIGINS` | ✗ | `localhost:3000,localhost:5173` | Origin yang diizinkan CORS |
| `ADMIN_EMAIL` | ✓ | — | Email admin yang di-seed |
| `ADMIN_PASSWORD` | ✓ | — | Password admin yang di-seed |

> **Catatan:** Saat `APP_ENV=development`, rate limit otomatis ditingkatkan (auth: 500 req/min, global: 1000 req/min) agar tidak mengganggu pengembangan.

---

## 🔐 Autentikasi

API menggunakan **JWT Bearer Token** dengan mekanisme **Refresh Token Rotation**.

### Flow Lengkap

```
1. Register   →  POST /api/v1/auth/register
2. Login      →  POST /api/v1/auth/login
                 ← { access_token, refresh_token, expires_at }
3. Gunakan access_token di setiap request:
                 Authorization: Bearer <access_token>
4. Token expired (1 jam) → POST /api/v1/auth/refresh-token
                           ← { access_token (baru), refresh_token (baru) }
5. Logout     →  POST /api/v1/auth/logout  (menginvalidasi refresh token)
```

### Level Akses

| Level | Keterangan |
|---|---|
| 🌐 **Public** | Tidak butuh token |
| 🔒 **Auth** | Butuh `Authorization: Bearer <token>` |
| 👑 **Admin** | Butuh token dengan `role: admin` |

### Kebijakan Password

Password harus mengandung minimal **8 karakter** dengan semua kriteria berikut:
- Huruf besar (A–Z)
- Huruf kecil (a–z)
- Angka (0–9)
- Karakter spesial (`!@#$%^&*`)

---

## 📦 Format Respons

Semua endpoint menggunakan format JSON yang konsisten.

### ✅ Success (2xx)

```json
{
  "success": true,
  "message": "Data berhasil diambil",
  "data": { ... },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### ✅ Success dengan Paginasi

```json
{
  "success": true,
  "message": "Data berhasil diambil",
  "data": [ ... ],
  "pagination": {
    "total": 150,
    "per_page": 10,
    "current_page": 1,
    "last_page": 15
  },
  "request_id": "..."
}
```

### ❌ Error (4xx / 5xx)

```json
{
  "success": false,
  "message": "Validasi gagal",
  "errors": [
    { "field": "email", "message": "harus berupa email yang valid" }
  ],
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Parameter Paginasi

Semua endpoint list mendukung query parameter:

| Parameter | Tipe | Default | Maks | Keterangan |
|---|---|---|---|---|
| `page` | integer | `1` | — | Halaman (mulai dari 1) |
| `per_page` | integer | `10` | `100` | Jumlah data per halaman |

---

## 🚦 Rate Limiting

| Endpoint | Limit | Window |
|---|---|---|
| `POST /auth/register` | 30 req | 60 detik |
| `POST /auth/login` | 30 req | 60 detik |
| `POST /kontak-form` | 30 req | 60 detik |
| `POST /auth/forgot-password` | 15 req | 60 detik |
| `POST /auth/reset-password` | 15 req | 60 detik |
| Semua endpoint lainnya | 100 req | 60 detik |

> **Development mode** (`APP_ENV=development`): Limit ditingkatkan otomatis → auth: 500, strict: 300, global: 1000.

Respons saat limit terlampaui:
```json
{
  "success": false,
  "message": "Too many requests, please try again later"
}
```

---

## 📡 Daftar Endpoint

> **Base URL:** `http://localhost:3000/api/v1`  
> **Alias:** `http://localhost:3000/api` (backward compatibility)

---

### 1. Health & Info

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/` | 🌐 Public | Informasi API & link dokumentasi |
| `GET` | `/health` | 🌐 Public | Health check (termasuk koneksi DB) |

**Contoh Respons `/health`:**
```json
{
  "success": true,
  "message": "KPM Academy API is running",
  "version": "v1",
  "time": "2026-09-23T08:00:00+07:00"
}
```

---

### 2. Auth

> Base: `/auth`

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `POST` | `/auth/register` | 🌐 Public | Daftar akun baru |
| `POST` | `/auth/login` | 🌐 Public | Login dan dapatkan token |
| `POST` | `/auth/refresh-token` | 🌐 Public | Perbarui access token |
| `POST` | `/auth/forgot-password` | 🌐 Public | Kirim link reset password ke email |
| `POST` | `/auth/reset-password` | 🌐 Public | Reset password dengan token |
| `POST` | `/auth/logout` | 🔒 Auth | Logout dan invalidasi refresh token |

**`POST /auth/register` — Request Body:**
```json
{
  "name": "Budi Santoso",
  "email": "budi@email.com",
  "password": "StrongPass123!",
  "phone": "081234567890"
}
```

**`POST /auth/login` — Request Body:**
```json
{
  "email": "budi@email.com",
  "password": "StrongPass123!"
}
```

**`POST /auth/login` — Response `data`:**
```json
{
  "user": { "id": "...", "name": "Budi Santoso", "email": "...", "role": "user" },
  "tokens": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_at": 1695000000
  }
}
```

**`POST /auth/refresh-token` — Request Body:**
```json
{ "refresh_token": "eyJ..." }
```

**`POST /auth/forgot-password` — Request Body:**
```json
{ "email": "budi@email.com" }
```

**`POST /auth/reset-password` — Request Body:**
```json
{
  "token": "ABCDEF123456",
  "email": "budi@email.com",
  "new_password": "NewPass456!"
}
```

---

### 3. Profile

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/profile` | 🔒 Auth | Lihat profil user yang login |
| `PUT` | `/profile` | 🔒 Auth | Update profil user |
| `PUT` | `/profile/change-password` | 🔒 Auth | Ganti password |

**`PUT /profile` — Request Body (semua opsional):**
```json
{
  "name": "Budi Santoso",
  "phone": "081234567890",
  "student_name": "Budi",
  "student_class": "XII IPA 1",
  "student_major": "IPA",
  "school_name": "SMAN 1 Jakarta",
  "address": "Jakarta",
  "gender": "male",
  "religion": "Islam"
}
```

**`PUT /profile/change-password` — Request Body:**
```json
{
  "old_password": "OldPass123!",
  "new_password": "NewPass456!"
}
```

---

### 4. Dashboard

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/dashboard` | 🔒 Auth | Statistik & ringkasan user (paket aktif, skor latihan, dll.) |

---

### 5. Notifikasi

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/notifications` | 🔒 Auth | Daftar notifikasi user |
| `GET` | `/notifications/unread-count` | 🔒 Auth | Jumlah notifikasi yang belum dibaca |
| `POST` | `/notifications/:id/read` | 🔒 Auth | Tandai satu notifikasi sebagai dibaca |
| `POST` | `/notifications/read-all` | 🔒 Auth | Tandai semua notifikasi sebagai dibaca |

---

### 6. Paket

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/packages` | 🌐 Public | Daftar semua paket aktif |
| `GET` | `/packages/:id` | 🌐 Public | Detail satu paket |

**Query Parameters untuk `GET /packages`:**

| Parameter | Tipe | Keterangan |
|---|---|---|
| `page` | integer | Halaman (default: 1) |
| `per_page` | integer | Data per halaman (default: 10) |
| `search` | string | Cari berdasarkan judul |
| `kelas` | string | Filter berdasarkan kelas |
| `jenjang` | string | Filter berdasarkan jenjang |

---

### 7. Order

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `POST` | `/orders` | 🔒 Auth | Buat order baru untuk sebuah paket |
| `GET` | `/orders` | 🔒 Auth | Daftar order milik user yang login |
| `GET` | `/orders/:id` | 🔒 Auth | Detail satu order |
| `GET` | `/orders/status` | 🔒 Auth | Cek status pembayaran |
| `POST` | `/orders/:id/pay` | 🔒 Auth | Simulasi pembayaran (dev only) |

**`POST /orders` — Request Body:**
```json
{
  "package_id": "550e8400-e29b-41d4-a716-446655440000",
  "total_price": 150000,
  "is_custom_amount": false
}
```

**Order Status Values:**

| Status | Keterangan |
|---|---|
| `pending` | Order dibuat, belum dibayar |
| `paid` | Pembayaran berhasil dikonfirmasi |
| `failed` | Pembayaran gagal |
| `refunded` | Dana dikembalikan |

---

### 8. Pembayaran (Webhook)

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `POST` | `/payment/notification` | 🌐 Public | Webhook notifikasi Midtrans (paket) |
| `POST` | `/payment/video-notification` | 🌐 Public | Webhook notifikasi Midtrans (video) |

> ⚠️ Endpoint ini dipanggil oleh server Midtrans, bukan oleh frontend. Tidak perlu token.

---

### 9. Latihan Soal

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `POST` | `/practice/start` | 🔒 Auth | Mulai sesi latihan baru |
| `POST` | `/practice/submit` | 🔒 Auth | Submit jawaban & selesaikan sesi |
| `GET` | `/practice/history` | 🔒 Auth | Riwayat semua sesi latihan user |
| `GET` | `/practice/statistics` | 🔒 Auth | Statistik latihan user |
| `GET` | `/practice/:id` | 🔒 Auth | Detail satu sesi latihan |

**`POST /practice/start` — Request Body:**
```json
{
  "package_id": "550e8400-e29b-41d4-a716-446655440000",
  "card_id": "card-1"
}
```

**`POST /practice/submit` — Request Body:**
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440001",
  "answers": "[{\"question_id\":\"q1\",\"answer\":\"A\"},{\"question_id\":\"q2\",\"answer\":\"B\"}]",
  "duration": 1800
}
```

---

### 10. Video

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/videos` | 🌐 Public | Daftar semua video aktif |
| `GET` | `/videos/:id` | 🌐 Public | Detail satu video |
| `POST` | `/videos/:id/order` | 🔒 Auth | Buat order untuk video berbayar |
| `POST` | `/videos/:id/pay/:orderId` | 🔒 Auth | Simulasi pembayaran video (dev) |

---

### 11. Testimoni

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/testimonials` | 🌐 Public | Daftar testimoni yang sudah diapprove |
| `POST` | `/testimonials` | 🔒 Auth | Buat testimoni baru (menunggu approval) |
| `GET` | `/testimonials/my` | 🔒 Auth | Lihat testimoni milik user yang login |

**`POST /testimonials` — Request Body:**
```json
{
  "content": "Platform ini sangat membantu persiapan UTBK saya!",
  "rating": 5
}
```

---

### 12. Kontak Form

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `POST` | `/kontak-form` | 🌐 Public | Kirim pesan kontak |
| `GET` | `/kontak-form` | 🔒 Auth | Daftar form kontak milik user |
| `GET` | `/kontak-form/:id` | 🔒 Auth | Detail satu form kontak milik user |

**`POST /kontak-form` — Request Body:**
```json
{
  "name": "Budi Santoso",
  "email": "budi@email.com",
  "phone": "081234567890",
  "subject": "Pertanyaan tentang paket",
  "message": "Apakah ada diskon untuk pelajar?"
}
```

---

### 13. Support Ticket

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `POST` | `/support/submit` | 🌐 Public | Kirim pertanyaan/tiket support |

**`POST /support/submit` — Request Body:**
```json
{
  "name": "Budi Santoso",
  "email": "budi@email.com",
  "question": "Bagaimana cara reset password saya?"
}
```

---

### 14. Chat

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `POST` | `/chat/send` | 🌐 Public | Kirim pesan chat |
| `GET` | `/chat/history` | 🔒 Auth | Ambil riwayat chat berdasarkan session_id |

**`POST /chat/send` — Request Body:**
```json
{
  "message": "Halo, saya butuh bantuan",
  "session_id": "CHAT-ABC12345"
}
```

> Jika `session_id` tidak dikirim, server akan membuat sesi baru secara otomatis.

**`GET /chat/history` — Query Parameter:**
```
GET /chat/history?session_id=CHAT-ABC12345
```

---

### 15. Admin — Dashboard

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/dashboard` | 👑 Admin | Statistik overview platform |

---

### 16. Admin — User

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/users` | 👑 Admin | Daftar semua user (paginasi + filter) |
| `GET` | `/admin/users/:id` | 👑 Admin | Detail satu user |
| `POST` | `/admin/users/:id/toggle-active` | 👑 Admin | Aktifkan / nonaktifkan user |
| `GET` | `/admin/login-logs` | 👑 Admin | Log login user |

---

### 17. Admin — Paket

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/packages` | 👑 Admin | Daftar semua paket (termasuk tidak aktif) |
| `POST` | `/admin/packages` | 👑 Admin | Buat paket baru |
| `GET` | `/admin/packages/:id` | 👑 Admin | Detail satu paket |
| `PUT` | `/admin/packages/:id` | 👑 Admin | Update paket |
| `DELETE` | `/admin/packages/:id` | 👑 Admin | Hapus paket (jika tidak ada order aktif) |
| `POST` | `/admin/packages/:id/cards` | 👑 Admin | Tambah kartu ke paket |
| `DELETE` | `/admin/packages/:id/cards/:cardId` | 👑 Admin | Hapus kartu dari paket |
| `POST` | `/admin/packages/:id/import-pdf` | 👑 Admin | Import soal dari file PDF |

**`POST /admin/packages` — Request Body:**
```json
{
  "title": "Paket UTBK 2024 — IPA",
  "description": "Paket latihan soal UTBK lengkap",
  "price": 150000,
  "discount_price": 120000,
  "is_discount_active": true,
  "membership_duration_days": 30,
  "kelas": "XII",
  "jenjang": "SMA",
  "hide_explanation": false,
  "time_limit_minutes": 90
}
```

---

### 18. Admin — Order

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/orders` | 👑 Admin | Daftar semua order |
| `GET` | `/admin/orders/:id` | 👑 Admin | Detail satu order |
| `POST` | `/admin/orders/:id/verify` | 👑 Admin | Verifikasi order secara manual |

---

### 19. Admin — Transaksi

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/transactions` | 👑 Admin | Daftar semua transaksi |
| `GET` | `/admin/transactions/stats` | 👑 Admin | Statistik transaksi (revenue, dll.) |
| `GET` | `/admin/transactions/:id` | 👑 Admin | Detail satu transaksi |

---

### 20. Admin — Enroll Keys

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/enroll-keys` | 👑 Admin | Daftar semua enrollment key |
| `GET` | `/admin/enroll-keys/:id` | 👑 Admin | Detail satu enrollment key |

---

### 21. Admin — Statistik Latihan

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/practice-statistics` | 👑 Admin | Statistik latihan semua user |
| `GET` | `/admin/practice-statistics/:id` | 👑 Admin | Detail statistik satu sesi latihan |

---

### 22. Admin — Laporan

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/reports` | 👑 Admin | Laporan umum platform |

---

### 23. Admin — Testimoni

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/testimonials` | 👑 Admin | Daftar semua testimoni |
| `POST` | `/admin/testimonials/:id/approve` | 👑 Admin | Approve testimoni user |
| `POST` | `/admin/testimonials/:id/toggle-active` | 👑 Admin | Aktifkan / nonaktifkan testimoni |
| `DELETE` | `/admin/testimonials/:id` | 👑 Admin | Hapus satu testimoni |
| `POST` | `/admin/testimonials/bulk-delete` | 👑 Admin | Hapus banyak testimoni sekaligus |

**`POST /admin/testimonials/bulk-delete` — Request Body:**
```json
{
  "ids": ["uuid-1", "uuid-2", "uuid-3"]
}
```

---

### 24. Admin — Support

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/support` | 👑 Admin | Daftar semua tiket support |
| `GET` | `/admin/support/:id` | 👑 Admin | Detail satu tiket |
| `POST` | `/admin/support/:id/answer` | 👑 Admin | Jawab tiket support |
| `PUT` | `/admin/support/:id/status` | 👑 Admin | Update status tiket |
| `DELETE` | `/admin/support/:id` | 👑 Admin | Hapus tiket |

**`POST /admin/support/:id/answer` — Request Body:**
```json
{
  "answer": "Untuk reset password, silakan klik tombol Lupa Password di halaman login."
}
```

**`PUT /admin/support/:id/status` — Request Body:**
```json
{
  "status": "closed"
}
```
> Status yang valid: `pending`, `answered`, `closed`

---

### 25. Admin — Kontak Form

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/kontak-form` | 👑 Admin | Daftar semua form kontak |
| `GET` | `/admin/kontak-form/:id` | 👑 Admin | Detail satu form kontak |
| `POST` | `/admin/kontak-form/:id/reply` | 👑 Admin | Balas form kontak via email |
| `PUT` | `/admin/kontak-form/:id/status` | 👑 Admin | Update status form kontak |
| `DELETE` | `/admin/kontak-form/:id` | 👑 Admin | Hapus form kontak |

**`POST /admin/kontak-form/:id/reply` — Request Body:**
```json
{
  "admin_reply": "Terima kasih atas pertanyaan Anda. Berikut jawabannya..."
}
```

**`PUT /admin/kontak-form/:id/status` — Request Body:**
```json
{
  "status": "replied"
}
```
> Status yang valid: `pending`, `read`, `replied`, `closed`

---

### 26. Admin — Video

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/videos` | 👑 Admin | Daftar semua video |
| `POST` | `/admin/videos` | 👑 Admin | Buat video baru |
| `GET` | `/admin/videos/:id` | 👑 Admin | Detail satu video |
| `PUT` | `/admin/videos/:id` | 👑 Admin | Update video |
| `DELETE` | `/admin/videos/:id` | 👑 Admin | Hapus video |
| `POST` | `/admin/videos/:id/toggle` | 👑 Admin | Aktifkan / nonaktifkan video |

---

### 27. Admin — Video Orders

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/video-orders` | 👑 Admin | Daftar semua order video |
| `POST` | `/admin/video-orders/:id/grant` | 👑 Admin | Grant akses video secara manual |

---

### 28. Admin — Notifikasi

| Method | Endpoint | Auth | Keterangan |
|---|---|---|---|
| `GET` | `/admin/notifications` | 👑 Admin | Daftar semua notifikasi sistem |

---

## 🔴 Error Codes

| HTTP Code | Keterangan |
|---|---|
| `200 OK` | Request berhasil |
| `201 Created` | Data berhasil dibuat |
| `204 No Content` | Berhasil, tidak ada respons body |
| `400 Bad Request` | Input tidak valid / validasi gagal |
| `401 Unauthorized` | Token tidak ada atau tidak valid |
| `403 Forbidden` | Tidak punya akses (role tidak mencukupi) |
| `404 Not Found` | Data atau route tidak ditemukan |
| `409 Conflict` | Data sudah ada (misalnya email duplikat) |
| `429 Too Many Requests` | Rate limit terlampaui |
| `503 Service Unavailable` | Database tidak dapat dijangkau |
| `504 Gateway Timeout` | Request melebihi batas waktu 30 detik |

---

## 📖 Dokumentasi Interaktif

Saat server berjalan, akses dokumentasi di:

| URL | Keterangan |
|---|---|
| [`/docs`](http://localhost:3000/docs) | Landing page dokumentasi |
| [`/docs/scalar`](http://localhost:3000/docs/scalar) | **Scalar UI** — Modern, recommended |
| [`/docs/swagger`](http://localhost:3000/docs/swagger) | **Swagger UI** — Classic |
| [`/docs/redoc`](http://localhost:3000/docs/redoc) | **Redoc** — Readable |
| [`/docs/openapi.json`](http://localhost:3000/docs/openapi.json) | Spec OpenAPI (JSON) |
| [`/docs/swagger.yaml`](http://localhost:3000/docs/swagger.yaml) | Spec OpenAPI (YAML) |

---

## 🧪 Testing Endpoint

Proyek ini dilengkapi script bash untuk testing **118 test cases** secara otomatis.

```bash
# Jalankan semua test (memerlukan MySQL aktif)
bash test_all_endpoints.sh

# Gunakan custom base URL
API_BASE_URL=http://localhost:8080/api bash test_all_endpoints.sh
```

**Cakupan Test:**
- ✅ Auth (register, login, logout, token refresh, forgot/reset password)
- ✅ Profile (get, update, change password)
- ✅ Order (create, list, detail, payment simulation)
- ✅ Practice (start, submit, history, statistics)
- ✅ Semua endpoint Admin (dashboard, users, packages, orders, dll.)
- ✅ Edge cases (invalid token, wrong role, rate limit, 404, dll.)

---

## 👥 Tim Pengembang

| Role | Kontak |
|---|---|
| Backend Engineering | engineering@kpmacademy.com |

---

*Dokumentasi ini diperbarui pada: September 2026*
