# Peta File Nyata - Paket 3 Fase 1

Dokumen ini memetakan file nyata yang perlu disentuh saat Paket 3 Fase 1 mulai dieksekusi.

## Tujuan

1. Menjadikan asset inventory dasar siap dipakai sebagai fondasi Bastion dan monitoring.
2. Memecah perubahan menjadi file nyata per lapisan.
3. Meminimalkan scope agar review, verifikasi, dan audit lebih mudah.

## File Bersama yang Hampir Pasti Tersentuh

- `obx_base/prisma/schema/ict_machine.prisma`
- `obx_base/prisma/schema/ict_website.prisma`
- `obx_base/prisma/schema/ict_security.prisma`
- `obx_base/prisma/schema/ict_mikrotik.prisma`
- `obx_base/prisma/schema/ict_monitoring.prisma`
- `obx_docs/autopilot/README.md`
- `obx_docs/blueprint/BASE/README.md`
- `obx_docs/guide/README.md`

## Asset Host / Machine

### Schema
- `obx_base/prisma/schema/ict_machine.prisma`

### Backend Reference
- `obx_rest/skeleton/XX99/template.go`
- `obx_rest/skeleton/XX99/repository.go`
- `obx_rest/skeleton/XX99/usecase.go`
- `obx_rest/skeleton/XX99/handler.go`

### Frontend Reference
- `obx_site/src/app/board/pages/XX99/page.tsx`
- `obx_site/src/app/board/model/module.ts`

### Documentation
- `obx_docs/blueprint/BASE/ict_machine.md`
- `obx_docs/guide/BASE/ict_machine.md`

## Asset Website / Service

### Schema
- `obx_base/prisma/schema/ict_website.prisma`

### Backend Reference
- `obx_rest/skeleton/XX99/template.go`
- `obx_rest/skeleton/XX99/repository.go`
- `obx_rest/skeleton/XX99/usecase.go`
- `obx_rest/skeleton/XX99/handler.go`

### Frontend Reference
- `obx_site/src/app/board/pages/XX99/page.tsx`
- `obx_site/src/app/board/model/module.ts`

### Documentation
- `obx_docs/blueprint/BASE/ict_website.md`
- `obx_docs/guide/BASE/ict_website.md`

## Asset Security / Endpoint Internal

### Schema
- `obx_base/prisma/schema/ict_security.prisma`

### Backend Reference
- `obx_rest/skeleton/XX99/template.go`
- `obx_rest/skeleton/XX99/repository.go`
- `obx_rest/skeleton/XX99/usecase.go`
- `obx_rest/skeleton/XX99/handler.go`

### Frontend Reference
- `obx_site/src/app/board/pages/XX99/page.tsx`
- `obx_site/src/app/board/model/module.ts`

### Documentation
- `obx_docs/blueprint/BASE/ict_security.md`
- `obx_docs/guide/BASE/ict_security.md`

## Asset MikroTik / Network

### Schema
- `obx_base/prisma/schema/ict_mikrotik.prisma`
- `obx_base/prisma/schema/ict_monitoring.prisma`

### Backend Reference
- `obx_rest/skeleton/XX99/template.go`
- `obx_rest/skeleton/XX99/repository.go`
- `obx_rest/skeleton/XX99/usecase.go`
- `obx_rest/skeleton/XX99/handler.go`

### Frontend Reference
- `obx_site/src/app/board/pages/XX99/page.tsx`
- `obx_site/src/app/board/model/module.ts`

### Documentation
- `obx_docs/blueprint/BASE/ict_mikrotik.md`
- `obx_docs/guide/BASE/ict_mikrotik.md`

## Urutan Eksekusi yang Disarankan

1. Review schema asset yang sudah ada di `obx_base/prisma/schema`.
2. Gunakan `XX99` sebagai bootstrap reference bila modul asset final belum dibentuk.
3. Sinkronkan dokumentasi BASE yang paling relevan.
4. Siapkan menu board dan frontend reference bila halaman asset mulai dipisah.
5. Jalankan verifikasi build dan lint setelah perubahan.

## Risiko

1. Asset inventory dibangun terlalu generik sehingga sulit dipakai ulang.
2. Referensi bootstrap XX99 dipakai terlalu lama tanpa modul final.
3. Dokumentasi tidak mengikuti perubahan schema asset.

## Mitigasi

1. Tetapkan field minimum yang benar-benar dipakai oleh Bastion dan monitoring.
2. Gunakan XX99 hanya sebagai referensi bootstrap sementara.
3. Update schema dan dokumentasi dalam batch yang sama.

## Referensi

- phase1_package3_file_task_plan.md
- phase1_backend_coding_plan.md
- phase1_backend_execution_backlog.md
- phase1_execution_backlog.md
- ai_runbook.md
