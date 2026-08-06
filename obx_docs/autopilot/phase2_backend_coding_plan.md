# Rencana Coding Backend Fase 2 - Bastion Core

Dokumen ini menurunkan backlog backend Fase 2 menjadi paket coding yang lebih konkret untuk dieksekusi di repo.

## Tujuan

1. Mengaktifkan domain bastion berbasis schema JMS yang sudah tersedia.
2. Menyediakan jalur akses terotorisasi melalui connect token dan approval.
3. Menyediakan audit trail sesi, command, transfer file, dan proxy access.
4. Menjaga setiap paket kecil, terverifikasi, dan mudah di-rollback.

## Prinsip Coding

- Satu paket fokus pada satu klaster domain bastion.
- Schema, backend API, frontend board, dan dokumentasi terkait harus selesai dalam batch yang sama bila memungkinkan.
- Jangan lanjut ke paket berikutnya sebelum verifikasi paket berjalan.

## Paket 1 - Asset dan Account Inventory

### Target Modul
- jms_asset_group
- jms_asset
- jms_asset_group_member
- jms_account
- jms_account_secret

### File yang Ditargetkan
- `obx_base/prisma/schema/jsm_stack.prisma`
- `obx_rest/skeleton/{KODE_BASTION_ASSET}/template.go`
- `obx_rest/skeleton/{KODE_BASTION_ASSET}/repository.go`
- `obx_rest/skeleton/{KODE_BASTION_ASSET}/usecase.go`
- `obx_rest/skeleton/{KODE_BASTION_ASSET}/handler.go`
- `obx_rest/backbone/routes.go`
- `obx_site/src/app/board/pages/{KODE_BASTION_ASSET}/page.tsx`
- `obx_site/src/app/board/model/module.ts`
- `obx_docs/blueprint/BASE/jsm_stack.md`
- `obx_docs/guide/BASE/jsm_stack.md`

### Aksi Coding
1. Finalisasi validasi field minimum asset, account, dan secret.
2. Pastikan company scope konsisten di seluruh query.
3. Terapkan masking/abstraksi secret di response publik.
4. Siapkan list endpoint untuk pilot inventory bastion.

### Verifikasi
- CRUD asset dan account berjalan.
- Secret reference tidak diekspos ke payload list/detail.
- Build backend lulus.

## Paket 2 - Session Core dan Audit

### Target Modul
- jms_session
- jms_session_event
- jms_session_command
- jms_audit_log

### File yang Ditargetkan
- `obx_base/prisma/schema/jsm_stack.prisma`
- `obx_rest/skeleton/{KODE_BASTION_SESSION}/template.go`
- `obx_rest/skeleton/{KODE_BASTION_SESSION}/repository.go`
- `obx_rest/skeleton/{KODE_BASTION_SESSION}/usecase.go`
- `obx_rest/skeleton/{KODE_BASTION_SESSION}/handler.go`
- `obx_rest/backbone/routes.go`
- `obx_site/src/app/board/pages/{KODE_BASTION_SESSION}/page.tsx`
- `obx_site/src/app/board/model/module.ts`
- `obx_docs/blueprint/BASE/jsm_stack.md`
- `obx_docs/guide/BASE/jsm_stack.md`

### Aksi Coding
1. Definisikan lifecycle session dan status transisi.
2. Pisahkan penyimpanan event dan command log.
3. Hubungkan event connect/disconnect ke audit log.
4. Pastikan actor, asset, dan account tercatat penuh.

### Verifikasi
- Session event dan command log tercatat saat simulasi.
- Revoke/disconnect tercatat konsisten.
- Build backend lulus.

## Paket 3 - Approval Access dan Connect Token

### Target Modul
- jms_connect_token
- jms_policy
- jms_approval

### File yang Ditargetkan
- `obx_base/prisma/schema/jsm_stack.prisma`
- `obx_rest/skeleton/{KODE_BASTION_ACCESS}/template.go`
- `obx_rest/skeleton/{KODE_BASTION_ACCESS}/repository.go`
- `obx_rest/skeleton/{KODE_BASTION_ACCESS}/usecase.go`
- `obx_rest/skeleton/{KODE_BASTION_ACCESS}/handler.go`
- `obx_rest/backbone/routes.go`
- `obx_site/src/app/board/pages/{KODE_BASTION_ACCESS}/page.tsx`
- `obx_site/src/app/board/model/module.ts`
- `obx_docs/blueprint/BASE/jsm_stack.md`
- `obx_docs/guide/BASE/jsm_stack.md`

### Aksi Coding
1. Implement token short-lived dan one-time use.
2. Terapkan validasi status token: ACTIVE, USED, EXPIRED, REVOKED.
3. Kaitkan policy dan approval ke issuance token.
4. Validasi scope user, company, asset, account, dan app.

### Verifikasi
- Token tidak bisa dipakai ulang.
- Approval flow dapat membuka atau menolak akses.
- Build backend lulus.

## Paket 4 - Web SSH Bridge

### Target Modul
- endpoint token SSH
- lifecycle session SSH
- event runtime SSH

### File yang Ditargetkan
- `obx_rest/skeleton/{KODE_BASTION_CONNECT}/template.go`
- `obx_rest/skeleton/{KODE_BASTION_CONNECT}/repository.go`
- `obx_rest/skeleton/{KODE_BASTION_CONNECT}/usecase.go`
- `obx_rest/skeleton/{KODE_BASTION_CONNECT}/handler.go`
- `obx_rest/backbone/routes.go`
- `obx_site/src/app/board/pages/{KODE_BASTION_CONNECT}/page.tsx`
- `obx_docs/blueprint/PLAN/README.md`

