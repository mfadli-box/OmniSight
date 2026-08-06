# Peta File Nyata per Modul Paket 4 Fase 2

Dokumen ini memecah Paket 4 Fase 2 menjadi peta file nyata per modul domain Web SSH Bridge bastion.

## Tujuan

1. Menjadikan implementasi Web SSH siap dieksekusi per modul.
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

## Token SSH

### Backend
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/template.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/repository.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/handler.go
- obx_rest/backbone/routes.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_CONNECT}/page.tsx
- obx_site/src/app/board/model/module.ts

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Session SSH

### Backend
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/template.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/repository.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_CONNECT}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Event SSH

### Backend
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/template.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/repository.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_CONNECT}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Urutan Eksekusi yang Disarankan

1. Kerjakan token SSH terlebih dahulu.
2. Lanjutkan session SSH untuk lifecycle koneksi.
3. Tambahkan event SSH dan forced disconnect.
4. Sinkronkan dokumentasi BASE dan PLAN pada batch yang sama.

## Risiko

1. Token SSH dipakai ulang di luar batas expiry.
2. Session SSH tidak menulis event penutup.
3. Menu board belum sinkron dengan modul connect bastion.

## Mitigasi

1. Wajibkan one-time token dengan validasi status.
2. Wajibkan event DISCONNECT atau TIMEOUT saat sesi berakhir.
3. Update module.ts bersamaan dengan page dan route.

## Referensi

- phase2_package4_file_task_plan.md
- phase2_package4_actual_file_map.md
- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
