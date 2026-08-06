# Definition of Done - AI Module

Checklist ini adalah gerbang akhir sebelum modul dianggap selesai.

## A. Spec dan Scope

- [ ] File spec modul tersedia dan valid sesuai module_spec_schema.md.
- [ ] module_code, route, table, dan privilege konsisten.
- [ ] Tidak ada requirement kritis yang masih placeholder.

## B. Database

- [ ] Schema Prisma diperbarui pada file yang tepat.
- [ ] Migrasi dibuat dengan nama yang jelas.
- [ ] Migrasi sukses dijalankan di local dev DB.
- [ ] Seed minimal tersedia jika dibutuhkan alur UI/testing.

## C. Backend

- [ ] Route terdaftar di group yang benar.
- [ ] Handler, usecase, dan repo mengikuti pattern yang ada.
- [ ] Error handling menggunakan helper standar.
- [ ] Logging write route menggunakan kode modul yang benar.

## D. Frontend

- [ ] Halaman terhubung ke route backend yang tepat.
- [ ] Field form sesuai spec dan validasi berjalan.
- [ ] Perilaku mobile sesuai pattern proyek.
- [ ] State loading/error/empty ditangani.

## E. Testing

- [ ] Test minimal backend terpenuhi.
- [ ] Test minimal frontend terpenuhi.
- [ ] Tidak ada regresi pada route/page terkait.
- [ ] Build dan typecheck lolos.

## F. Dokumentasi

- [ ] Blueprint teknis modul dibuat atau diperbarui.
- [ ] User guide modul dibuat atau diperbarui.
- [ ] Referensi di indeks dokumen sudah terdaftar.

## G. Operasional

- [ ] Catatan konfigurasi runtime ditulis jika ada env baru.
- [ ] Catatan rollback disiapkan untuk migration/route kritis.
- [ ] Risiko utama dan asumsi dicatat singkat.

## Status Akhir

Modul hanya boleh diberi status done jika seluruh poin wajib sudah centang.