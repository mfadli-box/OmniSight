# Rancangan Submodul: JMS - Adaptasi JumpServer ke Stack OmniSight

## 1. Overview

- Group: Security / Operations
- Icon: Shield
- Akses: Admin, Operator, Auditor
- Tujuan: mengadaptasi konsep inti JumpServer (asset, akun, session, audit, policy) ke arsitektur Gin + Prisma + Next.js yang sudah dipakai di repo ini
- Referensi skema: GitHub repository jumpserver/jumpserver, khusus area asset management, account management, session recording, audit log

## 2. Konsep Mapping ke Stack Saat Ini

### Backend
- Framework: Go + Gin di obx_rest
- Pola: skeleton repository/usecase/handler seperti modul SM01–SM05 dan XX99
- Routing: route halaman di /rest/pages/JMS dan route robot untuk webhook/agent di /rest/robot/JMS

### Database
- ORM: Prisma di obx_base/prisma
- Konvensi: UUID string, company_id untuk tenancy, created_at/updated_at untuk audit, status dan is_active untuk lifecycle

### Frontend
- Framework: Next.js di obx_site
- Pola: DataTable + DataDialog + module page berbasis board

## 3. Database Schema

### Tabel utama yang disarankan

- jms_asset
  - id, company_id, code, name, asset_type, address, protocol, port, platform, status, is_active
- jms_asset_group
  - id, company_id, code, name, parent_id
- jms_asset_group_member
  - id, asset_id, group_id
- jms_account
  - id, company_id, name, username, secret_ref, account_type, is_active
- jms_account_secret
  - id, company_id, account_id, secret_value, secret_mode, expires_at
- jms_session
  - id, company_id, asset_id, account_id, user_id, session_type, started_at, ended_at, status, connection_info
- jms_session_command
  - id, session_id, command, executed_at, user_id
- jms_audit_log
  - id, company_id, actor_id, action, target_type, target_id, detail, created_at
- jms_policy
  - id, company_id, name, policy_type, is_active, config_json
- jms_approval
  - id, company_id, request_type, request_id, approver_id, status, note
- jms_connect_token
  - id, company_id, user_id, asset_id, account_id, app_id, connect_type, token_hash, status, expires_at, used_at
- jms_web_app
  - id, company_id, code, name, target_url, strip_prefix, allowed_roles, is_active
- jms_session_event
  - id, session_id, user_id, event_type, message, command_text, payload_json, created_at
- jms_file_transfer
  - id, company_id, session_id, asset_id, user_id, transfer_type, source_path, target_path, file_name, file_size, status, created_at

### Relasi penting
- jms_asset -> jms_asset_group_member -> jms_asset_group
- jms_account -> jms_account_secret
- jms_session -> jms_asset, jms_account, dat_user
- jms_session_command -> jms_session
- jms_policy -> jms_audit_log (opsional untuk tracing)

## 4. API Endpoints

