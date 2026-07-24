# Todo: Rancangan Database Schema OmniSight

## Prioritas Pengembangan

### Fase 1 — Fondasi (MVP)

Modul-modul yang menjadi dasar dan paling sering diintegrasikan oleh modul lain.

| # | Modul | File | Estimasi | Status |
|---|-------|------|----------|--------|
| 1 | **Location Management** | `dat_location.prisma` | 5 tabel | ✅ |
| 2 | **Asset/Inventory** | `ict_asset.prisma` | 5 tabel | ✅ |
| 3 | **Vulnerability** | `ict_vulnerability.prisma` | 4 tabel | ⬜ |

**Alasan prioritas:**
- `ict_location` menjadi foreign key untuk modul asset, docker, service_request
- `ict_asset` menjadi foreign key untuk modul vuln, compliance, incident, risk, certificate, SIEM
- `ict_vulnerability` adalah fitur inti DevSecOps

---

### Fase 2 — Core Security

Modul yang berfokus pada pengelolaan keamanan langsung.

| # | Modul | File | Estimasi | Status |
|---|-------|------|----------|--------|
| 3 | **Incident Management** | `ict_incident.prisma` | 4 tabel | ⬜ |
| 4 | **SIEM / Log Monitoring** | `ict_siem.prisma` | 4 tabel | ⬜ |

**Alasan prioritas:**
- Incident adalah workflow utama DevSecOps
- SIEM menghubungkan log existing (`ict_web_security`) dengan detection rules

---

### Fase 3 — Governance

Modul terkait tata kelola dan kepatuhan.

| # | Modul | File | Estimasi | Status |
|---|-------|------|----------|--------|
| 5 | **Risk Assessment** | `ict_risk.prisma` | 4 tabel | ⬜ |
| 6 | **Compliance Monitoring** | `ict_compliance.prisma` | 4 tabel | ⬜ |

**Alasan prioritas:**
- Risk dan Compliance saling terkait
- Membutuhkan data asset yang sudah stabil

---

### Fase 4 — Supporting

Modul pendukung dan integrasi.

| # | Modul | File | Estimasi | Status |
|---|-------|------|----------|--------|
| 7 | **Certificate Management** | `ict_certificate.prisma` | 3 tabel | ⬜ |
| 8 | **Dashboard & Reporting** | `ict_dashboard.prisma` | 5 tabel | ⬜ |

---

### Fase 5 — Document Management

Modul pengelolaan dokumen, SOP, dan todo list.

| # | Modul | File | Estimasi | Status |
|---|-------|------|----------|--------|
| 9 | **Document Management (DMS)** | `ict_docman.prisma` | 10 tabel | ⬜ |

**Alasan prioritas:**
- Mendukung seluruh kegiatan dokumentasi IT & DevSecOps
- Terintegrasi dengan approval workflow (`dat_signature`)
- Fitur publish untuk distribusi dokumen

---

### Fase 6 — Docker Management & Monitoring

Modul manajemen dan monitoring Docker (inspirasi Dockge + Beszel).

| # | Modul | File | Estimasi | Status |
|---|-------|------|----------|--------|
| 10 | **Docker Management** | `ict_docker.prisma` | 14 tabel | ⬜ |
| 11 | **Docker Agent** | `ict_auto/ict_docker_agent` | Go CLI | ⬜ |

**Alasan prioritas:**
- Monitoring container sebagai bagian dari infrastruktur DevSecOps
- Manajemen Compose stack terintegrasi
- Agent-based untuk monitoring multi-server

---

### Fase 7 — Service Request Management

Modul pengelolaan permintaan layanan IT dari user non-IT.

| # | Modul | File | Estimasi | Status |
|---|-------|------|----------|--------|
| 12 | **Service Request** | `ict_service_request.prisma` | 8 tabel | ⬜ |

**Alasan prioritas:**
- Portal untuk user non-IT request layanan (server, akses, aset, aplikasi)
- Workflow approval terintegrasi
- Tracking SLA dan fulfilled request

---

### Fase 8 — Billing & Invoice Management

Modul pengelolaan tagihan yang harus dibayarkan, termasuk upload invoice dari vendor.

| # | Modul | File | Estimasi | Status |
|---|-------|------|----------|--------|
| 13 | **Billing & Invoice** | `ict_billing.prisma` | 3 tabel | ⬜ |

**Alasan prioritas:**
- Tracking tagihan vendor (cloud, license, hardware, internet)
- Upload invoice untuk dokumentasi
- Integrasi dengan Service Request dan Location

---

### Fase 9 — Preventive Maintenance (WiFi & Zoom)

Modul pengelolaan preventive maintenance untuk perangkat WiFi dan Zoom Meeting berdasarkan ruangan, dengan checklist harian, upload foto, dan kalender pengecekan.

| # | Modul | File | Estimasi | Status |
|---|-------|------|----------|--------|
| 14 | **Preventive Maintenance** | `ict_preventive.prisma` | 8 tabel | ⬜ |

**Alasan prioritas:**
- Checklist harian untuk perangkat WiFi (AP, Router, Switch)
- Checklist harian untuk perangkat Zoom (Camera, Speaker, Mic, Display)
- Upload foto/laporan setiap pengecekan
- Kalender penjadwalan pengecekan
- Issue tracking untuk kerusakan

---

### Fase 10 — ISO Compliance & Quality Management

Modul terkait kepatuhan ISO 27001/9001, CAPA, Audit Internal, Management Review, dan Quality Management System.

| # | Modul | File | Estimasi | Status |
|---|-------|------|----------|--------|
| 16 | **Statement of Applicability (SoA)** | `ict_compliance.prisma` | 1 tabel | ✅ |
| 17 | **CAPA** | `ict_compliance.prisma` | 3 tabel | ✅ |
| 18 | **Internal Audit** | `ict_compliance.prisma` | 3 tabel | ✅ |
| 19 | **Management Review** | `ict_compliance.prisma` | 3 tabel | ✅ |
| 20 | **Quality Objectives & KPIs** | `ict_qms.prisma` | 2 tabel | ✅ |
| 21 | **Supplier Management** | `ict_qms.prisma` | 2 tabel | ✅ |
| 22 | **Training & Competence** | `ict_qms.prisma` | 2 tabel | ✅ |
| 23 | **Calibration Management** | `ict_qms.prisma` | 2 tabel | ✅ |
| 24 | **Customer Feedback** | `ict_qms.prisma` | 1 tabel | ✅ |

**Alasan prioritas:**
- Dukung sertifikasi ISO 27001 (Information Security) dan ISO 9001 (Quality Management)
- CAPA dan Internal Audit adalah proses inti ISO
- Management Review diperlukan untuk evaluasi berkala
- Supplier Management dan Training mendukung operasional QMS
- Calibration penting untuk industri manufaktur/lab

---

### Fase 11 — Generic & Shared Tables

Shared tables yang digunakan oleh semua modul.

| # | Modul | File | Estimasi | Status |
|---|-------|------|----------|--------|
| 25 | **Generic** | `ict_generic.prisma` | 4 tabel | ⬜ |

**Alasan prioritas:**
- Audit trail untuk semua aktivitas
- File attachment untuk semua modul
- Approval workflow untuk semua modul
- Notification untuk semua modul

---

## Total Estimasi

| Fase | Jumlah Tabel | File Schema |
|------|-------------|-------------|
| Fase 1 | 13 | 3 |
| Fase 2 | 8 | 2 |
| Fase 3 | 8 | 2 |
| Fase 4 | 8 | 2 |
| Fase 5 | 10 | 1 |
| Fase 6 | 14 | 1 |
| Fase 7 | 8 | 1 |
| Fase 8 | 3 | 1 |
| Fase 9 | 8 | 1 |
| Fase 10 | 21 | 2 |
| Fase 11 | 4 | 1 |
| **Total** | **105 tabel** | **17 file schema** |

> **Catatan:** Setelah konsolidasi, tabel duplikat digantikan oleh shared tables di `ict_generic.prisma`.

---

## Cross-Cutting: Privilege System

Sistem akses yang berlaku untuk **semua modul** menggunakan `dat_user_privilege` + `dat_action_type`.

### Level Akses

| Level | Keterangan | Contoh Aksi |
|-------|------------|-------------|
| `hide` | Sembunyikan modul | Tidak bisa melihat menu |
| `view` | Lihat saja (read-only) | Lihat daftar, detail, dashboard |
| `book` | Kelola data (CRUD) | Tambah, edit, hapus data |
| `post` | Publish/Escalate | Publish, approve, close insiden |

### Mapping Level per Modul

| Modul | `hide` | `view` | `book` | `post` |
|-------|--------|--------|--------|--------|
| **asset** | — | Lihat aset | Tambah/edit/hapus | Ubah criticality, decommission |
| **vulnerability** | — | Lihat temuan | Scan, update status | Mark FP, accept risk |
| **incident** | — | Lihat insiden | Buat/update | Contain, close, eskalasi |
| **compliance** | — | Lihat assessment | Buat assessment | Approve, publish laporan |
| **capa** | — | Lihat CAPA | Buat/update CAPA | Verifikasi efektivitas |
| **internal_audit** | — | Lihat program audit | Buat/update audit | Finalisasi temuan |
| **management_review** | — | Lihat review | Buat/update review | Finalisasi keputusan |
| **qms** | — | Lihat data kualitas | Tambah/update data | Publish laporan kualitas |
| **risk** | — | Lihat register | Tambah/update | Approve treatment |
| **certificate** | — | Lihat sertifikat | Tambah/edit | Trigger renewal |
| **docker** | — | Lihat container | Start/stop, edit compose | Hapus, manage host |
| **docman** | — | Baca dokumen | Buat/edit | Publish, approve |
| **siem** | — | Lihat alert | Tambah/edit rules | Resolve alert |
| **reporting** | — | Lihat laporan | Generate laporan | — |
| **setting** | — | Lihat pengaturan | Kelola user | Hapus user |

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

