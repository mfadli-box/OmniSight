# Rencana Task File-by-File Paket 2 Fase 2

Dokumen ini memecah Paket 2 Fase 2 menjadi task file-by-file untuk klaster Session Core dan Audit bastion.

## Scope Paket 2

- Modul domain: jms_session, jms_session_event, jms_session_command, jms_audit_log.
- Tujuan: memastikan seluruh lifecycle sesi bastion tercatat konsisten dan dapat diaudit.

## Batch Task A - Finalisasi Model Session dan Event

### File Target
- `obx_base/prisma/schema/jsm_stack.prisma`
- `obx_base/prisma/schema/schema.prisma`

### Task
1. Review field lifecycle session: status, started_at, ended_at, connection_info.
2. Review event type yang dipakai lintas protokol (SSH, RDP, FTP, WEBAPP).
3. Pastikan relasi session ke asset, account, user, event, dan command valid.
4. Pastikan index query untuk histori sesi dan audit timeline memadai.

### Verifikasi
- `npx prisma generate` di `obx_base` lulus.

## Batch Task B - API Session dan Audit

### File Target
- `obx_rest/backbone/routes.go`
- `obx_rest/skeleton/{KODE_BASTION_SESSION}/template.go`
- `obx_rest/skeleton/{KODE_BASTION_SESSION}/repository.go`
- `obx_rest/skeleton/{KODE_BASTION_SESSION}/usecase.go`
- `obx_rest/skeleton/{KODE_BASTION_SESSION}/handler.go`

### Task
1. Definisikan DTO untuk create session, close session, log event, dan log command.
2. Implement endpoint list/detail session dengan filter company, user, status, waktu.
3. Implement endpoint event dan command log per session.
4. Pastikan jalur revoke/disconnect menutup session dan menulis event audit.
5. Pastikan write routes memakai logging yang konsisten.

### Verifikasi
- `go build ./...` di `obx_rest` lulus.
- Simulasi create session -> event -> command -> disconnect tercatat utuh.

## Batch Task C - Frontend Board Session Audit

### File Target
- `obx_site/src/app/board/pages/{KODE_BASTION_SESSION}/page.tsx`
- `obx_site/src/app/board/model/module.ts`

### Task
1. Buat halaman monitoring session berbasis DataTable + DataDialog.
2. Tampilkan event timeline dan command list per session.
3. Tambahkan aksi revoke session dari UI dengan konfirmasi.
4. Pastikan pencarian/filter berdasarkan status, user, asset, dan waktu.

### Verifikasi
- `npm run lint` di `obx_site` lulus.
- UI session audit menampilkan timeline dan command log sesuai backend.

## Batch Task D - Sinkronisasi Dokumentasi Session

### File Target
- `obx_docs/blueprint/BASE/jsm_stack.md`
- `obx_docs/guide/BASE/jsm_stack.md`
- `obx_docs/blueprint/PLAN/README.md`

### Task
1. Sinkronkan lifecycle session dan event type dengan schema aktual.
2. Tambahkan alur audit trail dari connect sampai disconnect.
3. Jelaskan batasan command logging untuk protokol non-shell.
4. Tautkan ke backlog Fase 2 paket session.

### Verifikasi
- Dokumen BASE dan PLAN konsisten dengan endpoint dan payload session.

## Urutan Eksekusi yang Disarankan

1. Batch A - Schema
2. Batch B - API
3. Batch C - Frontend
4. Batch D - Dokumentasi

## Risiko

1. Event model tidak konsisten lintas protokol.
2. Session tertutup tanpa event penutup.
3. UI menampilkan status yang terlambat dari backend.

## Mitigasi

1. Tetapkan event mandatory minimum: CONNECT dan DISCONNECT.
2. Paksa write event saat close/revoke session.
3. Gunakan refresh polling/trigger setelah aksi kritikal.

## Referensi

- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
- ai_runbook.md
