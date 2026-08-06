# Autopilot Big Pickle - Progress Tracker

Gunakan file ini untuk mencatat progres eksekusi harian per modul.

## Cara Pakai

1. Tambahkan blok tanggal baru setiap mulai sesi.
2. Isi status modul: `done`, `in-progress`, atau `blocked`.
3. Catat command validasi dan hasil ringkas.
4. Simpan blocker dan next action sebelum tutup sesi.

## Template Harian

```text
Tanggal: YYYY-MM-DD
Model: opencode/big-pickle
Executor: {nama}

Batch dikerjakan:
- F{n}-P{n} {module_code}: {done|in-progress|blocked}

Validation:
- cd obx_rest && go build ./... -> {pass|fail}
- cd obx_site && npx tsc --noEmit -> {pass|fail|n/a}

Ringkasan:
- Scope: ...
- Changes: ...
- Notes: ...

Blocker:
- ...

Next Action:
- ...
```

## Tracking Board

### Fase 1

- [ ] F1-P2 dat_request
- [ ] F1-P2 dat_signature
- [ ] F1-P2 dat_document
- [ ] F1-P2 dat_document_revision
- [ ] F1-P2 dat_document_evidence

### Fase 2

- [ ] F2-P1 jms_asset_group
- [ ] F2-P1 jms_asset
- [ ] F2-P1 jms_asset_group_member
- [ ] F2-P1 jms_account
- [ ] F2-P1 jms_account_secret
- [ ] F2-P2 jms_session
- [ ] F2-P2 jms_session_event
- [ ] F2-P2 jms_session_command
- [ ] F2-P2 jms_audit_log
- [ ] F2-P3 jms_connect_token
- [ ] F2-P3 jms_policy
- [ ] F2-P3 jms_approval
- [ ] F2-P4 ssh_token_endpoint
- [ ] F2-P4 ssh_session_lifecycle
- [ ] F2-P4 ssh_runtime_event
- [ ] F2-P5 jms_file_transfer
- [ ] F2-P5 jms_session_event
- [ ] F2-P6 jms_web_app
- [ ] F2-P6 jms_connect_token
- [ ] F2-P6 proxy_request_access_log
- [ ] F2-P7 rdp_token_endpoint
- [ ] F2-P7 rdp_session_revoke_control
- [ ] F2-P7 rdp_runtime_recovery_note

### Fase 3

- [ ] F3-P1 infra_host
- [ ] F3-P1 infra_stack
- [ ] F3-P1 vm_host
- [ ] F3-P2 web_site
- [ ] F3-P2 web_upstream
- [ ] F3-P2 web_certificate
- [ ] F3-P2 web_config_version
- [ ] F3-P2 web_reload_history
- [ ] F3-P3 net_device
- [ ] F3-P3 net_interface
- [ ] F3-P3 net_poll_sample
- [ ] F3-P4 infra_metric_sample
- [ ] F3-P4 infra_alert_rule
- [ ] F3-P4 net_alert_rule
- [ ] F3-P4 net_backup_job
- [ ] F3-P4 vm_resource_sample
- [ ] F3-P4 vm_permission
- [ ] F3-P4 vm_action_log
- [ ] F3-P5 sec_event
- [ ] F3-P5 sec_event_source
- [ ] F3-P5 sec_event_parser
- [ ] F3-P5 sec_rule
- [ ] F3-P5 sec_alert
- [ ] F3-P5 sec_incident

## Log Harian

### 2026-08-06

Model: opencode/big-pickle  
Executor: Owner

Batch dikerjakan:
- F0-P0 preparation-autopilot: done
- F1-P2 dat_request: in-progress

Validation:
- cd obx_rest && go build ./... -> pending
- cd obx_site && npx tsc --noEmit -> pass

Ringkasan:
- Scope: Menyusun prompt cepat autopilot dan tracker progres harian sampai final flow.
- Changes: Menambahkan prompt per fase/paket/modul dan checklist tracking harian.
- Notes: Backend runtime lokal belum stabil pada sesi ini, tetapi jalur proxy guest sudah merespons dari server frontend yang aktif.

Blocker:
- Backend endpoint lokal pada localhost:36665 belum konsisten tersedia saat diuji langsung.

Next Action:
- Jalankan backend lokal, verifikasi go build, lalu mulai eksekusi modul F1-P2 dat_request sebagai batch pertama implementasi.

### 2026-08-06 (Update 2)

Model: opencode/big-pickle  
Executor: Owner

Batch dikerjakan:
- F0-P0 preparation-autopilot: done
- F1-P2 dat_request: in-progress

Validation:
- cd obx_rest && go build ./... -> pass (EXIT:0)
- cd obx_site && npx tsc --noEmit -> pass
- curl http://localhost:36665/rest/guest/SP00 -> pass (HTTP 200)

Ringkasan:
- Scope: Verifikasi runtime backend dan kesiapan gate teknis sebelum eksekusi modul pertama.
- Changes: Tidak ada perubahan kode modul, hanya pembaruan tracker progres.
- Notes: Kondisi backend pada sesi ini terkonfirmasi reachable untuk endpoint guest.

Blocker:
- Tidak ada blocker aktif untuk memulai batch F1-P2 dat_request.

Next Action:
- Mulai implementasi modul F1-P2 dat_request dan catat output Scope/Changes/Validation/Notes setelah batch selesai.
