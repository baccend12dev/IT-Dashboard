# Frontend Architecture

## Project Overview

Frontend untuk Application Knowledge Management System dibangun menggunakan React dan TypeScript.

Frontend berfungsi sebagai client yang mengonsumsi REST API dari backend Golang.

Semua business logic, authentication, authorization, dan data management berada di backend Golang.

Frontend hanya bertanggung jawab untuk:

* Menampilkan data
* Mengelola state UI
* Mengelola routing
* Mengirim request ke backend API
* Menampilkan feedback kepada user

---

# Technology Stack

## Core Framework

* React
* TypeScript
* Vite

---

## Data Fetching

### TanStack Query (Required)

TanStack Query adalah standar wajib untuk seluruh data fetching dan state management yang berasal dari backend.

Gunakan TanStack Query untuk:

* GET requests
* POST requests
* PUT requests
* DELETE requests
* Cache management
* Background refetch
* Mutation handling

### Rules

DO:

* Gunakan useQuery untuk mengambil data.
* Gunakan useMutation untuk create, update, delete.
* Gunakan query invalidation setelah mutation berhasil.
* Gunakan query keys yang konsisten.

DON'T:

* Jangan menggunakan useEffect untuk fetch data API.
* Jangan menyimpan response API ke React Context.
* Jangan membuat state duplicate untuk data yang sudah dikelola TanStack Query.

Example Query Keys:

```text
systems
system-detail
servers
notes
documents
feature-requests
users
```

---

# Routing

## TanStack Router (Required)

Gunakan TanStack Router untuk seluruh routing aplikasi.

### Benefits

* Type-safe routing
* Nested layouts
* Better developer experience
* Strong TypeScript integration

### Rules

DO:

* Gunakan nested routes.
* Gunakan route params dari TanStack Router.

DON'T:

* Jangan menggunakan React Router.
* Jangan menggunakan hash routing.

---

# Data Table

## TanStack Table (Required)

Gunakan TanStack Table untuk seluruh tampilan data tabular.

### Features

* Sorting
* Filtering
* Pagination
* Column Visibility
* Global Search

### Usage

Gunakan untuk:

* Systems List
* Servers List
* Notes List
* Documents List
* Feature Requests List
* Users List

### Rules

DO:

* Gunakan server-side pagination jika data besar.
* Gunakan reusable DataTable component.

DON'T:

* Jangan membuat table custom berulang kali.

---

# Authentication

## JWT Authentication

Authentication menggunakan JWT yang diterbitkan oleh backend Golang.

### Login Endpoint

```http
POST /api/auth/login
```

### User Profile Endpoint

```http
GET /api/auth/me
```

### Logout Endpoint

```http
POST /api/auth/logout
```

---

# Auth State Management

## React Context API

Gunakan React Context hanya untuk:

* Current User
* Authentication State
* Login State

Jangan gunakan React Context untuk data API yang sudah dikelola TanStack Query.

### Context Example

```text
AuthContext

- user
- token
- login()
- logout()
- isAuthenticated
```

---

# Token Storage

## Recommended

Gunakan:

* Secure HTTP Cookie (Preferred)

atau

* localStorage (Development / Internal Use)

### Rules

DO:

* Simpan token secara terpusat.
* Bersihkan token saat logout.

DON'T:

* Jangan menyimpan token di component state.

---

# API Client

## Axios

Gunakan satu instance Axios untuk seluruh aplikasi.

### Requirements

* Base URL terpusat
* Authorization Header otomatis
* Error handling terpusat

### Authorization Header

```http
Authorization: Bearer <token>
```

Harus ditambahkan otomatis ke setiap request yang membutuhkan autentikasi.

---

# Backend Communication

Semua komunikasi menggunakan REST API.

Response format:

```json
{
  "success": true,
  "message": "Success",
  "data": {}
}
```

Error format:

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {}
}
```

---

# UI Principles

### Dashboard First

Fokus aplikasi adalah dashboard internal dan knowledge management.

Prioritas:

1. Readability
2. Fast Navigation
3. Searchability
4. Data Organization

---

# AI Development Rules

When generating code for this project:

1. Always use TypeScript.
2. Always use TanStack Query for API requests.
3. Always use TanStack Router for routing.
4. Always use TanStack Table for tabular data.
5. Never use React Router.
6. Never use useEffect for API fetching.
7. Never store API data in Context API.
8. Use reusable components whenever possible.
9. Follow feature-based folder structure.
10. Assume backend is a Golang REST API.

These rules must be followed consistently throughout the project.