### Aksi Coding
1. Sediakan endpoint pembuatan token untuk SSH.
2. Integrasikan bridge PTY/WebSocket runtime.
3. Simpan event connect/disconnect/timeout.
4. Tambahkan forced disconnect pada revoke.

### Verifikasi
- SSH via browser berjalan pada pilot host.
- Audit session SSH tercatat.

## Paket 5 - File Transfer Browser-Based

### Target Modul
- jms_file_transfer
- jms_session_event

### File yang Ditargetkan
- `obx_base/prisma/schema/jsm_stack.prisma`
- `obx_rest/skeleton/{KODE_BASTION_TRANSFER}/template.go`
- `obx_rest/skeleton/{KODE_BASTION_TRANSFER}/repository.go`
- `obx_rest/skeleton/{KODE_BASTION_TRANSFER}/usecase.go`
- `obx_rest/skeleton/{KODE_BASTION_TRANSFER}/handler.go`
- `obx_rest/backbone/routes.go`
- `obx_site/src/app/board/pages/{KODE_BASTION_TRANSFER}/page.tsx`
- `obx_docs/blueprint/BASE/jsm_stack.md`
- `obx_docs/guide/BASE/jsm_stack.md`

### Aksi Coding
1. Sediakan token transfer terpisah dari token SSH/RDP.
2. Implement upload/download stream terkontrol.
3. Simpan metadata transfer untuk audit.
4. Kaitkan transfer event dengan session.

### Verifikasi
- Transfer sukses dan gagal tercatat.
- Riwayat transfer dapat ditelusuri.

## Paket 6 - WebAppProxy

### Target Modul
- jms_web_app
- jms_connect_token
- request access log

### File yang Ditargetkan
- `obx_base/prisma/schema/jsm_stack.prisma`
- `obx_rest/skeleton/{KODE_BASTION_PROXY}/template.go`
- `obx_rest/skeleton/{KODE_BASTION_PROXY}/repository.go`
- `obx_rest/skeleton/{KODE_BASTION_PROXY}/usecase.go`
- `obx_rest/skeleton/{KODE_BASTION_PROXY}/handler.go`
- `obx_rest/backbone/routes.go`
- `obx_site/src/app/board/pages/{KODE_BASTION_PROXY}/page.tsx`
- `obx_docs/blueprint/BASE/jsm_stack.md`
- `obx_docs/guide/BASE/jsm_stack.md`

### Aksi Coding
1. Definisikan registry app dan allowlist target URL.
2. Terapkan guard permission per company/role.
3. Simpan request log per akses proxy.
4. Pastikan proxy menolak target di luar allowlist.

### Verifikasi
- App proxy hanya berjalan untuk target terdaftar.
- Log akses proxy dapat diinspeksi.

## Paket 7 - RDP Pilot dan Recovery Controls

### Target Modul
- endpoint token RDP
- session revoke controls
- runtime recovery note

### File yang Ditargetkan
- `obx_rest/skeleton/{KODE_BASTION_RDP}/template.go`
- `obx_rest/skeleton/{KODE_BASTION_RDP}/repository.go`
- `obx_rest/skeleton/{KODE_BASTION_RDP}/usecase.go`
- `obx_rest/skeleton/{KODE_BASTION_RDP}/handler.go`
- `obx_rest/backbone/routes.go`
- `obx_site/src/app/board/pages/{KODE_BASTION_RDP}/page.tsx`
- `obx_docs/blueprint/PLAN/README.md`

### Aksi Coding
1. Siapkan jalur token dan koneksi RDP untuk pilot.
2. Tambahkan timeout, reconnect, dan forced revoke.
3. Catat event RDP ke audit trail session.
4. Dokumentasikan recovery flow runtime.

### Verifikasi
- RDP pilot dapat diuji terbatas.
- Revoke memutus sesi aktif.

## Urutan Coding yang Disarankan

1. Paket 1 - Asset dan Account Inventory
2. Paket 2 - Session Core dan Audit
3. Paket 3 - Approval Access dan Connect Token
4. Paket 4 - Web SSH Bridge
5. Paket 5 - File Transfer Browser-Based
6. Paket 6 - WebAppProxy
7. Paket 7 - RDP Pilot dan Recovery Controls

## DoD Coding Fase 2

- Inventory asset/account bastion aktif.
- Session audit dan revoke konsisten.
- Connect token one-time berjalan aman.
- Web SSH, transfer file, dan WebAppProxy siap pilot.
- RDP pilot minimal siap uji terbatas.
- Dokumentasi BASE dan PLAN sinkron.

## Risiko

1. Scope runtime bridge terlalu cepat membesar.
2. Token dan approval flow tidak sinkron.
3. Session event tidak konsisten antar protokol.

## Mitigasi

1. Kunci urutan paket dan jangan lompat fase.
2. Verifikasi one-time token di setiap endpoint koneksi.
3. Standarkan event model untuk SSH, FTP, WEBAPP, dan RDP.

## Referensi

- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
- roadmap_platform_3_phase.md
- backend_module_matrix_platform_replacement.md
- ai_runbook.md
