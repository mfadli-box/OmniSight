# Rencana Task File-by-File Paket 7 Fase 2

Dokumen ini memecah Paket 7 Fase 2 menjadi task file-by-file untuk klaster RDP Pilot dan Recovery Controls bastion.

## Scope Paket 7

- Modul domain: token RDP, session revoke controls, dan runtime recovery note.
- Tujuan: menyediakan RDP browser access untuk pilot terbatas dengan kontrol pemulihan yang jelas.

## Batch Task A - Finalisasi Model dan Status RDP

### File Target
- obx_base/prisma/schema/jsm_stack.prisma
- obx_base/prisma/schema/schema.prisma

### Task
1. Review pemakaian connect_type RDP pada token.
2. Pastikan field session mendukung status reconnect dan revoke.
3. Tetapkan event minimum RDP untuk audit.
4. Pastikan index untuk query sesi RDP memadai.

### Verifikasi
- npx prisma generate di obx_base lulus.

## Batch Task B - API Token, Session, dan Revoke RDP

### File Target
- obx_rest/backbone/routes.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/template.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/repository.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/usecase.go
- obx_rest/skeleton/{KODE_BASTION_RDP}/handler.go

### Task
1. Definisikan DTO issue token RDP dan open session RDP.
2. Implement endpoint revoke dan forced disconnect.
3. Simpan event CONNECT, DISCONNECT, TIMEOUT, dan DENY.
4. Tambahkan status recovery note ketika gateway gagal.
5. Pastikan write routes memakai logging yang konsisten.

### Verifikasi
- go build ./... di obx_rest lulus.
- Revoke dapat memutus sesi RDP aktif pada simulasi.

## Batch Task C - Frontend RDP Pilot Board

### File Target
- obx_site/src/app/board/pages/{KODE_BASTION_RDP}/page.tsx
- obx_site/src/app/board/model/module.ts

### Task
1. Buat halaman RDP pilot berbasis DataTable dan DataDialog.
2. Tampilkan status koneksi, timeout, dan hasil revoke.
3. Tambahkan notifikasi recovery saat terjadi kegagalan gateway.
4. Pastikan tombol aksi dibuat touch-friendly dan mobile-friendly.

### Verifikasi
- npm run lint di obx_site lulus.
- UI dapat menampilkan status sesi RDP dan aksi revoke.

## Batch Task D - Sinkronisasi Dokumentasi RDP

### File Target
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md
- obx_docs/blueprint/PLAN/README.md

### Task
1. Sinkronkan alur RDP pilot dan batasannya.
2. Tambahkan recovery guidance saat gateway bermasalah.
3. Jelaskan perbedaan jalur RDP dibanding SSH.
4. Tautkan paket ini ke backlog Fase 2.

### Verifikasi
- Dokumen BASE dan PLAN konsisten dengan endpoint RDP.

## Urutan Eksekusi yang Disarankan

1. Batch A - Schema
2. Batch B - API
3. Batch C - Frontend
4. Batch D - Dokumentasi

## Risiko

1. RDP pilot belum stabil pada jaringan tertentu.
2. Revoke tidak selalu memutus koneksi runtime.
3. Recovery flow tidak terdokumentasi jelas.

## Mitigasi

1. Batasi pilot RDP pada scope kecil lebih dulu.
2. Uji forced disconnect secara berkala.
3. Simpan recovery note standar di dokumen operasional.

## Referensi

- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
- ai_runbook.md