| Role | Deskripsi | Service Access |
|------|-----------|---------------|
| `requester` | User umum | Submit request, lihat status sendiri |
| `department_head` | Kepala divisi | Submit + approve department |
| `finance` | Staff finance | Lihat request biaya, approve budget |
| `hrd` | Staff HRD | Request aset untuk karyawan baru |
| `project_manager` | Project manager | Request server/aplikasi project |

### Mapping Role → Service

| Layanan | `requester` | `department_head` | `project_manager` | `hrd` | `finance` | `staff` |
|---------|------------|-------------------|-------------------|-------|-----------|---------|
| Request Server | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| Request Akses | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| Request Aplikasi | ✅ | ✅ | ✅ | ✅ | ❌ | ✅ |
| Request Aset | ✅ | ✅ | ❌ | ✅ | ❌ | ✅ |
| Approve Department | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| Approve Budget | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ |

### Task: Privilege System

- [x] Buat middleware `USRole` di `ict_rest/backbone/memory.go`
- [ ] Buat komponen frontend privilege check
- [ ] Buat halaman manajemen privilege di `ict_site`

Lihat detail lengkap di `ict_rest.md`.

---

## Task: Login & Authentication System (SP00)

### Backend (ict_rest)
- [x] Buat `skeleton/SP00/template.go` — Interfaces (Repository, UseCase) + DTOs
  - [x] `HrisCompanyItem`, `LoginRequest`, `HrisSoapResponse`, `LoginResponse`, `ProfileData` structs
  - [x] `Repository` interface — DB methods (ListHrisCompany, GetCompanyHrisLink, FindUserByUsername, FindUserByUsernameAndCompany, CreateUserSession, UpdateUserKey, DeleteSession)
  - [x] `UseCase` interface — Business methods (ListHrisCompany, GetCompanyHrisLink, Login, LoginHris, Logout)
  - [x] Captcha: client-side only (tidak ada captcha di server)
- [x] Buat `skeleton/SP00/repository.go` — Database queries (NRepo constructor)
  - [x] `ListHrisCompany()` — Query dat_company WHERE is_active=true AND hris_link IS NOT NULL
  - [x] `GetCompanyHrisLink(companyID)` — Query single company hris_link
  - [x] `FindUserByUsername(username)` — Query dat_user (includes password_hash, company_id, is_admin, is_hris, is_active)
  - [x] `FindUserByUsernameAndCompany(username, companyID)` — Query dat_user WHERE is_hris=true AND company_id=$2
  - [x] `CreateUserSession(userID, token, ip, ua, expiresAt)` — INSERT dat_user_session (mechanic.NullableString untuk ip/ua)
  - [x] `UpdateUserKey(id, key)` — UPDATE dat_user.key (encrypted HRIS password)
  - [x] `DeleteSession(token)` — DELETE dat_user_session
- [x] Buat `skeleton/SP00/usecase.go` — Business logic (NCase constructor)
  - [x] `GenerateToken()` — crypto/rand 32 bytes hex
  - [x] `Login(username, password, ip, ua)` — Non-HRIS: find user by username, check company_id is empty, bcrypt.CompareHashAndPassword, create session
  - [x] `LoginHris(username, password, companyID, ip, ua)` — HRIS: find user by username+companyID, check is_hris=true, mechanic.Encrypt(password) → UpdateUserKey, create session
  - [x] `Logout(token)` — DeleteSession
- [x] Buat `skeleton/SP00/handler.go` — HTTP handlers (NHand constructor)
  - [x] `ListHrisCompany()` — GET /rest/guest/SP00 → List of HRIS companies
  - [x] `Login()` — POST /rest/guest/SP00 → Bind LoginRequest, route non-HRIS (company_id="" || "Non-HRIS") vs HRIS
  - [x] HRIS flow: GetCompanyHrisLink, replace `{$U}` and `{$P}` in URL template, HTTP GET external, XML Unmarshal HrisSoapResponse, check prefix "1|", then LoginHris
  - [x] `Logout()` — DELETE /rest/pages/SP00 → Extract token from Authorization header, Logout
- [x] Update `backbone/routes.go` — DI wiring
  - [x] Import `ict_rest/skeleton/SP00`
  - [x] `SP00R := SP00.NRepo(PgSQL)` → `SP00U := SP00.NCase(SP00R)` → `SP00H := SP00.NHand(SP00U)`
  - [x] `guest.GET("/SP00", SP00H.ListHrisCompany)`, `guest.POST("/SP00", SP00H.Login)`, `pages.DELETE("/SP00", SP00H.Logout)`
- [x] Hapus `backbone/handler.go` (old monolithic handler)
- [x] Hapus `skeleton/SP00.go` (old monolithic file)
- [x] Build test — `go build ./...` ✅ Berhasil

### Frontend (ict_site)
- [x] Buat `src/app/proxy/guest/[...path]/route.ts` — Guest proxy route (catch-all)
  - [x] GET handler (proxy ke backend, forward query params)
  - [x] POST handler (proxy ke backend, forward JSON body)
- [x] Buat `src/app/login/page.tsx` — Login page
  - [x] react-hook-form + zod schema validation (company_id, username, password, captcha)
  - [x] Company select (fetch dari GET /proxy/guest/SP00) + "(Non-HRIS)" empty option
  - [x] Username input (with Field/FieldLabel/FieldError)
  - [x] Password input (type=password)
  - [x] Captcha: client-side 3 random digits, InputGroup with hint + refresh button
  - [x] Submit: client-side captcha validation → POST ke /proxy/guest/SP00
  - [x] Session storage: localStorage + cookie (storageKey = "OmniSightMemory")
  - [x] Redirect ke /board on success
  - [x] Error handling + refresh captcha on fail

### Fitur Login
- Dua metode login: HRIS dan Non-HRIS
- HRIS: Validasi via external XML API (`{$U}`/`{$P}` template), find user by username+companyID, encrypt password → dat_user.key via mechanic.Encrypt
- Non-HRIS: company_id="" → find user by username WHERE company_id='', bcrypt.CompareHashAndPassword
- Captcha: Client-side only (3 random digits 0-9, tidak ada server-side captcha store)
- Session: Bearer token (hex 32 bytes), 24 jam expiry, dat_user_session
- Storage: Cookie `OmniSightMemory` (client) + localStorage
- Arsitektur: Linear Layer Execution (template → repository → usecase → handler)
- DI: `SP00.NRepo(PgSQL)` → `SP00.NCase(repo)` → `SP00.NHand(uc)`

Lihat detail lengkap di `omnisight.md`.

---

## Task: Session Profile / Board Data (SP01)

### Backend (ict_rest)
- [x] Buat `skeleton/SP01/template.go` — Interfaces (Repository, UseCase) + DTOs
  - [x] `UserCompanyItem` struct (ID, CompanyID, Name)
  - [x] `ModuleTreeNode` struct (ID, Code, Name, Path, IsPage, Children — recursive)
  - [x] `Repository` interface — `ListUserCompany`, `ListAllModuleTree`, `ListUserModule`
  - [x] `UseCase` interface — same methods
- [x] Buat `skeleton/SP01/repository.go` — Database queries (NRepo constructor)
  - [x] `ListUserCompany(userID)` — Query dat_company WHERE is_active=true, join dat_user_company OR dat_user.company_id
  - [x] `ListAllModuleTree()` — Query dat_module, build tree from flat rows (parent_id → children)
  - [x] `ListUserModule(userID, companyID)` — Recursive CTE: dat_user_company → dat_company_module → dat_user_privilege → dat_module, filter level <> 'hide', aggregate max privilege level
- [x] Buat `skeleton/SP01/usecase.go` — Business logic (NCase constructor)
  - [x] `ListUserCompany(userID)` — Delegate to repo
  - [x] `ListAllModuleTree()` — Delegate to repo
  - [x] `ListUserModule(userID, companyID)` — Return empty if companyID empty, else delegate to repo
- [x] Buat `skeleton/SP01/handler.go` — HTTP handlers (NHand constructor)
  - [x] `ListUserCompany()` — GET /rest/pages/SP01/company, extract userId from session (c.Get("userId"))
  - [x] `ListUserModule()` — GET /rest/pages/SP01/module?company_id=, admin bypass (isAdmin → ListAllModuleTree), non-admin → ListUserModule
- [x] Update `backbone/routes.go` — DI wiring
  - [x] Import `ict_rest/skeleton/SP01`
  - [x] `SP01R := SP01.NRepo(PgSQL)` → `SP01U := SP01.NCase(SP01R)` → `SP01H := SP01.NHand(SP01U)`
  - [x] `pages.GET("/SP01/company", SP01H.ListUserCompany)`, `pages.GET("/SP01/module", SP01H.ListUserModule)`
