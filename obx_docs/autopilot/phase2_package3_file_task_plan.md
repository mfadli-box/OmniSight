# Rencana Task File-by-File Paket 3 Fase 2

Dokumen ini memecah Paket 3 Fase 2 menjadi task file-by-file untuk klaster Approval Access dan Connect Token bastion.

## Scope Paket 3

- Modul domain: jms_connect_token, jms_policy, jms_approval.
- Tujuan: menerapkan akses bastion yang terotorisasi, short-lived, one-time, dan dapat diaudit.

## Batch Task A - Finalisasi Model Token, Policy, dan Approval

### File Target
- `obx_base/prisma/schema/jsm_stack.prisma`
- `obx_base/prisma/schema/schema.prisma`

### Task
1. Review field token: connect_type, token_hash, status, expires_at, used_at.
2. Review relasi token ke asset, account, dan web app.
3. Review model policy dan approval untuk workflow akses.
4. Pastikan index query untuk validasi token aktif dan expiry memadai.

### Verifikasi
- `npx prisma generate` di `obx_base` lulus.

## Batch Task B - API Approval Access dan Connect Token

### File Target
- `obx_rest/backbone/routes.go`
- `obx_rest/skeleton/{KODE_BASTION_ACCESS}/template.go`
- `obx_rest/skeleton/{KODE_BASTION_ACCESS}/repository.go`
- `obx_rest/skeleton/{KODE_BASTION_ACCESS}/usecase.go`
- `obx_rest/skeleton/{KODE_BASTION_ACCESS}/handler.go`

### Task
1. Definisikan DTO request akses, approval action, dan issue token.
2. Implement endpoint issue token dengan one-time dan short-lived policy.
3. Implement endpoint approve/reject request akses.
4. Validasi scope company, user, asset/account/app saat issue token.
5. Pastikan token status transisi konsisten: ACTIVE -> USED/EXPIRED/REVOKED.

### Verifikasi
- `go build ./...` di `obx_rest` lulus.
- Token tidak bisa dipakai ulang setelah status USED.
- Request approval dapat diproses sampai keputusan final.

## Batch Task C - Frontend Board Access Control

### File Target
- `obx_site/src/app/board/pages/{KODE_BASTION_ACCESS}/page.tsx`
- `obx_site/src/app/board/model/module.ts`

### Task
1. Buat halaman request akses dan approval berbasis DataTable + DataDialog.
2. Tampilkan status request dan status token secara jelas.
3. Sediakan aksi approve/reject/revoke dari UI sesuai privilege.
4. Tampilkan informasi masa berlaku token secara eksplisit.

### Verifikasi
- `npm run lint` di `obx_site` lulus.
- Alur request -> approve -> issue token berjalan di UI.

## Batch Task D - Sinkronisasi Dokumentasi Access Control

### File Target
- `obx_docs/blueprint/BASE/jsm_stack.md`
- `obx_docs/guide/BASE/jsm_stack.md`
- `obx_docs/blueprint/PLAN/README.md`

### Task
1. Sinkronkan definisi connect token dan approval dengan schema aktual.
2. Tambahkan alur pengguna untuk request dan approval access.
3. Jelaskan one-time token rule dan expiry behavior.
4. Tautkan ke risiko keamanan dan mitigasi Fase 2.

### Verifikasi
- Dokumen BASE dan PLAN konsisten dengan endpoint approval/token.

## Urutan Eksekusi yang Disarankan

1. Batch A - Schema
2. Batch B - API
3. Batch C - Frontend
4. Batch D - Dokumentasi

## Risiko

1. Token tidak benar-benar one-time.
2. Approval bypass pada jalur tertentu.
3. UI tidak menunjukkan expiry/status terbaru.

## Mitigasi

1. Validasi token status di semua endpoint konsumsi token.
2. Wajibkan gate approval untuk asset kritikal.
3. Refresh status token setelah aksi approval/revoke.

## Referensi

- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
- ai_runbook.md
