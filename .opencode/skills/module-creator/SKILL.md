---
name: module-creator
description: Primary skill for module-scoped work: new modules, module pages, and sub-entities. Triggers: "buat modul", "modul baru", "tambah halaman modul", "tambah sub-entity", "new module", "add module page", "create module". Uses a design-first flow with approval before implementation.
---

# Module Creator Skill

Automasi pembuatan modul OmniSight dengan pendekatan **design-first** yang tetap mengikuti prinsip **clean-execution**.

## Kapan dipakai

Gunakan skill ini jika request berkaitan dengan:

- pembuatan modul baru
- penambahan halaman pada modul yang sudah ada
- penambahan sub-entity atau child table pada halaman modul
- perubahan struktur modul yang membutuhkan rancangan sebelum implementasi

Jangan pakai skill ini jika task hanya berupa:

- refactor kecil pada halaman yang sudah ada
- perbaikan bug lokal tanpa perubahan struktur modul
- penyesuaian styling atau perilaku UI yang tidak mengubah desain modul
- pembuatan atau restrukturisasi satu halaman frontend ketika backend dan desain modul sudah jelas

Untuk task-task tersebut, pakai skill yang lebih sempit: `frontend-page-creator`, `backend-page-creator`, atau `clean-execution`.

## Referensi wajib

Sebelum membuat rancangan atau implementasi, baca dan ikuti referensi berikut:

- `obx_docs/blueprint/{KODE}/README.md` — dokumentasi teknis per sub-module (kode halaman, privilege, route, checklist)
- `obx_docs/guide/{KODE}.md` — dokumentasi user guide per sub-module
- `obx_docs/blueprint/SM/README.md` — referensi pola CRUD sederhana
- `obx_docs/blueprint/IM/README.md` — referensi pola multi-tab atau relasi yang lebih kompleks
- `obx_base/prisma/schema/` — naming convention dan relasi schema jika diperlukan

## Format prompt yang disarankan

Jika user ingin hasil cepat dan minim klarifikasi, dorong format input yang konsisten seperti ini:

```text
Scope: module
Goal: buat modul baru | tambah halaman modul | tambah sub-entity
Module code: {KODE}
Page code: {KODE}01 | {KODE}02 | ...
Entity name: {Nama}
Route group: admin | authu | guest | pages
Database table: {nama_tabel}
Prisma schema: sudah ada | belum ada
Fields:
- field_a: type, nullable/not null, default, unique?, relation?
- field_b: type, nullable/not null, default, unique?, relation?
Sub-entity: ada | tidak ada
Kebutuhan UI: CRUD standar | multi-tab | select endpoint | mobile card view
```

Contoh prompt siap pakai:

```text
Scope: module
Goal: buat modul baru
Module code: NM
Page code: NM01
Entity name: Notification
Route group: admin
Database table: dat_notification
Prisma schema: belum ada
Fields:
- id: text, primary key, default gen_random_uuid()
- company_id: text, not null, relation dat_company
- code: text, not null, unique per company
- message: text, not null
- severity: enum(critical, warning, info), default info
- is_active: boolean, default true
- created_at: timestamptz, default now()
- updated_at: timestamptz, default now()
Sub-entity: tidak ada
Kebutuhan UI: CRUD standar, select endpoint, mobile card view
```

## Mode eksekusi

- Pahami kebutuhan modul secara utuh sebelum menulis file.
- Ambil konteks minimum dari `obx_docs/blueprint/{KODE}/README.md`, `obx_docs/guide/{KODE}.md`, skeleton referensi, dan schema terkait.
- Gunakan pola repo yang sudah ada; jangan membuat variasi baru tanpa kebutuhan jelas.
- Kerjakan per fase: sinkronisasi konteks, rancangan, approval, implementasi, validasi.
- Setelah edit substantif pertama, langsung lakukan validasi terfokus.
- Jangan memperluas scope di luar modul atau halaman yang diminta.

## Output contract

Saat fase desain atau implementasi selesai, laporkan dengan prinsip `token-efficient`:

```text
Scope: {modul / halaman / sub-entity}
Changes: {file yang berubah}
Validation: {go build ./... / npx tsc --noEmit} -> {hasil}
Notes: {approval, blocker, atau next step}
```

