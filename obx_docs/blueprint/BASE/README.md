# Rancangan Submodul: BASE - obx_base Prisma Schema

## Ringkasan

- Scope: katalog database untuk obx_base
- Tujuan: referensi teknis tunggal untuk Prisma model, enum, relasi, dan konvensi schema
- Sumber utama: `obx_base/prisma/schema/*.prisma`

## Struktur Dokumen

- Summary teknis: file ini
- Detail page teknis: `obx_docs/blueprint/BASE/*.md`
- Summary user guide: `obx_docs/guide/BASE.md`
- Detail page user guide: `obx_docs/guide/BASE/*.md`

## Domain Ringkas

- dat: company, module, user
- jsm_stack: bastion, session access, web proxy, file transfer
- dat_signature: approval flow
- doc: document & request
- ict_machine: VM & Docker
- ict_mikrotik: network device
- ict_security: firewall, incident, vuln, FIM
- ict_website: Nginx, SSL, uptime

## Prisma Summary

- Generator: prisma-client-js
- Datasource: postgresql
- ID utama: String + uuid()
- Audit field umum: created_at, updated_at
- Tenancy umum: company_id

## Page Index

| Page | Detail Teknis | Detail Guide |
|---|---|---|
| all_enum_dat | all_enum_dat.md | all_enum_dat.md |
| dat_company | dat_company.md | dat_company.md |
| dat_module | dat_module.md | dat_module.md |
| dat_user | dat_user.md | dat_user.md |
| jsm_stack | jsm_stack.md | jsm_stack.md |
| dat_signature | dat_signature.md | dat_signature.md |
| dat_document | dat_document.md | dat_document.md |
| ict_machine | ict_machine.md | ict_machine.md |
| ict_monitoring | ict_monitoring.md | ict_monitoring.md |
| ict_mikrotik | ict_mikrotik.md | ict_mikrotik.md |
| ict_security | ict_security.md | ict_security.md |
| ict_website | ict_website.md | ict_website.md |

## Checklist

- [ ] Enum dan model terpetakan per schema file
- [ ] Relasi penting terdokumentasi
- [ ] Checklist migrasi dan generate tersedia

## Referensi

- `obx_base/prisma/schema/schema.prisma`
- `obx_base/prisma/schema/*.prisma`

### Domain: dat (Company, Module, User)

#### Tabel Inti dan Relasi

| Tabel | FK | Tabel Target | Keterangan |
|---|---|---|---|
| dat_company_module | company_id | dat_company | Modul yang diaktifkan per company |
| dat_company_module | module_id | dat_module | Referensi modul yang tersedia |
| dat_company_area | company_id | dat_company | Area/divisi fisik per company |
| dat_user_company | user_id | dat_user | Relasi user ke company |
| dat_user_company | company_id | dat_company | Relasi user ke company |
| dat_user_area | user_id | dat_user | Relasi user ke area |
| dat_user_area | company_area_id | dat_company_area | Referensi area |
| dat_user_privilege | user_company_id | dat_user_company | Scope privilege per user-company |
| dat_user_privilege | module_id | dat_module | Modul yang diberi akses |
| dat_user_session | user_id | dat_user | Sesi aktif per user |
| dat_module | parent_id | dat_module | Self-referential untuk hierarki menu |

#### Alur Akses (Access Chain)

```
dat_company
  â””â”€â”€ dat_company_module â”€â”€> dat_module (modul yang diaktifkan)
  â””â”€â”€ dat_company_area
        â””â”€â”€ dat_user_area â”€â”€> dat_user

dat_user
  â””â”€â”€ dat_user_company â”€â”€> dat_company
        â””â”€â”€ dat_user_privilege â”€â”€> dat_module (level: HIDE/VIEW/BOOK/POST)
  â””â”€â”€ dat_user_session (token, ip, expires)
  â””â”€â”€ dat_user_action (audit log)
```

#### Composite Unique

| Tabel | Constraint | Scope |
|---|---|---|
| dat_company_module | [company_id, module_id] | 1 modul per company |
| dat_company_area | [company_id, code] | kode area unik per company |
| dat_user_company | [user_id, company_id] | 1 pasang user-company |
| dat_user_area | [user_id, company_area_id] | 1 pasang user-area |
| dat_user_privilege | [user_company_id, module_id] | 1 privilege per user-company-modul |

