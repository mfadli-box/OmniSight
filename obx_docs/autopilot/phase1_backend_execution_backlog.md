# Backlog Backend Fase 1 - Governance dan Asset Fondasi

Dokumen ini memecah Fase 1 menjadi urutan kerja backend yang paling siap dieksekusi untuk memulai implementasi teknis.

## Tujuan

1. Menstabilkan master data dan session flow modul inti.
2. Menyediakan workflow approval dokumen terkendali.
3. Menyiapkan fondasi data asset untuk fase bastion dan monitoring.
4. Menjaga perubahan tetap kecil, dapat diuji, dan mudah diaudit.

## Urutan Backend

### Paket 1 - Master Data Inti

#### Modul
- SM01
- SM02
- SM03
- SM04
- SM05
- SP01
- SP02
- SP03

#### Aksi
1. Review route aktif dan privilege untuk modul inti.
2. Verifikasi session, profile, company, dan module mapping.
3. Pastikan write route memakai logging yang konsisten.
4. Sinkronkan blueprint REST dan SITE per modul inti.

#### Output
- Master data dan akses stabil.
- Jembatan ke tahap dokumen dan asset siap dipakai.

### Paket 2 - Workflow Dokumen Terkendali

#### Modul
- dat_request
- dat_signature
- dat_document
- dat_document_revision
- dat_document_evidence

#### Aksi
1. Finalisasi schema dan relasi dokumen terkendali.
2. Implement approval workflow berbasis signature type.
3. Simpan histori revisi dan bukti audit.
4. Pastikan status request dan actor audit konsisten.

#### Output
- Dokumen SOP/ISO dapat dibuat, direview, disetujui, dan ditelusuri.

### Paket 3 - Asset Inventory Dasar

#### Modul
- asset generik untuk host, service, app, dan endpoint internal
- seed data untuk pilot

#### Aksi
1. Tetapkan field minimum untuk inventory dasar.
2. Tentukan company scope dan tipe asset.
3. Tambahkan seed data minimum untuk smoke test.
4. Siapkan relasi untuk reuse di bastion dan monitoring.

#### Output
- Fondasi asset tersedia untuk fase berikutnya.

### Paket 4 - Dokumentasi Operasional Backend

#### Fokus
- blueprint teknis
- user guide
- indeks dokumen

#### Aksi
1. Perbarui blueprint REST/SITE untuk modul inti.
2. Perbarui guide user untuk alur yang berubah.
3. Tautkan sample SOP/ISO dan checklist rollout.
4. Pastikan catatan update terakhir jelas.

#### Output
- Fase 1 siap untuk audit dan handoff ke Fase 2.

## Dependensi

| Area | Dependensi Minimum |
|---|---|
| Master data inti | route aktif, privilege konsisten, session stabil |
| Workflow dokumen | dat_request, dat_signature, dat_document |
| Asset inventory | field minimum, seed data, company scope |
| Dokumentasi backend | blueprint dan guide sinkron |

## DoD Backend Fase 1

- Route inti stabil dan lulus build.
- Approval workflow dokumen aktif.
- Revision history dan evidence trail berfungsi.
- Asset inventory dasar tersedia.
- Dokumentasi teknis dan guide sinkron.

## Risiko

1. Master data belum stabil dan mengganggu workflow dokumen.
2. Asset inventory dibangun sebelum privilege dan session siap.
3. Dokumentasi tertinggal dari implementasi.

## Mitigasi

1. Gunakan batch kecil dan verifikasi per langkah.
2. Jangan lanjut ke Fase 2 sebelum DoD Fase 1 terpenuhi.
3. Update blueprint dan guide dalam batch yang sama dengan implementasi.

## Referensi

- roadmap_platform_3_phase.md
- phase1_execution_backlog.md
- backend_module_matrix_platform_replacement.md
- ai_runbook.md
