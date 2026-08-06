# Rencana Task File-by-File Paket 4 Fase 2

Dokumen ini memecah Paket 4 Fase 2 menjadi task file-by-file untuk klaster Web SSH Bridge bastion.

## Scope Paket 4

- Modul domain: endpoint token SSH, lifecycle session SSH, event runtime SSH.
- Tujuan: menyediakan akses SSH berbasis browser yang stabil dan teraudit.

## Batch Task A - Finalisasi Kontrak Token SSH

### File Target
- obx_base/prisma/schema/jsm_stack.prisma
- obx_base/prisma/schema/schema.prisma

### Task
1. Review pemakaian connect_type SSH pada token.
2. Pastikan relasi token ke asset dan account sesuai kebutuhan SSH.
3. Tetapkan field minimum context koneksi untuk audit.
4. Pastikan status token dan expiry mudah divalidasi di runtime.

### Verifikasi
- npx prisma generate di obx_base lulus.

## Batch Task B - API Token dan Session SSH

### File Target
- obx_rest/backbone/routes.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/template.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/repository.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_CONNECT}/handler.go

### Task
1. Definisikan DTO issue token SSH dan create session SSH.
2. Implement endpoint open session SSH dengan token valid.
3. Implement endpoint close session dan forced revoke.
4. Simpan event CONNECT, COMMAND, DISCONNECT, dan TIMEOUT.
5. Pastikan write routes memakai logging yang konsisten.

### Verifikasi
- go build ./... di obx_rest lulus.
- Simulasi issue token SSH, open session, close session tercatat utuh.

## Batch Task C - Frontend Web SSH Board

### File Target
- obx_site/src/app/board/pages/{KODE_BASTION_CONNECT}/page.tsx
- obx_site/src/app/board/model/module.ts

### Task
1. Buat halaman Web SSH berbasis DataTable dan DataDialog untuk daftar sesi/token.
2. Sediakan aksi open session, disconnect, dan revoke.
3. Tampilkan status token dan status sesi secara real-time minimum.
4. Pastikan state loading, error handling, dan validasi form konsisten.

### Verifikasi
- npm run lint di obx_site lulus.
- UI dapat melakukan alur token SSH dan menampilkan audit sesi.

## Batch Task D - Sinkronisasi Dokumentasi SSH

### File Target
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md
- obx_docs/blueprint/PLAN/README.md

### Task
1. Sinkronkan alur SSH token dan session ke dokumentasi domain JMS.
2. Tambahkan batasan timeout dan forced revoke.
3. Jelaskan event minimum untuk audit SSH.
4. Tautkan paket ini ke backlog Fase 2.

### Verifikasi
- Dokumen BASE dan PLAN konsisten dengan endpoint SSH.

## Urutan Eksekusi yang Disarankan

1. Batch A - Schema
2. Batch B - API
3. Batch C - Frontend
4. Batch D - Dokumentasi

## Risiko

1. Session SSH terbuka tanpa jalur close yang jelas.
2. Event COMMAND tidak konsisten formatnya.
3. UI tidak sinkron saat sesi diputus paksa.

## Mitigasi

1. Wajibkan endpoint close dan revoke terpisah.
2. Standarkan payload event COMMAND di backend.
3. Refresh data sesi setelah aksi disconnect atau revoke.

## Referensi

- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
- ai_runbook.md
