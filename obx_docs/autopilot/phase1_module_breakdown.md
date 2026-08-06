# Breakdown Fase 1 - Governance, SOP/ISO, dan Fondasi Asset

Dokumen ini memecah Fase 1 menjadi paket implementasi per modul agar dapat dieksekusi bertahap dengan risiko rendah.

## Tujuan Fase 1

- Menyediakan controlled documents untuk SOP/ISO
- Menstabilkan company, user, privilege, dan request workflow
- Menyediakan asset inventory dasar sebagai fondasi fase bastion dan monitoring

## Paket 1 - Core Master Data Stabilization

### Modul
- SM01
- SM02
- SM03
- SM04
- SM05
- SP01
- SP02
- SP03

### Fokus
- User
- Module
- Company
- Signature type
- Session management
- Profile, password, dan privilege current user

### Task
- Review route dan privilege eksisting
- Verifikasi company/module/area mapping
- Verifikasi audit write action
- Rapikan dokumentasi REST dan SITE per modul inti

### Output
- Fondasi identitas dan akses stabil
- Audit dan session management siap dipakai fase berikutnya

## Paket 2 - Controlled Document dan Approval Foundation

### Domain
- dat_document
- dat_signature
- dat_request

### Fokus
- Draft SOP/ISO
- Versioning
- Approval flow
- Request workflow formal

### Task
- Finalisasi struktur kategori dokumen
- Finalisasi metadata dokumen terkendali
- Mapping signature type untuk approval SOP/ISO
- Definisikan template dokumen: policy, SOP, work instruction, form, evidence

### Output
- Alur dokumen terkendali siap dipakai internal
- Approval dan histori perubahan terdokumentasi

## Paket 3 - Asset Inventory Dasar

### Domain target
- Asset generik untuk host, service, app, dan endpoint internal
- Dapat memakai XX99 sebagai referensi template implementasi awal sebelum modul final dibentuk

### Fokus
- Kode asset
- Company scope
- Tipe asset
- Status aktif
- Metadata alamat / host / port minimum

### Task
- Finalisasi struktur asset minimal
- Tentukan apakah implementasi awal memakai modul template atau langsung modul final
- Siapkan daftar field minimum lintas domain
- Siapkan seed data minimum untuk smoke test

### Output
- Inventory dasar siap dipakai untuk bastion dan monitoring tahap berikutnya

## Paket 4 - Dokumentasi Operasional dan Audit Readiness

### Fokus
- Template SOP/ISO
- Checklist bukti audit
- Link blueprint, guide, dan autopilot

### Task
- Buat template controlled document
- Buat template evidence checklist
- Buat daftar owner, reviewer, approver
- Pastikan semua dokumen punya update terakhir yang jelas

### Output
- Fase 1 siap dipakai untuk audit internal dasar

## Urutan Eksekusi Rekomendasi

1. Paket 1 - Core Master Data Stabilization
2. Paket 2 - Controlled Document dan Approval Foundation
3. Paket 3 - Asset Inventory Dasar
4. Paket 4 - Dokumentasi Operasional dan Audit Readiness

## Validasi Minimum Fase 1

- `cd obx_base; npx prisma generate`
- `cd obx_rest; go build ./...`
- `cd obx_site; npm run lint`
- Smoke test login, profile, company, privilege, session
- Review dokumen SOP/ISO sample dengan alur approval

## Exit Criteria Fase 1

- Master data stabil dan terdokumentasi.
- Dokumen SOP/ISO punya approval dan histori.
- Request workflow memakai `dat_request` jelas.
- Asset inventory dasar tersedia untuk dipakai fase 2.