- [x] Build test — `go build ./...` ✅ Berhasil

### Fitur SP01
- Multi-tenant: Query user company via dat_user_company junction table + dat_user.company_id
- Module tree: Recursive query dengan parent_id hierarchy, membangun tree di Go
- Privilege filtering: Recursive CTE menggabungkan dat_user_privilege (level: hide/view/book/post), aggregate max level per module
- Admin bypass: is_admin=true → tampilkan semua module tanpa filter privilege
- Endpoint protected: Bearer token required (via USLoad middleware)

---

## Task: Change Password (SP02)

### Backend (ict_rest)
- [x] Buat `skeleton/SP02/template.go` — Interfaces + DTOs
  - [x] `Repository` interface — `FindUserPasswordHash(id)`, `UpdateUserPassword(id, hashedPassword)`
  - [x] `UseCase` interface — `ChangePassword(userID, currentPassword, newPassword)`
- [x] Buat `skeleton/SP02/repository.go` — Database queries (NRepo constructor)
  - [x] `FindUserPasswordHash(id)` — Query dat_user.password by id
  - [x] `UpdateUserPassword(id, hashedPassword)` — UPDATE dat_user.password + updated_at
- [x] Buat `skeleton/SP02/usecase.go` — Business logic (NCase constructor)
  - [x] `ChangePassword` — Validasi input (userID, currentPassword, newPassword tidak kosong)
  - [x] bcrypt.CompareHashAndPassword untuk verify current password
  - [x] bcrypt.GenerateFromPassword untuk hash new password
- [x] Buat `skeleton/SP02/handler.go` — HTTP handlers (NHand constructor)
  - [x] `ChangePassword()` — PUT /rest/pages/SP02, extract userId dari session, bind JSON {current_password, new_password}
- [x] Update `backbone/routes.go` — DI wiring
  - [x] `SP02R := SP02.NRepo(PgSQL)` → `SP02U := SP02.NCase(SP02R)` → `SP02H := SP02.NHand(SP02U)`
  - [x] `pages.PUT("/SP02", SP02H.ChangePassword)`
- [x] Build test — `go build ./...` ✅ Berhasil

### Frontend (ict_site)
- [x] Buat `src/app/board/pages/SP02/page.tsx` — Change password page
  - [x] react-hook-form + zod schema (current_password, new_password, confirm_password)
  - [x] Zod validation: min 6 chars, uppercase, lowercase, number + confirm match
  - [x] HRIS check: Jika is_hris, tampilkan pesan "cannot change password"
  - [x] Submit: PUT /proxy/pages/SP02 dengan Bearer token
  - [x] Success/error alerts (Alert component)
  - [x] Session check: parseSession + isSessionExpired → redirect ke /login

### Fitur SP02
- Hanya untuk non-HRIS users (HRIS password dikelola oleh sistem HRIS)
- Current password diverifikasi dengan bcrypt sebelum update
- New password di-hash dengan bcrypt seimpan ke database
- Endpoint protected: Bearer token required (via USLoad middleware)

---

## Task: User Action History (SP03)

### Backend (ict_rest)
- [x] Buat `skeleton/SP03/template.go` — Interfaces + DTOs
  - [x] `UserActionItem` DTO — id, user_id, company_id, module_code, action, path, ip_address, user_agent, created_at
  - [x] `Repository` interface — `ListActions(userID, search, page, size, sortBy, sortOrder)` → `([]UserActionItem, GridMeta, error)`
  - [x] `UseCase` interface — same signature
- [x] Buat `skeleton/SP03/repository.go` — Database queries (NRepo constructor)
  - [x] `ListActions` — Dynamic WHERE (user_id filter, ILIKE search pada module_code/action/path/ip_address)
  - [x] Valid sort columns: created_at, action, module_code, ip_address (default: created_at desc)
  - [x] Pagination: LIMIT/OFFSET dengan COUNT(*) total → mechanic.GridMeta
- [x] Buat `skeleton/SP03/usecase.go` — Business logic (NCase constructor)
  - [x] `ListActions` — mechanic.CheckMeta untuk validasi page/size, delegate ke repo
- [x] Buat `skeleton/SP03/handler.go` — HTTP handlers (NHand constructor)
  - [x] `ListActions()` — GET /rest/pages/SP03, bind mechanic.ActionMeta dari query params, response {data, meta}
- [x] Update `backbone/routes.go` — DI wiring
  - [x] `SP03R := SP03.NRepo(PgSQL)` → `SP03U := SP03.NCase(SP03R)` → `SP03H := SP03.NHand(SP03U)`
  - [x] `pages.GET("/SP03", SP03H.ListActions)`
- [x] Build test — `go build ./...` ✅ Berhasil

### Frontend (ict_site)
- [x] Buat `src/app/board/pages/SP03/page.tsx` — User action history page
  - [x] DataTable dengan columns: Module, Action, Path, IP Address, Created At (hidden: id, user_id, company_id, user_agent)
  - [x] Sorting: module_code, action, ip_address, created_at
  - [x] Search + pagination (via DataTable built-in)
  - [x] DataDialog untuk detail view (read-only, no create/update/delete)
  - [x] Fetch: GET /proxy/pages/SP03 dengan Bearer token
  - [x] Session check: 401 → redirect ke /login

### Fitur SP03
- Menampilkan history aktivitas user (module_code, action, path, ip_address, user_agent)
- Filter per user (userID dari query param)
- Search: ILIKE pada module_code, action, path, ip_address
- Sorting: created_at (default), action, module_code, ip_address
- Pagination: page/size dengan GridMeta response
- Endpoint protected: Bearer token required (via USLoad middleware)

---

## Task: User Management (SM01)

### Backend (ict_rest)
- [x] Buat `skeleton/SM01/template.go` — Interfaces + DTOs
  - [x] `UserListItem` DTO — id, username, name, email, phone, is_hris, is_active, is_admin, company_id, company_name
  - [x] `UserCompanyItem` DTO — id, user_id, company_id, company_name, is_active
  - [x] `UserPrivilegeItem` DTO — id, module_id, module_name, module_code, level
  - [x] `HrisCompanyItem` DTO — id, company_name
  - [x] `Repository` interface — `ListUsers`, `FindUser`, `CreateUser`, `UpdateUser`, `DeleteUser`, `ListUserCompanies`, `AssignCompany`, `RemoveCompany`, `ListUserPrivileges`, `AssignPrivilege`, `UpdatePrivilege`, `RemovePrivilege`, `ListHrisCompanies`, `SearchModules`, `ListCompanies`
  - [x] `UseCase` interface — same methods
- [x] Buat `skeleton/SM01/repository.go` — Database queries (NRepo constructor, 16 methods)
- [x] Buat `skeleton/SM01/usecase.go` — Business logic (NCase constructor)
  - [x] `CreateUser` — bcrypt password hash sebelum insert
  - [x] `UpdateUser` — bcrypt password hash jika password baru tidak kosong
  - [x] `ListUserCompanies` — JOIN dat_company untuk nama
  - [x] `ListUserPrivileges` — JOIN dat_module untuk nama + code
- [x] Buat `skeleton/SM01/handler.go` — HTTP handlers (NHand constructor, 13 endpoints)
  - [x] `ListUsers()` — GET /rest/pages/admin/SM01 (search, page, size, sort)
  - [x] `FindUser()` — GET /rest/pages/admin/SM01/:id
  - [x] `CreateUser()` — POST /rest/pages/admin/SM01 (username, name, password, is_admin, is_hris, is_active)
  - [x] `UpdateUser()` — PUT /rest/pages/admin/SM01/:id
  - [x] `DeleteUser()` — DELETE /rest/pages/admin/SM01/:id
  - [x] `ListUserCompanies()` — GET /rest/pages/admin/SM01/:id/company
  - [x] `AssignCompany()` — POST /rest/pages/admin/SM01/:id/company
  - [x] `RemoveCompany()` — DELETE /rest/pages/admin/SM01/:id/company/:companyId
  - [x] `ListUserPrivileges()` — GET /rest/pages/admin/SM01/:id/privilege
  - [x] `AssignPrivilege()` — POST /rest/pages/admin/SM01/:id/privilege
  - [x] `UpdatePrivilege()` — PUT /rest/pages/admin/SM01/:id/privilege/:privilegeId
  - [x] `RemovePrivilege()` — DELETE /rest/pages/admin/SM01/:id/privilege/:privilegeId
  - [x] `ListHrisCompanies()` — GET /rest/pages/admin/SM01/hris-companies
- [x] Update `backbone/routes.go` — DI wiring
  - [x] `SM01R := SM01.NRepo(PgSQL)` → `SM01U := SM01.NCase(SM01R)` → `SM01H := SM01.NHand(SM01U)`
  - [x] `admin.GET("/SM01", ...)` — List, Find, ListHrisCompanies
  - [x] `admin.POST("/SM01", ...)` — Create
  - [x] `admin.PUT("/SM01/:id", ...)` — Update
  - [x] `admin.DELETE("/SM01/:id", ...)` — Delete
  - [x] Sub-routes: `/SM01/:id/company`, `/SM01/:id/privilege`
