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
