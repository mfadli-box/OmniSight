---
name: backend-page-creator
description: Skill for page-scoped backend work in `obx_rest/skeleton/{KODE}/`, especially template, repository, usecase, handler, and route wiring for one page. Triggers: "buat backend page", "backend page", "buat handler", "buat usecase", "buat repository", "skeleton page", "crud backend page", "rest page". If the request creates a new module, prefer `module-creator`.
---

# Backend Page Creator Skill

Skill ini dipakai untuk membuat atau merapikan implementasi backend untuk satu halaman di `obx_rest` dengan tetap mengikuti blueprint proyek.

## Kapan dipakai

Gunakan skill ini jika request berkaitan dengan:

- pembuatan backend untuk satu halaman `obx_rest/skeleton/{KODE}/`
- penambahan `template.go`, `repository.go`, `usecase.go`, `handler.go` untuk satu kode halaman
- penyesuaian route wiring atau usecase flow untuk satu halaman existing
- pekerjaan backend page-level yang tidak mengubah desain modul secara menyeluruh

Jangan pakai skill ini jika task berupa:

- pembuatan modul baru atau penambahan halaman modul yang butuh design doc dan approval formal
- perubahan frontend page-level sebagai fokus utama
- refactor kecil umum yang tidak butuh workflow backend page khusus

Untuk task-task tersebut:

- pakai `module-creator` untuk workflow modul
- pakai `frontend-page-creator` untuk page-level frontend
- pakai `clean-execution` untuk implementasi atau refactor umum

## Referensi wajib

Sebelum implementasi, baca dan ikuti referensi berikut:

- `obx_docs/blueprint/{KODE}/README.md` — blueprint teknis backend, route group, error handling, SQL pattern, dan checklist
- `obx_docs/guide/{KODE}.md` — user guide per sub-module
- skeleton pembanding di `obx_rest/skeleton/SM/`, `obx_rest/skeleton/IM/`, atau halaman serupa
- `obx_rest/backbone/routes.go` — pola DI dan registrasi route
- `obx_rest/mechanic/` — helper error dan meta utilities

## Format prompt yang disarankan

Jika user ingin hasil cepat dan minim klarifikasi, dorong format input yang konsisten seperti ini:

```text
Scope: backend-page
Goal: buat backend baru | ubah backend existing | rapikan flow backend
Page code: {KODE}
Entity name: {Nama}
Route group: admin | authu | guest | pages
Database table: {nama_tabel}
Schema/table: sudah ada | belum ada
Fields response:
- field_a: type
- field_b: type
Request fields:
- field_a: required/optional
- field_b: required/optional
Sub-entity: ada | tidak ada
Select endpoint: ya/tidak
Company scope: ya/tidak
```

Contoh prompt siap pakai:

```text
Scope: backend-page
Goal: buat backend baru
Page code: NM01
Entity name: Notification
Route group: admin
Database table: dat_notification
Schema/table: sudah ada
Fields response:
- id: string
- code: string
- message: string
- severity: string
- is_active: boolean
- created_at: string
Request fields:
- code: required
- message: required
- severity: required
- is_active: optional
Sub-entity: tidak ada
Select endpoint: ya
Company scope: ya
```

## Prinsip eksekusi

- Ambil konteks minimum yang benar-benar dibutuhkan.
- Fokus pada satu kode halaman atau satu slice backend terlebih dahulu.
- Ikuti blueprint teknis pada `obx_docs/blueprint/{KODE}/README.md`.
- Jangan memperluas scope ke frontend atau modul lain jika tidak diminta.
- Setelah edit substantif pertama, langsung lakukan validasi terfokus.

## Output contract

Setelah implementasi, laporkan hasil dengan format singkat dan prinsip `token-efficient`:

```text
Scope: {kode halaman / area backend}
Changes: {file yang berubah}
Validation: {go build ./...} -> {hasil}
Notes: {blocker atau catatan penting}
```

## Workflow

