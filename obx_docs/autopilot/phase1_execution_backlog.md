# Backlog Eksekusi Fase 1 - Governance, SOP/ISO, dan Fondasi Asset

Dokumen ini memecah Fase 1 menjadi urutan kerja yang siap dieksekusi, dengan fokus pada outcome yang bisa diverifikasi sebelum masuk Fase 2.

## Tujuan Eksekusi Fase 1

1. Menstabilkan master data identitas dan akses.
2. Mengaktifkan workflow approval dokumen terkendali.
3. Menyediakan asset inventory dasar untuk fondasi bastion dan monitoring.
4. Menyelaraskan dokumentasi teknis, user guide, dan audit readiness.

## Urutan Eksekusi

### Langkah 1 - Stabilkan Master Data Inti

#### Scope
- SM01
- SM02
- SM03
- SM04
- SM05
- SP01
- SP02
- SP03

#### Aksi
1. Verifikasi route aktif dan privilege setiap modul inti.
2. Pastikan company, module, user, dan session mapping konsisten.
3. Pastikan write route memakai logging yang benar.
4. Sinkronkan dokumentasi REST dan SITE per modul inti.

#### Output
- Master data stabil.
- Identity and access baseline siap dipakai.

#### Verifikasi
- Build backend lulus.
- Smoke test login, profile, company, privilege, session lulus.

### Langkah 2 - Aktifkan Workflow SOP/ISO

#### Scope
- dat_request
- dat_signature
- dat_document
- dat_document_revision
- dat_document_evidence

#### Aksi
1. Finalisasi struktur dokumen terkendali.
2. Tetapkan signature type untuk policy, SOP, work instruction, dan form/record.
3. Pastikan request workflow memakai status dan actor audit.
4. Tambahkan histori revisi dan evidence trail.

#### Output
- Dokumen terkendali bisa dibuat, direview, disetujui, dan direkam historinya.

#### Verifikasi
- Sample SOP/ISO berjalan melalui alur approval.
- Evidence checklist dapat diisi dan ditelusuri.

### Langkah 3 - Siapkan Asset Inventory Dasar

#### Scope
- asset generik untuk host, service, app, dan endpoint internal
- seed data minimum untuk pilot

#### Aksi
1. Finalisasi field minimum untuk inventory dasar.
2. Tentukan scope company dan tipe asset.
3. Siapkan seed data untuk smoke test dan pilot.
4. Pastikan asset dapat dipakai ulang oleh fase bastion dan monitoring.

#### Output
- Fondasi asset tersedia untuk tahap berikutnya.

#### Verifikasi
- Data asset dasar muncul di backend dan frontend bila sudah dipasang.

### Langkah 4 - Audit Readiness dan Dokumentasi Operasional

#### Scope
- template SOP/ISO
- template evidence checklist
- blueprint dan guide per domain inti

#### Aksi
1. Perbarui template dokumen terkendali.
2. Perbarui template bukti audit.
3. Daftarkan file sample dan checklist ke indeks dokumentasi.
4. Pastikan update terakhir terdokumentasi jelas.

#### Output
- Fase 1 siap untuk audit internal dasar.

#### Verifikasi
- Indeks dokumen memuat template, sample, dan checklist yang relevan.

## Dependensi Fase 1

| Area | Dependensi Minimum |
|---|---|
| Master data inti | route aktif, privilege konsisten, session stabil |
| Approval dokumen | dat_request, dat_signature, dat_document |
| Asset inventory | field minimum, seed data, company scope |
| Audit readiness | template, sample, dan indeks dokumen |

## DoD Fase 1

- Approval workflow dokumen terkendali aktif.
- Signature type untuk dokumen utama tersedia.
- Revision history dan evidence trail berfungsi.
- Asset inventory dasar tersedia.
- Dokumentasi teknis dan user guide sinkron.

## Risiko Utama

1. Scope dokumen terlalu longgar.
2. Master data belum stabil sehingga approval workflow terganggu.
3. Asset inventory dibangun terlalu dini sebelum privilege dan session siap.

## Mitigasi

1. Eksekusi dengan batch kecil.
2. Jangan lanjut ke Fase 2 sebelum DoD Fase 1 tercapai.
3. Gunakan roadmap dan matriks backend sebagai urutan modul.

## Referensi

- roadmap_platform_3_phase.md
- phase1_module_breakdown.md
- ai_runbook.md
- checklist_platform_replacement_execution.md
