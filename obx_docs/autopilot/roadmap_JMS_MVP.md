# Roadmap MVP: JMS (JumpServer Adaptation)

Dokumen ini adalah rencana eksekusi bertahap untuk membangun MVP JMS pada stack OmniSight:
- Database: Prisma (obx_base)
- Backend: Go + Gin (obx_rest)
- Frontend: Next.js board (obx_site)

## Target MVP

MVP dianggap selesai jika 3 kapabilitas utama berjalan end-to-end:
1. Asset inventory
2. Account management
3. Session dan audit log

Tambahan target akses web langsung:
4. Web SSH (browser terminal)
5. Web RDP (browser remote desktop)
6. FTP / file transfer browser-based
7. Web app proxy (akses aplikasi internal via browser)

## Prinsip Eksekusi

- Perubahan kecil per batch agar mudah rollback.
- Setiap batch wajib punya bukti verifikasi command.
- Route write wajib logging dan auth sesuai pola existing.

## Fase 0 - Foundation (0.5-1 hari)

### Tujuan
- Menetapkan kontrak data, route, dan boundary scope MVP.

### Output
- Blueprint JMS final untuk MVP.
- Daftar field minimal per entitas.
- Daftar endpoint prioritas MVP.

### Checklist
- [ ] Freeze scope tabel: jms_asset, jms_account, jms_session, jms_audit_log
- [ ] Freeze endpoint CRUD minimum
- [ ] Freeze UI minimum (table + form + detail)

## Fase 1 - Database Layer (1 hari)

### Tujuan
- Menyediakan schema Prisma yang cukup untuk 3 kapabilitas MVP.

### Implementasi
- Tambah model Prisma minimal:
  - jms_asset
  - jms_account
  - jms_session
  - jms_audit_log
- Tambah index untuk query list/sort/filter.
- Tambah seed minimal untuk smoke test JMS.

### Verifikasi wajib
- `cd obx_base; npx prisma migrate dev --name jms_mvp_init`
- `cd obx_base; npx prisma generate`
- `cd obx_base; npx prisma db seed`

### Exit criteria
- Migrasi sukses.
- Prisma client update.
- Data seed minimal tersedia.

### Draft skema Prisma minimal (tambahan Web SSH/RDP/Web App Proxy)

Tujuan draft ini adalah memberi kontrak data minimum untuk fase akses web tanpa menunggu semua fitur policy selesai.

```prisma
enum jms_connect_type {
  SSH
  RDP
  FTP
  WEBAPP
}

enum jms_token_status {
  ACTIVE
  USED
  EXPIRED
  REVOKED
}

enum jms_session_event_type {
  CONNECT
  DISCONNECT
  COMMAND
  FILE_UPLOAD
  FILE_DOWNLOAD
  PROXY_ACCESS
  DENY
  TIMEOUT
}

model jms_connect_token {
  id           String            @id @default(uuid())
  company_id   String
  user_id      String
  asset_id     String?
  account_id   String?
  app_id       String?
  connect_type jms_connect_type
  token_hash   String            @unique
  status       jms_token_status  @default(ACTIVE)
  client_ip    String?
  user_agent   String?
  expires_at   DateTime
  used_at      DateTime?
  created_at   DateTime          @default(now())
  updated_at   DateTime          @default(now())

  @@index([company_id])
  @@index([user_id])
  @@index([status])
  @@index([expires_at])
  @@index([connect_type])
}

model jms_web_app {
  id            String   @id @default(uuid())
  company_id    String
  code          String
  name          String
  target_url    String
  strip_prefix  Boolean  @default(true)
  allowed_roles String?
  is_active     Boolean  @default(true)
  created_at    DateTime @default(now())
  updated_at    DateTime @default(now())

  @@unique([company_id, code])
  @@index([company_id])
  @@index([is_active])
}

model jms_session_event {
  id           String                 @id @default(uuid())
  session_id   String
  user_id      String?
  event_type   jms_session_event_type
  message      String?
  command_text String?
  payload_json String?
  created_at   DateTime               @default(now())

  @@index([session_id])
  @@index([event_type])
  @@index([created_at])
}

model jms_file_transfer {
  id            String   @id @default(uuid())
  company_id    String
  session_id    String?
  asset_id      String?
  user_id       String
  transfer_type String
  source_path   String?
  target_path   String?
  file_name     String
  file_size     BigInt?
  status        String
  created_at    DateTime @default(now())

  @@index([company_id])
  @@index([session_id])
  @@index([user_id])
  @@index([created_at])
}
```

Catatan implementasi:
- `token_hash` menyimpan hash token, bukan token mentah.
- TTL token dipenuhi lewat `expires_at`; token one-time via `used_at` + `status`.
- `jms_web_app.target_url` wajib divalidasi allowlist sebelum dipakai proxy.
- `jms_session_event` dipakai untuk audit granular SSH/RDP/WEBAPP per sesi.
- `jms_file_transfer` dipakai untuk audit upload/download file dan quota/forensik dasar.

## Fase 2 - Backend CRUD Asset + Account (1-1.5 hari)

### Tujuan
- Menyediakan API utama asset dan account.

### Implementasi
- Buat skeleton baru:
  - obx_rest/skeleton/JMS/template.go
  - obx_rest/skeleton/JMS/repository.go
  - obx_rest/skeleton/JMS/usecase.go
  - obx_rest/skeleton/JMS/handler.go
- Daftarkan route di routes.go:
  - GET /rest/pages/JMS
  - GET /rest/pages/JMS/:id
  - POST /rest/pages/JMS
  - PUT /rest/pages/JMS/:id
  - DELETE /rest/pages/JMS/:id
  - GET /rest/pages/JMS/account
