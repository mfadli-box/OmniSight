# dat_document - Panduan Document Domain

## Ringkasan Fitur

Schema ini dipakai untuk dokumen, versioning, approval, dan request dokumen (`dat_request`).

## Struktur Utama

- `dat_document_category` untuk pengelompokan dokumen
- `dat_document` untuk entitas dokumen utama
- `dat_document_version` untuk riwayat versi
- `dat_document_approval` untuk jejak approval dokumen
- `dat_request` untuk alur pengajuan/request dokumen

## Langkah Penggunaan

1. Perbarui schema saat format document berubah.
2. Pastikan approval dan versioning tetap konsisten.
3. Pastikan alur request (`dat_request`) sinkron dengan approval type dan user relation.

## Validasi Hasil

- Kode dokumen dan kode request tetap unik per company.
- Relasi requester/completer dan status approval terbaca dengan benar.