---

### Domain: dat_signature (Approval Flow)

#### Tabel Inti dan Relasi

| Tabel | FK | Tabel Target | Keterangan |
|---|---|---|---|
| dat_approval_step | type_id | dat_signature_type | Step milik signature type tertentu |
| dat_approval_sign | step_id | dat_approval_step | Sign milik step tertentu |
| dat_approval_sign | user_id | dat_user | Approver per step |
| dat_signature_flag | form_id | dat_signature_form | Flag milik form tertentu |
| dat_signature_flag | user_id | dat_user | User yang memberikan flag |

#### Alur Approval (Approval Chain)

```
dat_signature_type (definisi flow)
  â””â”€â”€ dat_approval_step [step, condition: ALL/ANY_APPROVED]
        â””â”€â”€ dat_approval_sign â”€â”€> dat_user (siapa yang boleh approve di step ini)

dat_signature_form (instance form di satu request)
  â””â”€â”€ dat_signature_flag [status: approval_info] â”€â”€> dat_user (aksi approval aktual)
```

#### Enum yang Dipakai

- approval_flag: ALL_APPROVED (semua harus approve), ANY_APPROVED (cukup satu)
- approval_info: PENDING, IN_PROGRESS, APPROVED, REJECTED, COMPLETED

---

### Domain: doc (Document & Request)

#### Tabel Inti dan Relasi

| Tabel | FK | Tabel Target | Keterangan |
|---|---|---|---|
| dat_document | category_id | dat_document_category | Kategori opsional |
| dat_document_version | document_id | dat_document | Versi historis dokumen |
| dat_document_approval | document_id | dat_document | Log aksi approval dokumen |
| dat_request | company_id | dat_company | Request scoped per company |
| dat_request | type_id | dat_signature_type | Flow approval yang dipakai |
| dat_request | requester_id | dat_user | Pembuat request |
| dat_request | completed_by | dat_user | Penyelesai request |

#### Alur Dokumen

```
dat_document_category
  â””â”€â”€ dat_document [status: DRAFT/...]
        â””â”€â”€ dat_document_version (riwayat)
        â””â”€â”€ dat_document_approval (log aksi)

dat_request â”€â”€> dat_signature_type (flow approval)
dat_request â”€â”€> dat_user (requester, completer)
```

#### Composite Unique

| Tabel | Constraint | Scope |
|---|---|---|
| dat_document_category | [company_id, code] | kode kategori unik per company |
| dat_document | [company_id, code] | kode dokumen unik per company |
| dat_document_version | [document_id, version] | versi unik per dokumen |
| dat_request | [company_id, code] | kode request unik per company |

---

### Domain: ict_machine (VM & Docker)

#### Tabel Inti dan Relasi (ID-based, tidak ada FK Prisma eksplisit)

| Tabel | Referensi | Keterangan |
|---|---|---|
| ict_vm_host | company_id | Inventori host VM |
| ict_vm_stat | host_id (ict_vm_host) | Time-series stat host, index host_id+recorded_at |
| ict_docker_compose | host_id (ict_vm_host) | Compose file per host |
| ict_docker_container | host_id (ict_vm_host) | Definisi container per host |
| ict_docker_container_stat | host_id, container_name | Time-series stat container |
| ict_docker_image | host_id (ict_vm_host) | Image registry lokal per host |
| ict_docker_network | host_id (ict_vm_host) | Docker network per host |
| ict_docker_network_member | network_id, container_id | Mapping container ke network |
| ict_docker_backup | host_id (ict_vm_host) | Backup artifact per host |
| ict_docker_deploy | host_id (ict_vm_host) | Deploy job log per host |
| ict_host_group | company_id | Group inventori host |
| ict_host_group_member | group_id, host_id | Keanggotaan host dalam grup |
| ict_host_permission | user_id, ref_type, ref_id | RBAC akses ke host/group |
| ict_alert_rule | company_id, host_id? | Aturan alert monitoring |
| ict_alert_notif | company_id | Channel notifikasi |
| ict_alert_history | alert_rule_id, notification_id | Histori pengiriman alert |
| ict_update_job | company_id, host_id | Jadwal update otomatis |
| ict_update_history | host_id, job_id | Histori eksekusi update |
| ict_update_package | history_id | Paket yang diupdate per histori |
| ict_git_webhook | company_id | Webhook repo git |
| ict_git_deploy_mapping | webhook_id | Pemetaan webhook ke target deploy |
| ict_git_deploy_log | webhook_id, mapping_id | Log event push dan deploy |

