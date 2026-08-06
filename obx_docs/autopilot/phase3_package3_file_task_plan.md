# Rencana Task File-by-File Paket 3 Fase 3

Dokumen ini memecah Paket 3 Fase 3 menjadi task file-by-file untuk klaster Inventory MikroTik dan Polling.

## Scope Paket 3

- Modul domain: net_device, net_interface, net_poll_sample.
- Tujuan: monitoring network dasar siap dipakai dashboard dan alert.

## Batch Task A - Finalisasi Schema Network Monitoring

### File Target
- obx_base/prisma/schema/ict_mikrotik.prisma
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_base/prisma/schema/schema.prisma

### Verifikasi
- npx prisma generate di obx_base lulus.

## Batch Task B - API Device, Interface, Polling

### File Target
- obx_rest/backbone/routes.go
- obx_rest/skeleton/{KODE_NET_MON}/template.go
- obx_rest/skeleton/{KODE_NET_MON}/repository.go
- obx_rest/skeleton/{KODE_NET_MON}/usecase.go
- obx_rest/skeleton/{KODE_NET_MON}/handler.go

### Verifikasi
- go build ./... di obx_rest lulus.
- Device, interface, dan poll sample tersimpan.

## Batch Task C - Frontend Board Network Monitoring

### File Target
- obx_site/src/app/board/pages/{KODE_NET_MON}/page.tsx
- obx_site/src/app/board/model/module.ts

### Verifikasi
- npm run lint di obx_site lulus.
- UI menampilkan device, interface, dan polling sample.

## Batch Task D - Sinkronisasi Dokumentasi

### File Target
- obx_docs/blueprint/BASE/ict_mikrotik.md
- obx_docs/guide/BASE/ict_mikrotik.md
- obx_docs/blueprint/PLAN/README.md

### Verifikasi
- Dokumen konsisten dengan route dan payload API.

## Referensi

- phase3_backend_coding_plan.md
- phase3_backend_execution_backlog.md
- phase3_execution_backlog.md
- ai_runbook.md
