# Backlog Backend Fase 2 - Bastion Core

Dokumen ini memecah Fase 2 menjadi urutan kerja backend yang paling siap dieksekusi untuk membangun bastion core secara bertahap.

## Tujuan

1. Menyediakan model data bastion yang lengkap untuk asset, account, session, dan proxy target.
2. Menyediakan token akses sementara yang aman dan mudah diaudit.
3. Menyediakan session audit, file transfer, dan WebAppProxy yang konsisten.
4. Menjaga perubahan kecil, terukur, dan dapat diuji per paket.

## Urutan Backend

### Paket 1 - Asset dan Account Inventory

#### Modul
- jms_asset
- jms_account
- jms_account_secret

#### Aksi
1. Finalisasi schema asset bastion dan company scope.
2. Finalisasi schema account dan secret reference.
3. Tetapkan field minimum untuk status, tipe, dan metadata dasar.
4. Siapkan data pilot untuk host dan account referensi.

#### Output
- Asset dan account inventory bastion siap dipakai.

### Paket 2 - Session Core dan Audit

#### Modul
- jms_session
- jms_session_event
- jms_session_command

#### Aksi
1. Definisikan lifecycle session: connect, active, disconnect, revoke.
2. Simpan event sesi dan command log secara terpisah.
3. Pastikan actor user, asset, dan account tercatat.
4. Hubungkan session ke audit trail internal.

#### Output
- Audit sesi dasar berjalan.

### Paket 3 - Approval Access dan Connect Token

#### Modul
- jms_connect_token
- jms_policy
- jms_approval

#### Aksi
1. Buat token koneksi sementara yang short-lived dan one-time.
2. Kaitkan policy akses dengan approval workflow.
3. Pastikan token hanya valid untuk user, asset, dan account yang sesuai.
4. Simpan status token: pending, used, expired, revoked.

#### Output
- Mekanisme akses terotorisasi siap digunakan.

### Paket 4 - Web SSH Bridge

#### Modul
- jms_connect_token
- jms_session
- bridge runtime untuk SSH

#### Aksi
1. Sediakan endpoint token untuk SSH.
2. Siapkan bridge koneksi PTY/WebSocket.
3. Pastikan connect/disconnect masuk ke audit.
4. Tambahkan timeout dan forced disconnect.

#### Output
- SSH via browser berjalan untuk pilot.

### Paket 5 - File Transfer Browser-Based

#### Modul
- jms_file_transfer
- jms_session_event

#### Aksi
1. Sediakan token file transfer terpisah.
2. Implement upload/download stream terkontrol.
3. Catat sumber, tujuan, nama file, ukuran, dan status transfer.
4. Pastikan aktivitas transfer ikut masuk audit sesi.

#### Output
- FTP/file transfer browser-based berjalan untuk pilot.

### Paket 6 - WebAppProxy

#### Modul
- jms_web_app
- jms_connect_token

#### Aksi
1. Definisikan daftar aplikasi internal yang boleh diproxy.
2. Implement proxy path-based yang dibatasi allowlist.
3. Simpan access log per request.
4. Pastikan permission per company dan role tervalidasi.

#### Output
- Aplikasi internal dapat diakses via board.

### Paket 7 - RDP Pilot dan Recovery Controls

#### Modul
- bridge runtime untuk RDP
- jms_session
- jms_connect_token

#### Aksi
1. Siapkan client browser untuk RDP pilot.
2. Tambahkan batasan timeout dan reconnect.
3. Pastikan revoke dapat memutus sesi aktif.
4. Tambahkan recovery note jika gateway gagal.

#### Output
- RDP browser access tersedia untuk pilot terbatas.

## Dependensi

| Area | Dependensi Minimum |
|---|---|
| Asset dan account | Fase 1 stabil, asset dasar tersedia |
| Session audit | user/company/privilege stabil |
| Connect token | policy dan approval workflow aktif |
| Web SSH / RDP / FTP | gateway runtime dan session event siap |
| WebAppProxy | daftar app internal dan allowlist tersedia |

## DoD Backend Fase 2

- Asset dan account inventory bastion aktif.
- Session audit dan revoke berjalan.
- Connect token sekali pakai berfungsi.
- Web SSH, FTP/file transfer, dan WebAppProxy aktif untuk pilot.
- RDP pilot berjalan atau minimal siap uji terbatas.
- Dokumentasi runtime bastion diperbarui.

## Risiko

1. Token koneksi tidak benar-benar short-lived.
2. Audit sesi tidak konsisten dengan event runtime.
3. WebAppProxy membuka target yang tidak diinginkan.
4. RDP/gateway kompleks dan memperlambat pilot.

## Mitigasi

1. Wajibkan expiry dan one-time use pada token.
2. Simpan audit event di jalur yang sama dengan connect/disconnect.
3. Terapkan allowlist host/path untuk proxy.
4. Tunda perluasan fitur RDP sampai SSH, FTP, dan proxy stabil.

## Referensi

- roadmap_platform_3_phase.md
- phase2_execution_backlog.md
- backend_module_matrix_platform_replacement.md
- ai_runbook.md
