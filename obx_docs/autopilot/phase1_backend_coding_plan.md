# Rencana Coding Backend Fase 1 - Governance dan Asset Fondasi

Dokumen ini menurunkan backlog backend Fase 1 menjadi paket coding yang lebih konkret untuk dieksekusi di repo.

## Tujuan

1. Mengaktifkan master data dan session flow modul inti.
2. Membangun workflow approval dokumen terkendali.
3. Menyiapkan asset inventory dasar yang bisa dipakai ulang oleh Bastion dan monitoring.
4. Menjaga setiap paket kecil, terverifikasi, dan mudah di-rollback.

## Prinsip Coding

- Satu paket fokus pada satu klaster domain.
- Schema, backend, dan dokumentasi yang terkait harus selesai dalam batch yang sama bila memungkinkan.
- Jangan lanjut ke paket berikutnya sebelum verifikasi paket berjalan.

## Paket 1 - Master Data Inti

### Target Modul
- SM01
- SM02
- SM03
- SM04
- SM05
- SP01
- SP02
- SP03

### File yang Ditargetkan
- `obx_rest/skeleton/{KODE}/template.go`
- `obx_rest/skeleton/{KODE}/repository.go`
- `obx_rest/skeleton/{KODE}/usecase.go`
- `obx_rest/skeleton/{KODE}/handler.go`
- `obx_rest/backbone/routes.go`
- `obx_site/src/app/board/pages/{KODE}/page.tsx`
- `obx_docs/blueprint/REST/{KODE}.md`
- `obx_docs/guide/REST/{KODE}.md`
- `obx_docs/blueprint/SITE/{KODE}.md`
- `obx_docs/guide/SITE/{KODE}.md`

### Aksi Coding
1. Review route aktif dan privilege setiap modul inti.
2. Pastikan write route memakai logging yang konsisten.
3. Rapikan dokumentasi REST dan SITE untuk modul inti.
4. Verifikasi session, company, module, dan profile flow.

### Verifikasi
- `go build ./...` di `obx_rest` lulus.
- `npm run lint` di `obx_site` lulus.
- Smoke test login, profile, company, privilege, dan session lulus.

## Paket 2 - Workflow Dokumen Terkendali

### Target Modul
- dat_request
- dat_signature
- dat_document
- dat_document_revision
- dat_document_evidence

### File yang Ditargetkan
- `obx_base/prisma/schema/*.prisma`
- `obx_rest/skeleton/{KODE}/template.go`
- `obx_rest/skeleton/{KODE}/repository.go`
- `obx_rest/skeleton/{KODE}/usecase.go`
- `obx_rest/skeleton/{KODE}/handler.go`
- `obx_rest/backbone/routes.go`
- `obx_site/src/app/board/pages/{KODE}/page.tsx`
- `obx_docs/blueprint/BASE/{file}.md`
- `obx_docs/guide/BASE/{file}.md`
- `obx_docs/blueprint/REST/{KODE}.md`
- `obx_docs/guide/REST/{KODE}.md`
- `obx_docs/blueprint/SITE/{KODE}.md`
- `obx_docs/guide/SITE/{KODE}.md`

### Aksi Coding
1. Finalisasi schema dan relasi dokumen terkendali.
2. Implement approval workflow berbasis signature type.
3. Simpan histori revisi dan bukti audit.
4. Pastikan status request dan actor audit konsisten.
5. Tambahkan sample dokumen dan evidence checklist jika dibutuhkan untuk uji.

### Verifikasi
- Approval workflow dokumen aktif.
- Revision history dan evidence trail berfungsi.
- Sample SOP/ISO dapat melewati alur approval.

## Paket 3 - Asset Inventory Dasar

### Target Modul
- asset generik untuk host, service, app, dan endpoint internal
- seed data minimum untuk pilot

### File yang Ditargetkan
- `obx_base/prisma/schema/*.prisma`
- `obx_rest/skeleton/{KODE}/template.go`
- `obx_rest/skeleton/{KODE}/repository.go`
- `obx_rest/skeleton/{KODE}/usecase.go`
- `obx_rest/skeleton/{KODE}/handler.go`
- `obx_rest/backbone/routes.go`
- `obx_site/src/app/board/pages/{KODE}/page.tsx`

### Aksi Coding
1. Tetapkan field minimum untuk inventory dasar.
2. Tentukan company scope dan tipe asset.
3. Tambahkan seed data minimum untuk smoke test.
4. Siapkan relasi yang bisa dipakai ulang di Bastion dan monitoring.

### Verifikasi
- Data asset dasar dapat dibuat dan dibaca.
- Seed data pilot tersedia untuk smoke test.

## Paket 4 - Dokumentasi Operasional Backend

### Fokus
- blueprint teknis
- user guide
- indeks dokumen

### Aksi Coding
1. Perbarui blueprint REST/SITE untuk modul inti.
2. Perbarui guide user untuk alur yang berubah.
3. Tautkan sample SOP/ISO dan checklist rollout.
4. Pastikan catatan update terakhir jelas.

### Verifikasi
- Indeks dokumen memuat template, sample, dan checklist yang relevan.

## Urutan Coding yang Disarankan

1. Paket 1 - Master Data Inti
2. Paket 2 - Workflow Dokumen Terkendali
3. Paket 3 - Asset Inventory Dasar
4. Paket 4 - Dokumentasi Operasional Backend

## DoD Coding Fase 1

- Master data inti stabil.
- Dokumen SOP/ISO punya approval dan histori.
- Request workflow memakai `dat_request` jelas.
- Asset inventory dasar tersedia.
- Dokumentasi teknis dan user guide sinkron.

## Risiko

1. Paket dokumen dikerjakan sebelum master data stabil.
2. Seed asset dibuat tanpa skema minimum yang konsisten.
3. Dokumentasi tertinggal dari implementasi.

## Mitigasi

1. Eksekusi batch kecil dan lakukan verifikasi di setiap paket.
2. Jangan lanjut ke Paket 2 sebelum Paket 1 lulus build dan smoke test.
3. Update schema, backend, dan dokumentasi dalam batch yang sama bila memungkinkan.

## Referensi

- phase1_backend_execution_backlog.md
- phase1_execution_backlog.md
- roadmap_platform_3_phase.md
- backend_module_matrix_platform_replacement.md
- ai_runbook.md
