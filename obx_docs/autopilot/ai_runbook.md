# AI Runbook - Pembuatan Modul

Runbook ini mendefinisikan urutan kerja autopilot agar output AI konsisten.

## Mode Eksekusi

- Default mode: small-batch.
- Satu batch fokus pada satu module_code.
- Setiap batch harus lolos verifikasi sebelum lanjut.

## Alur Utama

1. Prepare Spec
2. Scaffold Backend
3. Scaffold Frontend
4. Implement Logic
5. Integrate Database
6. Test and Verify
7. Update Documentation
8. Final Review Gate

## Jalur Eksekusi untuk Platform Replacement

Jika tugas yang dikerjakan adalah bagian dari roadmap OmniSight sebagai platform pengganti, ikuti urutan fase berikut:

1. Fase 1: SOP/ISO governance dan fondasi asset.
2. Fase 2: Bastion access untuk SSH, RDP, FTP, dan WebAppProxy.
3. Fase 3: Inventory dan monitoring infra, lalu SIEM terakhir.

Gunakan dokumen pendukung berikut sebagai sumber urutan modul:

- [Roadmap prioritas 3 fase](roadmap_platform_3_phase.md)
- [Matriks aksi implementasi per fitur](../blueprint/PLAN/feature_action_matrix_platform_replacement.md)
- [Matriks modul backend per domain](../blueprint/PLAN/backend_module_matrix_platform_replacement.md)

### Batch Prioritas Backend

- Batch 1: `dat_request`, `dat_signature`, `dat_document`, `dat_document_revision`, `dat_document_evidence`
- Batch 2: `jms_asset`, `jms_account`, `jms_account_secret`, `jms_session`, `jms_connect_token`, `jms_file_transfer`, `jms_web_app`
- Batch 3: `infra_host`, `infra_stack`, `web_site`, `web_upstream`, `vm_host`, `net_device`
- Batch 4: `sec_event`, `sec_event_source`, `sec_event_parser`, `sec_rule`, `sec_alert`, `sec_incident`

### Gate per Fase

- Fase 1 tidak boleh lanjut ke Fase 2 sebelum approval workflow, signature type, dan evidence checklist stabil.
- Fase 2 tidak boleh lanjut ke Fase 3 sebelum session audit, revoke, dan file transfer stabil.
- Fase 3 tidak boleh dimulai jika inventory infra dan bastion core belum mencapai pilot.

## Detail Tahap

### 1) Prepare Spec

- Buat atau update file spec modul menggunakan schema standar.
- Validasi kelengkapan field wajib.
- Tetapkan scope: list/link/form/workflow.
- Jika scope berkaitan dengan platform replacement, pastikan fase dan batch backend sudah ditentukan terlebih dahulu.

Output:
- Spec valid dan siap eksekusi.

### 2) Scaffold Backend

- Tambah struktur route, handler, usecase, repo.
- Ikuti pattern yang sudah digunakan di project.
- Daftarkan route static sebelum dynamic route.

Output:
- Endpoint compile-ready.

### 3) Scaffold Frontend

- Tambah page berdasarkan page_code.
- Terapkan pola komponen yang sudah baku.
- Hubungkan request client ke endpoint spec.

Output:
- Halaman render dan bisa memanggil API.

### 4) Implement Logic

- Lengkapi query, validasi, mapping data, dan error path.
- Pastikan privilege dan logging diterapkan.

Output:
- CRUD/alur utama berjalan sesuai spec.

### 5) Integrate Database

- Update schema Prisma jika ada perubahan model.
- Buat migration dan jalankan ke dev DB.
- Tambahkan seed jika dibutuhkan oleh flow.

Output:
- Schema dan data dasar sinkron.

### 6) Test and Verify

- Jalankan build, typecheck, dan test minimal.
- Uji path sukses dan path error penting.

Output:
- Bukti verifikasi tersedia.

### 7) Update Documentation

- Update blueprint teknis modul.
- Update user guide modul.
- Daftarkan di indeks dokumen terkait.

Output:
- Dokumentasi sinkron dengan implementasi.

### 8) Final Review Gate

- Cocokkan hasil dengan definition_of_done.md.
- Jika ada poin gagal, kembali ke tahap terkait.

Output:
- Status done valid.

## Rule Fail Fast

- Jika migrasi gagal: hentikan tahap berikutnya.
- Jika build/typecheck gagal: jangan lanjut docs final.
- Jika privilege tidak sesuai spec: rollback patch terakhir dan perbaiki.

## Command Verifikasi per Tahap (Contoh XX99)

Urutan command di bawah ini dipakai sebagai baseline verifikasi implementasi XX99.

### Tahap Database (obx_base)

```sh
cd obx_base
npx prisma migrate dev --name xx99_init
npx prisma db seed
```

Expected:

- Migrasi sukses.
- Seed selesai tanpa error.

### Tahap Backend Build (obx_rest)

```sh
cd obx_rest
go mod tidy
go build ./...
```

Expected:

- Build berhasil tanpa compile error.

### Tahap Backend Run Smoke (obx_rest)

```sh
cd obx_rest
go run main.go
```

Smoke endpoint minimal:

- GET /
- GET /rest
- GET /rest/guest/

Expected:

- Service listen normal.
- Endpoint marker merespons.

### Tahap Frontend Check (obx_site)

```sh
cd obx_site
npm update
```

Expected:

- Dependency sinkron.

Catatan:

- Tambahkan command lint/typecheck/test sesuai script yang tersedia pada package frontend.

### Tahap Agent Dependency Sync (obx_auto)

```sh
cd obx_auto/agent_mx
go mod tidy
cd ../agent_vm
go mod tidy
cd ../agent_ws
go mod tidy
```

Expected:

- Seluruh agent tidy tanpa error.

## Catatan Audit

Setiap eksekusi autopilot wajib menyimpan:

1. module_code
2. daftar file berubah
3. command verifikasi dan hasil ringkas
4. risiko sisa