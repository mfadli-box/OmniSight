# Plan: Backend API (ict_rest)

Dokumen ini berisi rencana pengembangan backend API OmniSight menggunakan **Go 1.26.4** & **Gin Gonic** dengan arsitektur **Linear Layer Execution**: Template (Struct & Interface) → Repository → Usecase → Handler.

---

## Daftar Isi

1. [Arsitektur & Tech Stack](#1-arsitektur--tech-stack)
2. [Error Handling Terpadu](#2-error-handling-terpadu)
3. [Sistem Akses (Privilege)](#3-sistem-akses-privilege)
4. [API Endpoints yang Sudah Diimplementasi](#4-api-endpoints-yang-sudah-diimplementasi)
5. [API Endpoints yang Direncanakan](#5-api-endpoints-yang-direncanakan)
6. [Struktur File](#6-struktur-file)

---

## 1. Arsitektur & Tech Stack

### Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Bahasa | Go 1.26.4 |
| HTTP Framework | Gin Gonic v1.12.0 |
| Database | PostgreSQL (raw `database/sql` + `lib/pq`, tanpa ORM) |
| Logging | Zerolog v1.35.1 |
| Password Hashing | bcrypt (`golang.org/x/crypto`) |
| CORS | `gin-contrib/cors` |
| Config | `godotenv` (`.env` file) |

### Arsitektur Linear Layer Execution

```
Template (Struct & Interface)
    ↓
Repository (SQL queries, database interaction)
    ↓
Usecase (business logic, validation)
    ↓
Handler (HTTP binding, response)
```

Setiap modul halaman (`skeleton/{KODE}/`) mengikuti pola ini:
- `template.go` — mendefinisikan struct domain (request/response) + interface `Repository` dan `UseCase`
- `repository.go` — mengimplementasikan interface `Repository` dengan SQL mentah
- `usecase.go` — mengimplementasikan interface `UseCase` dengan logika bisnis
- `handler.go` — HTTP handler yang menerima request, memanggil usecase, mengembalikan JSON

### Konvensi Constructor

| Layer | Constructor | Parameter |
|-------|-------------|-----------|
| Repository | `NRepo(db *sql.DB)` | Database connection |
| UseCase | `NCase(r Repository)` | Repository interface |
| Handler | `NHand(u UseCase)` | UseCase interface |

### Struktur Folder

```
ict_rest/
├── main.go                    Entry point, graceful shutdown
├── go.mod / go.sum
├── Dockerfile                 Multi-stage build
├── backbone/                  Infrastructure layer
│   ├── cleanup.go             Session cleanup goroutine
│   ├── database.go            PostgreSQL connection setup
│   ├── logger.go              Zerolog init, RequestID & Logger middleware
│   ├── memory.go              Auth middleware (USLoad, USLock, USRole, USLogs)
│   ├── recovery.go            Custom panic recovery middleware
│   ├── routes.go              Router setup, CORS, DI wiring, all route registration
│   └── upload.go              File upload handler (belum diwire ke route)
├── mechanic/                  Utility functions
│   ├── crypto.go              AES-GCM encrypt/decrypt
│   ├── helper.go              AppError type, Error() responder, constructor functions
│   └── typography.go          Pagination types, parsing helpers, BuildMeta
└── skeleton/                  Page modules
    ├── XX99/                  Empty scaffold template
    ├── SP00/                  Login / Logout / HRIS company list
    ├── SP01/                  User companies & module tree
    ├── SP02/                  Change password
    ├── SP03/                  User action audit log
    ├── SM01/                  User management (CRUD, company, privilege, location)
    ├── SM02/                  Module management
    ├── SM03/                  Company management (CRUD, modules, location areas)
    ├── SM04/                  Signature/approval type management
    └── SM05/                  Session management (list, detail, revoke)
```

### Dependencies

| Library | Versi | Fungsi |
|---------|-------|--------|
| `github.com/gin-contrib/cors` | v1.7.7 | CORS middleware |
| `github.com/gin-gonic/gin` | v1.12.0 | HTTP framework |
| `github.com/google/uuid` | v1.6.0 | UUID generation |
| `github.com/joho/godotenv` | v1.5.1 | `.env` file loading |
| `github.com/lib/pq` | v1.12.3 | PostgreSQL driver |
| `github.com/rs/zerolog` | v1.35.1 | Structured logging |
| `golang.org/x/crypto` | v0.54.0 | bcrypt password hashing |

### Konfigurasi

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `PG_HOST` | `localhost` | Host PostgreSQL |
| `PG_PORT` | `5432` | Port PostgreSQL |
| `PG_USER` | `dbe` | User PostgreSQL |
| `PG_PASS` | `rahasia` | Password PostgreSQL |
| `PG_DATA` | `erp` | Nama database |
| `IS_POOL` | `false` | Gunakan connection pool (100 open / 10 idle) |
| `RE_SECRET` | `omnisight-default-key-32b!` | Secret key untuk AES-GCM encryption |

### Connection Pool

| Mode | Max Open | Max Idle | Lifetime |
|------|----------|----------|----------|
| `IS_POOL=true` | 100 | 10 | 5 menit |
| `IS_POOL=false` | 50 | 25 | 5 menit |

---

## 2. Error Handling Terpadu

> Status: **SELESAI**

### Arsitektur Error Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                        BACKEND (Go)                             │
│                                                                 │
│  Repository  →  raw sql error                                   │
│  UseCase     →  AppError (user) / wrapped (internal)            │
│  Handler     →  mechanic.Error(c, err)                          │
│                 ↳ auto-detect: AppError→status, others→500      │
│                 ↳ log structured: request_id, method, path      │
│                 ↳ response: { code, error, request_id }         │
└──────────────────────────┬──────────────────────────────────────┘
                           │ HTTP Response
                           │ Header: X-Request-ID: abc123
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                     PROXY LAYER (Next.js)                       │
│                                                                 │
│  apiFetch()  →  forward X-Request-ID header                     │
│  401         →  auto-redirect to /login                         │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                     FRONTEND (React)                            │
│                                                                 │
│  clientApi() →  unified fetch wrapper                           │
│                 ↳ checks res.ok                                 │
│                 ↳ throws ClientApiError(message, status, code)  │
│  Page        →  try/catch → toast.error                         │
│  ErrorBoundary → fallback UI on React crash                     │
└─────────────────────────────────────────────────────────────────┘
```

### AppError — Satu-Satunya Cara Melempar Error

Seluruh error di backend WAJIB menggunakan `mechanic.AppError` via helper constructor:

| Helper | Kode | HTTP Status | Penggunaan |
|--------|------|-------------|------------|
| `mechanic.ValidationError(msg)` | `VALIDATION_ERROR` | 400 | Input tidak valid, field wajib kosong |
| `mechanic.NotFound(msg)` | `NOT_FOUND` | 404 | Resource tidak ditemukan |
| `mechanic.Conflict(msg)` | `CONFLICT` | 409 | Data duplikat (username sudah ada) |
| `mechanic.Unauthorized(msg)` | `UNAUTHORIZED` | 401 | Token tidak valid/kedaluwarsa |
| `mechanic.Forbidden(msg)` | `FORBIDDEN` | 403 | Akses ditolak (non-admin ke route admin) |
| `mechanic.ExternalServiceError(msg, err)` | `EXTERNAL_SERVICE_ERROR` | 502 | Service eksternal tidak terjangkau |
| `mechanic.InternalError(msg, err)` | `INTERNAL_ERROR` | 500 | Panic recovery, error tak terduga |

### Handler — Selalu `mechanic.Error(c, err)`

```go
// ✅ BENAR
if err := h.usecase.DoSomething(req); err != nil {
    mechanic.Error(c, err)
    return
}

// ❌ SALAH
if err := h.usecase.DoSomething(req); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### Usecase — Return `error` biasi, bukan `AppError`

```go
// ✅ BENAR — usecase mengembalikan AppError atau error biasa
func (u *useCase) CreateItem(name string) error {
    if name == "" {
        return mechanic.ValidationError("Name is required")
    }
    if exists, _ := u.repo.FindByName(name); exists != nil {
        return mechanic.Conflict("Item already exists")
    }
    return u.repo.Insert(name)
}
```

### Repository — Return `error` biasa

Repository tidak perlu `mechanic.Error()`. Cukup return `error` dari query database, atau `mechanic.Conflict(...)` untuk kasus duplikat constraint.

### Response Format

```json
{
  "code": "VALIDATION_ERROR",
  "error": "Username is required",
  "request_id": "uuid-dari-request"
}
```

### Structured Logging

Gunakan `backbone.Log` (zerolog) untuk structured logging:

```go
backbone.Log.Info().Str("user_id", id).Msg("User created")
backbone.Log.Warn().Str("token", t).Msg("Token expiring soon")
backbone.Log.Error().Err(err).Msg("Database query failed")
```

### HTTP Status Code Rules

| Status | Kode Error | Kapan Digunakan |
|--------|-----------|-----------------|
| **400** | `VALIDATION_ERROR` | Input tidak valid |
| **401** | `UNAUTHORIZED` | Token tidak ada/expired |
| **403** | `FORBIDDEN` | Tidak punya hak akses |
| **404** | `NOT_FOUND` | Resource tidak ditemukan |
| **409** | `CONFLICT` | Data konflik |
| **500** | `INTERNAL_ERROR` | DB error, panic recovery |
| **502** | `EXTERNAL_SERVICE_ERROR` | Service eksternal gagal |

### Graceful Shutdown

```go
func main() {
    srv := &http.Server{Addr: ":36665", Handler: rest}
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            backbone.Log.Fatal().Err(err).Msg("Server failed to start")
        }
    }()
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
    backbone.PgSQL.Close()
}
```

### Session Cleanup

Goroutine `CleanupExpiredSessions()` berjalan setiap 1 jam dan menghapus session yang sudah kedaluwarsa dari `dat_user_session`.

---

## 3. Sistem Akses (Privilege)

### Level Akses

| Level | Keterangan | Hak Akses |
|-------|------------|-----------|
| `hide` | Sembunyikan modul | Tidak bisa melihat |
| `view` | Lihat | Read-only |
| `book` | Kelola | CRUD |
| `post` | Publish | Semua + Publish/Approve/Escalate |

### Middleware yang Sudah Diimplementasi

| Middleware | Fungsi | Keterangan |
|-----------|--------|------------|
| `USLoad()` | Autentikasi session token | Baca header `Authorization`, validasi token dari `dat_user_session`, set `userId`, `companyId`, `isAdmin`, `isHris` di context |
| `USLock()` | Admin-only gate | Cek `isAdmin` dari context. Return 403 jika bukan admin |
| `USRole(moduleCode, requiredLevel)` | Role-based access | Admin langsung lolos. Non-admin: query `dat_user_privilege` → cek hierarchy `hide < view < book < post` |
| `USLogs(moduleCode)` | Audit logger | Jalankan **setelah** handler. Insert ke `dat_user_action` dengan: user_id, company_id, module_code, action, path, ip_address, user_agent |
| `RequestID()` | Request ID | Baca/generate UUID untuk `X-Request-ID` header |
| `Logger()` | Request logging | Log setiap request: level Info (<400), Warn (400-499), Error (≥500) |
| `CustomRecovery()` | Panic recovery | Tangkap panic, log stack trace, return 500 INTERNAL_ERROR |

### Alur Privilege Check

```
Request masuk
  → RequestID()         : generate/forward X-Request-ID
  → Logger()            : log request
  → CustomRecovery()    : tangkap panic
  → USLoad()            : autentikasi session token
  → USLock()            : cek admin (hanya untuk admin routes)
  → USRole(kode, level) : cek dat_user_privilege (belum diwire ke route)
  → USLogs(kode)        : audit log (setelah handler)
  → Handler             : proses request
```

### Level Hierarchy

```
hide < view < book < post
```

- `GET` request memerlukan minimal `view`
- `POST` / `PUT` / `DELETE` request memerlukan `post`
- `book` untuk operasi reservasi/penandaan

### Aturan Middleware

| Tipe Halaman | Middleware | Privilege Level |
|--------------|------------|-----------------|
| Guest (public) | Tidak ada | Tidak ada autentikasi |
| Pages (auth) | `USLoad()` | Login required |
| Admin Only | `USLoad()` + `USLock()` | `isAdmin = true` |
| CRUD (post) | `USLoad()` + `USLock()` + `USLogs()` | Admin only (USRole belum diwire) |
| View Only | `USLoad()` + `USLock()` | Admin only (USRole belum diwire) |

### Catatan: `USRole` Belum Diwire

`USRole(moduleCode, requiredLevel)` sudah diimplementasi di `backbone/memory.go` tapi **belum digunakan** di registrasi route manapun. Saat ini semua route yang memerlukan autentikasi menggunakan `USLoad()` + `USLock()` (admin-only). Penerapan `USRole` akan dilakukan saat modul-modul non-admin diimplementasi.

### Role Templates (IT)

| Role | Default Level |
|------|---------------|
| `admin` | `post` semua modul |
| `manager` | `book`~`post` modul inti |
| `analyst` | `view`~`book` modul security |
| `operator` | `view`~`book` modul operasional |
| `staff` | `view` modul dasar |
| `viewer` | `view` semua modul |

### Role Templates (non-IT)

| Role | Service Access |
|------|---------------|
| `requester` | Submit request, lihat status sendiri |
| `department_head` | Submit + approve department |
| `finance` | Lihat request biaya, approve budget |
| `hrd` | Request aset untuk karyawan baru |
| `project_manager` | Request server/aplikasi project |

---

## 4. API Endpoints yang Sudah Diimplementasi

### Global Routes (tanpa auth)

| Method | Endpoint | Handler | Keterangan |
|--------|----------|---------|------------|
| GET | `/` | — | 200 OK (empty) |
| GET | `/rest` | — | `{"message":"rest"}` |
| GET | `/hook` | — | `{"message":"hook"}` |

### SP — Session Profile

#### Guest Routes (`/rest/guest`) — tanpa auth

| Method | Endpoint | Handler | Keterangan |
|--------|----------|---------|------------|
| GET | `/rest/guest/` | — | `{"message":"guest"}` |
| GET | `/rest/guest/SP00` | `SP00.ListHrisCompany` | List HRIS companies |
| POST | `/rest/guest/SP00` | `SP00.Login` | Login |

#### Pages Routes (`/rest/pages`) — USLoad() required

| Method | Endpoint | Handler | Keterangan |
|--------|----------|---------|------------|
| GET | `/rest/pages/` | — | `{"message":"pages"}` |
| DELETE | `/rest/pages/SP00` | `SP00.Logout` | Logout |
| GET | `/rest/pages/SP01/company` | `SP01.ListUserCompany` | User companies |
| GET | `/rest/pages/SP01/module` | `SP01.ListUserModule` | User module tree (query: `company_id`) |
| PUT | `/rest/pages/SP02` | `SP02.ChangePassword` | Change password |
| GET | `/rest/pages/SP03` | `SP03.ListActions` | User action audit log |

### SM — System Manager (Admin Only)

> Semua route SM di bawah ini menggunakan `USLoad()` + `USLock()` + `USLogs()` untuk operasi mutation.

#### SM01 — User Management

| Method | Endpoint | Handler | Audit | Keterangan |
|--------|----------|---------|-------|------------|
| GET | `/rest/pages/SM01` | `SM01.ListUser` | — | List users (query: `search`, `page`, `size`, `sort_by`, `sort_order`) |
| POST | `/rest/pages/SM01` | `SM01.CreateUser` | SM01 | Create user |
| PUT | `/rest/pages/SM01/:id` | `SM01.UpdateUser` | SM01 | Update user |
| GET | `/rest/pages/SM01/hris` | `SM01.ListHRISCompany` | — | List HRIS companies |
| GET | `/rest/pages/SM01/company` | `SM01.ListAllCompany` | — | List all companies |
| GET | `/rest/pages/SM01/module` | `SM01.ListAllModule` | — | List all modules |
| POST | `/rest/pages/SM01/:id/company` | `SM01.AssignCompany` | SM01 | Assign company to user (auto-create) |
| GET | `/rest/pages/SM01/:id/company` | `SM01.ListUserCompany` | — | List user companies |
| POST | `/rest/pages/SM01/:id/company/assign` | `SM01.CreateUserCompany` | SM01 | Create/assign user company |
| PUT | `/rest/pages/SM01/:id/company/:companyId` | `SM01.UpdateUserCompany` | SM01 | Update user company (is_active) |
| DELETE | `/rest/pages/SM01/:id/company/:companyId` | `SM01.DeleteUserCompany` | SM01 | Remove user company |
| GET | `/rest/pages/SM01/:id/privilege` | `SM01.ListUserPrivilege` | — | List user privileges |
| POST | `/rest/pages/SM01/:id/privilege` | `SM01.CreateUserPrivilege` | SM01 | Assign privilege |
| PUT | `/rest/pages/SM01/:id/privilege/:privilegeId` | `SM01.UpdateUserPrivilege` | SM01 | Update privilege |
| DELETE | `/rest/pages/SM01/:id/privilege/:privilegeId` | `SM01.DeleteUserPrivilege` | SM01 | Remove privilege |
| GET | `/rest/pages/SM01/type` | `SM01.ListLocationTypeByCompany` | — | List location types (query: `company_id`) |
| GET | `/rest/pages/SM01/:id/location` | `SM01.ListUserLocation` | — | List user locations |
| POST | `/rest/pages/SM01/:id/location` | `SM01.CreateUserLocation` | SM01 | Assign user location |
| PUT | `/rest/pages/SM01/:id/location/:locationId` | `SM01.UpdateUserLocation` | SM01 | Update user location |
| DELETE | `/rest/pages/SM01/:id/location/:locationId` | `SM01.DeleteUserLocation` | SM01 | Remove user location |

#### SM02 — Module Management

| Method | Endpoint | Handler | Audit | Keterangan |
|--------|----------|---------|-------|------------|
| GET | `/rest/pages/SM02` | `SM02.ListModule` | — | List modules (query: `search`, `page`, `size`, `sort_by`, `sort_order`) |
| POST | `/rest/pages/SM02` | `SM02.CreateModule` | SM02 | Create module |
| PUT | `/rest/pages/SM02/:id` | `SM02.UpdateModule` | SM02 | Update module |

#### SM03 — Company Management

| Method | Endpoint | Handler | Audit | Keterangan |
|--------|----------|---------|-------|------------|
| GET | `/rest/pages/SM03` | `SM03.ListCompany` | — | List companies (query: `search`, `page`, `size`, `sort_by`, `sort_order`) |
| POST | `/rest/pages/SM03` | `SM03.CreateCompany` | SM03 | Create company |
| PUT | `/rest/pages/SM03/:id` | `SM03.UpdateCompany` | SM03 | Update company |
| GET | `/rest/pages/SM03/module` | `SM03.ListModule` | — | List available modules (child modules only) |
| GET | `/rest/pages/SM03/:id/module` | `SM03.ListCompanyModule` | — | List company modules |
| POST | `/rest/pages/SM03/:id/module` | `SM03.AssignCompanyModule` | SM03 | Assign module to company |
| PUT | `/rest/pages/SM03/:id/module/:moduleId` | `SM03.UpdateCompanyModule` | SM03 | Update company module (is_active) |
| GET | `/rest/pages/SM03/:id/type` | `SM03.ListLocationType` | — | List location types |
| POST | `/rest/pages/SM03/:id/type` | `SM03.CreateLocationType` | SM03 | Create location type |
| PUT | `/rest/pages/SM03/:id/type/:typeId` | `SM03.UpdateLocationType` | SM03 | Update location type |
| DELETE | `/rest/pages/SM03/:id/type/:typeId` | `SM03.DeleteLocationType` | SM03 | Delete location type |

#### SM04 — Signature/Approval Type Management

| Method | Endpoint | Handler | Audit | Keterangan |
|--------|----------|---------|-------|------------|
| GET | `/rest/pages/SM04` | `SM04.ListSignatureType` | — | List signature types (dengan steps dan signs nested) |
| POST | `/rest/pages/SM04` | `SM04.CreateSignatureType` | SM04 | Create signature type + steps + signs |
| PUT | `/rest/pages/SM04/:id` | `SM04.UpdateSignatureType` | SM04 | Update signature type (full replace strategy) |
| GET | `/rest/pages/SM04/user` | `SM04.ListUser` | — | List users for signer selection (query: `company_id`, `search`) |

#### SM05 — Session Management

| Method | Endpoint | Handler | Audit | Keterangan |
|--------|----------|---------|-------|------------|
| GET | `/rest/pages/SM05` | `SM05.ListSession` | — | List sessions (query: `search`, `page`, `size`, `sort_by`, `sort_order`) |
| GET | `/rest/pages/SM05/:id` | `SM05.GetSessionDetail` | — | Detail session (full token, user info) |
| DELETE | `/rest/pages/SM05/:id` | `SM05.RevokeSession` | SM05 | Revoke session |

### Database Tables yang Dirujuk

| Table | Digunakan di |
|-------|-------------|
| `dat_user` | SP00, SP02, SM01, SM04, SM05, backbone/memory |
| `dat_user_session` | SP00, SM05, backbone/memory, backbone/cleanup |
| `dat_company` | SP00, SP01, SM01, SM03 |
| `dat_module` | SP01, SM01, SM02, SM03 |
| `dat_user_company` | SM01, SP01, backbone/memory |
| `dat_user_privilege` | SM01, SP01, backbone/memory |
| `dat_company_module` | SM03, SP01 |
| `dat_user_location` | SM01 |
| `dat_location_type` | SM01, SM03, SM06 |
| `dat_user_action` | SP03, backbone/memory |
| `dat_signature_type` | SM04 |
| `dat_approval_step` | SM04 |
| `dat_approval_sign` | SM04 |

### Utility: Upload Handler (Belum Diwire)

`backbone/upload.go` sudah diimplementasi tapi **belum didaftarkan** ke route. Endpoint yang direncanakan:

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/rest/pages/upload/:entityType/:entityId` | Upload file |

**Konfigurasi upload:**
- Max ukuran: 20MB
- Base path: `/root/files`
- Ekstensi yang diizinkan: `.pdf`, `.doc`, `.docx`, `.png`, `.jpg`, `.jpeg`, `.xlsx`, `.csv`, `.txt`
- Kategori: `documents`, `evidence`, `avatars`
- Penyimpanan: `/root/files/{category}/{id}/{uuid}_{timestamp}.{ext}`

---

## 5. API Endpoints yang Direncanakan

Modul-modul berikut belum diimplementasi dan akan dikembangkan di masa mendatang:

| Kode | Modul | Keterangan | Status |
|------|-------|------------|--------|
| LM | Location Management | Manajemen lokasi hierarki | ⬜ Direncanakan |
| AM | Asset Management | Manajemen aset IT | ⬜ Direncanakan |
| NK | Network Management | VLAN, subnet, route, NAT | ⬜ Direncanakan |
| MK | Mikrotik Device Management | Device, hotspot, DHCP, firewall, queues, port/VLAN groups, clients, access | ⬜ Direncanakan |
| DN | Domain Management | Domain & DNS records | ⬜ Direncanakan |
| IP | IP Public & Block Management | IP blocks, allocations, NAT mappings | ⬜ Direncanakan |
| VL | Vulnerability Management | Scanner, scan jobs, findings | ⬜ Direncanakan |
| IN | Incident Management | Insiden & response | ⬜ Direncanakan |
| SI | SIEM / Log Monitoring | Log sources, rules, alerts | ⬜ Direncanakan |
| RK | Risk Assessment | Risiko, assessment, treatment | ⬜ Direncanakan |
| CP | Compliance Monitoring | ISO 27001, kontrol, assessment | ⬜ Direncanakan |
| CA | CAPA | Corrective and Preventive Action | ⬜ Direncanakan |
| IA | Internal Audit | Program & instance audit | ⬜ Direncanakan |
| MR | Management Review | Review meeting & keputusan | ⬜ Direncanakan |
| QM | Quality Management System | QMS, KPI, supplier, pelatihan | ⬜ Direncanakan |
| CT | Certificate Management | Sertifikat SSL/TLS | ⬜ Direncanakan |
| DK | Docker Management | Host, container, compose, image, network, volume | ⬜ Direncanakan |
| DM | Document Management | Dokumen, editor, todos, approvals | ⬜ Direncanakan |
| SR | Service Request | Katalog layanan, request, approval | ⬜ Direncanakan |
| BL | Billing & Invoice | Tagihan, vendor, kategori | ⬜ Direncanakan |
| PM | Preventive Maintenance | Ruangan, perangkat, template, kalender | ⬜ Direncanakan |
| DB | Dashboard & Reporting | Overview, widget, report | ⬜ Direncanakan |
| WM | Web Monitoring | UptimeRobot monitors, incidents, SLA | ⬜ Direncanakan |
| WS | Web Security | Nginx logs, WAF, IP whitelist/blacklist | ⬜ Direncanakan |

### Generic Endpoints (Direncanakan)

| Method | Endpoint | Privilege | Keterangan |
|--------|----------|-----------|------------|
| POST | `/rest/pages/upload/:entityType/:entityId` | post | Upload file (handler sudah ada, belum diwire) |
| GET | `/rest/pages/files/:entityType/:entityId` | view | List file |
| DELETE | `/rest/pages/files/:id` | post | Hapus file |
| GET | `/rest/pages/audit/:entityType/:entityId` | view | Audit trail |
| POST | `/rest/pages/audit/:entityType/:entityId` | post | Tambah catatan |
| GET | `/rest/pages/notifications` | view | Notifikasi saya |
| PUT | `/rest/pages/notifications/:id/read` | view | Tandai dibaca |
| PUT | `/rest/pages/notifications/read-all` | view | Tandai semua dibaca |

---

## 6. Struktur File

### File yang Sudah Ada

| File | Keterangan |
|------|------------|
| `ict_rest/main.go` | Entry point, graceful shutdown |
| `ict_rest/Dockerfile` | Multi-stage build (golang:1.26.4-alpine → alpine:3.20) |
| `ict_rest/backbone/cleanup.go` | Session cleanup goroutine (setiap 1 jam) |
| `ict_rest/backbone/database.go` | PostgreSQL connection setup |
| `ict_rest/backbone/logger.go` | Zerolog init, RequestID & Logger middleware |
| `ict_rest/backbone/memory.go` | Auth middleware (USLoad, USLock, USRole, USLogs) |
| `ict_rest/backbone/recovery.go` | Custom panic recovery middleware |
| `ict_rest/backbone/routes.go` | Gin router, CORS, DI wiring, route registration |
| `ict_rest/backbone/upload.go` | File upload handler (belum diwire ke route) |
| `ict_rest/mechanic/crypto.go` | AES-GCM encrypt/decrypt |
| `ict_rest/mechanic/helper.go` | AppError type, Error() responder, constructor functions |
| `ict_rest/mechanic/typography.go` | Pagination types, parsing helpers, BuildMeta |

### File yang Perlu Dibuat untuk Modul Baru

Setiap modul baru memerlukan 4 file di `ict_rest/skeleton/{KODE}/`:

| File | Deskripsi |
|------|-----------|
| `template.go` | Struct domain (request/response) + interface Repository & UseCase |
| `repository.go` | Implementasi Repository dengan SQL mentah |
| `usecase.go` | Implementasi UseCase dengan logika bisnis |
| `handler.go` | HTTP handler dengan Gin |

### Catatan Arsitektur

- **Tanpa ORM** — seluruh SQL ditulis manual dengan `database/sql` dan `QueryRowContext`/`ExecContext`
- **Constructor naming** — `NRepo()`, `NCase()`, `NHand()` untuk setiap layer
- **Dependency injection** — setiap layer menerima interface dari layer sebelumnya via constructor
- **CORS** hanya mengizinkan `localhost:36666` dan `172.99.66.6:36666` (frontend app)
- **Session lifetime** — 24 jam (diset saat login). Cleanup berjalan setiap 1 jam
- **XX99** adalah scaffold kosong — template untuk membuat modul skeleton baru
