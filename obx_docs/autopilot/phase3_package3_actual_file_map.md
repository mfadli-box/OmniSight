# Peta File Nyata Paket 3 Fase 3

Dokumen ini memetakan file nyata untuk Paket 3 Fase 3 pada klaster Inventory MikroTik dan Polling.

## Scope

- net_device
- net_interface
- net_poll_sample

## Schema

- obx_base/prisma/schema/ict_mikrotik.prisma
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_base/prisma/schema/schema.prisma

## Backend

### Shared
- obx_rest/backbone/routes.go

### Skeleton Target
- obx_rest/skeleton/{KODE_NET_MON}/template.go
- obx_rest/skeleton/{KODE_NET_MON}/repository.go
- obx_rest/skeleton/{KODE_NET_MON}/usecase.go
- obx_rest/skeleton/{KODE_NET_MON}/handler.go

### Reference Pattern
- obx_rest/skeleton/SM03/template.go
- obx_rest/skeleton/SM03/repository.go
- obx_rest/skeleton/SM03/usecase.go
- obx_rest/skeleton/SM03/handler.go
- obx_rest/skeleton/XX99/template.go
- obx_rest/skeleton/XX99/repository.go
- obx_rest/skeleton/XX99/usecase.go
- obx_rest/skeleton/XX99/handler.go

## Frontend

### Target
- obx_site/src/app/board/pages/{KODE_NET_MON}/page.tsx
- obx_site/src/app/board/model/module.ts

### Reference Pattern
- obx_site/src/app/board/pages/SM03/page.tsx
- obx_site/src/app/board/pages/XX99/page.tsx

## Documentation

### BASE Domain
- obx_docs/blueprint/BASE/ict_mikrotik.md
- obx_docs/guide/BASE/ict_mikrotik.md

### Plan and Tracking
- obx_docs/blueprint/PLAN/README.md
- obx_docs/autopilot/phase3_execution_backlog.md
- obx_docs/autopilot/phase3_backend_execution_backlog.md
- obx_docs/autopilot/phase3_backend_coding_plan.md
- obx_docs/autopilot/phase3_package3_file_task_plan.md

## Referensi

- phase3_package3_file_task_plan.md
- phase3_backend_coding_plan.md
- phase3_backend_execution_backlog.md
- phase3_execution_backlog.md
- ai_runbook.md
