# Peta File Nyata per Modul Paket 7 Fase 2

Dokumen ini memecah Paket 7 Fase 2 menjadi peta file nyata per modul domain RDP Pilot dan Recovery Controls bastion.

## Tujuan

1. Menjadikan implementasi RDP pilot siap dieksekusi per modul.
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

## jms_connect_token untuk RDP

### Backend
- obx_rest/skeleton/{KODE_BASTION_RDP}/template.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/repository.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/handler.go
- obx_rest/backbone/routes.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_RDP}/page.tsx
- obx_site/src/app/board/model/module.ts

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_session untuk RDP

### Backend
- obx_rest/skeleton/{KODE_BASTION_RDP}/template.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/repository.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_RDP}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_session_event untuk RDP

### Backend
- obx_rest/skeleton/{KODE_BASTION_RDP}/template.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/repository.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_RDP}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Urutan Eksekusi yang Disarankan

1. Kerjakan token RDP terlebih dahulu.
2. Lanjutkan session RDP untuk lifecycle koneksi.
3. Tambahkan session event RDP dan recovery note.
4. Sinkronkan dokumentasi BASE dan PLAN pada batch yang sama.

## Risiko

1. Revoke tidak selalu memutus koneksi runtime RDP.
2. Event RDP tidak konsisten untuk timeout dan deny.
3. Menu board belum sinkron dengan modul RDP bastion.

## Mitigasi

1. Wajibkan jalur forced disconnect pada revoke.
2. Standarkan event CONNECT, DISCONNECT, TIMEOUT, DENY.
3. Update module.ts bersamaan dengan page dan route.

## Referensi

- phase2_package7_file_task_plan.md
- phase2_package7_actual_file_map.md
- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