#### Alur Deployment

```
ict_vm_host (inventori host)
  â””â”€â”€ ict_docker_compose / ict_docker_container (definisi workload)
  â””â”€â”€ ict_docker_deploy (job deploy/action)
  â””â”€â”€ ict_vm_stat / ict_docker_container_stat (monitoring time-series)

ict_git_webhook
  â””â”€â”€ ict_git_deploy_mapping â”€â”€> ref_id (compose/container)
  â””â”€â”€ ict_git_deploy_log (event + hasil)
```

---

### Domain: ict_mikrotik (Network Device)

#### Tabel Inti dan Relasi

| Tabel | Referensi | Keterangan |
|---|---|---|
| ict_mikrotik_device | company_id | Inventori perangkat MikroTik |
| ict_mikrotik_status | device_id | Time-series resource stat |
| ict_mikrotik_interface | device_id | Snapshot interface per device |
| ict_mikrotik_firewall_rule | device_id | Rule firewall yang disync |
| ict_mikrotik_address_list | device_id | Address list per device |
| ict_mikrotik_address_entry | address_list_id | Entry dalam address list |
| ict_mikrotik_backup | device_id | File backup config device |
| ict_mikrotik_backup_file | backup_id | Konten file backup |
| ict_mikrotik_log | device_id | Syslog dari device |

#### Alur Sinkronisasi

```
ict_mikrotik_device
  â””â”€â”€ ict_mikrotik_status (poll resource)
  â””â”€â”€ ict_mikrotik_interface (sync interface)
  â””â”€â”€ ict_mikrotik_firewall_rule (sync rules)
  â””â”€â”€ ict_mikrotik_address_list
        â””â”€â”€ ict_mikrotik_address_entry
  â””â”€â”€ ict_mikrotik_backup
        â””â”€â”€ ict_mikrotik_backup_file
  â””â”€â”€ ict_mikrotik_log
```

---

### Domain: ict_security (Firewall, Incident, Vuln, FIM)

#### Tabel Inti dan Relasi

| Tabel | Referensi | Keterangan |
|---|---|---|
| ict_firewall_rule | company_id, host_id | Rule nftables per host |
| ict_firewall_zone | company_id, host_id | Zone firewall per host |
| ict_firewall_zone_rule | zone_id, rule_id | Mapping rule ke zone |
| ict_incident | company_id | Inventori insiden keamanan |
| ict_incident_evidence | incident_id | Bukti/file per insiden |
| ict_incident_timeline | incident_id | Audit trail aksi per insiden |
| ict_incident_ioc | incident_id | IOC (indicator of compromise) |
| ict_vuln_scan_schedule | company_id, host_id | Jadwal scan kerentanan |
| ict_vuln_scan | schedule_id, host_id | Hasil satu sesi scan |
| ict_vuln_finding | scan_id | Temuan per scan |
| ict_fim_path | company_id, host_id | Path yang dipantau FIM |
| ict_fim_snapshot | path_id | Snapshot kondisi file |
| ict_fim_alert | path_id, host_id | Alert perubahan file |

#### Alur Incident Response

```
ict_incident
  â””â”€â”€ ict_incident_evidence
  â””â”€â”€ ict_incident_timeline
  â””â”€â”€ ict_incident_ioc

ict_vuln_scan_schedule â”€â”€> ict_vuln_scan â”€â”€> ict_vuln_finding

ict_fim_path â”€â”€> ict_fim_snapshot (baseline + periodik)
ict_fim_path â”€â”€> ict_fim_alert (delta alert)
```

---

### Domain: ict_website (Nginx, SSL, Uptime)

#### Tabel Inti dan Relasi

