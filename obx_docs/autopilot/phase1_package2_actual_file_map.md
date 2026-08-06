# Peta File Nyata - Paket 2 Fase 1

Dokumen ini memetakan file nyata yang perlu disentuh saat Paket 2 Fase 1 mulai dieksekusi.

## Tujuan

1. Menjadikan workflow dokumen terkendali langsung operasional di repo.
2. Memecah perubahan menjadi file nyata per lapisan.
3. Meminimalkan scope agar review, verifikasi, dan audit lebih mudah.

## File Bersama yang Hampir Pasti Tersentuh

- `obx_base/prisma/schema/schema.prisma`
- `obx_base/prisma/schema/dat_document.prisma`
- `obx_base/prisma/schema/dat_signature.prisma`
- `obx_docs/autopilot/README.md`
- `obx_docs/blueprint/BASE/README.md`
- `obx_docs/guide/README.md`
- `obx_docs/blueprint/PLAN/README.md`

## dat_request

### Schema
- `obx_base/prisma/schema/schema.prisma`
- `obx_base/prisma/schema/all_enum_dat.prisma`
- `obx_base/prisma/schema/dat_document.prisma`

### Backend
- `obx_rest/skeleton/SM04/template.go`
- `obx_rest/skeleton/SM04/repository.go`
- `obx_rest/skeleton/SM04/usecase.go`
- `obx_rest/skeleton/SM04/handler.go`
- `obx_rest/backbone/routes.go`

### Frontend
- `obx_site/src/app/board/pages/SM04/page.tsx`
- `obx_site/src/app/board/model/module.ts`

### Documentation
- `obx_docs/blueprint/BASE/dat_document.md`
- `obx_docs/guide/BASE/dat_document.md`
- `obx_docs/blueprint/BASE/dat_document.md`
- `obx_docs/guide/BASE/dat_document.md`

## dat_signature

### Schema
- `obx_base/prisma/schema/schema.prisma`
- `obx_base/prisma/schema/all_enum_dat.prisma`
- `obx_base/prisma/schema/dat_signature.prisma`

### Backend
- `obx_rest/skeleton/SM04/template.go`
- `obx_rest/skeleton/SM04/repository.go`
- `obx_rest/skeleton/SM04/usecase.go`
- `obx_rest/skeleton/SM04/handler.go`
- `obx_rest/backbone/routes.go`

### Frontend
- `obx_site/src/app/board/pages/SM04/page.tsx`
- `obx_site/src/app/board/model/module.ts`

### Documentation
- `obx_docs/blueprint/BASE/dat_signature.md`
- `obx_docs/guide/BASE/dat_signature.md`
- `obx_docs/blueprint/BASE/dat_document.md`
- `obx_docs/guide/BASE/dat_document.md`

## dat_document

### Schema
- `obx_base/prisma/schema/schema.prisma`
- `obx_base/prisma/schema/dat_document.prisma`
- `obx_base/prisma/schema/all_enum_dat.prisma`

### Backend
- `obx_rest/skeleton/SM04/template.go`
- `obx_rest/skeleton/SM04/repository.go`
- `obx_rest/skeleton/SM04/usecase.go`
- `obx_rest/skeleton/SM04/handler.go`
- `obx_rest/backbone/routes.go`

### Frontend
- `obx_site/src/app/board/pages/SM04/page.tsx`
- `obx_site/src/app/board/model/module.ts`

### Documentation
- `obx_docs/blueprint/BASE/dat_document.md`
- `obx_docs/guide/BASE/dat_document.md`
- `obx_docs/blueprint/BASE/dat_document.md`
- `obx_docs/guide/BASE/dat_document.md`

## dat_document_revision

### Schema
- `obx_base/prisma/schema/schema.prisma`
- `obx_base/prisma/schema/dat_document.prisma`

### Backend
- `obx_rest/skeleton/SM04/template.go`
- `obx_rest/skeleton/SM04/repository.go`
- `obx_rest/skeleton/SM04/usecase.go`
- `obx_rest/skeleton/SM04/handler.go`

### Documentation
- `obx_docs/blueprint/BASE/dat_document.md`
- `obx_docs/guide/BASE/dat_document.md`
- `obx_docs/blueprint/BASE/dat_document.md`
- `obx_docs/guide/BASE/dat_document.md`

## dat_document_evidence

### Schema
- `obx_base/prisma/schema/schema.prisma`
- `obx_base/prisma/schema/dat_document.prisma`

### Backend
- `obx_rest/skeleton/SM04/template.go`
- `obx_rest/skeleton/SM04/repository.go`
- `obx_rest/skeleton/SM04/usecase.go`
- `obx_rest/skeleton/SM04/handler.go`

### Documentation
- `obx_docs/blueprint/BASE/dat_document.md`
- `obx_docs/guide/BASE/dat_document.md`
- `obx_docs/autopilot/template_iso_evidence_checklist.md`
- `obx_docs/autopilot/template_approval_matrix.md`

## Urutan Eksekusi yang Disarankan

1. Update schema Prisma untuk dokumen dan signature.
2. Sinkronkan backend SM04 dan route terkait.
3. Sinkronkan page SM04 dan menu board.
4. Update dokumentasi BASE, template, dan guide.
5. Jalankan verifikasi build dan smoke check approval.

## Risiko

1. Relasi schema dokumen dan approval tidak konsisten.
2. Backend dan frontend SM04 berubah tanpa dokumentasi yang mengikuti.
3. Evidence trail tidak sinkron dengan alur approval.

## Mitigasi

1. Kerjakan schema, backend, frontend, dan dokumentasi dalam batch kecil.
2. Update file SM04 sebagai satu paket bersama dokumen BASE terkait.
3. Verifikasi setelah setiap batch perubahan.

## Referensi

- phase1_package2_file_task_plan.md
- phase1_backend_coding_plan.md
- phase1_backend_execution_backlog.md
- phase1_execution_backlog.md
- ai_runbook.md
