# User Guide BASE - Operasional obx_base Prisma

## Ringkasan

Panduan ini dipakai untuk workflow harian perubahan database obx_base berbasis Prisma.

## Struktur Dokumen

- Summary guide: file ini
- Detail page guide: `obx_docs/guide/BASE/*.md`
- Summary teknis: `obx_docs/blueprint/BASE/README.md`
- Detail teknis: `obx_docs/blueprint/BASE/*.md`

## Domain Ringkas

- dat: company, module, user
- jsm_stack: bastion, token koneksi, web proxy, file transfer
- dat_signature: approval flow
- doc: document & request (`dat_request`)
- ict_machine / ict_mikrotik / ict_security / ict_website: domain operasional infra

## Alur Kerja Inti

1. Pilih domain schema yang akan diubah.
2. Lakukan perubahan kecil dan spesifik.
3. Jalankan `npx prisma generate`.
4. Buat migrasi development dengan nama deskriptif.
5. Verifikasi status migrasi.
6. Perbarui detail page guide dan blueprint.

## Checklist

- [ ] Schema terkait sudah jelas
- [ ] Generate client lulus
- [ ] Migration status aman
- [ ] Detail page guide diperbarui

## Referensi

- `obx_base/prisma/schema/*.prisma`
- `obx_docs/blueprint/BASE/README.md`