- [x] Build test — `go build ./...` ✅ Berhasil

### Frontend (ict_site)
- [x] Buat `src/app/board/pages/SM01/page.tsx` — User management page
  - [x] DataTable: Username, Name, Email, Phone, HRIS badge, Active badge, Admin badge
  - [x] Create DataDialog: react-hook-form + zod (username, name, email, phone, password, is_admin, is_active)
  - [x] Edit DataDialog: react-hook-form + zod (semua field kecuali password)
  - [x] Delete confirmation dialog
  - [x] Tabs: Information (detail), Companies (assign/remove company), Privilege (assign/update/remove privilege)
  - [x] Company sub-dialog: select dari HRIS companies, assign ke user
  - [x] Privilege sub-dialog: select module + level (hide/view/book/post)
  - [x] Fetch: GET/POST/PUT/DELETE /proxy/pages/admin/SM01

### Fitur SM01
- CRUD user dengan role admin
- Assign/remove user ke multiple companies (multi-tenant)
- Assign/update/remove user privileges per module
- HRIS badge: user HRIS tidak bisa edit password
- Endpoint protected: Bearer token + USLock (admin session)

---

## Task: Module Management (SM02)

### Backend (ict_rest)
- [x] Buat `skeleton/SM02/template.go` — Interfaces + DTOs
  - [x] `ModuleListItem` DTO — id, code, name, path, is_page, parent_id, parent_name, children
  - [x] `ModuleCreateRequest` DTO — code, name, path, is_page, parent_id
  - [x] `ModuleUpdateRequest` DTO — same fields
  - [x] `Repository` interface — `ListModules`, `FindModule`, `CreateModule`, `UpdateModule`
  - [x] `UseCase` interface — same methods
- [x] Buat `skeleton/SM02/repository.go` — Database queries (NRepo constructor)
  - [x] `ListModules()` — Query dat_module, build tree dari flat rows
  - [x] `FindModule()` — Query by ID
  - [x] `CreateModule()` — INSERT dat_module
  - [x] `UpdateModule()` — UPDATE dat_module
- [x] Buat `skeleton/SM02/usecase.go` — Business logic (NCase constructor)
  - [x] `CreateModule` — Validasi code unik, parent exist jika parent_id != null
  - [x] `UpdateModule` — Validasi code unik (kecuali diri sendiri)
- [x] Buat `skeleton/SM02/handler.go` — HTTP handlers (NHand constructor, 3 endpoints)
  - [x] `ListModules()` — GET /rest/pages/admin/SM02
  - [x] `CreateModule()` — POST /rest/pages/admin/SM02
  - [x] `UpdateModule()` — PUT /rest/pages/admin/SM02/:id
- [x] Update `backbone/routes.go` — DI wiring
  - [x] `SM02R := SM02.NRepo(PgSQL)` → `SM02U := SM02.NCase(SM02R)` → `SM02H := SM02.NHand(SM02U)`
- [x] Build test — `go build ./...` ✅ Berhasil

### Frontend (ict_site)
- [x] Buat `src/app/board/pages/SM02/page.tsx` — Module management page
  - [x] DataTable: Code, Name, Path, Page badge, Parent name
  - [x] Create DataDialog: react-hook-form + zod (code, name, path, is_page, parent_id select)
  - [x] Edit DataDialog: react-hook-form + zod
  - [x] Tree view display (parent-child hierarchy)
  - [x] Parent select: select existing modules sebagai parent
  - [x] Fetch: GET/POST/PUT /proxy/pages/admin/SM02

### Fitur SM02
- CRUD modules dengan hierarchy parent-child
- Tree view display
- Parent module select untuk assign hierarchy
- is_page flag untuk module frontend pages
- Endpoint protected: Bearer token + USLock (admin session)

---

## Task: Company Management (SM03)

### Backend (ict_rest)
- [x] Buat `skeleton/SM03/template.go` — Interfaces + DTOs
  - [x] `CompanyListItem` DTO — id, code, name, hris_link, is_active, currency, timezone, phone, email, website, address, city, state, country, postal_code
  - [x] `CompanyModuleItem` DTO — id, module_id, module_name, module_code
  - [x] `CompanyCreateRequest` DTO — code, name, hris_link, is_active, currency, timezone, phone, email, website, address, city, state, country, postal_code
  - [x] `CompanyModuleAssignRequest` DTO — module_id
  - [x] `Repository` interface — `ListCompanies`, `FindCompany`, `CreateCompany`, `UpdateCompany`, `DeleteCompany`, `ListCompanyModules`, `AssignModule`, `RemoveModule`
  - [x] `UseCase` interface — same methods
- [x] Buat `skeleton/SM03/repository.go` — Database queries (NRepo constructor, 8 methods)
- [x] Buat `skeleton/SM03/usecase.go` — Business logic (NCase constructor)
  - [x] `CreateCompany` — Validasi code unik
  - [x] `UpdateCompany` — Validasi code unik (kecuali diri sendiri)
  - [x] `ListCompanyModules` — JOIN dat_module
- [x] Buat `skeleton/SM03/handler.go` — HTTP handlers (NHand constructor, 7 endpoints)
  - [x] `ListCompanies()` — GET /rest/pages/admin/SM03
  - [x] `FindCompany()` — GET /rest/pages/admin/SM03/:id
  - [x] `CreateCompany()` — POST /rest/pages/admin/SM03
  - [x] `UpdateCompany()` — PUT /rest/pages/admin/SM03/:id
  - [x] `DeleteCompany()` — DELETE /rest/pages/admin/SM03/:id
  - [x] `ListCompanyModules()` — GET /rest/pages/admin/SM03/:id/module
  - [x] `AssignModule()` — POST /rest/pages/admin/SM03/:id/module
  - [x] `RemoveModule()` — DELETE /rest/pages/admin/SM03/:id/module/:moduleId
- [x] Update `backbone/routes.go` — DI wiring
  - [x] `SM03R := SM03.NRepo(PgSQL)` → `SM03U := SM03.NCase(SM03R)` → `SM03H := SM03.NHand(SM03U)`
- [x] Build test — `go build ./...` ✅ Berhasil

### Frontend (ict_site)
- [x] Buat `src/app/board/pages/SM03/page.tsx` — Company management page
  - [x] DataTable: Code, Name, HRIS Link, Active badge, Currency
  - [x] Create DataDialog: react-hook-form + zod (code, name, hris_link, is_active, currency, timezone, phone, email, website, address, city, state, country, postal_code)
  - [x] Edit DataDialog: react-hook-form + zod
  - [x] Delete confirmation dialog
  - [x] Tabs: Information (detail), Modules (assign/remove modules)
  - [x] Module sub-dialog: select modules, assign ke company
  - [x] Valuta select: IDR, USD, EUR, SGD, JPY, AUD
  - [x] Fetch: GET/POST/PUT/DELETE /proxy/pages/admin/SM03

### Fitur SM03
- CRUD company dengan admin
- Multi-tenant support: company_code, hris_link, timezone, currency
- Module assignment: pilih modules yang tersedia untuk company
- Endpoint protected: Bearer token + USLock (admin session)

---

## Task: Document Type & Approval (SM04)

### Backend (ict_rest)
- [x] Buat `skeleton/SM04/template.go` — Interfaces + DTOs
  - [x] `SignatureTypeItem` DTO — id, company_id, company_name, name, description, category, is_active, approval_steps (children)
  - [x] `ApprovalStepI` DTO — id, name, step_order, is_required, step_type, category
  - [x] `ApprovalSignI` DTO — id, name, email
  - [x] `SignatureTypeCreateRequest` DTO — company_id, name, description, category, is_active, approval_steps
  - [x] `ApprovalStepCreateRequest` DTO — name, step_order, is_required, step_type, category
  - [x] `Repository` interface — `ListSignatureTypes`, `FindSignatureType`, `CreateSignatureType`, `UpdateSignatureType`, `DeleteSignatureType`, `ListApprovalSteps`, `CreateApprovalStep`, `UpdateApprovalStep`, `DeleteApprovalStep`, `ListApprovalSigns`, `SearchUsers`
  - [x] `UseCase` interface — same methods
- [x] Buat `skeleton/SM04/repository.go` — Database queries (NRepo constructor, 11 methods)
  - [x] `ListSignatureTypes` — JOIN dat_company + nested approval_steps
  - [x] `CreateSignatureType` — Transaction: INSERT dat_document_type + INSERT dat_document_type_approval (steps)
  - [x] `UpdateSignatureType` — Transaction: UPDATE + DELETE + re-INSERT approval steps
  - [x] `ListApprovalSteps` — Query by signature_type_id
  - [x] `ListApprovalSigns` — Query approval signs by step_id
  - [x] `SearchUsers` — Query dat_user WHERE name ILIKE or email ILIKE
- [x] Buat `skeleton/SM04/usecase.go` — Business logic (NCase constructor)
  - [x] All CRUD methods delegate ke repo
- [x] Buat `skeleton/SM04/handler.go` — HTTP handlers (NHand constructor, 4 endpoints)
  - [x] `ListSignatureTypes()` — GET /rest/pages/admin/SM04
  - [x] `CreateSignatureType()` — POST /rest/pages/admin/SM04
  - [x] `UpdateSignatureType()` — PUT /rest/pages/admin/SM04/:id
  - [x] `ListUsers()` — GET /rest/pages/admin/SM04/users (untuk assign signer)
