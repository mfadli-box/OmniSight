# Peta File Nyata per Modul Paket 3 Fase 2

Dokumen ini memecah Paket 3 Fase 2 menjadi peta file nyata per modul domain Approval Access dan Connect Token bastion.

## Tujuan

1. Menjadikan implementasi approval dan token siap dieksekusi per modul.
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

## jms_connect_token

### Backend
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/template.go
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/repository.go
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/handler.go
- obx_rest/backbone/routes.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_ACCESS}/page.tsx
- obx_site/src/app/board/model/module.ts

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_policy

### Backend
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/template.go
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/repository.go
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_ACCESS}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## jms_approval

### Backend
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/template.go
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/repository.go
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_ACCESS}/handler.go

### Frontend
- obx_site/src/app/board/pages/{KODE_BASTION_ACCESS}/page.tsx

### Documentation
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md

## Urutan Eksekusi yang Disarankan

1. Kerjakan jms_policy sebagai aturan akses.
2. Lanjutkan jms_approval untuk alur keputusan akses.
3. Implement jms_connect_token setelah approval stabil.
4. Sinkronkan dokumentasi BASE dan PLAN pada batch yang sama.

## Risiko

1. Token terbit tanpa validasi approval.
2. Status token tidak sinkron dengan status request.
3. Menu board belum sinkron dengan modul approval bastion.

## Mitigasi

1. Wajibkan gate approval untuk issuance token.
2. Validasi transisi status token pada semua endpoint.
3. Update module.ts bersamaan dengan page dan route.

## Referensi

- phase2_package3_file_task_plan.md
- phase2_package3_actual_file_map.md
- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
