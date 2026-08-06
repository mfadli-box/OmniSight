# Rencana Task File-by-File Paket 5 Fase 2

Dokumen ini memecah Paket 5 Fase 2 menjadi task file-by-file untuk klaster File Transfer browser-based bastion.

## Scope Paket 5

- Modul domain: jms_file_transfer dan jms_session_event.
- Tujuan: menyediakan upload dan download file yang terkontrol dan dapat diaudit.

## Batch Task A - Finalisasi Model Transfer

### File Target
- obx_base/prisma/schema/jsm_stack.prisma
- obx_base/prisma/schema/schema.prisma

### Task
1. Review field transfer_type, source_path, target_path, file_name, file_size, dan status.
2. Pastikan relasi transfer ke session dan asset konsisten.
3. Pastikan index query histori transfer memadai.
4. Tetapkan status transfer minimum untuk sukses dan gagal.

### Verifikasi
- npx prisma generate di obx_base lulus.

## Batch Task B - API Transfer dan Audit Event

### File Target
- obx_rest/backbone/routes.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/template.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/repository.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_TRANSFER}/handler.go

### Task
1. Definisikan DTO request transfer upload dan download.
2. Implement endpoint inisiasi transfer berbasis token.
3. Simpan metadata transfer dan status hasil.
4. Simpan event FILE_UPLOAD dan FILE_DOWNLOAD di session event.
5. Pastikan write routes memakai logging yang konsisten.

### Verifikasi
- go build ./... di obx_rest lulus.
- Simulasi upload dan download menghasilkan histori transfer yang bisa ditelusuri.

## Batch Task C - Frontend Transfer Board

### File Target
- obx_site/src/app/board/pages/{KODE_BASTION_TRANSFER}/page.tsx
- obx_site/src/app/board/model/module.ts

### Task
1. Buat halaman transfer berbasis DataTable dan DataDialog.
2. Tampilkan metadata transfer: file, ukuran, sumber, tujuan, status.
3. Sediakan filter status transfer dan rentang waktu.
4. Tampilkan error transfer melalui toast dan peta error code.

### Verifikasi
- npm run lint di obx_site lulus.
- UI menampilkan riwayat transfer sukses dan gagal dengan benar.

## Batch Task D - Sinkronisasi Dokumentasi Transfer

### File Target
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md
- obx_docs/blueprint/PLAN/README.md

### Task
1. Sinkronkan alur transfer file dan audit event ke dokumentasi.
2. Tambahkan aturan minimum metadata transfer untuk audit.
3. Jelaskan perbedaan token transfer dan token SSH.
4. Tautkan paket ini ke backlog Fase 2.

### Verifikasi
- Dokumen BASE dan PLAN konsisten dengan endpoint transfer.

## Urutan Eksekusi yang Disarankan

1. Batch A - Schema
2. Batch B - API
3. Batch C - Frontend
4. Batch D - Dokumentasi

## Risiko

1. Metadata transfer tidak lengkap untuk audit.
2. Transfer berjalan di luar session context.
3. Endpoint transfer rentan dipakai ulang.

## Mitigasi

1. Wajibkan field transfer minimum di backend.
2. Kaitkan transfer ke session aktif bila tersedia.
3. Gunakan token short-lived dengan status sekali pakai.

## Referensi

- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
- ai_runbook.md
