## Name : Rehan Athallah Azhar


# CineStream 

layanan RESTful API Backend berbasis Golang yang  mengelola sistem perbioskopan, mencakup otentikasi pengguna, manajemen bioskop, studio, data film (movies), dan penjadwalan pemutaran film (showtimes).

---

## Teknologi & Library

- **Bahasa Pemrograman**: Go (Golang) v1.25+
- **Framework Web**: [Echo v4](https://echo.labstack.com/)
- **Database**: PostgreSQL 16
- **Database Access & Code Generation**: [SQLC](https://sqlc.dev/) & `database/sql` (`lib/pq`)
- **Otentikasi & Keamanan**: JWT (`golang-jwt/jwt/v5`) & Bcrypt (`golang.org/x/crypto`)
- **Validasi Data**: `go-playground/validator/v10`
- **Database Migration**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Containerization**: Docker & Docker Compose

---

## Cara Menjalankan Menggunakan Docker

### 1. Prasyarat
Pastikan Anda sudah menginstal:
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### 2. Salin File Environment (`.env`)
Buat file `.env` berdasarkan contoh `.env.example`:
```bash
cp .env.example .env
```

### 3. Jalankan Service dengan Docker Compose
Jalankan command berikut untuk membuat dan menjalankan container API dan PostgreSQL:
```bash
docker-compose up -d --build
```

Service yang akan berjalan:
- **API Server**: `http://localhost:8080`
- **PostgreSQL**: `localhost:5433` (di dalam network docker terhubung ke port `5432`)

---

## Database Migrations

Project ini menggunakan `golang-migrate` untuk mengelola skema database dan data awal (seeding).

### Opsi A: Jalankan Migrasi Menggunakan Docker

#### 1. Menjalankan Migration UP (Membuat Tabel & Seed Data)
Jalankan skema database dasar (`000001_init_schema` dan `000002_seed_data`):

**Bash:**
```bash
docker run --rm -v $(pwd)/db/migrations:/migrations --network cinestream_network migrate/migrate -path=/migrations/ -database "postgres://postgres:postgres@postgres:5432/cinestream_db?sslmode=disable" up
```


#### 2. Menjalankan Migration DOWN (Rollback)
Untuk membatalkan/menghapus skema tabel:

**Linux / macOS (Bash):**
```bash
docker run --rm -v $(pwd)/db/migrations:/migrations --network cinestream_network migrate/migrate -path=/migrations/ -database "postgres://postgres:postgres@postgres:5432/cinestream_db?sslmode=disable" down -all
```


## Mengunduh Dependensi (`go get` Script)

### Menggunakan Script Otomatis

**Linux / macOS (Bash):**
```bash
chmod +x install_deps.sh
./install_deps.sh
```


**atau bisa Manual Command (`go get`)**

```bash
go get github.com/labstack/echo/v4
go get github.com/lib/pq
go get github.com/golang-jwt/jwt/v5
go get github.com/google/uuid
go get github.com/go-playground/validator/v10
go get github.com/joho/godotenv
go get golang.org/x/crypto

# Verifikasi & Tidy
go mod tidy
go mod download
```

---

## Pengujian API

Dokumentasi endpoint dan koleksi API dapat di-import langsung ke Postman:
- File Koleksi: [`CineStream.postman_collection.json`](./CineStream.postman_collection.json)

---

## Struktur Folder Project

```text
CineStream/
├── cmd/
│   └── server/          # Entry point aplikasi (main.go)
├── db/
│   ├── migrations/      # File SQL migration (.up.sql & .down.sql)
│   └── queries/         # SQL queries untuk generator SQLC
├── internal/
│   ├── auth/            # Modul otentikasi & user management
│   ├── config/          # Konfigurasi aplikasi & database
│   ├── database/        # Output kode tersinkronisasi SQLC
│   ├── middleware/      # Custom middleware (JWT, Error Handler, dll)
│   ├── response/        # Standard JSON response helper
│   └── showtime/        # Modul penjadwalan tayang film
├── Dockerfile           # Multi-stage Docker build config
├── docker-compose.yml   # Orchestration untuk API & Postgres container
├── install_deps.sh      # Bash script pengunduhan dependensi
├── install_deps.ps1     # PowerShell script pengunduhan dependensi
└── README.md            # Dokumentasi utama project
```
