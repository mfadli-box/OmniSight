# Rencana Task File-by-File Paket 5 Fase 3

Dokumen ini memecah Paket 5 Fase 3 menjadi task file-by-file untuk klaster Event Store dan SIEM Fondasi.

## Scope Paket 5

- Modul domain: sec_event, sec_event_source, sec_event_parser, sec_rule, sec_alert, sec_incident.
- Tujuan: fondasi ingest, normalisasi, korelasi, dan incident timeline aktif secara bertahap.

## Batch Task A - Finalisasi Schema SIEM Core

### File Target
- obx_base/prisma/schema/ict_security.prisma
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_base/prisma/schema/schema.prisma

### Verifikasi
- npx prisma generate di obx_base lulus.

## Batch Task B - API Event, Rule, Alert, Incident

### File Target
- obx_rest/backbone/routes.go
- obx_rest/skeleton/{KODE_SIEM_CORE}/template.go
- obx_rest/skeleton/{KODE_SIEM_CORE}/repository.go
- obx_rest/skeleton/{KODE_SIEM_CORE}/usecase.go
- obx_rest/skeleton/{KODE_SIEM_CORE}/handler.go

### Verifikasi
- go build ./... di obx_rest lulus.
- Event dari satu sumber bisa ingest dan membentuk alert.

## Batch Task C - Frontend Board SIEM Core

### File Target
- obx_site/src/app/board/pages/{KODE_SIEM_CORE}/page.tsx
- obx_site/src/app/board/model/module.ts

### Verifikasi
- npm run lint di obx_site lulus.
- UI menampilkan event, alert, dan incident timeline.

## Batch Task D - Sinkronisasi Dokumentasi

### File Target
- obx_docs/blueprint/BASE/ict_security.md
- obx_docs/guide/BASE/ict_security.md
- obx_docs/blueprint/PLAN/README.md

### Verifikasi
- Dokumen konsisten dengan route dan payload API.

## Referensi

- phase3_backend_coding_plan.md
- phase3_backend_execution_backlog.md
- phase3_execution_backlog.md
- ai_runbook.md
