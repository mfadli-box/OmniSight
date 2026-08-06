# Prompt Playbook - AI Module Autopilot

Dokumen ini berisi prompt siap pakai agar proses pembuatan modul lebih konsisten.

## Cara Pakai

1. Ganti placeholder yang diawali tanda kurung siku.
2. Tempel prompt per tahap ke agent.
3. Simpan hasil tiap tahap ke catatan audit runbook.

## Prompt 1 - Validate Spec

Tujuan: memastikan spec valid sebelum coding.

```text
Validasi file spec modul ini: [PATH_SPEC].
Periksa konsistensi module_code, route_base, page_code, primary_table, privilege_code, write_log_code, dan path dokumentasi.
Keluaran wajib:
1) daftar valid,
2) daftar mismatch,
3) rekomendasi perbaikan dengan patch minimal.
```

## Prompt 2 - Scaffold Backend

Tujuan: membuat skeleton backend dari spec.

```text
Gunakan spec di [PATH_SPEC] untuk membuat scaffold backend modul [MODULE_CODE] di obx_rest.
Ikuti pola Repo -> UseCase -> Handler yang sudah ada.
Daftarkan route static sebelum dynamic route.
Pastikan error handling memakai mechanic.Error dan internal error wrapping konsisten.
Keluaran wajib:
1) daftar file yang dibuat/diubah,
2) ringkasan endpoint,
3) command verifikasi backend.
```

## Prompt 3 - Scaffold Frontend

Tujuan: membuat struktur halaman frontend dari spec.

```text
Gunakan spec di [PATH_SPEC] untuk membuat draft halaman [PAGE_CODE] di obx_site.
Ikuti pola UI [UI_PATTERN] dan mobile behavior [MOBILE_BEHAVIOR].
Field form wajib mengikuti spec.
Keluaran wajib:
1) daftar file yang dibuat/diubah,
2) ringkasan alur user,
3) command verifikasi frontend.
```

## Prompt 4 - Implement Logic CRUD

Tujuan: melengkapi logika create, list, detail, update, delete.

```text
Lanjutkan implementasi modul [MODULE_CODE] berdasarkan scaffold yang sudah ada.
Lengkapi query parameterized, whitelist sort, validasi field, dan logging write action dengan USLogs("[MODULE_CODE]").
Keluaran wajib:
1) ringkasan logic selesai,
2) daftar risiko sisa,
3) command smoke test endpoint.
```

## Prompt 5 - Generate Docs

Tujuan: membuat dokumen blueprint dan guide dari spec final.

```text
Generate atau update dokumentasi modul dari [PATH_SPEC].
Target teknis: [BLUEPRINT_TARGET].
Target guide: [GUIDE_TARGET].
Pastikan indeks summary terkait juga diperbarui.
Keluaran wajib:
1) file docs yang diubah,
2) ringkasan perubahan,
3) gap informasi yang masih pending.
```

## Prompt 6 - Final DoD Gate

Tujuan: menutup pekerjaan hanya jika seluruh gate lolos.

```text
Lakukan final review modul [MODULE_CODE] menggunakan checklist di definition_of_done.md.
Tampilkan hasil per bagian: Spec, Database, Backend, Frontend, Testing, Dokumentasi, Operasional.
Jika ada poin gagal, berikan patch dan command verifikasi untuk menutup gap.
Jangan nyatakan selesai jika masih ada item wajib yang belum lolos.
```