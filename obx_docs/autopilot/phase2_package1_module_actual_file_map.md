# Peta File Nyata per Modul Paket 1 Fase 2

Dokumen ini memecah Paket 1 Fase 2 menjadi peta file nyata per modul domain inventory bastion.

## Tujuan

1. Menjadikan implementasi inventory bastion siap dieksekusi per modul.
2. Menyediakan daftar file nyata untuk schema, backend, frontend, dan dokumentasi.
3. Meminimalkan scope agar review dan verifikasi per modul lebih mudah.

## File Bersama yang Hampir Pasti Tersentuh

- obx_base/prisma/schema/jsm_stack.prisma
- obx_base/prisma/schema/schema.prisma
- obx_base/prisma/schema/all_enum_dat.prisma
- obx_rest/backbone/routes.go
- obx_site/src/app/board/model/module.ts
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md
- obx_docs/blueprint/PLAN/README.md
- obx_docs/autopilot/README.md

## jms_asset_group

### Backend
- obx_rest/skeleton/{KODE_BASTION_ASSET}/template.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/repository.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/handler.go
- obx_rest/backbone/routes.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_ASSET}/page.tsx
- obx_site/src/app/board/model/module.ts

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_asset

### Backend
- obx_rest/skeleton/{KODE_BASTION_ASSET}/template.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/repository.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/handler.go
- obx_rest/backbone/routes.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_ASSET}/page.tsx
- obx_site/src/app/board/model/module.ts

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_asset_group_member

### Backend
- obx_rest/skeleton/{KODE_BASTION_ASSET}/template.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/repository.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_ASSET}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_account

### Backend
- obx_rest/skeleton/{KODE_BASTION_ASSET}/template.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/repository.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/handler.go
- obx_rest/backbone/routes.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_ASSET}/page.tsx
- obx_site/src/app/board/model/module.ts

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_account_secret

### Backend
- obx_rest/skeleton/{KODE_BASTION_ASSET}/template.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/repository.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_ASSET}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_ASSET}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Urutan Eksekusi yang Disarankan

1. Kerjakan jms_asset_group dan jms_asset terlebih dahulu.
2. Lanjutkan jms_account setelah relasi asset stabil.
3. Tambahkan jms_account_secret dengan kontrol expose ketat.
4. Sinkronkan dokumentasi BASE dan PLAN pada batch yang sama.

## Risiko

1. Relasi asset ke account tidak konsisten.
2. Secret terekspos pada endpoint list atau detail.
3. Menu board belum sinkron dengan modul bastion asset.

## Mitigasi

1. Validasi relasi sejak layer schema dan repository.
2. Pisahkan jalur secret dari payload publik.
3. Update module.ts bersamaan dengan page dan route.

## Referensi

- phase2_package1_file_task_plan.md
- phase2_package1_actual_file_map.md
- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