- [x] Update `backbone/routes.go` — DI wiring
  - [x] `SM04R := SM04.NRepo(PgSQL)` → `SM04U := SM04.NCase(SM04R)` → `SM04H := SM04.NHand(SM04U)`
- [x] Build test — `go build ./...` ✅ Berhasil

### Frontend (ict_site)
- [x] Buat `src/app/board/pages/SM04/page.tsx` — Document type & approval management
  - [x] DataTable: Name, Description, Category, Company, Active badge
  - [x] Create DataDialog: react-hook-form + zod (company_id, name, description, category, is_active, dynamic approval steps)
  - [x] Edit DataDialog: react-hook-form + zod
  - [x] useFieldArray untuk dynamic approval steps (add/remove/reorder)
  - [x] Approval step fields: name, step_order, is_required, step_type, category
  - [x] Fetch: GET/POST/PUT /proxy/pages/admin/SM04

### Fitur SM04
- CRUD document type dengan approval workflow
- Dynamic approval steps (useFieldArray)
- Company assignment per document type
- Category-based classification (approval_type, signature_type, etc.)
- Endpoint protected: Bearer token + USLock (admin session)

---

---

## Task: Board Layout & Dashboard (Frontend)

### Frontend (ict_site)
- [x] Buat `src/app/board/layout.tsx` — Board layout
  - [x] SidebarProvider (defaultOpen dari sidebar_collapsible preference)
  - [x] AppSidebar (variant + collapsible dari preference store)
  - [x] Header (SidebarTrigger + Separator + BoardBreadcrumb)
  - [x] Content area (p-4 md:p-6, overflow handling)
- [x] Buat `src/app/board/page.tsx` — Board dashboard page
  - [x] Renders `<Board />` component
- [x] Buat `src/app/board/[...not-found]/page.tsx` — Catch-all not-found
  - [x] UnderDevelopment placeholder
- [x] Buat `src/app/board/model/module.ts` — Navigation model
  - [x] Types: `NavGroup`, `NavMainItem`, `NavMainLinkItem`, `NavMainParentItem`, `NavSubItem`, `NavBadge`
  - [x] `useModuleItem()` hook — fetch dari `/proxy/pages/SP01/module?company_id=`, build NavGroups
  - [x] Static items: `staticAdminItem` (SM01-SM05 System Manager), `staticSessionItem` (SP01-SP03 Session Profile)
  - [x] `mapIconByCode()` — Icon mapping berdasarkan module code prefix (NW, SM, SP, XX)
  - [x] `buildNavGroups()` — Convert ModuleTreeNode → NavMainItem, dedupe, append admin/session items
- [x] Buat `src/app/board/widget/sidebar.tsx` — AppSidebar
  - [x] Header: Logo + WS_CONF.name (link ke /board)
  - [x] Content: CompanyCombobox + NavMain
  - [x] Footer: User avatar dropdown (Profile, Password jika non-HRIS, History) + Logout
  - [x] Session check: parseSession + isSessionExpired → forceLogout
  - [x] Company validation: Fetch SP01/company on mount, forceLogout jika 401/403
- [x] Buat `src/app/board/widget/breadcrumb.tsx` — BoardBreadcrumb
  - [x] Dynamic breadcrumb dari pathname + NavGroups
  - [x] Fallback: Last path segment di-format sebagai label
- [x] Buat `src/app/board/widget/company.tsx` — CompanyCombobox
  - [x] Fetch dari `/proxy/pages/SP01/company` dengan Bearer token
  - [x] Zustand store: companyId, companyList dari usePreferencesStore
  - [x] Locked company: Jika session.company_id ada, kunci pilihan (tidak bisa diganti)
  - [x] upsertSessionCompanyName: Update localStorage session dengan nama company terpilih
- [x] Buat `src/app/board/widget/dashboard.tsx` — Board component
  - [x] User profile display (read-only Input fields)
  - [x] Filter: Sembunyikan field internal (password, company_id, dll) dan boolean status fields
  - [x] Status badges: is_hris (emerald/red badge)
  - [x] Session check: parseSession + isSessionExpired → forceLogout
- [x] Buat `src/app/board/widget/module.tsx` — NavMain
  - [x] `NavItem` — Route ke NavLinkItem (single link) atau NavDropdownItem/NavCollapsibleItem (dengan subItems)
  - [x] `NavLinkItem` — SidebarMenuButton dengan Link
  - [x] `NavDropdownItem` — DropdownMenu saat sidebar collapsed (desktop)
  - [x] `NavCollapsibleItem` — Collapsible tree saat sidebar expanded
  - [x] `CollapsedIconFallback` — Initial letter icon saat collapsed
  - [x] `NavItemBadge` — Badge "new"/"soon" pada nav items

---

## Task: SM01 Location Management (Backend + Frontend)

### Backend (ict_rest)
- [x] Buat `skeleton/SM01/template.go` — User location DTOs
  - [x] `UserLocationListItem` DTO — id, location_type_id, type_code, type_name, is_active, company_name
  - [x] `UserLocationCreateRequest` DTO — company_id, location_type_id, is_active
  - [x] `UserLocationUpdateRequest` DTO — location_type_id, is_active
  - [x] `LocationTypeSelectItem` DTO — id, code, name
- [x] Buat `skeleton/SM01/repository.go` — User location queries
  - [x] `ListUserLocation(userID)` — Query dat_user_location JOIN dat_location_type + dat_company
  - [x] `CreateUserLocation(userID, req)` — INSERT dat_user_location ON CONFLICT UPDATE
  - [x] `UpdateUserLocation(userID, locationID, req)` — UPDATE dat_user_location
  - [x] `DeleteUserLocation(userID, locationID)` — DELETE dat_user_location
  - [x] `ListLocationTypeByCompany(companyID)` — Query dat_location_type WHERE company_id
- [x] Buat `skeleton/SM01/usecase.go` — User location business logic
  - [x] All CRUD methods delegate ke repo
- [x] Buat `skeleton/SM01/handler.go` — User location HTTP handlers
  - [x] `ListUserLocation()` — GET /rest/pages/admin/SM01/:id/location
  - [x] `CreateUserLocation()` — POST /rest/pages/admin/SM01/:id/location
  - [x] `UpdateUserLocation()` — PUT /rest/pages/admin/SM01/:id/location/:locationId
  - [x] `DeleteUserLocation()` — DELETE /rest/pages/admin/SM01/:id/location/:locationId
  - [x] `ListLocationTypeByCompany()` — GET /rest/pages/admin/SM01/type?company_id=
- [x] Update `backbone/routes.go` — DI wiring + 5 routes
- [x] Build test — `go build ./...` ✅ Berhasil

### Frontend (ict_site)
- [x] Update `src/app/board/pages/SM01/page.tsx` — Location Tab
  - [x] `UserLocationListItem` interface (with company_name)
  - [x] `LocationTypeSelectItem` interface
  - [x] Zod schemas: `locationCreateSchema` (with company_id), `locationUpdateSchema`
  - [x] State management: `userLocations`, `locationFormOpen`, `locationEditing`, `locationTypeOptions`
  - [x] CRUD handlers: `fetchUserLocations`, `fetchLocationTypes`, `handleLocationCreate`, `handleLocationUpdate`, `handleLocationDelete`
  - [x] `useWatch` on `company_id` triggers `fetchLocationTypes` to filter types by company
  - [x] DataTable columns: Type Code, Type Name, Company, Active badge
  - [x] Create DataDialog: Company select + Location Type select (filtered by company) + Active checkbox
  - [x] Update DataDialog: Location Type select + Active checkbox
  - [x] Delete confirmation dialog
  - [x] `TabsContent("location")` in User Details modal
- [x] Build test — `npx tsc --noEmit` ✅ Berhasil

---

## Task: SM03 Location Type Management (Backend + Frontend)

### Backend (ict_rest)
- [x] Buat `skeleton/SM03/template.go` — Location type DTOs
  - [x] `LocationTypeItem` DTO — id, company_id, parent_id, code, name, description, icon, color, is_active
  - [x] `LocationTypeCreateRequest` DTO — parent_id, code, name, description, icon, color, is_active
  - [x] `LocationTypeUpdateRequest` DTO — same fields
- [x] Buat `skeleton/SM03/repository.go` — Location type queries
  - [x] `ListLocationType(companyID)` — Query dat_location_type WHERE company_id
  - [x] `CreateLocationType(req)` — INSERT dat_location_type
  - [x] `UpdateLocationType(id, req)` — UPDATE dat_location_type
  - [x] `DeleteLocationType(id)` — DELETE dat_location_type
- [x] Buat `skeleton/SM03/usecase.go` — Location type business logic
  - [x] All CRUD methods delegate ke repo
- [x] Buat `skeleton/SM03/handler.go` — Location type HTTP handlers
  - [x] `ListLocationType()` — GET /rest/pages/admin/SM03/:id/type
  - [x] `CreateLocationType()` — POST /rest/pages/admin/SM03/:id/type
  - [x] `UpdateLocationType()` — PUT /rest/pages/admin/SM03/:id/type/:typeId
  - [x] `DeleteLocationType()` — DELETE /rest/pages/admin/SM03/:id/type/:typeId
