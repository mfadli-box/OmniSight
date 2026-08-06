# Peta File Nyata Paket 5 Fase 2

Dokumen ini memetakan file nyata untuk Paket 5 Fase 2 pada klaster File Transfer browser-based bastion.

## Scope

- jms_file_transfer
- jms_session_event terkait transfer

## Schema

- obx_base/prisma/schema/jsm_stack.prisma
- obx_base/prisma/schema/schema.prisma
- obx_base/prisma/schema/all_enum_dat.prisma

## Backend

### Shared
- obx_rest/backbone/routes.go

### Skeleton Target
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/template.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/repository.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/handler.go

### Reference Pattern
- obx_rest/skeleton/SM03/template.go
- obx_rest/skeleton/SM03/repository.go
- obx_rest/skeleton/SM03/usecase.go
- obx_rest/skeleton/SM03/handler.go
- obx_rest/skeleton/SM04/template.go
- obx_rest/skeleton/SM04/repository.go
- obx_rest/skeleton/SM04/usecase.go
- obx_rest/skeleton/SM04/handler.go

## Frontend

### Target
- obx_site/src/app/board/pages/{KODE_BASTION_TRANSFER}/page.tsx
- obx_site/src/app/board/model/module.ts

### Reference Pattern
- obx_site/src/app/board/pages/SM03/page.tsx
- obx_site/src/app/board/pages/SM04/page.tsx
- obx_site/src/app/board/pages/XX99/page.tsx

## Documentation

### BASE Domain
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

### Plan and Tracking
- obx_docs/blueprint/PLAN/README.md
- obx_docs/autopilot/phase2_execution_backlog.md
- obx_docs/autopilot/phase2_backend_execution_backlog.md
- obx_docs/autopilot/phase2_backend_coding_plan.md
- obx_docs/autopilot/phase2_package5_file_task_plan.md

## Verifikasi yang Disarankan

1. npx prisma generate di obx_base.
2. go build ./... di obx_rest.
3. npm run lint di obx_site.
4. Simulasi upload dan download menghasilkan metadata audit utuh.

## Catatan

1. Placeholder KODE_BASTION_TRANSFER perlu ditetapkan ke kode modul final sebelum implementasi.
2. Transfer event harus terkait session bila sesi aktif tersedia.
3. Company scope wajib diterapkan konsisten pada query repository.

## Referensi

- phase2_package5_file_task_plan.md
- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
- ai_runbook.md
