# Daftar Perubahan Implementasi Pertama Paket 1 Fase 3

Dokumen ini menjadi jembatan dari perencanaan ke eksekusi coding untuk inventory host, stack, dan VM pada Paket 1 Fase 3.

## Prasyarat

1. Kode modul final untuk placeholder KODE_INFRA_CORE telah ditetapkan.
2. Schema infra terbaru sudah sinkron.
3. Dokumen task plan dan peta file Paket 1 Fase 3 telah ditinjau.

## Langkah Implementasi Pertama

### Langkah 1 - Finalisasi Kontrak Data Schema

#### File
- obx_base/prisma/schema/ict_machine.prisma
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_base/prisma/schema/schema.prisma

#### Acceptance Check
- npx prisma generate di obx_base lulus.

### Langkah 2 - Definisi DTO dan Interface Backend

#### File
- obx_rest/skeleton/{KODE_INFRA_CORE}/template.go

#### Acceptance Check
- Go compiler mengenali DTO dan interface tanpa error.

### Langkah 3 - Implement Repository Layer

#### File
- obx_rest/skeleton/{KODE_INFRA_CORE}/repository.go

#### Acceptance Check
- go build ./... di obx_rest lulus.

### Langkah 4 - Implement Usecase Layer

#### File
- obx_rest/skeleton/{KODE_INFRA_CORE}/usecase.go

#### Acceptance Check
- Validasi input, pagination, dan wrapping error berjalan konsisten.

### Langkah 5 - Implement Handler dan Route

#### File
- obx_rest/skeleton/{KODE_INFRA_CORE}/handler.go
- obx_rest/backbone/routes.go

#### Acceptance Check
- Smoke test API CRUD host, stack, dan VM lulus.

### Langkah 6 - Implement Halaman Board

#### File
- obx_site/src/app/board/pages/{KODE_INFRA_CORE}/page.tsx
- obx_site/src/app/board/model/module.ts

#### Acceptance Check
- npm run lint di obx_site lulus.
- UI list/create/update inventory pilot berjalan.

### Langkah 7 - Sinkronisasi Dokumentasi

#### File
- obx_docs/blueprint/BASE/ict_machine.md
- obx_docs/guide/BASE/ict_machine.md
- obx_docs/blueprint/PLAN/README.md

#### Acceptance Check
- Dokumen sinkron dengan implementasi aktual.

## Referensi

- phase3_package1_file_task_plan.md
- phase3_package1_actual_file_map.md
- phase3_package1_module_actual_file_map.md
- phase3_backend_coding_plan.md
- ai_runbook.md
