# Application Knowledge Management System

## Overview

Application Knowledge Management System adalah aplikasi yang digunakan untuk mengelola informasi seluruh aplikasi atau website perusahaan dalam satu portal terpusat.

Sistem ini bertujuan untuk membantu tim IT dan developer dalam menyimpan dokumentasi, catatan teknis, informasi server, serta riwayat pengembangan aplikasi sehingga seluruh pengetahuan terkait aplikasi dapat terdokumentasi dengan baik.

---

## Objectives

* Menyimpan daftar aplikasi atau website perusahaan.
* Menyimpan informasi server yang digunakan oleh setiap aplikasi.
* Menyimpan dokumentasi teknis dan bisnis.
* Menyimpan catatan developer.
* Menyimpan request fitur dan enhancement.
* Menyimpan riwayat perubahan aplikasi.
* Mempermudah onboarding developer baru.
* Menjadi pusat knowledge base untuk seluruh aplikasi perusahaan.

---

## Technology Stack

### Backend

* Golang
* Gin Framework
* GORM
* MySQL / MariaDB
* REST API

### Frontend (Planned)

* Laravel Blade
* Vue.js (Optional)

---

## Current Development Scope

### Systems

Master data aplikasi atau website.

Contoh:

* Portal Gudang
* Sistem Kalibrasi
* HRIS
* IT Ticketing

### Servers

Menyimpan informasi server yang digunakan oleh sistem.

Contoh:

* Application Server
* Database Server
* Backup Server

### Notes

Menyimpan catatan terkait sistem.

Contoh:

* Cara Deploy
* Konfigurasi Nginx
* Struktur Database
* Catatan Maintenance
* Informasi Integrasi API

---

## Planned Modules

### Feature Requests

Menyimpan daftar request pengembangan.

Status:

* Pending
* In Progress
* Testing
* Done
* Closed

### Documentation

Menyimpan dokumentasi aplikasi.

Kategori:

* Business Flow
* Technical Flow
* API Documentation
* Database Documentation
* Deployment Guide
* User Manual

### Changelog

Menyimpan riwayat perubahan aplikasi.

Contoh:

#### Version 1.0

* Login
* Dashboard
* User Management

#### Version 1.1

* Export Excel
* Upload Attachment

### User Management

Mengelola pengguna dan hak akses.

Role:

* Administrator
* Developer
* Viewer


---

## API Roadmap

### Systems

```http
GET    /api/systems
GET    /api/systems/{id}
POST   /api/systems
PUT    /api/systems/{id}
DELETE /api/systems/{id}
```

### Servers

```http
GET    /api/servers
GET    /api/servers/{id}
POST   /api/servers
PUT    /api/servers/{id}
DELETE /api/servers/{id}
```

### Notes

```http
```
