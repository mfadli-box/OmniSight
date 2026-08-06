# Peta File Nyata Paket 1 Fase 2

Dokumen ini memetakan file nyata untuk Paket 1 Fase 2 pada klaster Asset dan Account Inventory bastion.

## Scope

- jms_asset_group
- jms_asset
- jms_asset_group_member
- jms_account
- jms_account_secret

## Schema

- obx_base/prisma/schema/jsm_stack.prisma
- obx_base/prisma/schema/schema.prisma
- obx_base/prisma/schema/all_enum_dat.prisma

## Backend

### Shared
- obx_rest/backbone/routes.go

### Skeleton Target
- obx_rest/skeleton/{KODE_BASTION_ASSET}/template.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/repository.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/handler.go

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
- obx_site/src/app/board/pages/{KODE_BASTION_ASSET}/page.tsx
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
- obx_docs/autopilot/phase2_package1_file_task_plan.md

## Verifikasi yang Disarankan

1. npx prisma generate di obx_base.
2. go build ./... di obx_rest.
3. npm run lint di obx_site.
4. Smoke test list create update untuk asset dan account.

## Catatan

1. Placeholder KODE_BASTION_ASSET perlu ditetapkan ke kode modul final sebelum implementasi.
2. Secret value tidak boleh diekspos pada payload list atau detail publik.
3. Company scope wajib diterapkan konsisten pada query repository.

## Referensi

- phase2_package1_file_task_plan.md
- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
- ai_runbook.md
