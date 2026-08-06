# Peta File Nyata - Paket 1 Fase 1

Dokumen ini memetakan file nyata yang perlu disentuh saat Paket 1 Fase 1 mulai dieksekusi.

## Tujuan

1. Menjadikan Paket 1 langsung operasional di repo.
2. Memecah perubahan menjadi file nyata per modul.
3. Meminimalkan scope agar review dan verifikasi lebih mudah.

## File Bersama yang Hampir Pasti Tersentuh

- `obx_rest/backbone/routes.go`
- `obx_site/src/app/board/model/module.ts`
- `obx_docs/autopilot/README.md`
- `obx_docs/blueprint/README.md`
- `obx_docs/guide/README.md`

## SM01

### Backend
- `obx_rest/skeleton/SM01/template.go`
- `obx_rest/skeleton/SM01/repository.go`
- `obx_rest/skeleton/SM01/usecase.go`
- `obx_rest/skeleton/SM01/handler.go`

### Frontend
- `obx_site/src/app/board/pages/SM01/page.tsx`

### Documentation
- `obx_docs/blueprint/REST/SM01.md`
- `obx_docs/guide/REST/SM01.md`
- `obx_docs/blueprint/SITE/SM01.md`
- `obx_docs/guide/SITE/SM01.md`

## SM02

### Backend
- `obx_rest/skeleton/SM02/template.go`
- `obx_rest/skeleton/SM02/repository.go`
- `obx_rest/skeleton/SM02/usecase.go`
- `obx_rest/skeleton/SM02/handler.go`

### Frontend
- `obx_site/src/app/board/pages/SM02/page.tsx`

### Documentation
- `obx_docs/blueprint/REST/SM02.md`
- `obx_docs/guide/REST/SM02.md`
- `obx_docs/blueprint/SITE/SM02.md`
- `obx_docs/guide/SITE/SM02.md`

## SM03

### Backend
- `obx_rest/skeleton/SM03/template.go`
- `obx_rest/skeleton/SM03/repository.go`
- `obx_rest/skeleton/SM03/usecase.go`
- `obx_rest/skeleton/SM03/handler.go`

### Frontend
- `obx_site/src/app/board/pages/SM03/page.tsx`

### Documentation
- `obx_docs/blueprint/REST/SM03.md`
- `obx_docs/guide/REST/SM03.md`
- `obx_docs/blueprint/SITE/SM03.md`
- `obx_docs/guide/SITE/SM03.md`

## SM04

### Backend
- `obx_rest/skeleton/SM04/template.go`
- `obx_rest/skeleton/SM04/repository.go`
- `obx_rest/skeleton/SM04/usecase.go`
- `obx_rest/skeleton/SM04/handler.go`

### Frontend
- `obx_site/src/app/board/pages/SM04/page.tsx`

### Documentation
- `obx_docs/blueprint/REST/SM04.md`
- `obx_docs/guide/REST/SM04.md`
- `obx_docs/blueprint/SITE/SM04.md`
- `obx_docs/guide/SITE/SM04.md`

## SM05

### Backend
- `obx_rest/skeleton/SM05/template.go`
- `obx_rest/skeleton/SM05/repository.go`
- `obx_rest/skeleton/SM05/usecase.go`
- `obx_rest/skeleton/SM05/handler.go`

### Frontend
- `obx_site/src/app/board/pages/SM05/page.tsx`

### Documentation
- `obx_docs/blueprint/REST/SM05.md`
- `obx_docs/guide/REST/SM05.md`
- `obx_docs/blueprint/SITE/SM05.md`
- `obx_docs/guide/SITE/SM05.md`

## SP01

### Backend
- `obx_rest/skeleton/SP01/template.go`
- `obx_rest/skeleton/SP01/repository.go`
- `obx_rest/skeleton/SP01/usecase.go`
- `obx_rest/skeleton/SP01/handler.go`

### Frontend
- `obx_site/src/app/board/pages/SP01/page.tsx`

### Documentation
- `obx_docs/blueprint/REST/SP01.md`
- `obx_docs/guide/REST/SP01.md`
- `obx_docs/blueprint/SITE/SP01.md`
- `obx_docs/guide/SITE/SP01.md`

## SP02

### Backend
- `obx_rest/skeleton/SP02/template.go`
- `obx_rest/skeleton/SP02/repository.go`
- `obx_rest/skeleton/SP02/usecase.go`
- `obx_rest/skeleton/SP02/handler.go`

### Frontend
- `obx_site/src/app/board/pages/SP02/page.tsx`

### Documentation
- `obx_docs/blueprint/REST/SP02.md`
- `obx_docs/guide/REST/SP02.md`
- `obx_docs/blueprint/SITE/SP02.md`
- `obx_docs/guide/SITE/SP02.md`

## SP03

### Backend
- `obx_rest/skeleton/SP03/template.go`
- `obx_rest/skeleton/SP03/repository.go`
- `obx_rest/skeleton/SP03/usecase.go`
- `obx_rest/skeleton/SP03/handler.go`

### Frontend
- `obx_site/src/app/board/pages/SP03/page.tsx`

### Documentation
- `obx_docs/blueprint/REST/SP03.md`
- `obx_docs/guide/REST/SP03.md`
- `obx_docs/blueprint/SITE/SP03.md`
- `obx_docs/guide/SITE/SP03.md`

## Urutan Eksekusi yang Disarankan

1. Update `routes.go` dan file backend inti per modul.
2. Sinkronkan `page.tsx` dan `module.ts`.
3. Update blueprint dan guide per modul.
4. Jalankan verifikasi build dan lint.

## Risiko

1. Terlalu banyak file diubah sekaligus.
2. Route dan page tidak sinkron jika modul dikerjakan terpisah.
3. Dokumentasi tertinggal dari implementasi.

## Mitigasi

1. Kerjakan modul satu per satu dalam batch kecil.
2. Update backend, frontend, dan dokumentasi dalam satu paket bila memungkinkan.
3. Verifikasi setelah setiap batch.

## Referensi

- phase1_package1_file_task_plan.md
- phase1_backend_coding_plan.md
- phase1_backend_execution_backlog.md
- phase1_execution_backlog.md
- ai_runbook.md
