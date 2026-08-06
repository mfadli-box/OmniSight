# Peta File Nyata per Modul Paket 6 Fase 2

Dokumen ini memecah Paket 6 Fase 2 menjadi peta file nyata per modul domain WebAppProxy bastion.

## Tujuan

1. Menjadikan implementasi proxy aplikasi internal siap dieksekusi per modul.
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

## jms_web_app

### Backend
- obx_rest/skeleton/{KODE_BASTION_PROXY}/template.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/repository.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/handler.go
- obx_rest/backbone/routes.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_PROXY}/page.tsx
- obx_site/src/app/board/model/module.ts

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_connect_token untuk WEBAPP

### Backend
- obx_rest/skeleton/{KODE_BASTION_PROXY}/template.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/repository.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_PROXY}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Access Log PROXY_ACCESS

### Backend
- obx_rest/skeleton/{KODE_BASTION_PROXY}/template.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/repository.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_PROXY}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Urutan Eksekusi yang Disarankan

1. Kerjakan jms_web_app sebagai registry aplikasi.
2. Lanjutkan token WEBAPP untuk akses sementara.
3. Tambahkan access log PROXY_ACCESS.
4. Sinkronkan dokumentasi BASE dan PLAN pada batch yang sama.

## Risiko

1. Target proxy di luar allowlist tetap lolos.
2. Token WEBAPP tidak dibatasi one-time.
3. Menu board belum sinkron dengan modul proxy bastion.

## Mitigasi

1. Wajibkan allowlist host dan path di backend.
2. Terapkan validasi status token sebelum forward request.
3. Update module.ts bersamaan dengan page dan route.

## Referensi

- phase2_package6_file_task_plan.md
- phase2_package6_actual_file_map.md
- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
