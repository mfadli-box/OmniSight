---
name: frontend-page-creator
description: Skill for page-scoped frontend work in `obx_site`, especially CRUD pages, dialogs, tabs, tables, and mobile list views. Triggers: "buat halaman frontend", "buat page frontend", "crud page", "datatable", "datadialog", "mobile card view", "frontend page". If the request changes module structure, prefer `module-creator`.
---

# Frontend Page Creator Skill

Skill ini dipakai untuk membuat atau merapikan halaman frontend di `obx_site` dengan tetap mengikuti pola UI dan execution style proyek.

## Kapan dipakai

Gunakan skill ini jika request berkaitan dengan:

- pembuatan halaman frontend baru di `obx_site/src/app/board/pages/{KODE}/`
- penyesuaian besar pada satu halaman CRUD yang sudah ada
- pembuatan atau perapihan flow `DataTable + DataDialog`
- penambahan tab, form, filter, mobile card view, atau select flow pada satu halaman
- pekerjaan page-level yang tidak mengubah desain modul secara lintas backend dan frontend

Jangan pakai skill ini jika task berupa:

- pembuatan modul baru atau penambahan halaman modul yang butuh design doc dan approval formal
- perubahan backend skeleton, route group, atau struktur modul di `obx_docs/blueprint/{KODE}/README.md`
- refactor kecil yang hanya butuh execution discipline biasa

Untuk task-task tersebut:

- pakai `module-creator` untuk workflow modul
- pakai `clean-execution` untuk implementasi atau refactor umum

## Referensi wajib

Sebelum implementasi, baca dan ikuti referensi berikut:

- `obx_docs/blueprint/{KODE}/README.md` — standar frontend, mobile view, DataTable, DataDialog, SearchSelect, loading state, dan validation flow
- `obx_docs/guide/{KODE}.md` — user guide per sub-module
- halaman frontend serupa di `obx_site/src/app/board/pages/`
- komponen shared di `obx_site/src/uix/`

## Format prompt yang disarankan

Jika user ingin hasil cepat dan minim klarifikasi, dorong format input yang konsisten seperti ini:

```text
Scope: frontend-page
Goal: buat halaman baru | ubah halaman existing | rapikan flow UI
Page code: {KODE}
Page purpose: {deskripsi}
Page type: CRUD standar | multi-tab | detail only | dashboard
Primary data: {entity / endpoint}
Field form:
- field_a: input type, required/optional
- field_b: input type, required/optional
Kebutuhan UI:
- DataTable: ya/tidak
- DataDialog: ya/tidak
- Tabs: ya/tidak
- SearchSelect/Combobox: ya/tidak
- Mobile card view: ya/tidak
Backend endpoint: sudah ada | belum ada
```

Contoh prompt siap pakai:

```text
Scope: frontend-page
Goal: buat halaman baru
Page code: NM01
Page purpose: CRUD Notification
Page type: CRUD standar
Primary data: notification
Field form:
- code: text, required
- message: textarea, required
- severity: select, required
- is_active: checkbox, optional
Kebutuhan UI:
- DataTable: ya
- DataDialog: ya
- Tabs: tidak
- SearchSelect/Combobox: tidak
- Mobile card view: ya
Backend endpoint: sudah ada
```

## Prinsip eksekusi

- Ambil konteks minimum yang benar-benar dibutuhkan.
- Fokus pada satu halaman atau satu slice UI terlebih dahulu.
- Gunakan komponen dan pola yang sudah ada di `src/uix/`.
- Jangan memperluas scope ke backend atau modul lain jika tidak diminta.
- Setelah edit substantif pertama, langsung lakukan validasi terfokus.

## Output contract

Setelah implementasi, laporkan hasil dengan format singkat dan prinsip `token-efficient`:

```text
Scope: {kode halaman / area frontend}
Changes: {file yang berubah}
Validation: {npx tsc --noEmit atau npm run build} -> {hasil}
Notes: {blocker atau catatan penting}
```

## Workflow

### Phase 0: Sinkronisasi konteks

