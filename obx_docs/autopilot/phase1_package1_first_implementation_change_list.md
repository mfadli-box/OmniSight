# Daftar Perubahan Implementasi Pertama Paket 1 Fase 1

Dokumen ini menjadi jembatan dari perencanaan ke eksekusi coding untuk modul inti master data pada Paket 1 Fase 1.

## Prasyarat

1. Kode modul final untuk placeholder SM01, SM02, SM03, SM04, SM05, SP01, SP02, SP03 telah ditetapkan.
2. Rute aktif dan middleware dasar di backend sudah dipetakan.
3. Dokumen task plan dan peta file Paket 1 Fase 1 telah ditinjau.

## Langkah Implementasi Pertama

### Langkah 1 - Verifikasi Route dan Middleware Inti

#### File
- obx_rest/backbone/routes.go

#### Acceptance Check
- Urutan route static sebelum dynamic tetap aman.
- Write route memakai middleware logging yang konsisten.

### Langkah 2 - Definisi DTO dan Interface per Modul Inti

#### File
- obx_rest/skeleton/SM01/template.go
- obx_rest/skeleton/SM02/template.go
- obx_rest/skeleton/SM03/template.go
- obx_rest/skeleton/SM04/template.go
- obx_rest/skeleton/SM05/template.go
- obx_rest/skeleton/SP01/template.go
- obx_rest/skeleton/SP02/template.go
- obx_rest/skeleton/SP03/template.go

#### Acceptance Check
- Go compiler mengenali DTO dan interface tanpa error.

### Langkah 3 - Review Repository dan Usecase Batch Pertama

#### File
- obx_rest/skeleton/SM01/repository.go
- obx_rest/skeleton/SM01/usecase.go
- obx_rest/skeleton/SM02/repository.go
- obx_rest/skeleton/SM02/usecase.go
- obx_rest/skeleton/SM03/repository.go
- obx_rest/skeleton/SM03/usecase.go

#### Acceptance Check
- Query dan validasi input konsisten pada batch awal.
- go build ./... di obx_rest lulus.

### Langkah 4 - Sinkronisasi Handler dan Page Modul Pilot

#### File
- obx_rest/skeleton/SM01/handler.go
- obx_rest/skeleton/SM02/handler.go
- obx_site/src/app/board/pages/SM01/page.tsx
- obx_site/src/app/board/pages/SM02/page.tsx
- obx_site/src/app/board/model/module.ts

#### Acceptance Check
- Smoke test login, profile, dan company flow lulus.
- npm run lint di obx_site lulus.

### Langkah 5 - Sinkronisasi Dokumentasi Batch Pertama

#### File
- obx_docs/blueprint/REST/SM01.md
- obx_docs/guide/REST/SM01.md
- obx_docs/blueprint/SITE/SM01.md
- obx_docs/guide/SITE/SM01.md
- obx_docs/blueprint/REST/SM02.md
- obx_docs/guide/REST/SM02.md
- obx_docs/blueprint/SITE/SM02.md
- obx_docs/guide/SITE/SM02.md

#### Acceptance Check
- Dokumen teknis dan guide sinkron dengan implementasi batch pertama.

## Referensi

- phase1_package1_file_task_plan.md
- phase1_package1_actual_file_map.md
- phase1_package1_module_actual_file_map.md
- phase1_backend_coding_plan.md
- ai_runbook.md

## Status Eksekusi Saat Ini

1. Langkah 1 selesai: route SM01 diurutkan agar static route didaftarkan sebelum route dinamis parameter.
2. Validasi backend lulus: `cd obx_rest && go build ./...`.
3. Validasi frontend terfokus lulus: `npx eslint` untuk file page SM01, SM02, SM03, SP01, SP02, SP03 menghasilkan `SCOPED_LINT_OK`.
4. Catatan: lint global workspace masih memiliki technical debt di file lain di luar scope Paket 1 batch ini.
5. Review DTO/interface dan repository-usecase SM01-SM03 selesai.
6. Perbaikan backend diterapkan pada sort mapping `SM03` agar `ORDER BY` memakai kolom eksplisit (`m.code`, `m.name`, `cm.created_at`) dan menghindari ambigu.
7. Validasi backend pasca-perbaikan lulus: `cd obx_rest && go build ./...`.
8. Hardening handler SM01-SM03 selesai: validasi parameter path wajib (`id`, `companyId`, `moduleId`, `areaId`, `privilegeId`) ditambahkan pada handler.
9. Sinkronisasi dokumentasi batch pertama selesai untuk SM01 dan SM02 (REST + SITE, blueprint + guide).
10. Catatan: smoke test manual login/profile/company flow belum dijalankan pada sesi ini.
11. Hardening tambahan selesai pada handler SP01-SP03 untuk menghindari panic type assertion context session (`userId`, `isAdmin`).
12. Sinkronisasi dokumentasi lanjutan selesai untuk SM03 serta SP01, SP02, SP03 (REST + SITE, blueprint + guide).
13. Technical debt lint global frontend telah dibersihkan; `npm run lint` di `obx_site` lulus tanpa temuan.
14. Validasi frontend produksi lulus: `cd obx_site && npm run build`.
15. Sinkronisasi dokumentasi batch lanjutan selesai untuk SM04 dan SM05 (REST + SITE, blueprint + guide).
16. Smoke test autentik berhasil untuk login UI admin dengan akun `admin` dan password yang diberikan pada sesi ini.
17. Smoke test halaman inti berhasil: SP01, SP02, SP03, SM01, SM02, SM03, SM04, dan SM05 dapat diakses setelah session aktif.
18. Smoke test data demo berhasil: satu demo company dibuat pada SM03 dan satu demo user dibuat pada SM01, lalu relasi company user tervalidasi di UI.
19. Acceptance check Paket 1 terpenuhi: backend build lulus, frontend lint lulus, frontend build lulus, serta smoke test login, profile, company, privilege dasar, dan session berhasil.
20. Status akhir: implementasi pertama Paket 1 Fase 1 selesai secara teknis dan siap ditutup secara administratif.
