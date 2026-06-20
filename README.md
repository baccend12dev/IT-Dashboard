# IT Dashboard - Application Knowledge Management System

Aplikasi IT Dashboard adalah portal terpusat (Knowledge Base) yang digunakan untuk mengelola informasi seluruh aplikasi/website perusahaan, server penampung (host nodes), catatan teknis developer, dokumentasi sistem, serta pengajuan request fitur (*feature requests*).

Sistem ini membantu mempermudah tim IT dan developer dalam mendokumentasikan infrastruktur, proses onboarding developer baru, serta melacak riwayat pengembangan aplikasi.

---

## 🛠️ Tech Stack

| Layer | Technology |
| :--- | :--- |
| **Runtime & package manager** | Go Runtime, Node.js + npm |
| **Frontend** | React, Vite, TanStack Router, TanStack Query |
| **Backend** | Go, Gin Gonic |
| **Database & ORM** | SQLite / MySQL + GORM |
| **Auth** | JWT (JSON Web Tokens) |
| **UI** | Vanilla CSS (Premium Dark Mode) + Lucide Icons |
| **API Docs** | Swagger (gin-swagger) |
| **Testing** | Go testing (standard library) + `httptest` |

---

## ✨ Fitur Utama

1. **Dashboard Overview**: Menampilkan list seluruh sistem/aplikasi perusahaan beserta tipe teknologi, status operasional (Active, Maintenance, Offline), tautan akses, dan host server.
2. **Servers Manager**: Manajemen data server fisik atau virtual (VPS) lengkap dengan IP Address, Operating System, dan lokasinya.
3. **System Detail View**: Halaman detail yang menggabungkan:
   * **Informasi Sistem & Server**: Detail spesifikasi server tempat aplikasi di-deploy.
   * **Developer Notes**: Catatan operasional seperti Nginx Config, API Integration, atau cara deployment.
   * **System Documentations**: Dokumentasi yang dibagi menjadi beberapa kategori (*Business Flow, Technical Flow, API Documentation, Database Documentation, Deployment Guide, User Manual*).
   * **Feature Requests**: Form pengajuan fitur baru beserta status implementasinya.
4. **Pending Requests Menu (Sidebar)**: Halaman khusus lintas sistem yang menyaring dan menampilkan seluruh request fitur yang masih berstatus **"Pending"** untuk segera dieksekusi.
5. **Autentikasi & Role-based Access Control (RBAC)**:
   * Login JWT Session.
   * Role: **Administrator**, **Developer**, dan **Viewer** (Viewer hanya memiliki akses baca/read-only).
6. **Dokumentasi Swagger**: Akses API Docs interaktif untuk seluruh modul sistem backend.

---

## 🚀 Cara Menjalankan Project

### Prasyarat
* [Go](https://go.dev/doc/install) (versi 1.20 ke atas)
* [Node.js](https://nodejs.org/) (versi 18 ke atas) & npm

---

### 1. Menjalankan Backend (Go)

1. Masuk ke direktori `backend`:
   ```bash
   cd backend
   ```
2. Buat file `.env` di dalam folder `backend` (atau gunakan `.env.example` sebagai referensi):
   ```ini
   PORT=8080
   JWT_SECRET=rahasia_super_secure_anda
   ALLOWED_ORIGIN=http://localhost:5173
   # Jika menggunakan database SQLite, biarkan konfigurasi database kosong atau default.
   ```
3. Jalankan aplikasi backend. GORM secara otomatis akan melakukan migrasi database (membuat tabel) dan membuat **User Default** pertama kali:
   ```bash
   go run main.go
   ```
   * **Akun Default Admin**: 
     * **Username**: `admin`
     * **Password**: `admin123`
     * **Role**: `Administrator`

4. **Akses Swagger API Documentation**:
   Buka browser dan akses URL:
   [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

### 2. Menjalankan Frontend (React + Vite)

1. Masuk ke direktori `frontend`:
   ```bash
   cd ../frontend
   ```
2. Instal semua dependensi Node.js:
   ```bash
   npm install
   ```
3. Jalankan server development frontend:
   ```bash
   npm run dev
   ```
4. Buka browser dan buka alamat yang diberikan oleh Vite (biasanya `http://localhost:5173`).
5. Login menggunakan akun default admin (`admin` / `admin123`).

---

## 🧪 Pengujian (Testing) pada Backend

Backend telah dilengkapi dengan unit test dan integration test untuk memvalidasi fungsionalitas HTTP Handler (Controllers), Autentikasi JWT, serta operasi database menggunakan Database In-Memory (`sqlite :memory:`).

Untuk menjalankan seluruh test pada backend:

1. Masuk ke direktori `backend`:
   ```bash
   cd backend
   ```
2. Jalankan perintah test:
   ```bash
   go test ./...
   ```
3. Untuk melihat hasil test secara detail per modul (misalnya modul `controllers`):
   ```bash
   go test -v ./controllers/...
   ```

Cakupan pengujian meliputi:
* **Auth Controller**: Uji coba register user baru, login valid/invalid, serta pembuatan token JWT.
* **System Controller**: Uji coba pembuatan sistem baru, pengambilan data, validasi hak akses, dan penghapusan sistem.
* **Server Controller**: Uji coba operasi CRUD Server.
* **Feature Request Controller**: Uji coba penambahan request fitur per sistem, pengubahan status fitur, serta fitur penarikan seluruh data request fitur pending (`GET /api/feature-requests/pending`).
* **Documentation & Notes Controllers**: Validasi logika CRUD dokumentasi teknis dan catatan developer.

---

## 🐳 Docker Deployment (Opsional)

Project ini mendukung deployment menggunakan Docker & Docker Compose.

* **Mode Development/Production**:
  Jalankan perintah berikut di root folder project untuk menjalankan backend dan frontend sekaligus dalam kontainer Docker:
  ```bash
  docker-compose up --build
  ```
  Ini akan menjalankan:
  * Backend API pada port `8080`
  * Frontend Web pada port `80` (dilayani oleh Nginx di dalam Docker)
