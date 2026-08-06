# Rencana Task File-by-File Paket 6 Fase 2

Dokumen ini memecah Paket 6 Fase 2 menjadi task file-by-file untuk klaster WebAppProxy bastion.

## Scope Paket 6

- Modul domain: jms_web_app dan jms_connect_token.
- Tujuan: menyediakan akses aplikasi internal via proxy dengan allowlist dan audit access log.

## Batch Task A - Finalisasi Model Web App Proxy

### File Target
- obx_base/prisma/schema/jsm_stack.prisma
- obx_base/prisma/schema/schema.prisma

### Task
1. Review field web app: code, target_url, strip_prefix, allowed_roles, is_active.
2. Pastikan relasi token WEBAPP ke app valid.
3. Pastikan index query list app aktif memadai.
4. Tetapkan aturan validasi target_url dan role allowlist.

### Verifikasi
- npx prisma generate di obx_base lulus.

## Batch Task B - API Registry dan Proxy Guard

### File Target
- obx_rest/backbone/routes.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/template.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/repository.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_PROXY}/handler.go

### Task
1. Definisikan DTO registry app dan akses proxy.
2. Implement endpoint list/create/update app proxy target.
3. Implement guard permission berbasis company, role, dan token.
4. Pastikan deny terhadap host atau path di luar allowlist.
5. Simpan access log untuk setiap request proxy.

### Verifikasi
- go build ./... di obx_rest lulus.
- Request proxy ke target tidak terdaftar ditolak.
- Request ke target valid tercatat pada access log.

## Batch Task C - Frontend WebAppProxy Board

### File Target
- obx_site/src/app/board/pages/{KODE_BASTION_PROXY}/page.tsx
- obx_site/src/app/board/model/module.ts

### Task
1. Buat halaman registry WebAppProxy berbasis DataTable dan DataDialog.
2. Tampilkan daftar app, target_url, role, dan status aktif.
3. Sediakan aksi enable atau disable target app.
4. Tambahkan panel sederhana untuk inspeksi hasil akses proxy.

### Verifikasi
- npm run lint di obx_site lulus.
- UI dapat kelola registry app proxy dan melihat status akses.

## Batch Task D - Sinkronisasi Dokumentasi WebAppProxy

### File Target
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md
- obx_docs/blueprint/PLAN/README.md

### Task
1. Sinkronkan definisi app registry dan proxy guard.
2. Tambahkan aturan allowlist host dan path.
3. Jelaskan log akses dan mekanisme deny.
4. Tautkan paket ini ke backlog Fase 2.

### Verifikasi
- Dokumen BASE dan PLAN konsisten dengan endpoint proxy.

## Urutan Eksekusi yang Disarankan

1. Batch A - Schema
2. Batch B - API
3. Batch C - Frontend
4. Batch D - Dokumentasi

## Risiko

1. Proxy membuka akses ke target yang tidak diinginkan.
2. Role guard tidak sinkron dengan policy approval.
3. Access log tidak cukup detail untuk audit.

## Mitigasi

1. Terapkan validasi allowlist ketat pada backend.
2. Selaraskan guard proxy dengan policy akses.
3. Simpan log minimal: actor, target, waktu, status.

## Referensi

- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
- ai_runbook.md
