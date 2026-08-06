# dat_company - Panduan Company Domain

## Ringkasan Fitur

Schema ini mengatur company, module, area, dan relasi request yang dipakai banyak modul lain.

## Struktur Utama

- `dat_company` sebagai master company
- `dat_company_module` untuk aktivasi modul per company
- `dat_company_area` untuk area/divisi per company
- relasi `doc_requests` ke `dat_request`

## Aturan Penting

- `dat_company.code` bersifat **unik global** (bukan per-tenant).
- `dat_company_area` tetap unik per company melalui kombinasi `[company_id, code]`.
- Mapping modul tetap unik per company melalui kombinasi `[company_id, module_id]`.

## Langkah Penggunaan

1. Gunakan schema ini saat menambah company scope.
2. Pastikan mapping module dan area ikut diperbarui.
3. Jika alur request berubah, validasi relasi `dat_company` ke `dat_request`.

## Validasi Hasil

- Company yang baru bisa dipakai oleh user dan page terkait.
- Tidak ada duplikasi `code` pada `dat_company`.
- Relasi request per company tetap konsisten.