## Workflow

### Phase 0: Sinkronisasi konteks

- Tentukan apakah request adalah modul baru, halaman baru, atau sub-entity.
- Identifikasi anchor terdekat: kode halaman, folder skeleton, route, atau tabel yang ada.
- Validasi kode modul, route group, dan privilege terhadap `obx_docs/blueprint/{KODE}/README.md` (jika sudah tersedia).
- Jika konteks inti sudah ada, jangan bertanya ulang.
- Jika ada ambiguity yang memblokir desain, tanyakan hanya poin minimum.

### Phase 1: Kumpulkan input minimum

Kumpulkan hanya yang benar-benar dibutuhkan:

1. **Tipe**: modul baru / tambah halaman / tambah sub-entity
2. **Kode**: `{KODE}` atau `{KODE}{XX}`
3. **Nama entity**
4. **Route group**
5. **Database table**
6. **Prisma schema**: sudah ada / belum ada
7. **Fields** jika schema belum ada atau belum jelas
8. **Sub-entity**: ada / tidak ada

Aturan phase ini:

- Jangan minta data yang sudah bisa dibaca dari schema atau modul serupa.
- Jika hanya menambah satu halaman ke modul existing, fokus pada halaman itu.
- Jika kode modul belum terdaftar atau privilege tidak cocok dengan `obx_docs/blueprint/{KODE}/README.md`, minta konfirmasi eksplisit.

### Phase 2: Rancang design document

Buat `obx_docs/blueprint/{KODE}/README.md` dengan struktur:

```
# Rancangan Modul: {KODE} — {Nama}

## 1. Overview
- Group, Icon, Akses, Tujuan
- Referensi Skema
- Tabel Halaman

## 2. Database Schema
- Schema Prisma lengkap (jika belum ada, buat rancangan konkret, bukan sekadar ringkasan)
- DTOs (Go structs)

## 3. API Endpoints
- Per halaman: method, endpoint, handler, description
- Query parameters

## 4. Frontend Pages
- Per halaman: tabs, DataTable columns, Form fields
- Mobile behavior: desktop table vs mobile card, `useIsMobile()`, `DataTableCard` bila relevan
- Drill-down levels (jika ada)

## 5. Implementation Checklist
- Per halaman: 4 file backend + 1 file frontend + routes

## 6. Relationships
- Relasi antar tabel
- Dependencies antar halaman

## 7. Backend Pattern Reference
- DI wiring
- companyId extraction
- InternalError wrapping
- SQL patterns
- Response format
```

Aturan phase ini:

- Rancangan harus ringkas, konkret, dan langsung bisa diimplementasikan.
- Semua route backend dan frontend harus konsisten dengan `obx_docs/blueprint/{KODE}/README.md`.
- CRUD standar harus memakai pola `DataTable + DataDialog`.
- Jika perlu mobile list view, cantumkan `useIsMobile()` dan `DataTableCard`.
- Jika schema baru, tampilkan model Prisma yang cukup lengkap untuk di-review.
- Jika ada write flow di frontend, cantumkan loading state dan `DataDialog loading`.

### Phase 3: Review

- Tampilkan rancangan secara utuh atau cukup lengkap untuk di-review.
- Sorot keputusan penting: schema, privilege, sub-entity, mobile behavior, dan route pattern.
- Jangan mulai implementasi sebelum approval jelas.
- Jika revisi hanya lokal, edit bagian itu saja.

### Phase 4: Implementasi

Setelah user approve:

- Selalu lengkapi dokumentasi yang relevan yang dirujuk (khususnya `obx_docs/blueprint/{KODE}/README.md` dan `obx_docs/guide/{KODE}.md`) jika perubahan memengaruhi rancangan, route, privilege, atau struktur modul.

1. **Database** (jika belum ada):
   - Buat file Prisma
   - `npx prisma generate && npx prisma migrate dev --name {KODE}`