- Validasi kode halaman terhadap `obx_docs/blueprint/{KODE}/README.md`.
- Identifikasi apakah halaman ini baru atau perubahan atas halaman existing.
- Identifikasi anchor terdekat: file page saat ini, halaman pembanding, atau komponen UIX yang akan dipakai.
- Jika backend endpoint belum ada, berhenti dan arahkan ke `module-creator` atau minta konfirmasi scope.

### Phase 1: Kumpulkan input minimum

Kumpulkan hanya informasi yang benar-benar perlu:

1. **Kode halaman**: `{KODE}`
2. **Tujuan halaman**: CRUD / dashboard / detail / multi-tab
3. **Data source**: endpoint atau bentuk data utama
4. **Field form**: nama field, tipe input, required/optional
5. **UI pattern**: DataTable, DataDialog, Tabs, SearchSelect/Combobox, mobile card view
6. **Status backend**: sudah ada / belum ada

Aturan clean-execution untuk phase ini:

- Jangan bertanya ulang jika endpoint, field, atau pola UI sudah jelas dari halaman pembanding atau kode yang ada.
- Jika hanya satu dialog atau satu tab yang berubah, fokus di slice itu saja.
- Jika kebutuhan UI bertentangan dengan standar proyek, gunakan standar proyek dan minta konfirmasi hanya bila benar-benar perlu.

### Phase 2: Rancang struktur halaman

Tentukan secara singkat:

- apakah memakai `DataTable + DataDialog`
- apakah perlu `Tabs`
- field input utama dan komponen UI yang dipakai
- mobile view dengan `useIsMobile()` dan `DataTableCard` bila relevan
- loading, error, dan submit flow

Jika request cukup besar, laporkan rancangan singkat sebelum implementasi. Jika request lokal, boleh lanjut tanpa design doc formal.

### Phase 3: Implementasi

Implementasi WAJIB mengikuti standar berikut:

- Jika task memengaruhi halaman yang dijelaskan di dokumentasi referensi, lengkapi atau perbarui dokumen yang relevan (khususnya `obx_docs/blueprint/{KODE}/README.md` dan `obx_docs/guide/{KODE}.md`) bersama implementasi.

- Gunakan `clientApi()` untuk fetch
- Gunakan `toast` + `getErrorMessage()` untuk error user-facing
- Gunakan `zod` + `react-hook-form`
- Gunakan `defaultValues` di `useForm()`
- Gunakan `value={field.value ?? ""}` untuk `Input` dan `Textarea` optional
- Gunakan `DataTable` untuk listing dan `DataDialog` untuk create/update/detail
- Gunakan `FieldError` untuk validasi
- Gunakan `SearchSelect` untuk searchable dropdown; gunakan `Combobox` hanya jika memang sesuai pattern proyek
- Gunakan `useIsMobile()` dan `DataTableCard` untuk halaman list yang perlu mobile card view
- Terapkan loading state pattern untuk semua write action

### Phase 4: Verifikasi

Lakukan validasi sesempit mungkin lebih dulu:

1. `npx tsc --noEmit` dari `obx_site/` untuk perubahan terfokus
2. `npm run build` dari `obx_site/` jika perubahan cukup besar atau menyentuh shared UIX

Jangan mengklaim selesai tanpa hasil validasi yang segar.

## Gaya komunikasi

- Gunakan bahasa yang singkat dan operasional.
- Laporkan perubahan per slice, bukan per spekulasi.
- Saat ada blocker backend, jelaskan scope dengan tegas.

## Constraints

- Jangan membuat page CRUD terpisah jika pola `DataTable + DataDialog` sudah cukup.
- Jangan pakai elemen HTML mentah jika komponen di `src/uix/` sudah tersedia.
- Jangan abaikan mobile compatibility untuk halaman list.
- Jangan hardcode error handling atau fetch pattern di luar standar proyek.
- Jangan memperluas perubahan ke backend kecuali user memang meminta atau scope sudah mencakupnya.
- Jika ternyata endpoint backend belum ada, hentikan implementasi frontend-only dan jelaskan bahwa scope perlu dipindah ke `module-creator` atau task backend terpisah.