| Method | Endpoint | Handler | Keterangan |
|---|---|---|---|
| GET | /rest/pages/JMS | JMSHandler.ListAsset | Daftar asset |
| GET | /rest/pages/JMS/:id | JMSHandler.DetailAsset | Detail asset |
| POST | /rest/pages/JMS | JMSHandler.CreateAsset | Buat asset |
| PUT | /rest/pages/JMS/:id | JMSHandler.UpdateAsset | Ubah asset |
| DELETE | /rest/pages/JMS/:id | JMSHandler.DeleteAsset | Hapus asset |
| GET | /rest/pages/JMS/account | JMSHandler.ListAccount | Daftar akun |
| POST | /rest/pages/JMS/session | JMSHandler.CreateSession | Buat catatan sesi |
| GET | /rest/pages/JMS/session | JMSHandler.ListSession | Daftar sesi |
| GET | /rest/pages/JMS/session/:id | JMSHandler.DetailSession | Detail sesi |
| POST | /rest/robot/JMS/heartbeat | JMSHandler.Heartbeat | Heartbeat agent |
| POST | /rest/robot/JMS/session | JMSHandler.ReceiveSession | Sinkronisasi sesi dari agent |
| POST | /rest/pages/JMS/connect/token | JMSHandler.CreateConnectToken | Buat token koneksi sementara |
| GET | /rest/pages/JMS/connect/ws/:sessionId | JMSHandler.ProxyConnectWS | WebSocket proxy untuk SSH/RDP bridge |
| GET | /rest/pages/JMS/proxy/:appId/*path | JMSHandler.WebAppProxy | Reverse proxy web app internal |
| POST | /rest/pages/JMS/file/token | JMSHandler.CreateFileTransferToken | Buat token file transfer sementara |
| GET | /rest/pages/JMS/file/ws/:sessionId | JMSHandler.ProxyFileTransferWS | WebSocket/file stream proxy untuk FTP/SFTP browser-based |

## 5. Frontend Pages

- Halaman utama: /board/pages/JMS
- Subpage: /board/pages/JMS/assets
- Subpage: /board/pages/JMS/accounts
- Subpage: /board/pages/JMS/sessions
- Subpage: /board/pages/JMS/audit
- Subpage: /board/pages/JMS/connect
- Subpage: /board/pages/JMS/files
- Komponen utama: DataTable, DataDialog, detail drawer, filter panel, action menu
- Field form utama: name, asset_type, address, protocol, port, platform, username, account_type, status
- Mobile behavior: daftar ringkas, drawer detail untuk sesi dan audit

## 5A. Web Access Feature (RDP, SSH, FTP, Web App Proxy)

### Tujuan
- Menyediakan akses jarak jauh langsung dari browser tanpa perlu client lokal.

### Arsitektur ringkas
- Frontend Next.js membuka koneksi ke endpoint backend JMS.
- Backend Gin mengeluarkan token sesi sementara dan meneruskan koneksi ke service bridge (gateway).
- Gateway mendukung SSH terminal stream, RDP stream, file transfer browser-based, dan reverse proxy aplikasi web internal.

### Kapabilitas
- SSH via web terminal (WebSocket stream + PTY bridge).
- RDP via browser canvas/HTML5 client melalui gateway.
- FTP/SFTP browser-based untuk upload/download file terkontrol.
- Web app proxy berbasis path untuk dashboard internal (Grafana, Portainer, Kibana, dan sejenisnya).

### Kebijakan keamanan minimum
- Token akses bersifat short-lived dan one-time.
- RBAC per user-company-asset wajib diverifikasi sebelum connect.
- Semua aktivitas connect/disconnect dicatat ke jms_audit_log.
- Aktivitas upload/download file dicatat ke `jms_file_transfer` dan `jms_session_event`.
- Pembatasan allowlist destination host dan port.

## 6. Implementasi Rekomendasi

### Backend
- Buat skeleton baru di obx_rest/skeleton/JMS dengan pola:
  - repository.go
  - usecase.go
  - handler.go
- Tambahkan route baru di obx_rest/backbone/routes.go
- Gunakan middleware yang sudah ada: USAuth, USLock, USLogs("JMS"), USRobot

### Database
- Tambahkan model Prisma di obx_base/prisma/schema
- Jalankan migrate dan generate client
- Pastikan company_id dan relasi tenancy konsisten dengan modul lain

### Frontend
- Tambahkan page di obx_site/src/app/board/pages/JMS
- Tambahkan entry menu di obx_site/src/app/board/model/module.ts
- Pakai endpoint proxy /proxy/pages/JMS untuk menghindari CORS dan konsistensi auth

## 7. Relationship dengan Modul Eksisting

- Terhubung ke modul user/company dari obx_base: dat_user, dat_company
- Memanfaatkan mekanisme session/auth yang sudah ada di obx_rest/backbone/memory.go
- Dapat dipakai bersama modul AUTO untuk agent/agent_vm/agent_mx sebagai sumber data host dan status

## 8. Implementation Checklist

- [ ] Prisma schema siap
- [ ] Backend skeleton JMS selesai
- [ ] Route terdaftar di obx_rest/backbone/routes.go
- [ ] Frontend page dan menu selesai
- [ ] Endpoint proxy frontend siap
- [ ] Test login, CRUD, dan session ingestion lulus
- [ ] Web SSH dari browser lulus
- [ ] Web RDP dari browser lulus
- [ ] FTP / file transfer browser-based lulus
- [ ] Web app proxy dari browser lulus
- [ ] Dokumentasi user guide diperbarui

## 9. Backend Pattern Reference

- Error handling: gunakan mechanic.Error dan mechanic.InternalError
- Query: parameterized query dan whitelist sort
- Routing: static route sebelum dynamic route
- Logging: semua route write wajib memakai USLogs("JMS")
- Auth: route halaman pakai USAuth/USLock, route robot pakai USRobot

## 10. Rekomendasi Scope Awal

Untuk tahap MVP, fokus pada 3 fitur pertama:
1. Asset inventory
2. Account management
3. Session/audit log

Setelah itu baru tambah fitur:
4. Policy approval
5. Secret rotation
6. Session recording integration
