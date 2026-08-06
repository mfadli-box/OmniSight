# Rencana Task File-by-File - Paket 2 Fase 1

Dokumen ini menurunkan Paket 2 Fase 1 menjadi daftar task file-by-file yang siap dipakai untuk implementasi workflow dokumen terkendali.

## Ruang Lingkup

- Modul domain: dat_request, dat_signature, dat_document, dat_document_revision, dat_document_evidence
- Fokus: approval workflow, versioning, evidence trail, dan sinkronisasi dokumentasi

## Urutan Task

### 1. Schema and Relation Review

#### Target File
- `obx_base/prisma/schema/dat_request.prisma`
- `obx_base/prisma/schema/dat_signature.prisma`
- `obx_base/prisma/schema/dat_document.prisma`
- `obx_base/prisma/schema/dat_document_revision.prisma`
- `obx_base/prisma/schema/dat_document_evidence.prisma`

#### Task
1. Finalisasi struktur schema dan relasi antar tabel.
2. Pastikan company scope, status, dan actor audit konsisten.
3. Validasi field minimum untuk approval workflow dan evidence trail.

#### Output
- Schema dokumen terkendali siap dipakai backend.

### 2. Backend Repository and Usecase Review

#### Target File
- `obx_rest/skeleton/{KODE}/template.go`
- `obx_rest/skeleton/{KODE}/repository.go`
- `obx_rest/skeleton/{KODE}/usecase.go`
- `obx_rest/skeleton/{KODE}/handler.go`
- `obx_rest/backbone/routes.go`

#### Task
1. Review query path dan error handling.
2. Implement approval workflow berbasis signature type.
3. Simpan histori revisi dan bukti audit.
4. Pastikan status request dan actor audit konsisten.

#### Output
- Workflow backend dokumen terkendali aktif.

### 3. Frontend Page Review

#### Target File
- `obx_site/src/app/board/pages/{KODE}/page.tsx`
- `obx_site/src/app/board/model/module.ts`

#### Task
1. Pastikan page terhubung ke endpoint backend yang benar.
2. Verifikasi form, loading state, dan mobile behavior.
3. Sinkronkan menu board dan nama modul.

#### Output
- Halaman dokumen terkendali siap dipakai user.

### 4. Documentation Synchronization

#### Target File
- `obx_docs/blueprint/BASE/dat_request.md`
- `obx_docs/blueprint/BASE/dat_signature.md`
- `obx_docs/blueprint/BASE/dat_document.md`
- `obx_docs/blueprint/BASE/dat_document.md`
- `obx_docs/guide/BASE/dat_request.md`
- `obx_docs/guide/BASE/dat_signature.md`
- `obx_docs/guide/BASE/dat_document.md`
- `obx_docs/guide/BASE/dat_document.md`
- `obx_docs/autopilot/template_iso_policy.md`
- `obx_docs/autopilot/template_iso_sop.md`
- `obx_docs/autopilot/template_iso_work_instruction.md`
- `obx_docs/autopilot/template_iso_form_record.md`
- `obx_docs/autopilot/template_iso_evidence_checklist.md`
- `obx_docs/autopilot/template_approval_matrix.md`

#### Task
1. Update blueprint teknis sesuai state source code terbaru.
2. Update user guide sesuai alur approval dokumen.
3. Pastikan template, sample, dan checklist tetap konsisten.

#### Output
- Dokumentasi workflow dokumen terkendali sinkron dengan implementasi.

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
- Buat request dokumen
- Approve / reject request
- Tambah revision history
- Simpan evidence trail

#### Output
- Bukti verifikasi Paket 2 tersedia.

## Risiko

1. Schema relasi dokumen tidak konsisten dengan workflow backend.
2. Dokumentasi tidak mengikuti perubahan approval flow.
3. Smoke test gagal karena role atau status belum stabil.

## Mitigasi

1. Kerjakan schema dan backend dalam batch kecil.
2. Update dokumentasi dalam batch yang sama dengan perubahan utama.
3. Jalankan verifikasi setelah setiap perubahan besar.

## Referensi

- phase1_backend_coding_plan.md
- phase1_backend_execution_backlog.md
- phase1_execution_backlog.md
- ai_runbook.md
