# Peta File Nyata per Modul Paket 2 Fase 2

Dokumen ini memecah Paket 2 Fase 2 menjadi peta file nyata per modul domain Session Core dan Audit bastion.

## Tujuan

1. Menjadikan implementasi session audit siap dieksekusi per modul.
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

## jms_session

### Backend
- obx_rest/skeleton/{KODE_BASTION_SESSION}/template.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/repository.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/handler.go
- obx_rest/backbone/routes.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_SESSION}/page.tsx
- obx_site/src/app/board/model/module.ts

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_session_event

### Backend
- obx_rest/skeleton/{KODE_BASTION_SESSION}/template.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/repository.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_SESSION}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_session_command

### Backend
- obx_rest/skeleton/{KODE_BASTION_SESSION}/template.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/repository.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_SESSION}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_audit_log

### Backend
- obx_rest/skeleton/{KODE_BASTION_SESSION}/template.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/repository.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_SESSION}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_SESSION}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Urutan Eksekusi yang Disarankan

1. Kerjakan jms_session sebagai fondasi lifecycle.
2. Lanjutkan jms_session_event dan jms_session_command.
3. Tambahkan jms_audit_log untuk sinkronisasi jejak audit.
4. Sinkronkan dokumentasi BASE dan PLAN pada batch yang sama.

## Risiko

1. Event audit tidak sinkron dengan status session.
2. Command log tidak selalu terhubung ke session aktif.
3. Menu board belum sinkron dengan modul session bastion.

## Mitigasi

1. Wajibkan event CONNECT dan DISCONNECT untuk setiap sesi.
2. Validasi foreign key session sebelum menyimpan command.
3. Update module.ts bersamaan dengan page dan route.

## Referensi

- phase2_package2_file_task_plan.md
- phase2_package2_actual_file_map.md
- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