- [x] Update `backbone/routes.go` — DI wiring + 4 routes
- [x] Build test — `go build ./...` ✅ Berhasil

### Frontend (ict_site)
- [x] Update `src/app/board/pages/SM03/page.tsx` — Location Type Tab
  - [x] `LocationTypeItem` interface
  - [x] Zod schemas: `locationTypeCreateSchema`, `locationTypeUpdateSchema`
  - [x] State management: `locationTypes`, `locationTypeFormOpen`, `locationTypeEditing`
  - [x] CRUD handlers: `fetchLocationTypes`, `handleLocationTypeCreate`, `handleLocationTypeUpdate`, `handleLocationTypeDelete`
  - [x] DataTable columns: Code, Name, Description, Active badge
  - [x] Create DataDialog: react-hook-form + zod (code, name, description, is_active)
  - [x] Update DataDialog: react-hook-form + zod
  - [x] Delete confirmation dialog
  - [x] `TabsContent("location")` in Company Details modal

---

## Task: Session Management (SM05)

### Backend (ict_rest)
- [x] Buat `skeleton/SM05/template.go` — Interfaces + DTOs
  - [x] `SessionListItem` DTO — id, user_id, username, name, company_name, ip_address, user_agent, created_at, expires_at
  - [x] `Repository` interface — `ListSessions`, `FindSession`, `DeleteSession`
  - [x] `UseCase` interface — same methods
- [x] Buat `skeleton/SM05/repository.go` — Database queries (NRepo constructor)
  - [x] `ListSessions(search, page, size, sortBy, sortOrder)` — Query dat_user_session JOIN dat_user + dat_company
  - [x] `FindSession(id)` — Query single session
  - [x] `DeleteSession(id)` — DELETE dat_user_session
- [x] Buat `skeleton/SM05/usecase.go` — Business logic (NCase constructor)
- [x] Buat `skeleton/SM05/handler.go` — HTTP handlers (NHand constructor, 3 endpoints)
  - [x] `ListSessions()` — GET /rest/pages/admin/SM05
  - [x] `FindSession()` — GET /rest/pages/admin/SM05/:id
  - [x] `DeleteSession()` — DELETE /rest/pages/admin/SM05/:id
- [x] Update `backbone/routes.go` — DI wiring
- [x] Build test — `go build ./...` ✅ Berhasil

### Frontend (ict_site)
- [x] Buat `src/app/board/pages/SM05/page.tsx` — Session management page
  - [x] DataTable: Username, Name, Company, IP Address, User Agent, Created At, Expires At
  - [x] Search + pagination + sorting
  - [x] Delete confirmation dialog (revoke session)
  - [x] Fetch: GET/DELETE /proxy/pages/admin/SM05

### Fitur SM05
- Lihat semua active sessions (admin only)
- Revoke/invalidate session user
- Endpoint protected: Bearer token + USLock (admin session)

---

## Task: Bug Fixes

### Cast Fix: SM01 `::dat_action_type` → `::action_type`
- [x] Fix `ict_rest/skeleton/SM01/repository.go` — Changed 3 occurrences of `::dat_action_type` → `::action_type`

### Cast Fix: SM04 `::dat_approval_flag` → `::approval_flag`
- [x] Fix `ict_rest/skeleton/SM04/repository.go` — Changed 1 occurrence of `::dat_approval_flag` → `::approval_flag`

### Fix: SP03 userId from Session Context
- [x] Fix `ict_rest/skeleton/SP03/handler.go` — Changed to get `userId` from session context instead of query param

### Fix: SM04 Frontend Null-Safety
- [x] Fix `ict_site/src/app/board/pages/SM04/page.tsx` — Changed `row.steps.map(...)` → `(row.steps ?? []).map(...)`
- [x] Fix `ict_site/src/app/board/pages/SM04/page.tsx` — Changed `s.signers.map(...)` → `(s.signers ?? []).map(...)`

---

## Task Detail

> ~~**Prioritas: Error Handling Terpadu**~~ ✅ Selesai — Lihat `ict_rest.md` dan `ict_site.md` untuk detail.

### Error Handling — Backend (ict_rest)
- [x] Buat `mechanic/helper.go` — AppError struct + ValidationError/NotFound/Conflict/Unauthorized/Forbidden/ExternalServiceError/InternalError + Error() handler
- [x] Install zerolog, buat `backbone/logger.go` — structured logging + request ID middleware
- [x] Buat `backbone/recovery.go` — custom panic recovery dengan logging
- [x] Update `main.go` — graceful shutdown (http.Server + signal.Notify)
- [x] Update `backbone/routes.go` — tambah RequestID() + CustomRecovery() middleware
- [x] Hapus ValidError, ErrorNumber, ShowError, ParseError dari mechanic/helper.go
- [x] Update mechanic/typography.go — CheckMeta return *AppError
- [x] Migrate SP00 handler/usecase → mechanic.Error + AppError
- [x] Migrate SP01 handler → mechanic.Error
- [x] Migrate SP02 handler → mechanic.Error
- [x] Migrate SP03 handler → mechanic.Error
- [x] Migrate SM01 handler/usecase/repo → mechanic.Error + AppError
- [x] Migrate SM02 handler → mechanic.Error
- [x] Migrate SM03 handler → mechanic.Error
- [x] Migrate SM04 handler → mechanic.Error
- [x] Build test — `go build ./...`

### Error Handling — Frontend (ict_site)
- [x] Buat `src/uix/error-boundary.tsx` — ErrorBoundary class + ErrorFallback
- [x] Buat `src/lib/error-message.ts` — error code → user message mapping
- [x] Update `src/lib/client-api.ts` — tambah code + requestId ke ClientApiError
- [x] Tambah ErrorBoundary ke layout.tsx (root + board)
- [x] Migrate SM01–SM04 → clientApi + toast + hapus silent catches
- [x] Migrate SP01–SP03 → clientApi + toast
- [x] Fix sidebar/dashboard/login catch blocks
- [x] Build test — `npm run build`

### Error Handling — Proxy Layer
- [x] Update pages proxy — error transform + request ID forwarding
- [x] Update guest proxy — error transform + network error handling

Lihat detail lengkap di `ict_rest.md` dan `ict_site.md`.

---

### Fase 1: Fondasi

#### 1. Asset/Inventory Management
- [x] Buat `ict_base/prisma/schema/ict_asset.prisma`
- [x] Definisikan enum `asset_criticality`, `asset_status`
- [x] Buat tabel `ict_asset_category`
- [x] Buat tabel `ict_asset`
- [x] Buat tabel `ict_asset_ip`
- [x] Buat tabel `ict_asset_relation`
- [x] Buat tabel `ict_asset_assignment`
- [ ] Jalankan `npx prisma migrate dev --name add_asset`
- [ ] Update `prisma/seed.ts` dengan seed data aset
- [ ] Buat AM01-AM06 skeleton (template, repository, usecase, handler)
- [ ] Buat AM01-AM06 frontend pages

#### 2. Vulnerability Management
- [ ] Buat `ict_base/prisma/schema/ict_vulnerability.prisma`
- [ ] Definisikan enum `vuln_severity`, `vuln_status`, `scan_status`
- [ ] Buat tabel `ict_vuln_scan_type`
- [ ] Buat tabel `ict_vuln_scan`
- [ ] Buat tabel `ict_vuln_finding`
- [ ] Buat tabel `ict_vuln_remediation`
- [ ] Jalankan `npx prisma migrate dev --name add_vulnerability`
- [ ] Update `prisma/seed.ts` dengan seed scan_type
- [ ] Buat repository skeleton di `ict_rest`
- [ ] Buat usecase skeleton di `ict_rest`
- [ ] Buat handler skeleton di `ict_rest`
- [ ] Buat halaman DataTable di `ict_site`

---

### Fase 2: Core Security

#### 3. Incident Management
- [ ] Buat `ict_base/prisma/schema/ict_incident.prisma`
- [ ] Definisikan enum `incident_priority`, `incident_status`
- [ ] Buat tabel `ict_incident_type`
- [ ] Buat tabel `ict_incident` (termasuk `siem_alert_id` FK)
- [ ] Buat tabel `ict_incident_asset`
- [ ] Buat tabel `ict_incident_user`
- [ ] Jalankan `npx prisma migrate dev --name add_incident`
- [ ] Buat repository/usecase/handler di `ict_rest`
- [ ] Buat halaman di `ict_site`

#### 4. SIEM / Log Monitoring
- [ ] Buat `ict_base/prisma/schema/ict_siem.prisma`
- [ ] Definisikan enum `siem_alert_status`
- [ ] Buat tabel `ict_siem_log_source`
- [ ] Buat tabel `ict_siem_rule`
- [ ] Buat tabel `ict_siem_alert`
- [ ] Buat tabel `ict_siem_correlation`
- [ ] Jalankan `npx prisma migrate dev --name add_siem`
- [ ] Buat repository/usecase/handler di `ict_rest`
- [ ] Buat agent `ict_siem_evaluator` di `ict_auto`
- [ ] Buat agent `ict_siem_notifier` di `ict_auto`
- [ ] Buat halaman di `ict_site`

---

### Fase 3: Governance

