# Autopilot Big Pickle

Prompt cepat ini dipakai untuk eksekusi autopilot bertahap dengan model `opencode/big-pickle`.

Progress harian: lihat `autopilot-big-pickle-progress.md` pada folder yang sama.
Command Opencode CLI: lihat `opencode-cli-big-pickle-commands.md` pada folder yang sama.

## Tujuan
- Menjalankan task per modul dengan mode small-batch.
- Menjaga perubahan tetap kecil, terukur, dan mudah diverifikasi.
- Memastikan output konsisten dengan format parity Opencode.

## Prompt Template Utama

```text
Task: Eksekusi autopilot {fase} paket {paket} modul {module_code}
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
- stop if build or migration fails
Output:
- Scope
- Changes
- Validation
- Notes
```

## Prompt Cepat Fase 1 Paket 2

### 1) dat_request

```text
Task: Eksekusi autopilot Fase 1 Paket 2 modul dat_request
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) dat_signature

```text
Task: Eksekusi autopilot Fase 1 Paket 2 modul dat_signature
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) dat_document

```text
Task: Eksekusi autopilot Fase 1 Paket 2 modul dat_document
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 4) dat_document_revision

```text
Task: Eksekusi autopilot Fase 1 Paket 2 modul dat_document_revision
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 5) dat_document_evidence

```text
Task: Eksekusi autopilot Fase 1 Paket 2 modul dat_document_evidence
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Validasi Minimum

```text
cd obx_rest
go build ./...
```

Jika ada perubahan frontend:

```text
cd obx_site
npx tsc --noEmit
```

## Prompt Cepat Fase 2 Paket 1

### 1) jms_asset_group

```text
Task: Eksekusi autopilot Fase 2 Paket 1 modul jms_asset_group
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) jms_asset

```text
Task: Eksekusi autopilot Fase 2 Paket 1 modul jms_asset
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) jms_asset_group_member

```text
Task: Eksekusi autopilot Fase 2 Paket 1 modul jms_asset_group_member
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 4) jms_account

```text
Task: Eksekusi autopilot Fase 2 Paket 1 modul jms_account
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 5) jms_account_secret

```text
Task: Eksekusi autopilot Fase 2 Paket 1 modul jms_account_secret
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Prompt Cepat Fase 3 Paket 1

### 1) infra_host

```text
Task: Eksekusi autopilot Fase 3 Paket 1 modul infra_host
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) infra_stack

```text
Task: Eksekusi autopilot Fase 3 Paket 1 modul infra_stack
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) vm_host

```text
Task: Eksekusi autopilot Fase 3 Paket 1 modul vm_host
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Prompt Cepat Fase 2 Paket 2

### 1) jms_session

```text
Task: Eksekusi autopilot Fase 2 Paket 2 modul jms_session
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) jms_session_event

```text
Task: Eksekusi autopilot Fase 2 Paket 2 modul jms_session_event
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) jms_session_command

```text
Task: Eksekusi autopilot Fase 2 Paket 2 modul jms_session_command
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 4) jms_audit_log

```text
Task: Eksekusi autopilot Fase 2 Paket 2 modul jms_audit_log
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Prompt Cepat Fase 2 Paket 3

### 1) jms_connect_token

```text
Task: Eksekusi autopilot Fase 2 Paket 3 modul jms_connect_token
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) jms_policy

```text
Task: Eksekusi autopilot Fase 2 Paket 3 modul jms_policy
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) jms_approval

```text
Task: Eksekusi autopilot Fase 2 Paket 3 modul jms_approval
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Prompt Cepat Fase 2 Paket 4

### 1) ssh_token_endpoint

```text
Task: Eksekusi autopilot Fase 2 Paket 4 modul ssh_token_endpoint
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) ssh_session_lifecycle

```text
Task: Eksekusi autopilot Fase 2 Paket 4 modul ssh_session_lifecycle
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) ssh_runtime_event

```text
Task: Eksekusi autopilot Fase 2 Paket 4 modul ssh_runtime_event
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Prompt Cepat Fase 2 Paket 5

### 1) jms_file_transfer

```text
Task: Eksekusi autopilot Fase 2 Paket 5 modul jms_file_transfer
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) jms_session_event

```text
Task: Eksekusi autopilot Fase 2 Paket 5 modul jms_session_event
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Prompt Cepat Fase 2 Paket 6

### 1) jms_web_app

```text
Task: Eksekusi autopilot Fase 2 Paket 6 modul jms_web_app
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) jms_connect_token

