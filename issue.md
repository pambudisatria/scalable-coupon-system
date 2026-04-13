# Project Planning: Go Clean Architecture Skeleton (Coupon System)

## 1. Ringkasan Tugas
Tugas ini bertujuan untuk membangun struktur dasar (*skeleton*) backend menggunakan bahasa Go (Golang). Target utamanya adalah membiasakan proyek dengan pola desain **Clean Architecture** demi menunjang *modularity*, *testability*, dan *separation of concerns*.

📌 **ATURAN UTAMA (SCOPE LIMITATION):**
- **TIDAK PERLU** mengimplementasikan logika bisnis yang rumit.
- **TIDAK PERLU** menggunakan transaksi database (*DB transactions*).
- **TIDAK PERLU** menulis aksi spesifik (endpoint HTTP handler aktif) — cukup berikan *placeholder* atau kerangka kosong.

---

## 2. Deliverables & Requirement

### A. Struktur Folder (Direktori)
Buat rancangan direktori utama sebagai berikut:
```text
.
├── cmd/
│   └── main.go       (entrypoint aplikasi)
├── internal/
│   ├── domain/       (berisi definisi struct Entity dan interface layer)
│   ├── usecase/      (interface untuk business logic, tanpa implementasi konkret yang rumit)
│   ├── repository/   (adapter interaksi DB, tanpa business logic)
│   ├── delivery/     (layer / HTTP handler placeholder menggunakan Fiber)
│   └── pkg/          (kumpulan shared utilities)
```

### B. Definisi Entity & Database Schema (GORM)
Definisikan entitas di layer `domain` yang mewakili bentuk dan sifat data bisnis yang selaras dengan GORM:

1. **Entity `Coupon` (Tabel `coupons`)**:
   - `name` (string, primary key / unik)
   - `amount` (int)
   - `remaining_amount` (int)

2. **Entity `Claim` (Tabel `claims`)**:
   - `user_id` (string)
   - `coupon_name` (string)
   - `created_at` (tanggal waktu pembuatan otomatis GORM)
   - ⚠️ **Krusial**: Wajib ada properti *Composite Unique Constraint* pada tingkat skema (Tabel) untuk kombinasi kolom `user_id` dan `coupon_name` agar satu pengguna tidak bisa me-redeem kupon yang sama berulang kali.

### C. Definisi Interface Layer
Simpan rancangan perantara komponen pada file antarmuka (interface) di direktori `domain`:

- `CouponRepository`:
  - `Create(coupon)`
  - `FindByName(name)`
  - `UpdateStock(...)`
- `ClaimRepository`:
  - `Create(claim)`
  - `Exists(user_id, coupon_name)`

### D. Setup Dependensi Utama (*Library*)
Harus menggunakan spesifikasi paket modul berikut dan buat file `go.mod` (melalui perintah standar `go mod init`).
- **Web Framework**: `github.com/gofiber/fiber/v2`
- **ORM Toolkit**: `gorm.io/gorm`
- **Postgres driver**: `gorm.io/driver/postgres`

### E. Configuration & Dependency Injection
- Buat sebuah pola peracikan file konfigurasi yang bersumber pada file *environment* (misalnya pembacaan URL dan *password* langsung dari *env-vars*).
- Pada `cmd/main.go`, implementasikan ***Manual Dependency Injection*** sederhana. Lakukan koneksi instance db (GORM), *passing* instance db ke Repository, lalu serahkan pada Service/Usecase dan Fiber secara berurutan.

### F. Infrastruktur Dockerisasi
Buat file persiapan peluncuran berikut:
1. **`Dockerfile`**: Resep multi-stage build yang mengambil sumber Go, melakukan compile `main.go`, dan mengeksekusi *binary*-nya pada *lightweight image*.
2. **`docker-compose.yml`**: Rancang dua service utama:
   - `app` : Menjalankan build konfigurasi dari Dockerfile.
   - `postgres` : Service penyedia database. Sambungkan environment DB ke app-service.

---
**Catatan untuk Eksekutor (Junior Programmer / Low-Tier AI Model):**  
Dokumen ini difungsikan sebagai cetak biru murni untuk pola desain arsitektural. Fokuskan pengerjaan pada **kekomprehensifan setup Go, migrasi skema tabel via GORM, perangkaian awal Fiber (meski tanpa respons data final), serta kelengkapan Docker agar bisa di-build dengan lancar.**