#### 5. Risk Assessment
- [ ] Buat `ict_base/prisma/schema/ict_risk.prisma`
- [ ] Definisikan enum `risk_likelihood`, `risk_impact`, `risk_level`, `risk_status`
- [ ] Buat tabel `ict_risk_category`
- [ ] Buat tabel `ict_risk`
- [ ] Buat tabel `ict_risk_assessment`
- [ ] Buat tabel `ict_risk_treatment`
- [ ] Jalankan `npx prisma migrate dev --name add_risk`
- [ ] Buat repository/usecase/handler di `ict_rest`
- [ ] Buat halaman di `ict_site` (termasuk risk matrix 5×5)

#### 6. Compliance Monitoring
- [ ] Buat `ict_base/prisma/schema/ict_compliance.prisma`
- [ ] Definisikan enum `assessment_status`, `compliance_status`
- [ ] Buat tabel `ict_comp_standard`
- [ ] Buat tabel `ict_comp_control`
- [ ] Buat tabel `ict_comp_assessment`
- [ ] Buat tabel `ict_comp_finding`
- [ ] Jalankan `npx prisma migrate dev --name add_compliance`
- [ ] Buat repository/usecase/handler di `ict_rest`
- [ ] Buat halaman di `ict_site`

---

### Fase 4: Supporting

#### 7. Certificate Management
- [ ] Buat `ict_base/prisma/schema/ict_certificate.prisma`
- [ ] Definisikan enum `cert_status`
- [ ] Buat tabel `ict_cert_type`
- [ ] Buat tabel `ict_certificate`
- [ ] Buat tabel `ict_cert_renewal_log`
- [ ] Jalankan `npx prisma migrate dev --name add_certificate`
- [ ] Buat repository/usecase/handler di `ict_rest`
- [ ] Buat halaman di `ict_site`

#### 8. Dashboard & Reporting
- [ ] Buat `ict_base/prisma/schema/ict_dashboard.prisma`
- [ ] Buat tabel `ict_dash_widget`
- [ ] Buat tabel `ict_dash_layout`
- [ ] Buat tabel `ict_report_template`
- [ ] Buat tabel `ict_report_instance`
- [ ] Buat tabel `ict_dash_bookmark`
- [ ] Jalankan `npx prisma migrate dev --name add_dashboard`
- [ ] Buat widget data API di `ict_rest`
- [ ] Buat halaman dashboard di `ict_site`

---

### Existing Agents (ict_auto)

#### Log Rotation Agent (`ict_log_rotate`)
- [x] PostgreSQL connection
- [x] Retention config (normal/attack log days)
- [x] Archive to CSV + compress
- [x] Delete archived records
- [x] Dockerfile

#### Nginx Log Agent (`ict_log_nginx`)
- [x] Elasticsearch connection
- [x] PostgreSQL connection
- [x] Nginx log parser (JSON format)
- [x] Attack detection (HTTP flood, scanner, SQL injection, XSS)
- [x] Store to PostgreSQL (`ict_nginx_log`, `ict_nginx_atc`)
- [x] Dockerfile

---

### Fase 5: Document Management

#### 9. Document Management System (DMS)
- [ ] Buat `ict_base/prisma/schema/ict_docman.prisma`
- [ ] Definisikan enum `doc_status`, `doc_visibility`, `doc_content_format`, `doc_share_permission`, `doc_todo_priority`, `doc_approval_status`
- [ ] Buat tabel `ict_doc_type`
- [ ] Buat tabel `ict_doc_category`
- [ ] Buat tabel `ict_doc`
- [ ] Buat tabel `ict_doc_version`
- [ ] Buat tabel `ict_doc_share`
- [ ] Buat tabel `ict_doc_publish`
- [ ] Buat tabel `ict_doc_todo`
- [ ] Buat tabel `ict_doc_template`
- [ ] Buat tabel `ict_doc_approval`
- [ ] Buat tabel `ict_doc_approval_step`
- [ ] Jalankan `npx prisma migrate dev --name add_docman`
- [ ] Update `prisma/seed.ts` dengan seed doc_type
- [ ] Buat repository skeleton di `ict_rest`
- [ ] Buat usecase skeleton di `ict_rest`
- [ ] Buat handler skeleton di `ict_rest`
- [ ] Buat halaman DataTable di `ict_site`
- [ ] Buat halaman editor di `ict_site`
- [ ] Buat halaman publik view di `ict_site`
- [ ] Buat fitur export (PDF, Word, Markdown)

---

### Fase 6: Docker Management & Monitoring

#### 10. Docker Management
- [ ] Buat `ict_base/prisma/schema/ict_docker.prisma`
- [ ] Definisikan enum `docker_host_connection`, `docker_host_status`, `docker_container_status`, `docker_compose_status`
- [ ] Buat tabel `ict_docker_host` (termasuk `location_id` FK)
- [ ] Buat tabel `ict_docker_container`
- [ ] Buat tabel `ict_docker_container_port`
- [ ] Buat tabel `ict_docker_container_env`
- [ ] Buat tabel `ict_docker_container_mount`
- [ ] Buat tabel `ict_docker_container_network`
- [ ] Buat tabel `ict_docker_image`
- [ ] Buat tabel `ict_docker_network`
- [ ] Buat tabel `ict_docker_volume`
- [ ] Buat tabel `ict_docker_compose`
- [ ] Buat tabel `ict_docker_compose_service`
- [ ] Buat tabel `ict_docker_stats`
- [ ] Buat tabel `ict_docker_host_stats`
- [ ] Buat tabel `ict_docker_image_pull_log`
- [ ] Jalankan `npx prisma migrate dev --name add_docker`
- [ ] Buat repository skeleton di `ict_rest`
- [ ] Buat usecase skeleton di `ict_rest`
- [ ] Buat handler skeleton di `ict_rest`
- [ ] Buat Docker API client di `ict_rest`
- [ ] Buat halaman Dashboard di `ict_site`
- [ ] Buat halaman Containers di `ict_site`
- [ ] Buat halaman Compose Stacks di `ict_site`
- [ ] Buat halaman Images/Networks/Volumes di `ict_site`
- [ ] Buat halaman Monitoring di `ict_site`
- [ ] Buat komponen DockerStatusBadge
- [ ] Buat komponen ContainerStatsChart
- [ ] Buat komponen ComposeEditor
- [ ] Buat komponen LogViewer

#### 11. Docker Agent (ict_auto/ict_docker_agent)
- [ ] Buat modul `ict_auto/ict_docker_agent/main.go`
- [ ] Implementasi Docker SDK connection
- [ ] Implementasi container stats collector
- [ ] Implementasi host stats collector
- [ ] Implementasi container list sync
- [ ] Implementasi image list sync
- [ ] Implementasi network/volume sync
- [ ] Implementasi metric sender ke hub API
- [ ] Implementasi heartbeat mechanism
- [ ] Buat Dockerfile untuk agent
- [ ] Buat konfigurasi agent (env vars)
- [ ] Test agent dengan Docker host

---

### Fase 7: Service Request Management

#### 12. Service Request Management
- [ ] Buat `ict_base/prisma/schema/ict_service_request.prisma`
- [ ] Definisikan enum `sr_priority`, `sr_status`, `sr_ticket_status`, `sr_cab_status`, `sr_cab_decision`
- [ ] Buat tabel `ict_sr_type`
- [ ] Buat tabel `ict_sr_ticket`
- [ ] Buat tabel `ict_sr_ticket_assignee`
- [ ] Buat tabel `ict_sr_cab`
- [ ] Buat tabel `ict_sr_cab_vote`
- [ ] Buat tabel `ict_sr_comment`
- [ ] Buat tabel `ict_sr_history`
- [ ] Jalankan `npx prisma migrate dev --name add_service_request`
- [ ] Update `prisma/seed.ts` dengan seed tipe
- [ ] Buat repository skeleton di `ict_rest`
- [ ] Buat usecase skeleton di `ict_rest`
- [ ] Buat handler skeleton di `ict_rest`
- [ ] Buat halaman User Portal di `ict_site`
- [ ] Buat halaman Admin Portal di `ict_site`
- [ ] Buat komponen dynamic form dari form_schema
- [ ] Buat halaman katalog layanan di `ict_site`

---

## ict_site — UIX Components

### Status Komponen
| Kategori | Jumlah | Keterangan |
|----------|--------|------------|
| Existing (shadcn/ui) | 62 | Sudah ada di src/uix/ |
| New Reusable | 12 | Komponen reusable untuk semua modul |
| Location Components | 3 | LocationPicker, LocationBadge, LocationTree |
| Utility Components | 4 | ConfirmDialog, CopyButton, ErrorBoundary, Spinner |
| **Total** | **81** | — |

### Komponen Existing (62)
- accordion, alert, alert-dialog, aspect-ratio, attachment, avatar
- badge, breadcrumb, bubble, button, button-group, calendar
- card, carousel, chart, checkbox, collapsible, command
- combobox, context-menu, datadialog, datatable, dialog
- direction, drawer, dropdown-menu, empty, field
- hover-card, input, input-group, input-otp
- item, kbd, label, marker, menubar
- message, message-scroller, native-select, navigation-menu
- pagination, popover, progress, radio-group
- resizable, scroll-area, select, separator, sheet
- sidebar, skeleton, slider, sonner, spinner
- switch, table, tabs, textarea, toggle
- toggle-group, tooltip

