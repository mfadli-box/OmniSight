# Rencana Task File-by-File Paket 1 Fase 3

Dokumen ini memecah Paket 1 Fase 3 menjadi task file-by-file untuk klaster Inventory Host, Stack, dan VM.

## Scope Paket 1

- Modul domain: infra_host, infra_stack, vm_host.
- Tujuan: inventory infra inti siap dipakai dashboard dan collector.

## Batch Task A - Finalisasi Schema Infra Core

### File Target
- obx_base/prisma/schema/ict_machine.prisma
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_base/prisma/schema/schema.prisma

### Verifikasi
- npx prisma generate di obx_base lulus.

## Batch Task B - API Inventory Infra Core

### File Target
- obx_rest/backbone/routes.go
- obx_rest/skeleton/{KODE_INFRA_CORE}/template.go
- obx_rest/skeleton/{KODE_INFRA_CORE}/repository.go
- obx_rest/skeleton/{KODE_INFRA_CORE}/usecase.go
- obx_rest/skeleton/{KODE_INFRA_CORE}/handler.go

### Verifikasi
- go build ./... di obx_rest lulus.
- CRUD host, stack, dan VM berjalan.

## Batch Task C - Frontend Board Infra Core

### File Target
- obx_site/src/app/board/pages/{KODE_INFRA_CORE}/page.tsx
- obx_site/src/app/board/model/module.ts

### Verifikasi
- npm run lint di obx_site lulus.
- UI list, create, update data inventory pilot.

## Batch Task D - Sinkronisasi Dokumentasi

### File Target
- obx_docs/blueprint/BASE/ict_machine.md
- obx_docs/guide/BASE/ict_machine.md
- obx_docs/blueprint/PLAN/README.md

### Verifikasi
- Dokumen konsisten dengan route dan payload API.

## Referensi

- phase3_backend_coding_plan.md
- phase3_backend_execution_backlog.md
- phase3_execution_backlog.md
- ai_runbook.md
