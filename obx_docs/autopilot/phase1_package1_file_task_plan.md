# Rencana Task File-by-File - Paket 1 Fase 1

Dokumen ini menurunkan Paket 1 Fase 1 menjadi daftar task file-by-file yang siap dipakai untuk implementasi teknis.

## Ruang Lingkup

- Modul inti: SM01, SM02, SM03, SM04, SM05, SP01, SP02, SP03
- Fokus: master data, session, privilege, dan sinkronisasi dokumentasi teknis

## Urutan Task

### 1. Backend Route and Handler Review

#### Target File
- `obx_rest/backbone/routes.go`
- `obx_rest/skeleton/SM01/handler.go`
- `obx_rest/skeleton/SM02/handler.go`
- `obx_rest/skeleton/SM03/handler.go`
- `obx_rest/skeleton/SM04/handler.go`
- `obx_rest/skeleton/SM05/handler.go`
- `obx_rest/skeleton/SP01/handler.go`
- `obx_rest/skeleton/SP02/handler.go`
- `obx_rest/skeleton/SP03/handler.go`

#### Task
1. Verifikasi route aktif dan urutan static sebelum dynamic.
2. Pastikan handler write route memakai logging yang benar.
3. Validasi privilege dan middleware yang dipakai.

#### Output
- Route inti dan handler inti stabil.

### 2. Backend Usecase and Repository Review

#### Target File
- `obx_rest/skeleton/SM01/repository.go`
- `obx_rest/skeleton/SM01/usecase.go`
- `obx_rest/skeleton/SM02/repository.go`
- `obx_rest/skeleton/SM02/usecase.go`
- `obx_rest/skeleton/SM03/repository.go`
- `obx_rest/skeleton/SM03/usecase.go`
- `obx_rest/skeleton/SM04/repository.go`
- `obx_rest/skeleton/SM04/usecase.go`
- `obx_rest/skeleton/SM05/repository.go`
- `obx_rest/skeleton/SM05/usecase.go`
- `obx_rest/skeleton/SP01/repository.go`
- `obx_rest/skeleton/SP01/usecase.go`
- `obx_rest/skeleton/SP02/repository.go`
- `obx_rest/skeleton/SP02/usecase.go`
- `obx_rest/skeleton/SP03/repository.go`
- `obx_rest/skeleton/SP03/usecase.go`

#### Task
1. Review query path dan error handling.
2. Pastikan usecase memvalidasi input dan wrap error dengan helper standar.
3. Sinkronkan alur list, create, update, delete, dan session/profile flow.

#### Output
- Logic backend inti konsisten dan siap diverifikasi.

### 3. Frontend Page Review

#### Target File
- `obx_site/src/app/board/pages/SM01/page.tsx`
- `obx_site/src/app/board/pages/SM02/page.tsx`
- `obx_site/src/app/board/pages/SM03/page.tsx`
- `obx_site/src/app/board/pages/SM04/page.tsx`
- `obx_site/src/app/board/pages/SM05/page.tsx`
- `obx_site/src/app/board/pages/SP01/page.tsx`
- `obx_site/src/app/board/pages/SP02/page.tsx`
- `obx_site/src/app/board/pages/SP03/page.tsx`
- `obx_site/src/app/board/model/module.ts`

#### Task
1. Pastikan page terhubung ke endpoint backend yang benar.
2. Verifikasi form, loading state, dan mobile behavior.
3. Sinkronkan nama modul dan menu board.

#### Output
- Halaman inti tetap konsisten dengan backend.

### 4. Document Synchronization

#### Target File
- `obx_docs/blueprint/REST/SM01.md`
- `obx_docs/blueprint/REST/SM02.md`
- `obx_docs/blueprint/REST/SM03.md`
- `obx_docs/blueprint/REST/SM04.md`
- `obx_docs/blueprint/REST/SM05.md`
- `obx_docs/blueprint/REST/SP01.md`
- `obx_docs/blueprint/REST/SP02.md`
- `obx_docs/blueprint/REST/SP03.md`
- `obx_docs/guide/REST/SM01.md`
- `obx_docs/guide/REST/SM02.md`
- `obx_docs/guide/REST/SM03.md`
- `obx_docs/guide/REST/SM04.md`
- `obx_docs/guide/REST/SM05.md`
- `obx_docs/guide/REST/SP01.md`
- `obx_docs/guide/REST/SP02.md`
- `obx_docs/guide/REST/SP03.md`
- `obx_docs/blueprint/SITE/SM01.md`
- `obx_docs/blueprint/SITE/SM02.md`
- `obx_docs/blueprint/SITE/SM03.md`
- `obx_docs/blueprint/SITE/SM04.md`
- `obx_docs/blueprint/SITE/SM05.md`
- `obx_docs/blueprint/SITE/SP01.md`
- `obx_docs/blueprint/SITE/SP02.md`
- `obx_docs/blueprint/SITE/SP03.md`
- `obx_docs/guide/SITE/SM01.md`
- `obx_docs/guide/SITE/SM02.md`
- `obx_docs/guide/SITE/SM03.md`
- `obx_docs/guide/SITE/SM04.md`
- `obx_docs/guide/SITE/SM05.md`
- `obx_docs/guide/SITE/SP01.md`
- `obx_docs/guide/SITE/SP02.md`
- `obx_docs/guide/SITE/SP03.md`

#### Task
1. Update blueprint teknis sesuai state source code terbaru.
2. Update user guide sesuai alur pengguna.
3. Pastikan indeks blueprint dan guide tetap konsisten.

#### Output
- Dokumen teknis dan guide inti sinkron dengan implementasi.

### 5. Verification Batch

#### Command
```bash
cd obx_rest
go build ./...
cd ../obx_site
npm run lint
```

#### Smoke Check
- Login
- Profile
- Company
- Privilege
- Session

#### Output
- Bukti verifikasi Paket 1 tersedia.

## Risiko

1. Perubahan route inti berdampak ke banyak page sekaligus.
2. Dokumentasi dan kode tidak di-update dalam batch yang sama.
3. Smoke test gagal karena session atau privilege belum stabil.

## Mitigasi

1. Kerjakan per klaster modul, bukan seluruh workspace sekaligus.
2. Update backend, frontend, dan dokumentasi dalam satu paket bila memungkinkan.
3. Jalankan verifikasi setelah setiap perubahan besar.

## Referensi

- phase1_backend_coding_plan.md
- phase1_backend_execution_backlog.md
- phase1_execution_backlog.md
- ai_runbook.md