```text
Task: Eksekusi autopilot Fase 2 Paket 6 modul jms_connect_token
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) proxy_request_access_log

```text
Task: Eksekusi autopilot Fase 2 Paket 6 modul proxy_request_access_log
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Prompt Cepat Fase 2 Paket 7

### 1) rdp_token_endpoint

```text
Task: Eksekusi autopilot Fase 2 Paket 7 modul rdp_token_endpoint
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) rdp_session_revoke_control

```text
Task: Eksekusi autopilot Fase 2 Paket 7 modul rdp_session_revoke_control
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) rdp_runtime_recovery_note

```text
Task: Eksekusi autopilot Fase 2 Paket 7 modul rdp_runtime_recovery_note
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Prompt Cepat Fase 3 Paket 2

### 1) web_site

```text
Task: Eksekusi autopilot Fase 3 Paket 2 modul web_site
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) web_upstream

```text
Task: Eksekusi autopilot Fase 3 Paket 2 modul web_upstream
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) web_certificate

```text
Task: Eksekusi autopilot Fase 3 Paket 2 modul web_certificate
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 4) web_config_version

```text
Task: Eksekusi autopilot Fase 3 Paket 2 modul web_config_version
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 5) web_reload_history

```text
Task: Eksekusi autopilot Fase 3 Paket 2 modul web_reload_history
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Prompt Cepat Fase 3 Paket 3

### 1) net_device

```text
Task: Eksekusi autopilot Fase 3 Paket 3 modul net_device
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) net_interface

```text
Task: Eksekusi autopilot Fase 3 Paket 3 modul net_interface
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) net_poll_sample

```text
Task: Eksekusi autopilot Fase 3 Paket 3 modul net_poll_sample
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Prompt Cepat Fase 3 Paket 4

### 1) infra_metric_sample

```text
Task: Eksekusi autopilot Fase 3 Paket 4 modul infra_metric_sample
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) infra_alert_rule

```text
Task: Eksekusi autopilot Fase 3 Paket 4 modul infra_alert_rule
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) net_alert_rule

```text
Task: Eksekusi autopilot Fase 3 Paket 4 modul net_alert_rule
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 4) net_backup_job

```text
Task: Eksekusi autopilot Fase 3 Paket 4 modul net_backup_job
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 5) vm_resource_sample

```text
Task: Eksekusi autopilot Fase 3 Paket 4 modul vm_resource_sample
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 6) vm_permission

```text
Task: Eksekusi autopilot Fase 3 Paket 4 modul vm_permission
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 7) vm_action_log

```text
Task: Eksekusi autopilot Fase 3 Paket 4 modul vm_action_log
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Prompt Cepat Fase 3 Paket 5

### 1) sec_event

```text
Task: Eksekusi autopilot Fase 3 Paket 5 modul sec_event
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 2) sec_event_source

```text
Task: Eksekusi autopilot Fase 3 Paket 5 modul sec_event_source
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 3) sec_event_parser

```text
Task: Eksekusi autopilot Fase 3 Paket 5 modul sec_event_parser
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 4) sec_rule

```text
Task: Eksekusi autopilot Fase 3 Paket 5 modul sec_rule
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 5) sec_alert

```text
Task: Eksekusi autopilot Fase 3 Paket 5 modul sec_alert
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

### 6) sec_incident

```text
Task: Eksekusi autopilot Fase 3 Paket 5 modul sec_incident
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
Output: Scope, Changes, Validation, Notes
```

## Urutan Eksekusi Sampai Final

1. Fase 1 Paket 2: dat_request -> dat_signature -> dat_document -> dat_document_revision -> dat_document_evidence
2. Fase 2 Paket 1: jms_asset_group -> jms_asset -> jms_asset_group_member -> jms_account -> jms_account_secret
3. Fase 2 Paket 2: jms_session -> jms_session_event -> jms_session_command -> jms_audit_log
4. Fase 2 Paket 3: jms_connect_token -> jms_policy -> jms_approval
5. Fase 2 Paket 4: ssh_token_endpoint -> ssh_session_lifecycle -> ssh_runtime_event
6. Fase 2 Paket 5: jms_file_transfer -> jms_session_event
7. Fase 2 Paket 6: jms_web_app -> jms_connect_token -> proxy_request_access_log
8. Fase 2 Paket 7: rdp_token_endpoint -> rdp_session_revoke_control -> rdp_runtime_recovery_note
9. Fase 3 Paket 1: infra_host -> infra_stack -> vm_host
10. Fase 3 Paket 2: web_site -> web_upstream -> web_certificate -> web_config_version -> web_reload_history
11. Fase 3 Paket 3: net_device -> net_interface -> net_poll_sample
12. Fase 3 Paket 4: infra_metric_sample -> infra_alert_rule -> net_alert_rule -> net_backup_job -> vm_resource_sample -> vm_permission -> vm_action_log
13. Fase 3 Paket 5: sec_event -> sec_event_source -> sec_event_parser -> sec_rule -> sec_alert -> sec_incident

## Gate Final Minimum

```text
cd obx_rest
go build ./...

