# Matriks Modul Backend - Platform Replacement OmniSight

Dokumen ini memetakan fitur prioritas menjadi modul backend yang harus dibangun, termasuk urutan eksekusi, dependensi, dan keluaran yang diharapkan.

## Tujuan

Memberikan jalur implementasi backend yang jelas untuk domain:
- SOP / ISO governance
- Bastion / JumpServer
- Docker / Beszel
- Nginx / nginx-ui
- MikroTik / LibreNMS
- VM monitoring
- SIEM / Wazuh

## Prinsip Implementasi

1. Mulai dari modul yang memberi nilai governance paling cepat.
2. Bangun backend yang bisa diaudit sebelum memperluas UI.
3. Prioritaskan data model, API, lalu workflow.
4. Simpan SIEM sebagai fase terakhir karena paling berat.

## 1. SOP / ISO Governance

| Modul Backend | Fungsi | Dependensi | Output |
|---|---|---|---|
| dat_document | controlled document master | dat_user, dat_company | dokumen terkendali |
| dat_signature | approval signature type | dat_user, dat_request | konfigurasi approval |
| dat_request | workflow request | dat_user, dat_company | request approval berjalan |
| dat_document_revision | histori revisi | dat_document | versioning dokumen |
| dat_document_evidence | bukti audit | dat_document, dat_request | evidence trail |

### Urutan Eksekusi

1. dat_request.
2. dat_signature.
3. dat_document.
4. revision dan evidence.

### Catatan Implementasi

- Gunakan pola approval bertingkat.
- Simpan status request dan actor audit.
- Pastikan document owner, reviewer, approver, dan archive owner bisa dilacak.

## 2. Bastion / JumpServer

| Modul Backend | Fungsi | Dependensi | Output |
|---|---|---|---|
| jms_asset | inventory asset | dat_company | daftar host/target |
| jms_account | inventory account | jms_asset, dat_company | akun terkelola |
| jms_account_secret | vault/secret ref | jms_account | secret aman |
| jms_policy | access approval policy | dat_request, dat_signature | policy akses |
| jms_connect_token | token koneksi sementara | jms_asset, jms_account, dat_user | access token |
| jms_session | session master | jms_asset, jms_account, dat_user | catatan sesi |
| jms_session_event | event sesi | jms_session | audit sesi |
| jms_session_command | command log | jms_session | histori command |
| jms_file_transfer | transfer audit | jms_session | audit file |
| jms_web_app | proxy target web app | dat_company | daftar app proxy |

### Urutan Eksekusi

1. Asset dan account.
2. Session dan session event.
3. Connect token.
4. File transfer.
5. WebAppProxy.

### Catatan Implementasi

- Session token harus short-lived dan one-time.
- Semua akses connect/disconnect wajib masuk audit log.
- File transfer perlu event yang terpisah dari session connect.

## 3. Docker / Dockge - Beszel

| Modul Backend | Fungsi | Dependensi | Output |
|---|---|---|---|
| infra_host | inventory host | dat_company | host terdaftar |
| infra_stack | compose/project inventory | infra_host | stack terkelola |
| infra_deploy_history | riwayat deploy | infra_stack, dat_user | audit deploy |
| infra_metric_sample | metrics host/container | infra_host | metrik berkala |
| infra_alert_rule | alert threshold | infra_host | rule monitoring |

### Urutan Eksekusi

1. infra_host.
2. infra_stack.
3. infra_deploy_history.
4. infra_metric_sample.
5. infra_alert_rule.

## 4. Nginx / nginx-ui

| Modul Backend | Fungsi | Dependensi | Output |
|---|---|---|---|
| web_site | site inventory | dat_company | virtual host |
| web_upstream | upstream inventory | dat_company | upstream target |
| web_certificate | SSL inventory | dat_company | sertifikat terkelola |
| web_config_version | config versioning | web_site, web_upstream | rollback config |
| web_reload_history | reload audit | web_config_version, dat_user | histori reload |

### Urutan Eksekusi

1. web_site dan web_upstream.
2. web_certificate.
3. web_config_version.
4. web_reload_history.

## 5. MikroTik / LibreNMS

| Modul Backend | Fungsi | Dependensi | Output |
|---|---|---|---|
| net_device | inventory perangkat | dat_company | daftar device |
| net_interface | inventory interface | net_device | port/interface |
| net_poll_sample | hasil polling | net_device | data health |
| net_alert_rule | threshold alert | net_device | alert resource |
| net_backup_job | backup config | net_device, dat_user | arsip backup |

### Urutan Eksekusi

1. net_device.
2. net_interface.
3. net_poll_sample.
4. net_alert_rule.
5. net_backup_job.

## 6. VM Monitoring / Management

| Modul Backend | Fungsi | Dependensi | Output |
|---|---|---|---|
| vm_host | inventory host/cluster | dat_company | VM inventory |
| vm_resource_sample | resource metrics | vm_host | status kesehatan |
| vm_permission | akses grup | vm_host, dat_user | permission model |
| vm_action_log | histori tindakan | vm_host, dat_user | audit perubahan |

### Urutan Eksekusi

1. vm_host.
2. vm_resource_sample.
3. vm_permission.
4. vm_action_log.

## 7. SIEM / Wazuh

| Modul Backend | Fungsi | Dependensi | Output |
|---|---|---|---|
| sec_event | event master | dat_company | event terpusat |
| sec_event_source | sumber event | sec_event | klasifikasi sumber |
| sec_event_parser | parsing/normalisasi | sec_event_source | event standar |
| sec_rule | correlation rule | sec_event_parser | rule correlation |
| sec_alert | alert hasil rule | sec_rule | alert keamanan |
| sec_incident | incident workflow | sec_alert, dat_user | timeline insiden |

### Urutan Eksekusi

1. sec_event dan source.
2. sec_event_parser.
3. sec_rule.
4. sec_alert.
5. sec_incident.

## Prioritas Implementasi Backend

| Prioritas | Modul Utama | Alasan |
|---|---|---|
| 1 | dat_request, dat_signature, dat_document | quick win governance |
| 2 | jms_asset, jms_account, jms_session | bastion core |
| 3 | jms_connect_token, jms_file_transfer, jms_web_app | akses web-based |
| 4 | infra_host, web_site, vm_host | inventory infra |
| 5 | net_device | network monitoring |
| 6 | sec_event | SIEM pipeline |

## Keluaran Minimum per Tahap

### Tahap Governance
- Approval workflow aktif.
- Dokumen terkendali bisa dibuat dan direvisi.

### Tahap Bastion
- Asset, account, session, dan file transfer tercatat.

### Tahap Infra Inventory
- Host, site, dan VM dapat diinventaris.

### Tahap Monitoring
- Metric sample dan alert rule aktif.

### Tahap SIEM
- Event ingest dan correlation rule berjalan.

## Kesimpulan

Backend replacement paling realistis dimulai dari governance dan bastion core. Setelah itu baru perluas ke inventory infra, monitoring, lalu SIEM. Urutan ini menjaga agar setiap tahap punya nilai mandiri dan tetap bisa diaudit.
