# Peta File Nyata Paket 2 Fase 3

Dokumen ini memetakan file nyata untuk Paket 2 Fase 3 pada klaster Inventory Nginx.

## Scope

- web_site
- web_upstream
- web_certificate
- web_config_version
- web_reload_history

## Schema

- obx_base/prisma/schema/ict_website.prisma
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_base/prisma/schema/schema.prisma

## Backend

### Shared
- obx_rest/backbone/routes.go

### Skeleton Target
- obx_rest/skeleton/{KODE_WEB_MGMT}/template.go
- obx_rest/skeleton/{KODE_WEB_MGMT}/repository.go
- obx_rest/skeleton/{KODE_WEB_MGMT}/usecase.go
- obx_rest/skeleton/{KODE_WEB_MGMT}/handler.go

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
- obx_site/src/app/board/pages/{KODE_WEB_MGMT}/page.tsx
- obx_site/src/app/board/model/module.ts

### Reference Pattern
- obx_site/src/app/board/pages/SM03/page.tsx
- obx_site/src/app/board/pages/XX99/page.tsx

## Documentation

### BASE Domain
- obx_docs/blueprint/BASE/ict_website.md
- obx_docs/guide/BASE/ict_website.md

### Plan and Tracking
- obx_docs/blueprint/PLAN/README.md
- obx_docs/autopilot/phase3_execution_backlog.md
- obx_docs/autopilot/phase3_backend_execution_backlog.md
- obx_docs/autopilot/phase3_backend_coding_plan.md
- obx_docs/autopilot/phase3_package2_file_task_plan.md

## Referensi

- phase3_package2_file_task_plan.md
- phase3_backend_coding_plan.md
- phase3_backend_execution_backlog.md
- phase3_execution_backlog.md
- ai_runbook.md