cd obx_site
npx tsc --noEmit
```

Checklist lulus final:
- Semua modul pada urutan eksekusi sudah selesai per small-batch.
- Tidak ada bloker build backend.
- Tidak ada error typecheck frontend.
- Output tiap batch tersedia: Scope, Changes, Validation, Notes.

## One-Liner Super Ringkas

Gunakan pola ini untuk semua modul:

```text
Task: Eksekusi autopilot {fase}-{paket} modul {module_code}; Mode: backend; Scope: satu modul small-batch; Model: opencode/big-pickle; Constraints: keep changes minimal, follow existing patterns, validate after changes; Output: Scope, Changes, Validation, Notes.
```

### Fase 1 Paket 2

```text
Task: Eksekusi autopilot F1-P2 modul dat_request; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F1-P2 modul dat_signature; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F1-P2 modul dat_document; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F1-P2 modul dat_document_revision; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F1-P2 modul dat_document_evidence; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 2 Paket 1

```text
Task: Eksekusi autopilot F2-P1 modul jms_asset_group; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P1 modul jms_asset; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P1 modul jms_asset_group_member; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P1 modul jms_account; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P1 modul jms_account_secret; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 2 Paket 2

```text
Task: Eksekusi autopilot F2-P2 modul jms_session; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.

## Checklist Eksekusi Harian

Gunakan checklist ini untuk tracking cepat. Ubah [ ] menjadi [x] saat modul selesai dan lolos validasi minimum.

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

### Rekap Status Cepat

- [ ] Build backend lulus: `cd obx_rest && go build ./...`
- [ ] Typecheck frontend lulus: `cd obx_site && npx tsc --noEmit`
- [ ] Semua batch punya output: Scope, Changes, Validation, Notes
Task: Eksekusi autopilot F2-P2 modul jms_session_event; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P2 modul jms_session_command; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P2 modul jms_audit_log; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 2 Paket 3

```text
Task: Eksekusi autopilot F2-P3 modul jms_connect_token; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P3 modul jms_policy; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P3 modul jms_approval; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 2 Paket 4

```text
Task: Eksekusi autopilot F2-P4 modul ssh_token_endpoint; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P4 modul ssh_session_lifecycle; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P4 modul ssh_runtime_event; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 2 Paket 5

```text
Task: Eksekusi autopilot F2-P5 modul jms_file_transfer; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P5 modul jms_session_event; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 2 Paket 6

```text
Task: Eksekusi autopilot F2-P6 modul jms_web_app; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P6 modul jms_connect_token; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P6 modul proxy_request_access_log; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 2 Paket 7

```text
Task: Eksekusi autopilot F2-P7 modul rdp_token_endpoint; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P7 modul rdp_session_revoke_control; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F2-P7 modul rdp_runtime_recovery_note; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 3 Paket 1

```text
Task: Eksekusi autopilot F3-P1 modul infra_host; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P1 modul infra_stack; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P1 modul vm_host; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 3 Paket 2

```text
Task: Eksekusi autopilot F3-P2 modul web_site; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P2 modul web_upstream; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P2 modul web_certificate; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P2 modul web_config_version; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P2 modul web_reload_history; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 3 Paket 3

```text
Task: Eksekusi autopilot F3-P3 modul net_device; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P3 modul net_interface; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P3 modul net_poll_sample; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 3 Paket 4

```text
Task: Eksekusi autopilot F3-P4 modul infra_metric_sample; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P4 modul infra_alert_rule; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P4 modul net_alert_rule; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P4 modul net_backup_job; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P4 modul vm_resource_sample; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P4 modul vm_permission; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P4 modul vm_action_log; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```

### Fase 3 Paket 5

```text
Task: Eksekusi autopilot F3-P5 modul sec_event; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P5 modul sec_event_source; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P5 modul sec_event_parser; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P5 modul sec_rule; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P5 modul sec_alert; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
Task: Eksekusi autopilot F3-P5 modul sec_incident; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes.
```
