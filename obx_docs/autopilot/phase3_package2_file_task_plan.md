# Rencana Task File-by-File Paket 2 Fase 3

Dokumen ini memecah Paket 2 Fase 3 menjadi task file-by-file untuk klaster Inventory Nginx.

## Scope Paket 2

- Modul domain: web_site, web_upstream, web_certificate, web_config_version, web_reload_history.
- Tujuan: inventory site dan versioning config siap untuk manajemen runtime.

## Batch Task A - Finalisasi Schema Web Management

### File Target
- obx_base/prisma/schema/ict_website.prisma
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_base/prisma/schema/schema.prisma

### Verifikasi
- npx prisma generate di obx_base lulus.

## Batch Task B - API Inventory dan History

### File Target
- obx_rest/backbone/routes.go
- obx_rest/skeleton/{KODE_WEB_MGMT}/template.go
- obx_rest/skeleton/{KODE_WEB_MGMT}/repository.go
- obx_rest/skeleton/{KODE_WEB_MGMT}/usecase.go
- obx_rest/skeleton/{KODE_WEB_MGMT}/handler.go

### Verifikasi
- go build ./... di obx_rest lulus.
- Site, upstream, certificate, version, dan reload history tercatat.

## Batch Task C - Frontend Board Web Management

### File Target
- obx_site/src/app/board/pages/{KODE_WEB_MGMT}/page.tsx
- obx_site/src/app/board/model/module.ts

### Verifikasi
- npm run lint di obx_site lulus.
- UI menampilkan inventory dan histori reload.

## Batch Task D - Sinkronisasi Dokumentasi

### File Target
- obx_docs/blueprint/BASE/ict_website.md
- obx_docs/guide/BASE/ict_website.md
- obx_docs/blueprint/PLAN/README.md

### Verifikasi
- Dokumen konsisten dengan route dan payload API.

## Referensi

- phase3_backend_coding_plan.md
- phase3_backend_execution_backlog.md
- phase3_execution_backlog.md
- ai_runbook.md
