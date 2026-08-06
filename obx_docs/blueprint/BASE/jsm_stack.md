# jsm_stack - Bastion dan Web Access Domain

## Ringkasan

- File schema: `obx_base/prisma/schema/jsm_stack.prisma`
- Isi utama: asset bastion, account, session, audit, token koneksi, web app proxy, dan file transfer

## Model Inti

- jms_asset_group
- jms_asset
- jms_asset_group_member
- jms_account
- jms_account_secret
- jms_session
- jms_session_command
- jms_audit_log
- jms_policy
- jms_approval
- jms_web_app
- jms_connect_token
- jms_session_event
- jms_file_transfer

## Enum Inti

- jms_connect_type: SSH, RDP, FTP, WEBAPP
- jms_token_status: ACTIVE, USED, EXPIRED, REVOKED
- jms_session_event_type: CONNECT, DISCONNECT, COMMAND, FILE_UPLOAD, FILE_DOWNLOAD, PROXY_ACCESS, DENY, TIMEOUT

## Relasi Penting

- jms_asset_group_member -> jms_asset + jms_asset_group
- jms_account -> jms_asset
- jms_account_secret -> jms_account
- jms_session -> jms_asset + jms_account
- jms_session_command -> jms_session
- jms_session_event -> jms_session
- jms_file_transfer -> jms_session + jms_asset
- jms_connect_token -> jms_asset + jms_account + jms_web_app

## Constraint dan Index Penting

- `jms_asset_group`: unique `[company_id, code]`
- `jms_asset`: unique `[company_id, code]`
- `jms_asset_group_member`: unique `[asset_id, group_id]`
- `jms_policy`: unique `[company_id, name]`
- `jms_web_app`: unique `[company_id, code]`
- `jms_connect_token.token_hash`: unique global

## Catatan Teknis

- `company_id` dan `user_id` tetap dipakai sebagai scope tenancy dan audit.
- Token koneksi disimpan sebagai hash, bukan token mentah.
- File transfer dipisahkan ke `jms_file_transfer` agar audit upload/download mudah dilacak.
- Enum pada domain `jsm_stack` (dengan prefiks `jms_`) diletakkan langsung di file schema domain agar perubahan tetap lokal dan tidak mengganggu enum lain.

## Checklist

- [ ] Generate Prisma lulus
- [ ] Migrasi schema dibuat
- [ ] Blueprint PLAN jsm_stack sinkron dengan schema aktual
- [ ] Guide BASE jsm_stack tersedia