| Tabel | Referensi | Keterangan |
|---|---|---|
| ict_nginx_site | company_id, host_id | Konfigurasi virtual host |
| ict_nginx_site | upstream_id | Referensi ke ict_nginx_upstream |
| ict_nginx_upstream | company_id | Definisi upstream pool |
| ict_nginx_upstream_server | upstream_id | Server dalam upstream pool |
| ict_nginx_config_ver | site_id | Riwayat versi config site |
| ict_nginx_config_file | host_id | File konfigurasi mentah per host |
| ict_ssl_certificate | company_id | Sertifikat SSL per domain |
| ict_ssl_acme_account | company_id | Akun ACME per provider |
| ict_ssl_site_link | cert_id, site_id | Mapping sertifikat ke site |
| ict_nginx_log | host, timestamp | Access log Nginx |
| ict_nginx_app | host, timestamp | App request log |
| ict_nginx_atc | host, timestamp | Attack/threat log |
| ict_nginx_atc_sum | date, client_ip | Agregasi serangan harian |
| ict_nginx_sla | date | SLA harian Nginx |
| ict_uptimerobot_log | monitor_id, alert_datetime | Raw webhook UptimeRobot |
| ict_uptimerobot_sum | date, monitor_id | Agregasi uptime per monitor/hari |
| ict_uptimerobot_sla | date | SLA harian agregat semua monitor |
| ict_ip_whitelist | ip_or_cidr | IP/CIDR yang dikecualikan WAF |
| ict_ip_blacklist | ip | IP yang diblokir WAF |
| ict_waf_bypass_rule | domain, url_path | Bypass rule WAF per path |

#### Alur Nginx Site

```
ict_nginx_upstream
  â””â”€â”€ ict_nginx_upstream_server (member pool)

ict_nginx_site â”€â”€> ict_nginx_upstream (proxy pass pool)
ict_nginx_site â”€â”€> ict_ssl_certificate (via ict_ssl_site_link)
ict_nginx_site â”€â”€> ict_nginx_config_ver (riwayat config)
```

#### Alur Uptime Monitoring

```
UptimeRobot webhook â”€â”€> ict_uptimerobot_log (raw event)
                    â”€â”€> ict_uptimerobot_sum (agregasi per hari per monitor)
                    â”€â”€> ict_uptimerobot_sla (SLA harian global)
```

## 6. Constraints and Index Pattern

- Composite unique dipakai luas untuk tenancy/scope:
  - dat_company: [code] (unique global)
  - dat_company_module: [company_id, module_id]
  - dat_user_company: [user_id, company_id]
  - dat_user_privilege: [user_company_id, module_id]
  - dat_document: [company_id, code]
- Time-series/log table umumnya memiliki index timestamp/recorded_at untuk query monitoring
- Banyak domain memakai status string default untuk lifecycle state

## 7. Implementation Checklist (Database-first)

- [ ] Ubah model hanya di file schema yang sesuai domain
- [ ] Pastikan unique/index tetap konsisten dengan tenancy (company_id)
- [ ] Jalankan prisma generate setelah perubahan schema
- [ ] Jalankan migration dev dengan nama yang deskriptif
- [ ] Perbarui blueprint BASE ini jika ada model/enum baru
- [ ] Perbarui user guide BASE jika workflow operasional berubah

## 8. Validation Commands

- npx prisma generate
- prisma migrate dev (via tool prisma-migrate-dev)
- prisma migrate status (via tool prisma-migrate-status)

## 9. Referensi

- obx_base/prisma/schema/schema.prisma
- obx_base/prisma/schema/all_enum_dat.prisma
- obx_base/prisma/schema/dat_company.prisma
- obx_base/prisma/schema/dat_module.prisma
- obx_base/prisma/schema/dat_user.prisma
- obx_base/prisma/schema/dat_signature.prisma
- obx_base/prisma/schema/jsm_stack.prisma
- obx_base/prisma/schema/dat_document.prisma
- obx_base/prisma/schema/ict_machine.prisma
- obx_base/prisma/schema/ict_mikrotik.prisma
- obx_base/prisma/schema/ict_security.prisma
- obx_base/prisma/schema/ict_website.prisma