2. **Backend** (per halaman):
   - Buat `obx_rest/skeleton/{KODE}{XX}/template.go`
   - Buat `obx_rest/skeleton/{KODE}{XX}/repository.go`
   - Buat `obx_rest/skeleton/{KODE}{XX}/usecase.go`
   - Buat `obx_rest/skeleton/{KODE}{XX}/handler.go`
   - Register DI + routes di `obx_rest/backbone/routes.go`

3. **Frontend** (per halaman):
   - Buat `obx_site/src/app/board/pages/{KODE}{XX}/page.tsx`
   - Terapkan `DataTable + DataDialog`, loading state pattern, dan mobile behavior sesuai blueprint teknis sub-module

4. **Documentation**:
   - Update `obx_docs/blueprint/{KODE}/README.md` untuk dokumentasi teknis
   - Update `obx_docs/guide/{KODE}.md` untuk dokumentasi user guide

5. **Build test**:
   - `go build ./...`
   - `npx tsc --noEmit`

Aturan phase ini:

1. Kerjakan per irisan kecil: database bila perlu, backend satu halaman, frontend halaman yang sama.
2. Setelah edit substantif pertama, jalankan validasi paling murah yang relevan.
3. Jangan lanjut ke halaman berikutnya sebelum halaman saat ini konsisten dengan blueprint.
4. Jika validasi gagal, perbaiki root cause pada irisan yang sama.
5. Jangan membuat file dokumentasi di luar struktur `obx_docs/blueprint/` dan `obx_docs/guide/`.

## Constraints

- Ikuti blueprint teknis di `obx_docs/blueprint/{KODE}/README.md` secara ketat, terutama checklist, route group, dan standar frontend
- Error: `mechanic.Error(c, err)` di handler, `mechanic.InternalError()` di usecase
- SQL: parameterized queries, `validSort` whitelist, `rows.Err()` setelah loop
- Frontend: `clientApi()`, `toast.error(getErrorMessage())`, `defaultValues`, `Controller`
- **Loading state WAJIB**: `useState(false)` → `setLoading(true)` diawal → `finally { setLoading(false) }` — untuk semua write operation (create, update, delete). Pass `loading` prop ke DataDialog
- DataTable: wajib support sorting, searching, pagination
- DataDialog: wajib dipakai untuk create, update, dan detail; jangan buat page CRUD terpisah
- Optional input field: wajib controlled dengan `value={field.value ?? ""}` untuk `Input` dan `Textarea`
- Mobile page: gunakan `useIsMobile()` dan `DataTableCard` bila halaman menampilkan list tabular
- Review design doc harus cukup detail untuk mengunci schema, route, mobile behavior, dan loading state sebelum implementasi
- Static route seperti `/{KODE}/select` harus didaftarkan sebelum dynamic route seperti `/{KODE}/:id`
- Semua write route backend wajib memakai `USLogs("{KODE}")`
- Bahasa: kode = Inggris, dokumentasi = Indonesia
- Design doc harus di-approve SEBELUM implementasi
- Jangan modifikasi file Prisma yang sudah ada
- Jangan commit
- Jangan menambah fitur di luar cakupan request user
- Jangan mengubah banyak file jika satu file atau satu slice sudah cukup

## Validasi minimum

- Backend: `go build ./...` dari `obx_rest/`
- Frontend: `npx tsc --noEmit` atau `npm run build` dari `obx_site/`
- Modul penuh: validasi backend lalu frontend
- Jangan mengklaim selesai tanpa hasil validasi yang segar.

## Gaya komunikasi

- Gunakan bahasa yang singkat dan operasional.
- Saat di fase desain, sorot keputusan yang memang perlu approval user.
- Saat di fase implementasi, laporkan progres per halaman atau per irisan.

## File Structure

```
obx_docs/blueprint/{KODE}/
├── README.md              ← technical design document

obx_rest/skeleton/{KODE}/
├── {KODE}01/
│   ├── template.go
│   ├── repository.go
│   ├── usecase.go
│   └── handler.go
├── {KODE}02/
│   ├── template.go
│   ├── repository.go
│   ├── usecase.go
│   └── handler.go
└── ...

obx_site/src/app/board/pages/
├── {KODE}01/
│   └── page.tsx
├── {KODE}02/
│   └── page.tsx
└── ...
```
