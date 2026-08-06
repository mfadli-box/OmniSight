# Rencana Task File-by-File - Paket 3 Fase 1

Dokumen ini menurunkan Paket 3 Fase 1 menjadi daftar task file-by-file yang siap dipakai untuk implementasi asset inventory dasar.

## Ruang Lingkup

- Domain: asset generik untuk host, service, app, dan endpoint internal
- Fokus: field minimum, company scope, seed data, dan reuse untuk Bastion serta monitoring

## Urutan Task

### 1. Schema Review for Asset Inventory

#### Target File
- `obx_base/prisma/schema/*.prisma`
- `obx_base/prisma/schema/dat_company.prisma`
- `obx_base/prisma/schema/dat_module.prisma`
- `obx_base/prisma/schema/dat_user.prisma`

#### Task
1. Tetapkan struktur asset minimum untuk host, service, app, dan endpoint internal.
2. Pastikan company scope dan status aktif konsisten.
3. Validasi field minimum untuk host, address, port, type, dan metadata.

#### Output
- Schema asset dasar siap dipakai backend.

### 2. Backend Inventory Review

#### Target File
- `obx_rest/skeleton/{KODE}/template.go`
- `obx_rest/skeleton/{KODE}/repository.go`
- `obx_rest/skeleton/{KODE}/usecase.go`
- `obx_rest/skeleton/{KODE}/handler.go`
- `obx_rest/backbone/routes.go`

#### Task
1. Review query path dan error handling.
2. Implement create, list, update, delete untuk asset inventory dasar.
3. Pastikan company scope dan session audit konsisten.
4. Siapkan reuse di Bastion dan monitoring.

#### Output
- Backend inventory dasar aktif.

### 3. Frontend Page Review

#### Target File
- `obx_site/src/app/board/pages/{KODE}/page.tsx`
- `obx_site/src/app/board/model/module.ts`

#### Task
1. Pastikan page terhubung ke endpoint backend yang benar.
2. Verifikasi form, loading state, dan mobile behavior.
3. Sinkronkan nama modul dan menu board.

#### Output
- Halaman asset inventory siap dipakai user.

### 4. Documentation Synchronization

#### Target File
- `obx_docs/blueprint/BASE/dat_company.md`
- `obx_docs/blueprint/BASE/dat_module.md`
- `obx_docs/blueprint/BASE/dat_user.md`
- `obx_docs/blueprint/BASE/dat_document.md`
- `obx_docs/guide/BASE/dat_company.md`
- `obx_docs/guide/BASE/dat_module.md`
- `obx_docs/guide/BASE/dat_user.md`
- `obx_docs/guide/BASE/dat_document.md`
- `obx_docs/autopilot/template_iso_evidence_checklist.md`
- `obx_docs/autopilot/template_approval_matrix.md`

#### Task
1. Update blueprint teknis sesuai state source code terbaru.
2. Update user guide sesuai alur asset inventory.
3. Pastikan template, sample, dan checklist tetap konsisten bila asset dipakai di workflow dokumen.

#### Output
- Dokumentasi asset inventory sinkron dengan implementasi.

### 5. Verification Batch

#### Command
```bash
cd obx_base
npx prisma generate
cd ../obx_rest
go build ./...
cd ../obx_site
npm run lint
```

#### Smoke Check
- Buat asset host
- Buat asset service
- Buat asset app
- Buat asset endpoint internal

#### Output
- Bukti verifikasi Paket 3 tersedia.

## Risiko

1. Struktur asset terlalu generik sehingga tidak cocok untuk Bastion atau monitoring.
2. Dokumentasi tidak mengikuti field minimum yang dipakai backend.
3. Smoke test gagal karena company scope belum stabil.

## Mitigasi

1. Tetapkan field minimum lebih awal dan gunakan company scope secara konsisten.
2. Update backend dan dokumentasi dalam batch yang sama.
3. Jalankan verifikasi setelah setiap perubahan besar.

## Referensi

- phase1_backend_coding_plan.md
- phase1_backend_execution_backlog.md
- phase1_execution_backlog.md
- ai_runbook.md