### Komponen Baru (12) — ✅ Semua Sudah Dibuat
| Komponen | File | Kebutuhan Modul | Status |
|----------|------|-----------------|--------|
| `status-badge` | `status-badge.tsx` | Semua modul (severity, status, level) | ✅ Dibuat |
| `timeline` | `timeline.tsx` | Incident, SIEM, SR, Docman | ✅ Dibuat |
| `tree-view` | `tree-view.tsx` | Location, Compliance, Docman, Asset | ✅ Dibuat |
| `file-upload` | `file-upload.tsx` | Billing, Docman, PM, Compliance | ✅ Dibuat |
| `stat-card` | `stat-card.tsx` | Semua dashboard | ✅ Dibuat |
| `approval-steps` | `approval-steps.tsx` | SR, Billing, Docman | ✅ Dibuat |
| `heatmap` | `heatmap.tsx` | Risk, SIEM, Dashboard | ✅ Dibuat |
| `kanban-board` | `kanban-board.tsx` | Docman, Incident | ✅ Dibuat |
| `event-calendar` | `event-calendar.tsx` | Preventive, Dashboard | ✅ Dibuat |
| `log-viewer` | `log-viewer.tsx` | Docker | ✅ Dibuat |
| `photo-upload` | `photo-upload.tsx` | Preventive | ✅ Dibuat |
| `markdown-preview` | `markdown-preview.tsx` | Docman | ✅ Dibuat |

### Komponen Location (3) — ✅ Semua Sudah Dibuat
| Komponen | File | Kebutuhan Modul | Status |
|----------|------|-----------------|--------|
| `location-picker` | `location-picker.tsx` | SM01, SM03, AM | ✅ Dibuat |
| `location-badge` | `location-badge.tsx` | Semua modul (location display) | ✅ Dibuat |
| `location-tree` | `location-tree.tsx` | SM03, AM, DK | ✅ Dibuat |

### Komponen Utility (4) — ✅ Semua Sudah Dibuat
| Komponen | File | Kebutuhan Modul | Status |
|----------|------|-----------------|--------|
| `confirm-dialog` | `confirm-dialog.tsx` | Semua modul (delete confirmation) | ✅ Dibuat |
| `copy-button` | `copy-button.tsx` | Semua modul (copy to clipboard) | ✅ Dibuat |
| `error-boundary` | `error-boundary.tsx` | Root layout + board layout | ✅ Dibuat |
| `spinner` | `spinner.tsx` | Loading states | ✅ Dibuat |

---

## ict_site — Custom React Hooks (Lib)

### Status Hooks
| Hook | File | Keterangan | Status |
|------|------|------------|--------|
| `useDebounce` | `use-debounce.ts` | Debounce value dengan configurable delay | ✅ Dibuat |
| `useEventListener` | `use-event-listener.ts` | Attach event listener ke element/window | ✅ Dibuat |
| `useWindowSize` | `use-window-size.ts` | Track window dimensions | ✅ Dibuat |
| `useLocalStorage` | `use-local-storage.ts` | Persistent state ke localStorage | ✅ Dibuat |
| `useMediaQuery` | `use-media-query.ts` | CSS media query matching | ✅ Dibuat |
| `usePrevious` | `use-previous.ts` | Track nilai sebelumnya | ✅ Dibuat |
| `useUpdateEffect` | `use-update-effect.ts` | Effect skip initial render | ✅ Dibuat |
| `useIsomorphicLayoutEffect` | `use-isomorphic-layout-effect.ts` | Layout effect SSR-safe | ✅ Dibuat |

### Utility Modules
| Module | File | Keterangan | Status |
|--------|------|------------|--------|
| `storage` | `storage.ts` | LocalStorage wrapper + `getPreference()` | ✅ Dibuat |
| `serverApi` | `server-api.ts` | Server-side fetch (server components) | ✅ Dibuat |

---

## ict_site — Package Management

### Status Package
| Kategori | Jumlah | Keterangan |
|----------|--------|------------|
| Digunakan aktif | 24 | Sudah di-import di source code (termasuk react-hook-form, zod, @hookform/resolvers) |
| Direncanakan | 10 | Untuk modul DMS, forms, export |
| Dihapus | 3 | `axios`, `@excalidraw/excalidraw`, `mermaid` (unused + security risk) |
| **Total dependencies** | **34** | — |
| **npm audit** | **0** | ✅ Bersih dari vulnerabilities |

### Package yang Digunakan Aktif
| Package | Keterangan |
|---------|------------|
| `@base-ui/react` | Base UI primitives (47 imports) |
| `@shadcn/react` | shadcn primitives |
| `class-variance-authority` | Variant styling (16 imports) |
| `clsx` | Class name utility |
| `cmdk` | Command palette |
| `embla-carousel-react` | Carousel |
| `geist` | Font |
| `input-otp` | OTP input |
| `lucide-react` | Icons (22 imports) |
| `next` | Framework |
| `next-themes` | Theme switching |
| `react` | UI library |
| `react-day-picker` | Calendar |
| `react-dom` | React DOM |
| `react-resizable-panels` | Resizable panels |
| `recharts` | Charts |
| `shadcn` | CSS theme |
| `sonner` | Toast notifications |
| `tailwind-merge` | Tailwind merge |
| `vaul` | Drawer |
| `zustand` | State management |

### Package yang Direncanakan (Belum Digunakan)
| Package | Modul | Keterangan |
|---------|-------|------------|
| `@tiptap/*` (9 packages) | DMS | Rich text editor |
| `date-fns` | Semua | Date utility library |
| `docx` | DMS | Word document export |
| `html2canvas` | DMS | HTML capture untuk PDF |
| `jspdf` | DMS | PDF document export (v4.2.1, updated for security) |
| `lowlight` | DMS | Code syntax highlighting |

---

## Catatan Teknis

### Pola Konsisten
- Primary key: `String @id @default(uuid())`
- Timestamps: `created_at DateTime @default(now())` + `updated_at DateTime @default(now())`
- Soft delete: gunakan `is_active` atau `status` (tidak hard delete)
- Multi-tenant: semua tabel domain menggunakan `company_id`
- Unique constraint: `@@unique([company_id, ...])` untuk isolasi data per perusahaan
- Index: `@@index` untuk kolom yang sering di-query

### Naming Convention
- Tabel: `ict_{modul}` atau `dat_{entitas}`
- Enum: `snake_case` sesuai Prisma convention
- Field: `snake_case` (konsisten dengan existing)

### Relasi Utama
```
dat_company
├── ict_location (Location Management)
│   ├── ict_location_detail
│   ├── ict_location_contact
│   ├── ict_location_asset
│   └── ict_asset (FK location_id)
├── ict_asset
│   ├── ict_vuln_finding
│   ├── ict_incident_asset
│   ├── ict_certificate
│   └── ict_risk
├── ict_vuln_scan
│   └── ict_vuln_finding
├── ict_comp_assessment
│   └── ict_comp_finding
├── ict_incident (FK siem_alert_id → ict_siem_alert)
│   ├── ict_incident_asset
│   └── ict_incident_user
├── ict_risk (FK comp_finding_id → ict_comp_finding)
│   ├── ict_risk_assessment
│   └── ict_risk_treatment
├── ict_siem_alert (FK incident_id → ict_incident)
│   └── ict_siem_correlation
├── ict_doc (Document Management)
│   ├── ict_doc_version
│   ├── ict_doc_share
│   ├── ict_doc_publish
│   ├── ict_doc_todo
│   ├── ict_doc_template
│   └── ict_doc_approval
│       └── ict_doc_approval_step
├── ict_docker (Docker Management)
│   ├── ict_docker_container
│   │   ├── ict_docker_container_port
│   │   ├── ict_docker_container_env
│   │   ├── ict_docker_container_mount
│   │   └── ict_docker_container_network
│   ├── ict_docker_image
│   ├── ict_docker_network
│   ├── ict_docker_volume
│   ├── ict_docker_compose
│   │   └── ict_docker_compose_service
│   ├── ict_docker_stats
│   ├── ict_docker_host_stats
│   └── ict_docker_image_pull_log
├── ict_sr (Service Request)
│   ├── ict_sr_type
│   ├── ict_sr_ticket
│   │   ├── ict_sr_ticket_assignee
│   │   ├── ict_sr_comment
│   │   └── ict_sr_history
│   ├── ict_sr_cab
│   │   └── ict_sr_cab_vote
│   └── ict_sr_fulfillment
├── ict_audit_trail (shared across modules)
├── ict_file_attachment (shared across modules)
├── ict_approval (shared across modules)
└── ict_notification (shared across modules)
```

---

## Update Terakhir

- **Tanggal**: 22 Juli 2026
- **Oleh**: System
- **Keterangan**: Sinkronisasi documentation dengan source code:
  - Hapus duplikasi section (SM01 Location ×3, SM03 Location Area ×3, Bug Fixes ×2)
  - Tambah Task: Session Management (SM05) — backend + frontend
  - Update UIX component count: 74 → 81 (tambah Location Components: 3, Utility Components: 4)
  - Tambah section Custom React Hooks (8 hooks) + Utility Modules (2 modules)
  - Total lib files: 17
