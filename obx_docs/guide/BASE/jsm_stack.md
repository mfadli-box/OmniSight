# jsm_stack - Panduan Bastion dan Web Access Domain

## Ringkasan Fitur

Schema ini dipakai untuk fondasi bastion internal berbasis web, termasuk asset, account, session audit, web app proxy, dan file transfer.

## Struktur Utama

- `jms_asset` untuk inventori target akses
- `jms_account` dan `jms_account_secret` untuk akun akses
- `jms_session` untuk sesi koneksi
- `jms_connect_token` untuk token akses sementara
- `jms_web_app` untuk registry aplikasi internal yang diproxy
- `jms_file_transfer` untuk audit upload/download file

## Aturan Penting

- Token koneksi wajib short-lived dan one-time.
- `token_hash` menyimpan hash token, bukan token mentah.
- Semua aktivitas koneksi dan transfer file harus bisa ditelusuri.
- Scope data mengikuti `company_id`.

## Langkah Penggunaan

1. Tambahkan asset target terlebih dahulu.
2. Tambahkan account akses jika diperlukan.
3. Buat token koneksi sementara untuk SSH, RDP, FTP, atau WEBAPP.
4. Simpan event sesi dan transfer file untuk audit.

## Validasi Hasil

- Asset dan account dapat di-query berdasarkan company.
- Token koneksi dapat dibedakan status aktif, terpakai, kedaluwarsa, atau dicabut.
- Upload/download file tercatat pada histori transfer.