### Phase 0: Sinkronisasi konteks

- Validasi kode halaman terhadap `obx_docs/blueprint/{KODE}/README.md`.
- Identifikasi apakah backend page ini baru atau perubahan atas page existing.
- Identifikasi anchor terdekat: folder skeleton, route yang ada, atau halaman backend pembanding.
- Jika desain modul, schema, atau privilege belum jelas, berhenti dan arahkan ke `module-creator` atau minta konfirmasi scope.

### Phase 1: Kumpulkan input minimum

Kumpulkan hanya informasi yang benar-benar perlu:

1. **Kode halaman**: `{KODE}`
2. **Nama entity**: nama utama untuk DTO dan message
3. **Route group**: `guest` / `authu` / `admin` / `pages` / `robot`
4. **Database table**: tabel utama
5. **Field response dan request**
6. **Sub-entity**: ada / tidak ada
7. **Select endpoint**: ya / tidak
8. **Company scope**: ya / tidak

Aturan clean-execution untuk phase ini:

- Jangan bertanya ulang jika struktur tabel atau pola endpoint sudah jelas dari schema dan halaman pembanding.
- Jika hanya satu file backend yang perlu disesuaikan, fokus di file itu dulu.
- Jika route group atau privilege bertentangan dengan `obx_docs/blueprint/{KODE}/README.md`, hentikan asumsi dan minta konfirmasi eksplisit.

### Phase 2: Rancang struktur backend page

Tentukan secara singkat:

- DTO di `template.go`
- method repository yang dibutuhkan
- validasi usecase dan wrapping `mechanic.InternalError()`
- response pattern di handler
- route wiring, urutan static route, dan kebutuhan select/sub-entity

Jika request cukup besar, laporkan rancangan singkat sebelum implementasi. Jika request lokal, boleh lanjut tanpa design doc formal.

### Phase 3: Implementasi

Implementasi WAJIB mengikuti standar berikut:

- Jika task memengaruhi backend page yang dijelaskan di dokumentasi referensi, lengkapi atau perbarui dokumen yang relevan (khususnya `obx_docs/blueprint/{KODE}/README.md` dan `obx_docs/guide/{KODE}.md`) bersama implementasi.

- Gunakan pola yang sudah ditetapkan di blueprint teknis sub-module
- Handler wajib memakai `mechanic.Error(c, err)`
- Usecase wajib wrap error repository dengan `mechanic.InternalError()`
- Gunakan query parameterized dan whitelist sort field
- Cek `rows.Err()` setelah iterasi rows
- Gunakan `USLogs("{KODE}")` untuk semua write route
- Daftarkan static route seperti `/select` sebelum route dinamis
- Pertahankan nama type, field, dan response dalam bahasa Inggris

### Phase 4: Verifikasi

Lakukan validasi sesempit mungkin lebih dulu:

1. `go build ./...` dari `obx_rest/`
2. Jika perubahan menyentuh wiring lintas lebih luas, ulangi `go build ./...` setelah seluruh slice selesai

Jangan mengklaim selesai tanpa hasil validasi yang segar.

## Gaya komunikasi

- Gunakan bahasa yang singkat dan operasional.
- Laporkan perubahan per file atau per slice backend yang jelas.
- Saat ada blocker schema atau privilege, jelaskan scope dengan tegas.

## Constraints

- Jangan return raw error dari usecase.
- Jangan pakai `c.JSON(... gin.H{"error": ...})` untuk error handler.
- Jangan buat route yang menyimpang dari kode halaman resmi di `obx_docs/blueprint/{KODE}/README.md`.
- Jangan melewati `USLogs("{KODE}")` untuk write route.
- Jangan memperluas perubahan ke frontend kecuali user memang meminta atau scope sudah mencakupnya.
- Jika ternyata desain modul atau schema belum jelas, hentikan implementasi backend-only dan jelaskan bahwa scope perlu dipindah ke `module-creator` atau task desain terpisah.
