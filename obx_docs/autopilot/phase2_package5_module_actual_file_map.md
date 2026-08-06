# Peta File Nyata per Modul Paket 5 Fase 2

Dokumen ini memecah Paket 5 Fase 2 menjadi peta file nyata per modul domain File Transfer bastion.

## Tujuan

1. Menjadikan implementasi transfer file siap dieksekusi per modul.
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

## jms_file_transfer

### Backend
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/template.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/repository.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/handler.go
- obx_rest/backbone/routes.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_TRANSFER}/page.tsx
- obx_site/src/app/board/model/module.ts

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Event FILE_UPLOAD

### Backend
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/template.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/repository.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_TRANSFER}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Event FILE_DOWNLOAD

### Backend
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/template.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/repository.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_TRANSFER}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Urutan Eksekusi yang Disarankan

1. Kerjakan jms_file_transfer sebagai data inti histori.
2. Lanjutkan event FILE_UPLOAD.
3. Tambahkan event FILE_DOWNLOAD.
4. Sinkronkan dokumentasi BASE dan PLAN pada batch yang sama.

## Risiko

1. Metadata transfer tidak lengkap untuk audit.
2. Event transfer tidak terhubung ke session.
3. Menu board belum sinkron dengan modul transfer bastion.

## Mitigasi

1. Wajibkan field file_name, transfer_type, status, dan waktu.
2. Kaitkan event transfer ke session aktif bila tersedia.
3. Update module.ts bersamaan dengan page dan route.

## Referensi

- phase2_package5_file_task_plan.md
- phase2_package5_actual_file_map.md
- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
