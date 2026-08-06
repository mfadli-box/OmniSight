# Backlog Eksekusi Fase 2 - Bastion Access dan Operasi Web-Based

Dokumen ini memecah Fase 2 menjadi urutan kerja yang siap dieksekusi untuk membangun bastion core secara bertahap.

## Tujuan Eksekusi Fase 2

1. Menyediakan akses web-based untuk SSH, RDP, FTP/file transfer, dan WebAppProxy.
2. Menjaga semua sesi dan aktivitas akses tetap teraudit.
3. Menyediakan approval access minimum untuk aset kritikal.
4. Menyiapkan fondasi runtime yang aman sebelum masuk Fase 3.

## Urutan Eksekusi

### Langkah 1 - Finalisasi Asset dan Account Inventory

#### Scope
- jms_asset
- jms_account
- jms_account_secret

#### Aksi
1. Finalisasi struktur asset bastion dan relasi company scope.
2. Finalisasi struktur account dan secret reference.
3. Tentukan status aktif, tipe asset, tipe account, dan metadata minimum.
4. Pastikan data bisa dipakai oleh approval access dan connect token.

#### Output
- Asset dan account inventory bastion siap dipakai.

#### Verifikasi
- Daftar asset dan account bisa dibuat dan dibaca.
- Secret reference tidak diekspos ke output publik.

### Langkah 2 - Aktifkan Session Audit Core

#### Scope
- jms_session
- jms_session_event
- jms_session_command

#### Aksi
1. Definisikan lifecycle session: connect, active, disconnect, revoke.
2. Simpan event sesi dan command log secara terpisah.
3. Pastikan actor user dan asset/account tercatat.
4. Hubungkan session ke audit trail internal.

#### Output
- Audit sesi dasar berjalan.

#### Verifikasi
- Session event dan command log muncul saat simulasi akses.
- Revoke session dapat tercatat.

### Langkah 3 - Bangun Connect Token dan Approval Access

#### Scope
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

#### Verifikasi
- Token tidak bisa dipakai ulang.
- Approval access minimum berhasil diverifikasi.

### Langkah 4 - Implement Web SSH dan Bridge Dasar

#### Scope
- jms_connect_token
- jms_session
- gateway/bridge runtime untuk SSH

#### Aksi
1. Sediakan endpoint pembuatan token untuk SSH.
2. Siapkan bridge koneksi PTY/WebSocket.
3. Pastikan connect/disconnect masuk ke audit.
4. Tambahkan handling timeout dan forced disconnect.

#### Output
- SSH via browser berjalan untuk pilot.

#### Verifikasi
- Web SSH bisa membuka sesi ke host pilot.
- Session audit tercatat penuh.

### Langkah 5 - Implement File Transfer Browser-Based

#### Scope
- jms_file_transfer
- jms_session_event

#### Aksi
1. Sediakan token file transfer terpisah.
2. Implement upload/download stream terkontrol.
3. Catat sumber, tujuan, file name, ukuran, dan status transfer.
4. Pastikan aktivitas transfer ikut masuk audit sesi.

#### Output
- FTP/file transfer browser-based berjalan untuk pilot.

#### Verifikasi
- Upload/download tercatat dengan benar.
- Transfer gagal/berhasil bisa ditelusuri.

### Langkah 6 - Implement WebAppProxy

#### Scope
- jms_web_app
- jms_connect_token

#### Aksi
1. Definisikan daftar aplikasi internal yang boleh diproxy.
2. Implement proxy path-based yang dibatasi allowlist.
3. Simpan access log per request.
4. Pastikan permission per company/role tervalidasi.

#### Output
- Aplikasi internal dapat diakses via board.

#### Verifikasi
- Proxy hanya mengizinkan target yang terdaftar.
- Request log dapat dilacak.

### Langkah 7 - Rapikan RDP Pilot dan Recovery Controls

#### Scope
- bridge runtime untuk RDP
- session audit
- approval access

#### Aksi
1. Siapkan client browser untuk RDP pilot.
2. Tambahkan batasan timeout dan reconnect.
3. Pastikan revoke dapat memutus sesi aktif.
4. Tambahkan recovery note jika gateway gagal.

#### Output
- RDP browser access tersedia untuk pilot terbatas.

#### Verifikasi
- RDP session dapat dicatat dan diputus saat revoke.

## Dependensi Fase 2

| Area | Dependensi Minimum |
|---|---|
| Inventory bastion | Fase 1 stabil, asset dasar tersedia |
| Session audit | user/company/privilege stabil |
| Connect token | policy dan approval workflow aktif |
| Web SSH / RDP / FTP | gateway runtime dan session event siap |
| WebAppProxy | daftar app internal dan allowlist tersedia |

## DoD Fase 2

- Asset dan account inventory bastion aktif.
- Session audit dan revoke berjalan.
- Connect token sekali pakai berfungsi.
- Web SSH, FTP/file transfer, dan WebAppProxy aktif untuk pilot.
- RDP pilot berjalan atau minimal siap uji terbatas.
- Dokumentasi runtime bastion diperbarui.

## Risiko Utama

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
- backend_module_matrix_platform_replacement.md
- feature_action_matrix_platform_replacement.md
- checklist_platform_replacement_execution.md