- Terapkan middleware:
  - USAuth + USLock untuk pages
  - USLogs("JMS") untuk write route

### Verifikasi wajib
- `cd obx_rest; go build ./...`
- `cd obx_rest; go run main.go`
- Uji endpoint dengan token login (smoke GET/POST/PUT/DELETE)

### Exit criteria
- Build backend lulus.
- CRUD asset lulus.
- List account lulus.

## Fase 3 - Backend Session + Audit (1 hari)

### Tujuan
- Menyediakan pelacakan sesi dan jejak audit operasional.

### Implementasi
- Tambah endpoint session:
  - POST /rest/pages/JMS/session
  - GET /rest/pages/JMS/session
  - GET /rest/pages/JMS/session/:id
- Tambah endpoint audit:
  - GET /rest/pages/JMS/audit
- Integrasi logging action ke jms_audit_log pada event write penting.

### Verifikasi wajib
- `cd obx_rest; go build ./...`
- Uji create session -> list session -> detail session
- Uji audit list setelah operasi write

### Exit criteria
- Session create/list/detail lulus.
- Audit trail terbaca dan konsisten.

## Fase 4 - Frontend Board JMS (1-1.5 hari)

### Tujuan
- Menyediakan UI operasional minimum untuk Asset, Account, Session, Audit.

### Implementasi
- Tambah menu JMS pada board module (admin).
- Buat page utama JMS dan sub-tab:
  - Assets
  - Accounts
  - Sessions
  - Audit
- Gunakan pola DataTable + DataDialog yang konsisten dengan modul existing.
- Gunakan proxy endpoint untuk akses API.

### Verifikasi wajib
- `cd obx_site; npm run dev`
- `cd obx_site; npm run lint`
- Uji flow UI: create asset, update asset, lihat session, lihat audit

### Exit criteria
- Halaman load tanpa error runtime.
- Flow UI utama berjalan.
- Lint frontend lulus.

## Fase 5 - Integration Gate (0.5 hari)

### Tujuan
- Menjamin alur end-to-end stabil untuk rilis MVP internal.

### Skenario E2E minimum
1. Login admin
2. Create asset
3. Create account
4. Create session record
5. Validasi audit log muncul
6. Update asset
7. Delete asset

### Verifikasi wajib
- `cd obx_base; npx prisma db seed`
- `cd obx_rest; go build ./...`
- `cd obx_site; npm run lint`
- Smoke test API + UI sesuai skenario

### Exit criteria
- Semua skenario E2E minimum lulus.
- Tidak ada blocker P0/P1.

## Fase 6 - Web Access Gateway (1-1.5 hari)

### Tujuan
- Menyediakan akses SSH, RDP, FTP/file transfer, dan web app proxy langsung dari web board JMS.

### Implementasi
- Backend:
  - POST /rest/pages/JMS/connect/token
  - GET /rest/pages/JMS/connect/ws/:sessionId
  - POST /rest/pages/JMS/file/token
  - GET /rest/pages/JMS/file/ws/:sessionId
  - GET /rest/pages/JMS/proxy/:appId/*path
- Tambahkan validasi RBAC dan short-lived token sebelum membuka koneksi.
- Integrasikan gateway bridge untuk stream SSH/RDP, file transfer browser-based, dan reverse proxy web app.
- Simpan log connect/disconnect ke jms_audit_log.
- Simpan event upload/download ke jms_file_transfer dan jms_session_event.

### Verifikasi wajib
- `cd obx_rest; go build ./...`
- Uji login -> create connect token -> buka web SSH
- Uji login -> create connect token -> buka web RDP
- Uji login -> create file token -> upload/download file dari browser
- Uji login -> akses web app internal melalui endpoint proxy

### Exit criteria
- SSH web terminal terbuka dan interaktif.
- RDP web client terbuka dan stabil.
- Upload/download file dari browser berjalan dan tercatat.
- Web app proxy merespons dengan header keamanan yang sesuai.

## Fase 7 - Hardening & Policy (0.5-1 hari)

### Tujuan
- Mengamankan akses remote agar siap dipakai internal secara terkendali.

### Implementasi
- Tambahkan allowlist host/port per company atau per asset group.
- Batasi TTL token koneksi dan one-time usage.
- Tambahkan idle timeout session dan revoke paksa dari UI admin.
- Tambahkan audit event minimum: connect, disconnect, deny, timeout.
- Tambahkan kebijakan upload/download: ukuran file, ekstensi, dan tujuan direktori.

### Verifikasi wajib
- Uji token expired ditolak.
- Uji user tanpa privilege ditolak.
- Uji revoke session memutus koneksi aktif.
- Uji file yang tidak sesuai policy ditolak.

### Exit criteria
- Kontrol akses tervalidasi.
- Audit trail lengkap untuk kejadian utama koneksi.

## RACI Ringkas

- Owner teknis backend: obx_rest maintainer
- Owner teknis database: obx_base maintainer
- Owner teknis frontend: obx_site maintainer
- QA gate: reviewer lint/build/smoke

## Risiko dan Mitigasi

- Risiko: mismatch schema vs query backend
  - Mitigasi: freeze field + review query sebelum coding
- Risiko: route auth tidak konsisten
  - Mitigasi: checklist middleware per endpoint
- Risiko: UI lebih cepat dari API
  - Mitigasi: kerjakan backend selesai per fase sebelum UI integrasi

## Definition of Done MVP JMS

- Database migrasi dan seed stabil.
- API Asset/Account/Session/Audit berfungsi.
- UI board JMS dapat dipakai operator internal.
- Web SSH/RDP/FTP/Web App Proxy dapat diakses dari board JMS sesuai RBAC.
- Bukti command verifikasi tercatat di audit template.
