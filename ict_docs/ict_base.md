# Plan: Database Schemas (ict_base)

Dokumen ini berisi seluruh rancangan schema database untuk modul-modul OmniSight menggunakan **Prisma ORM 7.8.0** dengan teknik **Multi-file schema layout** pada **PostgreSQL 18**.

---

## Daftar Isi

1. [Shared / Generic Tables](#1-shared--generic-tables)
2. [Core Data (User, Company, Module, Signature)](#2-core-data-user-company-module-signature)
3. [Location Management](#3-location-management)
4. [Asset / Inventory Management](#4-asset--inventory-management)
5. [Vulnerability Management](#5-vulnerability-management)
6. [Incident Management](#6-incident-management)
7. [SIEM / Log Monitoring](#7-siem--log-monitoring)
8. [Risk Assessment](#8-risk-assessment)
9. [Compliance Monitoring](#9-compliance-monitoring)
10. [Certificate Management](#10-certificate-management)
11. [Docker Management](#11-docker-management)
12. [Document Management](#12-document-management)
13. [Service Request Management](#13-service-request-management)
14. [Billing & Invoice](#14-billing--invoice)
15. [Preventive Maintenance](#15-preventive-maintenance)
16. [Dashboard & Reporting](#16-dashboard--reporting)
17. [Web Security (Nginx/WAF)](#17-web-security-nginxwaf)
18. [Web Monitoring (UptimeRobot)](#18-web-monitoring-uptimerobot)
19. [Asset Ownership & Assignment](#19-asset-ownership--assignment)
20. [Network Management (VLAN/Subnet/Route)](#20-network-management-vlansubnetroute)
21. [Mikrotik Device Management](#21-mikrotik-device-management)
22. [Domain Management](#22-domain-management)
23. [IP Public & Block Management](#23-ip-public--block-management)
24. [Monitor Agent](#24-monitor-agent)
25. [Monitor Metric](#25-monitor-metric)
26. [Monitor Log](#26-monitor-log)
27. [Monitor FIM](#27-monitor-fim)
28. [Monitor SCA](#28-monitor-sca)
29. [Monitor Discovery](#29-monitor-discovery)
30. [Monitor Alert](#30-monitor-alert)
31. [Monitor Docker Extensions](#31-monitor-docker-extensions)
32. [ISO 27001 & SOP](#32-iso-27001--sop)
33. [Naming Convention](#33-naming-convention)

---

## 1. Shared / Generic Tables

> File: `ict_base/prisma/schema/ict_generic.prisma`

Tabel-tabel generik yang digunakan bersama oleh beberapa modul untuk menghindari duplikasi struktur.

### Shared Enum

#### `severity_level` — Severity level universal

Digunakan oleh: SIEM, Vulnerability, Compliance, Docker Alert, Preventive Issue, Incident.

```prisma
enum severity_level {
  critical
  high
  medium
  low
  info
}
```

#### `approval_status` — Status approval universal

```prisma
enum approval_status {
  pending
  approved
  rejected
}
```

### `ict_audit_trail` — Timeline/audit trail universal

Menggantikan `ict_incident_timeline`, `ict_siem_alert_timeline`, dan `ict_sr_timeline`.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `entity_type` | String | incident, siem_alert, sr_request, bill_invoice, doc, risk, certificate, docker_host, pm_checklist |
| `entity_id` | String | ID entitas terkait |
| `user_id` | String | FK → dat_user |
| `action` | String | status_change, note_added, assignment_changed, escalation |
| `old_value` | String? | Nilai sebelumnya |
| `new_value` | String? | Nilai sesudahnya |
| `note` | String? | Catatan |
| `attachment_url` | String? | Lampiran |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([company_id, entity_type, entity_id])`, `@@index([user_id])`, `@@index([created_at])`

### `ict_file_attachment` — File attachment universal

Menggantikan `ict_doc_file`, `ict_pm_checklist_file`, `ict_bill_invoice_file`, dan `ict_comp_evidence`.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `entity_type` | String | doc, pm_checklist, bill_invoice, comp_finding, sr_request |
| `entity_id` | String | ID entitas terkait |
| `file_name` | String | Nama file asli |
| `file_path` | String | Path penyimpanan |
| `file_size` | Int | Size dalam bytes |
| `file_type` | String | MIME type |
| `file_hash` | String? | SHA-256 hash untuk integrity check |
| `description` | String? | Deskripsi file |
| `is_original` | Boolean | Default true (file asli dari vendor) |
| `upload_type` | String? | photo, document, report |
| `uploaded_by` | String | FK → dat_user |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([company_id, entity_type, entity_id])`, `@@index([file_hash])`, `@@index([uploaded_by])`

### `ict_approval` — Approval workflow universal

Menggantikan `ict_sr_approval` dan `ict_bill_approval`.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `entity_type` | String | sr_request, bill_invoice |
| `entity_id` | String | ID entitas terkait |
| `step` | Int | Urutan approval (1, 2, 3) |
| `approver_id` | String | FK → dat_user |
| `status` | enum | pending, approved, rejected |
| `comment` | String? | Catatan approval/reject |
| `acted_at` | DateTime? | Waktu diproses |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([entity_type, entity_id, step])`, `@@index([approver_id])`, `@@index([status])`

### `ict_notification` — Notification log universal

Menggantikan `ict_cert_notification` dan `ict_bill_notification`.

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `entity_type` | String | certificate, bill_invoice, sr_request |
| `entity_id` | String | ID entitas terkait |
| `user_id` | String | FK → dat_user (penerima) |
| `type` | String | reminder, overdue, approved, paid, expiry_warning, renewal_success, renewal_failed |
| `message` | String | Pesan notifikasi |
| `is_read` | Boolean | Default false |
| `sent_at` | DateTime | Default now() |
| `read_at` | DateTime? | Waktu dibaca |

**Index:** `@@index([user_id])`, `@@index([company_id, entity_type, entity_id])`, `@@index([is_read])`

### Mapping Modul → Generic Table

| Modul | Tabel Asli (Dihapus) | Digantikan Oleh | entity_type |
|-------|----------------------|-----------------|-------------|
| Incident | `ict_incident_timeline` | `ict_audit_trail` | `incident` |
| SIEM | `ict_siem_alert_timeline` | `ict_audit_trail` | `siem_alert` |
| Service Request | `ict_sr_timeline` | `ict_audit_trail` | `sr_request` |
| Service Request | `ict_sr_approval` | `ict_approval` | `sr_request` |
| Billing | `ict_bill_approval` | `ict_approval` | `bill_invoice` |
| Billing | `ict_bill_notification` | `ict_notification` | `bill_invoice` |
| Billing | `ict_bill_invoice_file` | `ict_file_attachment` | `bill_invoice` |
| Docman | `ict_doc_file` | `ict_file_attachment` | `doc` |
| Preventive | `ict_pm_checklist_file` | `ict_file_attachment` | `pm_checklist` |
| Compliance | `ict_comp_evidence` | `ict_file_attachment` | `comp_finding` |
| Certificate | `ict_cert_notification` | `ict_notification` | `certificate` |

### Storage Structure Universal

```
storage/
├── document/
│   └── {company_id}/
│       └── {entity_type}/{entity_id}/
├── billing/
│   └── {company_id}/
│       └── {invoice_id}/
├── preventive/
│   └── {company_id}/
│       └── {checklist_id}/
└── compliance/
    └── {company_id}/
        └── {finding_id}/
```

---

## 2. Core Data (User, Company, Module, Signature)

### 2.1 User Management

> File: `ict_base/prisma/schema/dat_user.prisma`

#### `dat_user` — Data pengguna utama

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `username` | String | Username login |
| `email` | String | Email |
| `password` | String | Hashed password |
| `fullname` | String | Nama lengkap |
| `phone` | String? | Nomor telepon |
| `company_id` | String | Default: "" |
| `employee_id` | String? | ID karyawan |
| `location_id` | String? | FK → dat_location |
| `department_id` | String? | ID departemen |
| `division_id` | String? | ID divisi |
| `role` | String | Default: "staff" |
| `job` | String | Default: "" |
| `key` | String | Default: "" |
| `is_admin` | Boolean | Default false |
| `is_hris` | Boolean | Default false |
| `is_active` | Boolean | Default false |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, username])`

#### `dat_user_company` — Relasi user ke perusahaan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `user_id` | String | FK → dat_user |
| `company_id` | String | FK → dat_company |
| `is_active` | Boolean | Default false |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([user_id, company_id])`

#### `dat_user_location` — Akses lokasi/area user

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `user_id` | String | FK → dat_user |
| `location_type_id` | String | FK → dat_location_type |
| `is_active` | Boolean | Default false |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([user_id, location_type_id])`

#### `dat_user_privilege` — Hak akses user per modul

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `user_company_id` | String | FK → dat_user_company |
| `module_id` | String | FK → dat_module |
| `level` | enum | Default: hide |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([user_company_id, module_id])`

#### `dat_user_session` — Sesi login

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `user_id` | String | FK → dat_user |
| `token` | String | Unique |
| `ip_address` | String? | IP client |
| `user_agent` | String? | User agent browser |
| `created_at` | DateTime | Default now() |
| `expires_at` | DateTime | Waktu kedaluwarsa |

#### `dat_user_action` — Log aktivitas user

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `user_id` | String | FK → dat_user |
| `company_id` | String? | |
| `module_code` | String? | Kode modul |
| `action` | String | Jenis aksi |
| `path` | String | Path yang diakses |
| `ip_address` | String? | |
| `user_agent` | String? | |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([user_id])`, `@@index([company_id])`

#### Enum

```prisma
enum action_type {
  hide
  view
  book
  post
}
```

### 2.2 Company Management

> File: `ict_base/prisma/schema/dat_company.prisma`

#### `dat_company` — Data perusahaan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `code` | String | Unique |
| `name` | String | Nama perusahaan |
| `vat_id` | String? | NPWP |
| `reg_no` | String? | No. registrasi |
| `address` | String? | Alamat |
| `valuta` | String | Default: "IDR" |
| `hris_link` | String? | Link HRIS |
| `is_active` | Boolean | Default false |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

#### `dat_company_module` — Modul yang aktif per perusahaan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `module_id` | String | FK → dat_module |
| `is_active` | Boolean | Default false |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, module_id])`

### 2.3 Module Management

> File: `ict_base/prisma/schema/dat_module.prisma`

#### `dat_module` — Daftar modul/halaman

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `parent_id` | String? | FK self-ref (hierarki) |
| `code` | String | Unique, kode halaman (AS01, VL02, dll) |
| `name` | String | Nama modul |
| `path` | String | Route path |
| `is_page` | Boolean | Default true |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

### 2.4 Signature & Approval

> File: `ict_base/prisma/schema/dat_signature.prisma`

#### `dat_signature_type` — Tipe tanda tangan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `code` | String | Unique |
| `name` | String | Nama tipe |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

#### `dat_approval_step` — Step approval dalam tipe

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `type_id` | String | FK → dat_signature_type |
| `step` | Int | Urutan step |
| `condition` | enum | ALL_APPROVED, ANY_APPROVED |

**Index:** `@@unique([type_id, step])`

#### `dat_approval_sign` — Penandatangan per step

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `step_id` | String | FK → dat_approval_step |
| `user_id` | String | FK → dat_user |

**Index:** `@@unique([step_id, user_id])`

#### `dat_signature_form` — Form signature yang sedang berjalan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `step` | Int | Step saat ini |
| `request_id` | String | ID request terkait |
| `condition` | enum | ALL_APPROVED, ANY_APPROVED |
| `status` | enum | Default: PENDING |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([request_id, step])`

#### `dat_signature_flag` — Status tanda tangan per user

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `form_id` | String | FK → dat_signature_form |
| `user_id` | String | FK → dat_user |
| `status` | enum | Default: PENDING |
| `comment` | String? | Catatan |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([form_id, user_id])`

#### Enum

```prisma
enum approval_flag {
  ALL_APPROVED
  ANY_APPROVED
}

enum approval_info {
  PENDING
  IN_PROGRESS
  APPROVED
  REJECTED
  COMPLETED
}
```

---

## 3. Location Management

> File: `ict_base/prisma/schema/dat_location.prisma`
> Status: **SELESAI**

### `dat_location_type` — Tipe lokasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `parent_id` | String? | FK self-ref (hierarki tipe) |
| `code` | String | Unique per company |
| `name` | String | Nama tipe |
| `description` | String? | Deskripsi |
| `icon` | String? | Icon untuk UI |
| `color` | String? | Warna untuk UI (hex) |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([parent_id])`

### `dat_location` — Data lokasi utama

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `location_type_id` | String | FK → dat_location_type |
| `parent_id` | String? | FK self-ref (hierarki) |
| `code` | String | Unique per company |
| `name` | String | Nama lokasi |
| `description` | String? | Deskripsi |
| `address` | String? | Alamat lengkap |
| `city` | String? | Kota |
| `province` | String? | Provinsi |
| `country` | String | Default: Indonesia |
| `postal_code` | String? | Kode pos |
| `latitude` | Decimal? | Latitude GPS |
| `longitude` | Decimal? | Longitude GPS |
| `phone` | String? | Nomor telepon |
| `email` | String? | Email kontak |
| `timezone` | String? | Default: Asia/Jakarta |
| `status` | enum | active, inactive, maintenance, closed |
| `is_active` | Boolean | Default true |
| `metadata` | Json? | Data tambahan per tipe |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([location_type_id])`, `@@index([parent_id])`

### `dat_location_detail` — Detail spesifik per tipe lokasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `location_id` | String | FK → dat_location (unique) |
| `detail_type` | String | datacenter, warehouse, store, office, production |
| `capacity` | Int? | Kapasitas (rack/server/orang) |
| `power_capacity` | String? | Kapasitas listrik (VA/kW) |
| `cooling_type` | String? | AC, Precision Cooling, Free Cooling |
| `floor_area` | Decimal? | Luas area (m²) |
| `floor_level` | String? | Lantai (GF, 1, 2, B1) |
| `rack_count` | Int? | Jumlah rack |
| `operating_hours` | String? | Jam operasional |
| `security_level` | String? | low, medium, high, critical |
| `has_generator` | Boolean | Default false |
| `has_ups` | Boolean | Default false |
| `has_fire_suppression` | Boolean | Default false |
| `has_cctv` | Boolean | Default false |
| `has_access_control` | Boolean | Default false |
| `contact_person` | String? | Nama PIC |
| `contact_phone` | String? | HP PIC |
| `contact_email` | String? | Email PIC |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

### `dat_location_contact` — Kontak person lokasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `location_id` | String | FK → dat_location |
| `name` | String | Nama kontak |
| `role` | String | Manager, Staff, Security, Maintenance |
| `phone` | String? | Nomor telepon |
| `email` | String? | Email |
| `is_primary` | Boolean | Default false |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([location_id, name])`

### `dat_location_asset` — Mapping aset ke lokasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `location_id` | String | FK → dat_location |
| `asset_id` | String | FK → ict_asset |
| `rack_position` | String? | Posisi rack (A1, B2) |
| `port` | String? | Port assignment |
| `install_date` | DateTime? | Tanggal instalasi |
| `notes` | String? | Catatan |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([location_id, asset_id])`

### Enum

```prisma
enum location_status {
  active
  inactive
  maintenance
  closed
}
```

### Seed Data: `dat_location_type`

| Code | Name | Icon | Color |
|------|------|------|-------|
| `store` | Store | store | #3B82F6 |
| `regional_office` | Regional Office | office | #8B5CF6 |
| `head_office` | Head Office | building | #10B981 |
| `warehouse` | Warehouse | warehouse | #F59E0B |
| `production` | Production | factory | #EF4444 |
| `datacenter` | Data Center | server | #6366F1 |

---

## 4. Asset / Inventory Management

> File: `ict_base/prisma/schema/ict_asset.prisma`

### `ict_asset_category` — Kategori aset (hierarchical)

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `parent_id` | String? | FK self-ref |
| `code` | String | Unique per company |
| `name` | String | Nama kategori |
| `description` | String? | Deskripsi |
| `icon` | String? | Nama icon untuk UI |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`

### `ict_asset` — Data aset utama

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `category_id` | String | FK → ict_asset_category |
| `location_id` | String? | FK → dat_location |
| `supplier_id` | String? | FK → ict_supplier |
| `code` | String | Unique per company |
| `name` | String | Nama aset |
| `description` | String? | Deskripsi |
| `asset_type` | String | server, vm, container, website, api, database, firewall, router, switch, ap, workstation, laptop, mobile |
| `ip_address` | String? | IP utama |
| `hostname` | String? | Hostname |
| `fqdn` | String? | Fully qualified domain name |
| `os_type` | String? | linux, windows, macos, ios, android, firmware |
| `os_version` | String? | Versi OS |
| `manufacturer` | String? | Vendor/pabrikan |
| `model` | String? | Model perangkat |
| `serial_number` | String? | Nomor serial |
| `owner_id` | String? | FK → dat_user |
| `criticality` | enum | low, medium, high, critical |
| `status` | enum | active, inactive, maintenance, decommissioned, planned |
| `purchase_date` | DateTime? | Tanggal beli |
| `warranty_end` | DateTime? | Akhir garansi |
| `last_scan_date` | DateTime? | Terakhir di-scan |
| `tags` | String[] | Tags fleksibel |
| `metadata` | Json? | Data tambahan |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([asset_type])`, `@@index([status])`, `@@index([criticality])`, `@@index([supplier_id])`

### `ict_asset_ip` — IP addresses aset

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `asset_id` | String | FK → ict_asset |
| `ip_address` | String | IP address |
| `ip_type` | String | ipv4, ipv6, public, private |
| `is_primary` | Boolean | Apakah IP utama |
| `interface_name` | String? | eth0, wlan0 |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([asset_id, ip_address])`, `@@index([ip_address])`

### `ict_asset_relation` — Relasi dependensi antar aset

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `source_asset_id` | String | FK → ict_asset |
| `target_asset_id` | String | FK → ict_asset |
| `relation_type` | String | depends_on, connects_to, hosted_on, protects, monitors |
| `description` | String? | Deskripsi relasi |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([source_asset_id, target_asset_id, relation_type])`

### Enum

```prisma
enum asset_criticality {
  low
  medium
  high
  critical
}

enum asset_status {
  active
  inactive
  maintenance
  decommissioned
  planned
}

enum assignment_status {
  assigned
  returned
  pending_return
  lost
  damaged
}

enum assignment_condition {
  excellent
  good
  fair
  poor
  damaged
  non_functional
}
```

---

## 5. Vulnerability Management

> File: `ict_base/prisma/schema/ict_vulnerability.prisma`

### `ict_vuln_scan_type` — Tipe scanner

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `code` | String | Unique: nessus, qualys, owasp_zap, nuclei, openvas, manual |
| `name` | String | Nama tipe scan |
| `description` | String? | |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |

### `ict_vuln_scan` — Instance scan (job)

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `scan_type_id` | String | FK → ict_vuln_scan_type |
| `name` | String | Nama scan job |
| `status` | enum | pending, running, completed, failed, cancelled |
| `started_at` | DateTime? | Waktu mulai |
| `finished_at` | DateTime? | Waktu selesai |
| `target_count` | Int | Jumlah target |
| `finding_count` | Int | Default 0 |
| `critical_count` | Int | Default 0 |
| `high_count` | Int | Default 0 |
| `medium_count` | Int | Default 0 |
| `low_count` | Int | Default 0 |
| `info_count` | Int | Default 0 |
| `scan_config` | Json? | Konfigurasi scan |
| `raw_result_url` | String? | URL/path ke hasil raw |
| `triggered_by` | String? | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, name])`, `@@index([company_id])`, `@@index([status])`

### `ict_vuln_finding` — Temuan vulnerability

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `scan_id` | String | FK → ict_vuln_scan |
| `asset_id` | String? | FK → ict_asset |
| `company_id` | String | FK → dat_company |
| `plugin_id` | String? | Plugin ID dari scanner |
| `cve_id` | String? | CVE identifier |
| `cwe_id` | String? | CWE identifier |
| `cvss_v2` | Decimal? | CVSS v2 score |
| `cvss_v3` | Decimal? | CVSS v3 score |
| `severity` | enum | Menggunakan `severity_level` (shared) |
| `title` | String | Judul temuan |
| `description` | String | Deskripsi detail |
| `solution` | String? | Saran remediasi |
| `reference_urls` | String[] | URL referensi |
| `status` | enum | open, confirmed, in_progress, remediated, false_positive, accepted, deferred |
| `verified` | Boolean | Default false |
| `first_seen` | DateTime | Pertama kali terdeteksi |
| `last_seen` | DateTime | Terakhir terdeteksi |
| `remediation_deadline` | DateTime? | Batas waktu remediasi |
| `remediation_notes` | String? | Catatan remediasi |
| `assigned_to` | String? | FK → dat_user |
| `evidence` | Json? | Proof/screenshot |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([scan_id, asset_id, plugin_id])`, `@@index([company_id])`, `@@index([asset_id])`, `@@index([severity])`, `@@index([status])`, `@@index([cve_id])`

### `ict_vuln_remediation` — Log remediasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `finding_id` | String | FK → ict_vuln_finding |
| `action` | String | patch, config_change, workaround, accept_risk, false_positive |
| `description` | String | Deskripsi tindakan |
| `evidence` | Json? | Bukti remediasi |
| `performed_by` | String | FK → dat_user |
| `verified_by` | String? | FK → dat_user |
| `verified_at` | DateTime? | Waktu verifikasi |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([finding_id])`

### Enum

```prisma
enum vuln_severity {
  critical
  high
  medium
  low
  info
}

enum vuln_status {
  open
  confirmed
  in_progress
  remediated
  false_positive
  accepted
  deferred
}

enum scan_status {
  pending
  running
  completed
  failed
  cancelled
}
```

---

## 6. Incident Management

> File: `ict_base/prisma/schema/ict_incident.prisma`

### `ict_incident_type` — Tipe insiden

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Unique per company |
| `name` | String | data_breach, malware, phishing, ddos, unauthorized_access, insider_threat |
| `description` | String? | |
| `default_priority` | enum | Default priority saat type dipilih |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`

### `ict_incident` — Data insiden utama

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `type_id` | String | FK → ict_incident_type |
| `code` | String | Unique per company, format: INC-YYYYMMDD-XXXX |
| `title` | String | Judul insiden |
| `description` | String | Deskripsi detail |
| `severity` | enum | Menggunakan `severity_level` (shared) |
| `priority` | enum | p1, p2, p3, p4 |
| `status` | enum | detected, triaged, contained, eradicated, recovered, closed |
| `detection_source` | String? | siem, user_report, vulnerability_scan, external_report, monitoring |
| `siem_alert_id` | String? | FK → ict_siem_alert |
| `reported_by` | String | FK → dat_user |
| `assigned_to` | String? | FK → dat_user |
| `sla_deadline` | DateTime? | Deadline SLA |
| `affected_users_count` | Int? | Jumlah user terdampak |
| `data_compromised` | Boolean | Default false |
| `data_type_compromised` | String? | pii, phi, financial, credentials |
| `root_cause` | String? | Diisi saat remediasi |
| `lessons_learned` | String? | Diisi saat close |
| `started_at` | DateTime | Waktu kejadian |
| `detected_at` | DateTime | Waktu terdeteksi |
| `contained_at` | DateTime? | Waktu contained |
| `eradicated_at` | DateTime? | Waktu eradicated |
| `recovered_at` | DateTime? | Waktu recovered |
| `closed_at` | DateTime? | Waktu closed |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([status])`, `@@index([severity])`, `@@index([priority])`

### `ict_incident_asset` — Aset terdampak

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `incident_id` | String | FK → ict_incident |
| `asset_id` | String | FK → ict_asset |
| `impact_description` | String? | Deskripsi dampak |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([incident_id, asset_id])`

### `ict_incident_user` — Tim penanganan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `incident_id` | String | FK → ict_incident |
| `user_id` | String | FK → dat_user |
| `role` | String | reporter, handler, escalator, approver |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([incident_id, user_id, role])`

### Enum

```prisma
enum incident_priority {
  p1
  p2
  p3
  p4
}

enum incident_status {
  detected
  triaged
  contained
  eradicated
  recovered
  closed
}
```

### SLA Matrix

| Priority | Deteksi ke Triage | Triage ke Contain | Total Resolution |
|----------|-------------------|-------------------|------------------|
| P1 (Critical) | 15 menit | 1 jam | 4 jam |
| P2 (High) | 30 menit | 4 jam | 24 jam |
| P3 (Medium) | 2 jam | 8 jam | 72 jam |
| P4 (Low) | 8 jam | 24 jam | 168 jam |

---

## 7. SIEM / Log Monitoring

> File: `ict_base/prisma/schema/ict_siem.prisma`

### Model yang Sudah Ada (`ict_web_security.prisma`)

| Model | Keterangan |
|-------|------------|
| `ict_nginx_log` | Log akses nginx per-domain |
| `ict_nginx_app` | Log aplikasi nginx |
| `ict_nginx_atc` | Log ATC (Attack Traffic Counter) |
| `ict_nginx_atc_sum` | Ringkasan ATC per hari |
| `ict_nginx_sla` | SLA nginx per hari |
| `ict_ip_whitelist` | IP whitelist |
| `ict_ip_blacklist` | IP blacklist |
| `ict_waf_bypass_rule` | Rule bypass WAF |

### `ict_siem_source` — Sumber log

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Unique per company |
| `name` | String | Nama sumber log |
| `source_type` | String | syslog, filebeat, api, database, agent |
| `endpoint` | String? | URL/host sumber |
| `config` | Json? | Konfigurasi koneksi |
| `is_enabled` | Boolean | Default true |
| `last_received_at` | DateTime? | Terakhir terima log |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`

### `ict_siem_rule` — Detection rules

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Unique per company |
| `name` | String | Nama rule |
| `description` | String? | |
| `rule_type` | String | threshold, pattern, anomaly, compound |
| `query` | Json | Rule definition |
| `severity` | enum | Menggunakan `severity_level` (shared) |
| `mitre_tactic` | String? | MITRE ATT&CK tactic |
| `mitre_technique` | String? | MITRE ATT&CK technique |
| `mitre_technique_id` | String? | T1595, T1190, dll |
| `is_enabled` | Boolean | Default true |
| `cooldown_minutes` | Int | Default 15 |
| `trigger_count` | Int | Default 0 |
| `last_triggered_at` | DateTime? | |
| `created_by` | String? | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([severity])`, `@@index([is_enabled])`

### `ict_siem_alert` — Alert dari rule matching

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `rule_id` | String | FK → ict_siem_rule |
| `company_id` | String | FK → dat_company |
| `title` | String | Judul alert |
| `description` | String | Deskripsi alert |
| `severity` | enum | Menggunakan `severity_level` (shared) |
| `status` | enum | new, acknowledged, investigating, resolved, false_positive |
| `source_ip` | String? | IP sumber |
| `source_host` | String? | Host sumber |
| `affected_asset_id` | String? | FK → ict_asset |
| `evidence` | Json? | Log/bukti |
| `mitre_tactic` | String? | |
| `mitre_technique` | String? | |
| `assigned_to` | String? | FK → dat_user |
| `acknowledged_at` | DateTime? | |
| `investigating_at` | DateTime? | |
| `resolved_at` | DateTime? | |
| `resolution_notes` | String? | |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([company_id])`, `@@index([rule_id])`, `@@index([status])`, `@@index([severity])`, `@@index([created_at])`, `@@index([source_ip])`

### `ict_siem_correlation` — Correlation events

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `name` | String | Nama correlation |
| `description` | String? | |
| `rule_ids` | String[] | Array of rule IDs |
| `time_window_minutes` | Int | Jendela waktu korelasi |
| `match_count` | Int | Jumlah match minimum |
| `is_enabled` | Boolean | Default true |
| `created_at` | DateTime | Default now() |

### Enum

```prisma
enum alert_status {
  new
  acknowledged
  investigating
  resolved
  false_positive
}
```

---

## 8. Risk Assessment

> File: `ict_base/prisma/schema/ict_risk.prisma`

### `ict_risk_category` — Kategori risiko (hierarchical)

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `parent_id` | String? | FK self-ref |
| `code` | String | Unique per company |
| `name` | String | operational, financial, strategic, compliance, reputational |
| `description` | String? | |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`

### `ict_risk` — Register risiko

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `category_id` | String | FK → ict_risk_category |
| `asset_id` | String? | FK → ict_asset |
| `comp_finding_id` | String? | FK → ict_comp_finding |
| `code` | String | Unique per company, format: RSK-XXXX |
| `title` | String | Judul risiko |
| `description` | String | Deskripsi detail |
| `source` | String? | Sumber identifikasi |
| `likelihood` | enum | rare(1), unlikely(2), possible(3), likely(4), almost_certain(5) |
| `impact` | enum | negligible(1), minor(2), moderate(3), major(4), severe(5) |
| `risk_score` | Int | Auto: likelihood × impact (1-25) |
| `risk_level` | enum | critical(20-25), high(15-19), medium(10-14), low(5-9), very_low(1-4) |
| `inherent_score` | Int? | Score inherent |
| `residual_score` | Int? | Score residual |
| `status` | enum | identified, assessed, treating, monitored, closed, accepted |
| `owner_id` | String? | FK → dat_user |
| `identified_date` | DateTime | Tanggal identifikasi |
| `last_review_date` | DateTime? | Terakhir di-review |
| `next_review_date` | DateTime? | Review berikutnya |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([risk_level])`, `@@index([status])`

### `ict_risk_assessment` — Penilaian risiko periodik

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `risk_id` | String | FK → ict_risk |
| `assessor_id` | String | FK → dat_user |
| `assessment_date` | DateTime | Tanggal penilaian |
| `new_likelihood` | enum | Likelihood terbaru |
| `new_impact` | enum | Impact terbaru |
| `new_score` | Int | Score terbaru |
| `notes` | String? | Catatan |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([risk_id])`

### `ict_risk_treatment` — Treatment plan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `risk_id` | String | FK → ict_risk |
| `treatment_type` | String | mitigate, transfer, avoid, accept |
| `description` | String | Deskripsi treatment |
| `action_items` | Json? | Daftar action items |
| `cost_estimate` | Decimal? | Estimasi biaya |
| `target_date` | DateTime? | Target selesai |
| `owner_id` | String? | FK → dat_user |
| `status` | String | planned, in_progress, completed, cancelled |
| `effectiveness` | String? | effective, partially_effective, ineffective |
| `completed_at` | DateTime? | |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([risk_id])`, `@@index([status])`

### Enum

```prisma
enum risk_likelihood {
  rare
  unlikely
  possible
  likely
  almost_certain
}

enum risk_impact {
  negligible
  minor
  moderate
  major
  severe
}

enum risk_level {
  critical
  high
  medium
  low
  very_low
}

enum risk_status {
  identified
  assessed
  treating
  monitored
  closed
  accepted
}
```

### Risk Matrix (5×5)

```
Impact →     Negligible(1)  Minor(2)  Moderate(3)  Major(4)  Severe(5)
Likelihood ↓
Almost Certain(5)    5           10          15           20        25
Likely(4)            4            8          12           16        20
Possible(3)          3            6           9           12        15
Unlikely(2)          2            4           6            8        10
Rare(1)              1            2           3            4         5
```

---

## 9. Compliance Monitoring

> File: `ict_base/prisma/schema/ict_compliance.prisma`

### Enum

```prisma
enum assessment_status {
  planned
  in_progress
  completed
  archived
}

enum compliance_status {
  compliant
  non_compliant
  partial
  not_applicable
}

enum soa_status {
  implemented
  partial
  excluded
}

enum capa_source {
  incident
  audit
  complaint
  risk
  improvement
  nonconformity
}

enum capa_status {
  open
  investigating
  root_cause
  corrective_action
  verification
  effectiveness
  closed
  reopened
}

enum capa_severity {
  critical
  major
  minor
  observation
}

enum finding_type {
  nonconformity
  observation
  opportunity_for_improvement
  major_nonconformity
  minor_nonconformity
}

enum audit_type {
  internal
  external
  supplier
  certification
}

enum audit_status {
  planned
  in_progress
  completed
  cancelled
}

enum review_status {
  scheduled
  in_progress
  completed
  deferred
}
```

---

### `ict_comp_standard` — Standar compliance

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Kode standar (unique per company) |
| `name` | String | Nama standar |
| `version` | String? | Versi standar |
| `description` | String? | Deskripsi singkat |
| `website` | String? | URL referensi resmi |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |

**Relasi:** `controls` → ict_comp_control[], `assessments` → ict_comp_assessment[], `ictCompSoas` → ict_comp_soa[]

**Index:** `@@unique([company_id, code])`

### `ict_comp_control` — Kontrol/requirement (hierarchical)

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `standard_id` | String | FK → ict_comp_standard |
| `parent_id` | String? | FK self-ref → ict_comp_control |
| `code` | String | Kode kontrol (unique per standard) |
| `name` | String | Nama kontrol |
| `description` | String? | Deskripsi kontrol |
| `requirement_text` | String | Teks requirement resmi |
| `implementation_guidance` | String? | Panduan implementasi |
| `is_mandatory` | Boolean | Default true |
| `sort_order` | Int | Urutan tampilan |
| `created_at` | DateTime | Default now() |

**Relasi:** `standard` → ict_comp_standard, `parent` → ict_comp_control? (self-ref), `children` → ict_comp_control[], `findings` → ict_comp_finding[], `ictCompSoas` → ict_comp_soa[]

**Index:** `@@unique([standard_id, code])`, `@@index([standard_id])`

### `ict_comp_assessment` — Hasil assessment compliance

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `standard_id` | String | FK → ict_comp_standard |
| `name` | String | Nama assessment (unique per company) |
| `assessment_date` | DateTime | Tanggal assessment |
| `status` | assessment_status | Default planned |
| `assessor_name` | String? | Nama auditor/assessor |
| `assessor_org` | String? | Organisasi auditor |
| `total_controls` | Int | Default 0, jumlah total kontrol |
| `compliant_count` | Int | Default 0, jumlah sesuai |
| `non_compliant_count` | Int | Default 0, jumlah tidak sesuai |
| `partial_count` | Int | Default 0, jumlah sebagian |
| `not_applicable_count` | Int | Default 0, jumlah tidak berlaku |
| `compliance_score` | Decimal? | Persentase skor kepatuhan |
| `report_url` | String? | URL laporan assessment |
| `notes` | String? | Catatan tambahan |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Relasi:** `standard` → ict_comp_standard, `findings` → ict_comp_finding[]

**Index:** `@@unique([company_id, name])`, `@@index([company_id])`, `@@index([standard_id])`

### `ict_comp_finding` — Temuan hasil assessment

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `assessment_id` | String | FK → ict_comp_assessment |
| `control_id` | String | FK → ict_comp_control |
| `company_id` | String | FK → dat_company |
| `status` | compliance_status | Status kepatuhan temuan |
| `finding_text` | String | Deskripsi temuan |
| `evidence_text` | String? | Bukti kepatuhan |
| `evidence_files` | String[] | Daftar path file bukti |
| `severity` | severity_level | Menggunakan shared enum `severity_level` |
| `remediation_plan` | String? | Rencana perbaikan |
| `target_date` | DateTime? | Target tanggal perbaikan |
| `assigned_to` | String? | FK → dat_user |
| `completed_at` | DateTime? | Waktu perbaikan selesai |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Relasi:** `assessment` → ict_comp_assessment, `control` → ict_comp_control, `capas` → ict_capa[]

**Index:** `@@unique([assessment_id, control_id])`, `@@index([company_id])`, `@@index([status])`, `@@index([severity])`

### `ict_comp_soa` — Statement of Applicability

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `standard_id` | String | FK → ict_comp_standard |
| `control_id` | String | FK → ict_comp_control |
| `is_applicable` | Boolean | Default true, apakah kontrol berlaku |
| `status` | soa_status | Default implemented |
| `implementation_notes` | String? | Catatan implementasi |
| `evidence_urls` | String[] | Daftar URL bukti implementasi |
| `exclusion_justification` | String? | Justifikasi jika tidak berlaku |
| `reviewed_by` | String? | FK → dat_user |
| `reviewed_at` | DateTime? | Tanggal review |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Relasi:** `standard` → ict_comp_standard, `control` → ict_comp_control, `reviewer` → dat_user?

**Index:** `@@unique([company_id, standard_id, control_id])`, `@@index([company_id])`, `@@index([standard_id])`, `@@index([is_applicable])`

### `ict_capa` — Corrective and Preventive Action

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Kode CAPA unik (unique per company) |
| `source_type` | capa_source | Sumber asal CAPA |
| `source_id` | String? | ID sumber (FK ke finding/audit, polymorphic) |
| `title` | String | Judul CAPA |
| `description` | String | Deskripsi detail |
| `immediate_action` | String? | Tindakan segera/awal |
| `root_cause` | String? | Hasil analisis akar masalah |
| `root_cause_method` | String? | Metode analisis (5-Why, Fishbone, FMEA, dll) |
| `severity` | capa_severity | Default minor |
| `status` | capa_status | Default open |
| `assigned_to` | String? | FK → dat_user |
| `target_date` | DateTime? | Target penyelesaian |
| `effectiveness` | String? | Hasil evaluasi efektivitas |
| `effectiveness_notes` | String? | Catatan evaluasi |
| `verified_by` | String? | FK → dat_user |
| `verified_at` | DateTime? | Tanggal verifikasi |
| `closed_at` | DateTime? | Tanggal ditutup |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Relasi:** `assignee` → dat_user?, `verifier` → dat_user?, `actions` → ict_capa_action[], `verification` → ict_capa_verification?, `finding` → ict_comp_finding?

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([status])`, `@@index([severity])`, `@@index([source_type, source_id])`, `@@index([assigned_to])`

### `ict_capa_action` — Tindakan detail CAPA

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `capa_id` | String | FK → ict_capa (onDelete: Cascade) |
| `action_type` | String | Jenis tindakan (corrective, preventive, dll) |
| `description` | String | Deskripsi tindakan |
| `responsible_id` | String? | FK → dat_user |
| `due_date` | DateTime? | Jatuh tempo |
| `completed_at` | DateTime? | Waktu penyelesaian |
| `evidence` | String? | Bukti pelaksanaan |
| `status` | String | Default "pending" |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Relasi:** `capa` → ict_capa, `responsible` → dat_user?

**Index:** `@@index([capa_id])`, `@@index([status])`

### `ict_capa_verification` — Verifikasi efektivitas CAPA

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `capa_id` | String | FK → ict_capa (unique, satu verifikasi per CAPA) |
| `verifier_id` | String | FK → dat_user |
| `verification_date` | DateTime | Tanggal verifikasi |
| `is_effective` | Boolean | Apakah CAPA efektif |
| `notes` | String? | Catatan verifikasi |
| `follow_up_required` | Boolean | Default false |
| `follow_up_date` | DateTime? | Tanggal follow-up |
| `created_at` | DateTime | Default now() |

**Relasi:** `capa` → ict_capa, `verifier` → dat_user

**Index:** `@@unique([capa_id])`

### `ict_audit_program` — Program audit tahunan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `year` | Int | Tahun program audit |
| `title` | String | Judul program |
| `description` | String? | Deskripsi program |
| `status` | String | Default "planned" |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |
| `dat_userId` | String? | FK → dat_user |

**Relasi:** `audits` → ict_audit[], `datUser` → dat_user?

**Index:** `@@unique([company_id, year])`, `@@index([company_id])`

### `ict_audit` — Data audit individual

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `program_id` | String | FK → ict_audit_program |
| `company_id` | String | FK → dat_company |
| `title` | String | Judul audit |
| `scope` | String | Ruang lingkup audit |
| `objectives` | String? | Tujuan audit |
| `audit_type` | audit_type | Default internal |
| `criteria` | String? | Kriteria/basis audit |
| `auditor_id` | String | FK → dat_user |
| `audit_date` | DateTime | Tanggal pelaksanaan audit |
| `planned_duration_days` | Int? | Durasi rencana (hari) |
| `status` | audit_status | Default planned |
| `finding_count` | Int | Default 0, jumlah temuan total |
| `nonconformity_count` | Int | Default 0, jumlah ketidaksesuaian |
| `observation_count` | Int | Default 0, jumlah observasi |
| `improvement_count` | Int | Default 0, jumlah peluang perbaikan |
| `report_url` | String? | URL laporan audit |
| `notes` | String? | Catatan |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Relasi:** `program` → ict_audit_program, `auditor` → dat_user, `findings` → ict_audit_finding[]

**Index:** `@@index([program_id])`, `@@index([company_id])`, `@@index([status])`, `@@index([audit_date])`

### `ict_audit_finding` — Temuan audit

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `audit_id` | String | FK → ict_audit (onDelete: Cascade) |
| `company_id` | String | FK → dat_company |
| `finding_type` | finding_type | Tipe temuan |
| `severity` | capa_severity | Tingkat keparahan |
| `clause` | String? | Klausul standar yang terdampak |
| `description` | String | Deskripsi temuan |
| `evidence` | String? | Bukti temuan |
| `requirement_reference` | String? | Referensi requirement |
| `root_cause` | String? | Akar masalah |
| `assigned_to` | String? | FK → dat_user |
| `target_date` | DateTime? | Target perbaikan |
| `status` | String | Default "open" |
| `verified_by` | String? | FK → dat_user |
| `verified_at` | DateTime? | Tanggal verifikasi |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Relasi:** `audit` → ict_audit, `assignee` → dat_user?, `verifier` → dat_user?

**Index:** `@@index([audit_id])`, `@@index([company_id])`, `@@index([finding_type])`, `@@index([status])`, `@@index([severity])`

### `ict_mgmt_review` — Management review

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `title` | String | Judul review |
| `review_date` | DateTime | Tanggal review |
| `review_type` | String | Jenis review (quarterly, annual, ad-hoc, dll) |
| `status` | review_status | Default scheduled |
| `attendees` | String[] | Daftar nama/ID peserta |
| `meeting_minutes_url` | String? | URL notulensi rapat |
| `summary` | String? | Ringkasan hasil review |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Relasi:** `inputs` → ict_mgmt_review_input[], `outputs` → ict_mgmt_review_output[]

**Index:** `@@index([company_id])`, `@@index([status])`, `@@index([review_date])`

### `ict_mgmt_review_input` — Input data untuk management review

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `review_id` | String | FK → ict_mgmt_review (onDelete: Cascade) |
| `category` | String | Kategori input (compliance, risk, incident, dll) |
| `title` | String | Judul input |
| `data_source` | String? | Sumber data |
| `summary` | String | Ringkasan input |
| `metrics` | Json? | Data metrik dalam format JSON |
| `created_at` | DateTime | Default now() |

**Relasi:** `review` → ict_mgmt_review

**Index:** `@@index([review_id])`

### `ict_mgmt_review_output` — Keputusan/tindakan dari management review

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `review_id` | String | FK → ict_mgmt_review (onDelete: Cascade) |
| `decision` | String | Keputusan hasil review |
| `action_required` | Boolean | Default false, apakah perlu tindakan lanjut |
| `assigned_to` | String? | FK → dat_user |
| `target_date` | DateTime? | Target pelaksanaan |
| `status` | String | Default "pending" |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Relasi:** `review` → ict_mgmt_review, `assignee` → dat_user?

**Index:** `@@index([review_id])`, `@@index([status])`

---

## 10. Certificate Management

> File: `ict_base/prisma/schema/ict_certificate.prisma`

### `ict_cert_type` — Tipe sertifikat

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `code` | String | Unique: ssl_tls, code_signing, client_cert, email_cert |
| `name` | String | |
| `description` | String? | |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |

### `ict_certificate` — Data sertifikat

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `type_id` | String | FK → ict_cert_type |
| `asset_id` | String? | FK → ict_asset |
| `name` | String | Nama sertifikat |
| `domain` | String | Domain utama |
| `san` | String[] | Subject Alternative Names |
| `issuer` | String | Penerbit |
| `serial_number` | String? | |
| `fingerprint_sha256` | String? | SHA-256 fingerprint |
| `not_before` | DateTime | Berlaku dari |
| `not_after` | DateTime | Berlaku sampai |
| `key_type` | String | RSA, ECDSA |
| `key_size` | Int? | 2048, 4096, 256, 384 |
| `auto_renew` | Boolean | Default false |
| `renewal_method` | String? | certbot, acme, manual |
| `notification_days` | Int | Default 30 |
| `status` | enum | active, expired, revoked, pending_renewal, suspended |
| `last_checked_at` | DateTime? | Terakhir dicek |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, domain])`, `@@index([company_id])`, `@@index([status])`, `@@index([not_after])`

### `ict_cert_renewal_log` — Log perpanjangan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `certificate_id` | String | FK → ict_certificate |
| `action` | String | renewed, failed, manual_update |
| `old_not_after` | DateTime? | Tanggal expired lama |
| `new_not_after` | DateTime? | Tanggal expired baru |
| `performed_by` | String? | FK → dat_user atau "system" |
| `notes` | String? | |
| `error_message` | String? | Jika gagal |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([certificate_id])`

### Enum

```prisma
enum cert_status {
  active
  expired
  revoked
  pending_renewal
  suspended
}
```

---

## 11. Docker Management

> File: `ict_base/prisma/schema/ict_docker.prisma`

### `ict_docker_host` — Docker host/server

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `location_id` | String? | FK → dat_location |
| `name` | String | Nama host |
| `hostname` | String | Hostname/IP |
| `port` | Int | Docker API port (default 2375) |
| `use_tls` | Boolean | Default false |
| `tls_cert_path` | String? | Path sertifikat TLS |
| `socket_path` | String? | Unix socket path |
| `connection_type` | enum | socket, tcp |
| `status` | enum | connected, disconnected, error |
| `last_seen_at` | DateTime? | Terakhir terhubung |
| `os_type` | String? | linux, windows, macos |
| `os_version` | String? | Versi OS |
| `docker_version` | String? | Versi Docker Engine |
| `cpu_count` | Int? | Jumlah CPU |
| `memory_total` | BigInt? | Total memory (bytes) |
| `disk_total` | BigInt? | Total disk (bytes) |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, name])`, `@@index([company_id])`, `@@index([status])`

### `ict_docker_container` — Container

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `host_id` | String | FK → ict_docker_host |
| `company_id` | String | FK → dat_company |
| `container_id` | String | Docker container ID (short) |
| `name` | String | Container name |
| `image` | String | Image name:tag |
| `image_id` | String | Image ID |
| `status` | enum | running, stopped, paused, restarting, created, removed |
| `state` | String | Container state detail |
| `created_at_docker` | DateTime | Created time dari Docker |
| `started_at` | DateTime? | Waktu start terakhir |
| `finished_at` | DateTime? | Waktu stop terakhir |
| `restart_count` | Int | Default 0 |
| `platform` | String? | linux/amd64 |
| `pid` | Int? | Process ID |
| `exit_code` | Int? | Exit code |
| `compose_stack` | String? | Nama compose stack |
| `compose_service` | String? | Nama service di compose |
| `labels` | Json? | Docker labels |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([host_id, container_id])`, `@@index([host_id])`, `@@index([company_id])`, `@@index([status])`, `@@index([compose_stack])`

### `ict_docker_container_port` — Port mapping

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `container_id` | String | FK → ict_docker_container |
| `private_port` | Int | Port di dalam container |
| `public_port` | Int? | Port di host |
| `protocol` | String | tcp, udp |
| `ip` | String? | IP binding |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([container_id])`, `@@unique([container_id, private_port, protocol])`

### `ict_docker_container_env` — Environment variables

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `container_id` | String | FK → ict_docker_container |
| `key` | String | Nama variabel |
| `value` | String | Nilai |
| `is_secret` | Boolean | Default false |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([container_id])`, `@@unique([container_id, key])`

### `ict_docker_container_mount` — Volume mounts

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `container_id` | String | FK → ict_docker_container |
| `source` | String | Host path atau volume name |
| `destination` | String | Container path |
| `type` | String | bind, volume, tmpfs |
| `mode` | String? | rw, ro |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([container_id])`

### `ict_docker_container_network` — Network attachment

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `container_id` | String | FK → ict_docker_container |
| `network_id` | String | FK → ict_docker_network |
| `ip_address` | String? | IP di network |
| `gateway` | String? | Gateway |
| `mac_address` | String? | MAC address |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([container_id])`, `@@index([network_id])`, `@@unique([container_id, network_id])`

### `ict_docker_image` — Image

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `host_id` | String | FK → ict_docker_host |
| `image_id` | String | Docker image ID |
| `repository` | String | Repository name |
| `tag` | String | Tag |
| `size` | BigInt | Size in bytes |
| `created_at_docker` | DateTime | Created time |
| `container_count` | Int | Default 0 |
| `is_dangling` | Boolean | Default false |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([host_id, image_id])`, `@@index([host_id])`, `@@index([repository])`

### `ict_docker_network` — Network

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `host_id` | String | FK → ict_docker_host |
| `network_id` | String | Docker network ID |
| `name` | String | Network name |
| `driver` | String | bridge, overlay, host, macvlan, none |
| `scope` | String | local, global, swarm |
| `subnet` | String? | Subnet CIDR |
| `gateway` | String? | Gateway IP |
| `ip_range` | String? | IP range |
| `is_internal` | Boolean | Default false |
| `container_count` | Int | Default 0 |
| `created_at_docker` | DateTime | Created time |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([host_id, network_id])`, `@@index([host_id])`

### `ict_docker_volume` — Volume

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `host_id` | String | FK → ict_docker_host |
| `volume_id` | String | Docker volume ID |
| `name` | String | Volume name |
| `driver` | String | local, nfs, etc |
| `mountpoint` | String | Host mount point |
| `scope` | String | local, global, swarm |
| `labels` | Json? | Volume labels |
| `options` | Json? | Driver options |
| `container_count` | Int | Default 0 |
| `size` | BigInt? | Size in bytes |
| `created_at_docker` | DateTime | Created time |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([host_id, volume_id])`, `@@index([host_id])`

### `ict_docker_compose` — Compose stack

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `host_id` | String | FK → ict_docker_host |
| `company_id` | String | FK → dat_company |
| `name` | String | Nama stack |
| `status` | enum | running, stopped, error, partial |
| `compose_content` | String | Isi compose.yaml |
| `env_content` | String? | Isi .env |
| `workdir` | String? | Working directory di host |
| `project_name` | String? | Docker Compose project name |
| `service_count` | Int | Default 0 |
| `running_count` | Int | Default 0 |
| `last_deployed_at` | DateTime? | Terakhir di-deploy |
| `last_action_by` | String? | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([host_id, name])`, `@@index([host_id])`, `@@index([company_id])`, `@@index([status])`

### `ict_docker_compose_service` — Service dalam stack

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `compose_id` | String | FK → ict_docker_compose |
| `container_id` | String? | FK → ict_docker_container |
| `name` | String | Nama service |
| `image` | String | Image yang digunakan |
| `status` | enum | running, stopped, error, unknown |
| `replicas` | Int | Default 1 |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([compose_id, name])`, `@@index([compose_id])`

### `ict_docker_stats` — Container stats (historical)

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `container_id` | String | FK → ict_docker_container |
| `host_id` | String | FK → ict_docker_host |
| `timestamp` | DateTime | Waktu pengukuran |
| `cpu_percent` | Decimal | CPU usage % |
| `memory_usage` | BigInt | Memory usage (bytes) |
| `memory_limit` | BigInt | Memory limit (bytes) |
| `memory_percent` | Decimal | Memory usage % |
| `network_rx` | BigInt | Network received (bytes) |
| `network_tx` | BigInt | Network transmitted (bytes) |
| `block_read` | BigInt | Block I/O read (bytes) |
| `block_write` | BigInt | Block I/O write (bytes) |
| `pids` | Int | Jumlah proses |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([container_id, timestamp])`, `@@index([host_id, timestamp])`, `@@index([timestamp])`

### `ict_docker_host_stats` — Host system stats

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `host_id` | String | FK → ict_docker_host |
| `timestamp` | DateTime | Waktu pengukuran |
| `cpu_percent` | Decimal | CPU usage % |
| `memory_usage` | BigInt | Memory usage (bytes) |
| `memory_limit` | BigInt | Memory limit (bytes) |
| `memory_percent` | Decimal | Memory usage % |
| `disk_usage` | BigInt | Disk usage (bytes) |
| `disk_total` | BigInt | Disk total (bytes) |
| `disk_percent` | Decimal | Disk usage % |
| `network_rx` | BigInt | Network received (bytes) |
| `network_tx` | BigInt | Network transmitted (bytes) |
| `load_avg_1` | Decimal? | Load average 1 min |
| `load_avg_5` | Decimal? | Load average 5 min |
| `load_avg_15` | Decimal? | Load average 15 min |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([host_id, timestamp])`, `@@index([timestamp])`

### `ict_docker_image_pull_log` — Image pull history

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `host_id` | String | FK → ict_docker_host |
| `image_name` | String | Nama image |
| `status` | enum | success, failed |
| `pulled_by` | String? | FK → dat_user |
| `error_message` | String? | |
| `duration_ms` | Int? | Durasi pull (ms) |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([host_id])`

### Enum

```prisma
enum host_connection {
  socket
  tcp
}

enum host_status {
  connected
  disconnected
  error
}

enum container_status {
  running
  stopped
  paused
  restarting
  created
  removed
}

enum compose_status {
  running
  stopped
  error
  partial
}
```

---

## 12. Document Management

> File: `ict_base/prisma/schema/ict_document.prisma`

### `ict_doc_type` — Tipe dokumen

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `code` | String | Unique: doc, sop, todo, template, policy |
| `name` | String | Nama tipe |
| `description` | String? | |
| `icon` | String? | Icon untuk UI |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |

### `ict_doc_category` — Kategori dokumen (hierarchical)

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `parent_id` | String? | FK self-ref |
| `code` | String | Unique per company |
| `name` | String | Nama kategori |
| `description` | String? | |
| `sort_order` | Int | Default 0 |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`

### `ict_doc` — Dokumen utama

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `type_id` | String | FK → ict_doc_type |
| `category_id` | String? | FK → ict_doc_category |
| `parent_id` | String? | FK self-ref |
| `code` | String | Unique per company |
| `title` | String | Judul dokumen |
| `slug` | String | URL-friendly slug |
| `description` | String? | Ringkasan singkat |
| `content` | String | Konten utama (Markdown/HTML) |
| `content_format` | String | markdown, html, json |
| `version` | Int | Default 1 |
| `status` | enum | draft, in_review, approved, published, archived |
| `visibility` | enum | private, internal, public |
| `author_id` | String | FK → dat_user |
| `owner_id` | String | FK → dat_user |
| `approval_required` | Boolean | Default false |
| `is_pinned` | Boolean | Default false |
| `is_controlled` | Boolean | Default false |
| `effective_date` | DateTime? | Tanggal berlaku dokumen |
| `expiry_date` | DateTime? | Tanggal kedaluwarsa dokumen |
| `review_cycle_months` | Int? | Siklus review (bulan) |
| `next_review_date` | DateTime? | Review berikutnya |
| `last_reviewed_by` | String? | FK → dat_user |
| `last_reviewed_at` | DateTime? | Waktu review terakhir |
| `department` | String? | Departemen terkait |
| `classification` | String? | Klasifikasi dokumen |
| `view_count` | Int | Default 0 |
| `last_viewed_at` | DateTime? | |
| `published_at` | DateTime? | |
| `archived_at` | DateTime? | |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@unique([company_id, slug])`, `@@index([company_id])`, `@@index([type_id])`, `@@index([status])`, `@@index([author_id])`, `@@index([owner_id])`, `@@index([is_controlled])`, `@@index([next_review_date])`

### `ict_doc_version` — Version history

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `doc_id` | String | FK → ict_doc |
| `version` | Int | Nomor versi |
| `title` | String | Judul saat itu |
| `content` | String | Konten saat itu |
| `change_note` | String? | Catatan perubahan |
| `author_id` | String | FK → dat_user |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([doc_id, version])`, `@@index([doc_id])`

### `ict_doc_share` — Sharing settings

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `doc_id` | String | FK → ict_doc |
| `user_id` | String? | FK → dat_user |
| `share_token` | String? | Unique token untuk link sharing |
| `permission` | enum | view, comment, edit |
| `expires_at` | DateTime? | Kedaluwarsa link |
| `created_by` | String | FK → dat_user |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([doc_id, user_id])`, `@@unique([share_token])`, `@@index([doc_id])`

### `ict_doc_publish` — Public page publish

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `doc_id` | String | FK → ict_doc |
| `publish_token` | String | Unique token untuk URL publik |
| `custom_slug` | String? | Custom slug |
| `is_enabled` | Boolean | Default true |
| `password` | String? | Password protection (hashed) |
| `allow_indexing` | Boolean | Default false |
| `view_count` | Int | Default 0 |
| `published_at` | DateTime | Default now() |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([doc_id])`, `@@unique([publish_token])`, `@@unique([custom_slug])`

### `ict_doc_todo` — Todo list items

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `doc_id` | String | FK → ict_doc |
| `parent_id` | String? | FK self-ref |
| `content` | String | Deskripsi tugas |
| `is_checked` | Boolean | Default false |
| `assigned_to` | String? | FK → dat_user |
| `due_date` | DateTime? | Deadline |
| `priority` | enum | low, medium, high, urgent |
| `sort_order` | Int | Default 0 |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([doc_id])`, `@@index([assigned_to])`, `@@index([due_date])`

### `ict_doc_template` — Template metadata

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `doc_id` | String | FK → ict_doc |
| `fields` | Json | Definisi field template |
| `instructions` | String? | Petunjuk pengisian |
| `is_repeatable` | Boolean | Default false |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([doc_id])`

### `ict_doc_approval` — Approval workflow

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `doc_id` | String | FK → ict_doc |
| `signature_type_id` | String? | FK → dat_signature_type |
| `status` | enum | pending, in_progress, approved, rejected |
| `requested_by` | String | FK → dat_user |
| `requested_at` | DateTime | Default now() |
| `completed_at` | DateTime? | |
| `notes` | String? | |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([doc_id])`, `@@index([status])`

### `ict_doc_approval_step` — Detail step approval

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `approval_id` | String | FK → ict_doc_approval |
| `step` | Int | Urutan step |
| `user_id` | String | FK → dat_user |
| `status` | enum | pending, approved, rejected |
| `comment` | String? | |
| `acted_at` | DateTime? | |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([approval_id, step])`, `@@index([user_id])`

### `ict_doc_read_receipt` — Tanda baca SOP (Acknowledgment)

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `doc_id` | String | FK → ict_doc |
| `user_id` | String | FK → dat_user |
| `version_read` | Int | Versi yang dibaca |
| `read_at` | DateTime | Default now() |
| `acknowledged` | Boolean | Default false |
| `acknowledged_at` | DateTime? | Waktu konfirmasi baca |
| `ip_address` | String? | IP client |
| `user_agent` | String? | User agent |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([doc_id, user_id])`, `@@index([doc_id])`, `@@index([user_id])`, `@@index([acknowledged])`

### Enum

```prisma
enum document_status {
  draft
  in_review
  approved
  published
  archived
}

enum document_visibility {
  private
  internal
  public
}

enum document_format {
  markdown
  html
  json
}

enum document_permission {
  view
  comment
  edit
}

enum document_priority {
  low
  medium
  high
  urgent
}

enum document_approval {
  pending
  in_progress
  approved
  rejected
}
```

---

## 13. Service Request Management

> File: `ict_base/prisma/schema/ict_service_request.prisma`

### `ict_form_category` — Kategori layanan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `parent_id` | String? | FK self-ref |
| `code` | String | Unique per company |
| `name` | String | Nama kategori |
| `description` | String? | |
| `icon` | String? | Icon untuk UI |
| `sort_order` | Int | Default 0 |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`

### `ict_form_service` — Katalog layanan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `category_id` | String | FK → ict_form_category |
| `code` | String | Unique per company |
| `name` | String | Nama layanan |
| `description` | String? | Deskripsi layanan |
| `form_schema` | Json | Definisi form fields |
| `approval_type` | enum | none, single, multi, manager, custom |
| `sla_hours` | Int? | Target penyelesaian (jam) |
| `estimated_cost` | Decimal? | Estimasi biaya |
| `required_role` | String? | Role yang boleh request |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([category_id])`

### `ict_form_request` — Data permintaan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `service_id` | String | FK → ict_form_service |
| `requester_id` | String | FK → dat_user |
| `code` | String | Unique per company, format: SR-YYYYMMDD-XXXX |
| `title` | String | Judul permintaan |
| `description` | String? | Deskripsi tambahan |
| `form_data` | Json | Data form yang diisi |
| `priority` | enum | low, medium, high, urgent |
| `status` | enum | draft, submitted, approved, rejected, in_progress, fulfilled, closed, cancelled |
| `assigned_to` | String? | FK → dat_user |
| `approved_by` | String? | FK → dat_user |
| `approved_at` | DateTime? | |
| `reject_reason` | String? | Alasan penolakan |
| `fulfilled_at` | DateTime? | Waktu selesai |
| `closed_at` | DateTime? | Waktu ditutup |
| `feedback_rating` | Int? | Rating 1-5 |
| `feedback_comment` | String? | |
| `due_date` | DateTime? | Deadline |
| `started_at` | DateTime? | Mulai dikerjakan |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([requester_id])`, `@@index([assigned_to])`, `@@index([status])`, `@@index([service_id])`, `@@index([created_at])`

### `ict_form_asset_request` — Detail request aset

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `request_id` | String | FK → ict_form_request |
| `asset_type` | String | laptop, desktop, monitor, printer, phone, tablet, accessories |
| `brand` | String? | Merek |
| `specification` | String? | Spesifikasi yang diminta |
| `quantity` | Int | Default 1 |
| `estimated_price` | Decimal? | Estimasi harga satuan |
| `justification` | String? | Alasan kebutuhan |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([request_id])`

### `ict_form_server_request` — Detail request server

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `request_id` | String | FK → ict_form_request |
| `server_type` | String | physical, vm, container, cloud |
| `os_type` | String | linux, windows |
| `os_version` | String? | Versi OS |
| `cpu_cores` | Int | Jumlah CPU core |
| `ram_gb` | Int | RAM dalam GB |
| `storage_gb` | Int | Storage dalam GB |
| `storage_type` | String? | ssd, hdd, nvme |
| `purpose` | String? | Tujuan server |
| `software_needed` | String? | Software yang dibutuhkan |
| `estimated_cost` | Decimal? | |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([request_id])`

### `ict_form_access_request` — Detail request akses

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `request_id` | String | FK → ict_form_request |
| `access_type` | String | vpn, ssh, rdp, database, application, api, firewall |
| `target_system` | String | Sistem yang diminta akses |
| `target_host` | String? | IP/hostname |
| `target_port` | Int? | Port |
| `target_database` | String? | Nama database |
| `target_application` | String? | Nama aplikasi |
| `access_level` | String | read_only, read_write, admin, superadmin |
| `duration_type` | String | permanent, temporary |
| `duration_days` | Int? | Durasi (hari) |
| `justification` | String? | Alasan |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([request_id])`

### `ict_form_app_request` — Detail request aplikasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `request_id` | String | FK → ict_form_request |
| `app_type` | String | new_development, modification, integration, license |
| `app_name` | String | Nama aplikasi |
| `app_description` | String? | Deskripsi aplikasi |
| `business_requirement` | String? | Kebutuhan bisnis |
| `target_users` | String? | Target pengguna |
| `estimated_timeline` | String? | Estimasi waktu |
| `budget_code` | String? | Kode budget |
| `priority_level` | String? | Urgency dari bisnis |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([request_id])`

### `ict_form_fulfillment` — Log pelaksanaan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `request_id` | String | FK → ict_form_request |
| `user_id` | String | FK → dat_user |
| `action` | String | started, in_progress, completed, blocked |
| `note` | String? | Catatan pengerjaan |
| `attachment_url` | String? | Bukti pengerjaan |
| `duration_minutes` | Int? | Durasi pengerjaan (menit) |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([request_id])`

### Enum

```prisma
enum form_approval {
  none
  single
  multi
  manager
  custom
}

enum form_priority {
  low
  medium
  high
  urgent
}

enum form_status {
  draft
  submitted
  approved
  rejected
  in_progress
  fulfilled
  closed
  cancelled
}
```

### SLA Matrix

| Priority | Approval | Fulfillment | Total |
|----------|----------|-------------|-------|
| Urgent | 1 jam | 4 jam | 5 jam |
| High | 4 jam | 24 jam | 28 jam |
| Medium | 24 jam | 72 jam | 96 jam |
| Low | 72 jam | 168 jam | 240 jam |

---

## 14. Billing & Invoice

> File: `ict_base/prisma/schema/ict_billing.prisma`

### `ict_bill_vendor` — Vendor/Supplier

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Unique per company |
| `name` | String | Nama vendor |
| `contact_person` | String? | Nama PIC vendor |
| `email` | String? | Email vendor |
| `phone` | String? | Telepon vendor |
| `address` | String? | Alamat vendor |
| `payment_terms` | String? | NET 30, NET 60, COD |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`

### `ict_bill_category` — Kategori tagihan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `parent_id` | String? | FK self-ref |
| `code` | String | Unique per company |
| `name` | String | Nama kategori |
| `description` | String? | |
| `gl_account` | String? | Akun GL |
| `budget_code` | String? | Kode budget |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`

### `ict_bill_invoice` — Invoice/Tagihan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `vendor_id` | String | FK → ict_bill_vendor |
| `category_id` | String | FK → ict_bill_category |
| `location_id` | String? | FK → dat_location |
| `sr_request_id` | String? | FK → ict_sr_request |
| `code` | String | Unique per company, format: INV-YYYYMMDD-XXXX |
| `vendor_invoice_no` | String? | Nomor invoice dari vendor |
| `title` | String | Judul/Deskripsi tagihan |
| `description` | String? | Detail tambahan |
| `invoice_date` | DateTime | Tanggal invoice |
| `due_date` | DateTime | Jatuh tempo |
| `total_amount` | Decimal | Total yang harus dibayar |
| `currency` | String | Default: IDR |
| `status` | enum | draft, submitted, approved, rejected, paid, overdue, cancelled |
| `priority` | enum | low, medium, high, urgent |
| `approved_by` | String? | FK → dat_user |
| `approved_at` | DateTime? | |
| `reject_reason` | String? | Alasan penolakan |
| `paid_at` | DateTime? | Waktu pembayaran dikonfirmasi |
| `notes` | String? | Catatan |
| `tags` | String[] | Tags fleksibel |
| `metadata` | Json? | Data tambahan |
| `created_by` | String | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([vendor_id])`, `@@index([status])`, `@@index([due_date])`

### `ict_bill_invoice_item` — Detail item tagihan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `invoice_id` | String | FK → ict_bill_invoice |
| `description` | String | Deskripsi item |
| `quantity` | Decimal | Default 1 |
| `unit_price` | Decimal | Harga satuan |
| `amount` | Decimal | quantity × unit_price |
| `sort_order` | Int | Default 0 |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([invoice_id])`

### Enum

```prisma
enum invoice_status {
  draft
  submitted
  approved
  rejected
  paid
  overdue
  cancelled
}

enum billing_priority {
  low
  medium
  high
  urgent
}
```

---

## 15. Preventive Maintenance

> File: `ict_base/prisma/schema/ict_preventive.prisma`

### `ict_pm_equipment_type` — Tipe perangkat

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Unique per company |
| `name` | String | Nama tipe perangkat |
| `category` | enum | wifi, zoom, other |
| `description` | String? | |
| `icon` | String? | Icon untuk UI |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`

### `ict_pm_room` — Ruangan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `location_id` | String | FK → dat_location |
| `code` | String | Unique per company |
| `name` | String | Nama ruangan |
| `floor` | String? | Lantai |
| `capacity` | Int? | Kapasitas orang |
| `has_wifi` | Boolean | Default false |
| `has_zoom` | Boolean | Default false |
| `description` | String? | |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([location_id])`

### `ict_pm_equipment` — Perangkat di ruangan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `room_id` | String | FK → ict_pm_room |
| `equipment_type_id` | String | FK → ict_pm_equipment_type |
| `asset_id` | String? | FK → ict_asset |
| `code` | String | Unique per company |
| `name` | String | Nama perangkat |
| `brand` | String? | Merek |
| `model` | String? | Model |
| `serial_number` | String? | Nomor serial |
| `ip_address` | String? | IP address |
| `install_date` | DateTime? | Tanggal instalasi |
| `warranty_end` | DateTime? | Akhir garansi |
| `status` | enum | active, inactive, maintenance, damaged |
| `notes` | String? | Catatan |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([room_id])`, `@@index([status])`

### `ict_pm_checklist_template` — Template checklist harian

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `equipment_type_id` | String | FK → ict_pm_equipment_type |
| `name` | String | Nama template |
| `description` | String? | |
| `items` | Json | Array checklist items |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, name])`, `@@index([company_id])`, `@@index([equipment_type_id])`

### `ict_pm_schedule` — Jadwal pengecekan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `room_id` | String | FK → ict_pm_room |
| `template_id` | String | FK → ict_pm_checklist_template |
| `name` | String | Nama jadwal |
| `frequency` | enum | daily, weekly, biweekly, monthly, quarterly |
| `day_of_week` | Int? | 0-6 (Minggu-Sabtu) |
| `day_of_month` | Int? | 1-31 |
| `time_start` | String | Waktu mulai (HH:MM) |
| `time_end` | String? | Waktu selesai (HH:MM) |
| `assigned_to` | String? | FK → dat_user |
| `next_check_date` | DateTime? | Tanggal pengecekan berikutnya |
| `last_check_date` | DateTime? | Tanggal pengecekan terakhir |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, name])`, `@@index([company_id])`, `@@index([room_id])`, `@@index([next_check_date])`

### `ict_pm_checklist` — Checklist yang sudah diisi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `schedule_id` | String | FK → ict_pm_schedule |
| `equipment_id` | String | FK → ict_pm_equipment |
| `room_id` | String | FK → ict_pm_room |
| `check_date` | DateTime | Tanggal pengecekan |
| `status` | enum | pending, in_progress, completed, overdue |
| `result` | enum? | pass, fail, needs_repair |
| `items` | Json | Hasil checklist |
| `notes` | String? | Catatan umum |
| `started_at` | DateTime? | Waktu mulai |
| `completed_at` | DateTime? | Waktu selesai |
| `checked_by` | String? | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, check_date, equipment_id])`, `@@index([company_id])`, `@@index([schedule_id])`, `@@index([check_date])`, `@@index([status])`

### `ict_pm_issue` — Masalah/kerusakan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `checklist_id` | String? | FK → ict_pm_checklist |
| `equipment_id` | String | FK → ict_pm_equipment |
| `room_id` | String | FK → ict_pm_room |
| `title` | String | Judul masalah |
| `description` | String? | Deskripsi masalah |
| `severity` | enum | Menggunakan `severity_level` (shared) |
| `status` | enum | open, in_progress, resolved, escalated |
| `assigned_to` | String? | FK → dat_user |
| `resolved_by` | String? | FK → dat_user |
| `resolved_at` | DateTime? | |
| `resolution_notes` | String? | |
| `sr_request_id` | String? | FK → ict_sr_request |
| `created_by` | String | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([company_id])`, `@@index([equipment_id])`, `@@index([room_id])`, `@@index([status])`, `@@index([severity])`

### `ict_pm_calendar_event` — Event kalender

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `schedule_id` | String? | FK → ict_pm_schedule |
| `checklist_id` | String? | FK → ict_pm_checklist |
| `title` | String | Judul event |
| `description` | String? | |
| `event_date` | DateTime | Tanggal event |
| `time_start` | String | Waktu mulai (HH:MM) |
| `time_end` | String? | Waktu selesai (HH:MM) |
| `event_type` | enum | scheduled, completed, overdue, cancelled |
| `room_id` | String? | FK → ict_pm_room |
| `assigned_to` | String? | FK → dat_user |
| `color` | String? | Warna event (hex) |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([company_id])`, `@@index([event_date])`, `@@index([event_type])`, `@@index([assigned_to])`

### Enum

```prisma
enum equipment_category {
  wifi
  zoom
  other
}

enum equipment_status {
  active
  inactive
  maintenance
  damaged
}

enum frequency {
  daily
  weekly
  biweekly
  monthly
  quarterly
}

enum checklist_status {
  pending
  in_progress
  completed
  overdue
}

enum checklist_result {
  pass
  fail
  needs_repair
}

enum issue_status {
  open
  in_progress
  resolved
  escalated
}

enum event_type {
  scheduled
  completed
  overdue
  cancelled
}
```

---

## 16. Dashboard & Reporting

> File: `ict_base/prisma/schema/ict_dashboard.prisma`

### Shared Enum

#### `report_category` — Kategori laporan

```prisma
enum report_category {
  general
  iso_compliance
  sop_execution
  security
  monitoring
  incident
  vulnerability
  audit
  risk
  management_review
}
```

### `ict_dash_widget` — Konfigurasi widget

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `code` | String | Unique: asset_summary, vuln_by_severity, comp_score, risk_matrix, cert_expiring, incident_trend |
| `name` | String | Nama widget |
| `type` | String | chart, stat, table, map, heatmap |
| `module_source` | String | asset, vulnerability, compliance, incident, risk, certificate, siem |
| `data_query` | String | Query/data source |
| `config` | Json | Konfigurasi default |
| `description` | String? | |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

### `ict_dash_layout` — Layout widget per user

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `user_id` | String | FK → dat_user |
| `widget_id` | String | FK → ict_dash_widget |
| `position_x` | Int | Posisi X (grid) |
| `position_y` | Int | Posisi Y (grid) |
| `width` | Int | Lebar (grid columns) |
| `height` | Int | Tinggi (grid rows) |
| `config_override` | Json? | Override config per user |
| `is_visible` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([user_id, widget_id])`

### `ict_report_template` — Template laporan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Unique per company |
| `name` | String | Nama laporan |
| `description` | String? | |
| `category` | report_category | Default "general" |
| `module_source` | String | Sumber data modul |
| `template_config` | Json | Layout/format laporan |
| `schedule` | String? | Legacy schedule |
| `schedule_cron` | String? | Cron expression untuk penjadwalan |
| `recipients` | Json? | Array of user_ids |
| `format` | String | Default "pdf" — pdf, xlsx, csv |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([category])`

### `ict_report_instance` — Instance laporan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `template_id` | String | FK → ict_report_template |
| `company_id` | String | FK → dat_company |
| `generated_by` | String | FK → dat_user atau "system" |
| `status` | String | generating, completed, failed |
| `file_path` | String? | Path file laporan |
| `file_type` | String | pdf, xlsx, csv |
| `file_size` | Int? | Ukuran byte |
| `parameters` | Json? | Parameter filter |
| `error_message` | String? | Jika gagal |
| `generation_time` | Int? | Waktu generasi dalam milidetik |
| `generated_at` | DateTime | Default now() |

**Index:** `@@index([template_id])`, `@@index([company_id])`, `@@index([status])`, `@@index([generated_at])`

### `ict_report_schedule` — Penjadwalan laporan otomatis

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `template_id` | String | FK → ict_report_template |
| `company_id` | String | FK → dat_company |
| `name` | String | Nama jadwal |
| `cron_expression` | String | Cron expression untuk penjadwalan |
| `parameters` | Json? | Parameter default untuk laporan |
| `next_run_at` | DateTime? | Waktu berikutnya dijalankan |
| `last_run_at` | DateTime? | Waktu terakhir dijalankan |
| `is_enabled` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([template_id])`, `@@index([company_id])`, `@@index([next_run_at])`, `@@index([is_enabled])`

### `ict_dash_bookmark` — Bookmark dashboard user

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `user_id` | String | FK → dat_user |
| `name` | String | Nama bookmark |
| `filters` | Json | Filter aktif |
| `is_default` | Boolean | Default false |
| `created_at` | DateTime | Default now() |

---

## 17. Web Security (Nginx/WAF)

> File: `ict_base/prisma/schema/ict_web_security.prisma`

### `ict_nginx_log` — Log akses nginx per-domain

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `timestamp` | DateTime | Default now() |
| `host` | String | Host server |
| `server_ip` | String | IP server |
| `client_ip` | String | IP client |
| `country_iso` | String? | Kode negara ISO |
| `xff` | String? | X-Forwarded-For header |
| `traffic_type` | String | Tipe traffic |
| `domain` | String | Domain target |
| `url` | String | URL yang diakses |
| `referer` | String? | Referer header |
| `args` | String? | Query parameters |
| `upstreamtime` | String? | Waktu upstream |
| `responsetime` | String? | Waktu respons |
| `request_method` | String | HTTP method |
| `status` | String | HTTP status code |
| `size` | String? | Response size |
| `request_body` | String? | Request body |
| `request_length` | Int? | Request length |
| `protocol` | String? | HTTP protocol |
| `upstreamhost` | String? | Upstream host |
| `file_dir` | String? | File directory |
| `http_user_agent` | String? | User agent |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([client_ip])`, `@@index([country_iso])`, `@@index([timestamp])`

### `ict_nginx_app` — Log aplikasi nginx

Struktur kolom sama dengan `ict_nginx_log`.

**Index:** `@@index([client_ip])`, `@@index([timestamp])`, `@@index([country_iso])`

### `ict_nginx_atc` — Log ATC (Attack Traffic Counter)

Struktur kolom sama dengan `ict_nginx_log`.

**Index:** `@@index([client_ip])`, `@@index([timestamp])`, `@@index([country_iso])`

### `ict_nginx_atc_sum` — Ringkasan ATC per hari

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `date` | String | Tanggal (YYYY-MM-DD) |
| `client_ip` | String | IP client |
| `traffic_type` | String | Tipe traffic |
| `target_domain` | String | Domain target |
| `total_hits` | BigInt | Default 1 |
| `last_seen` | DateTime | Default now() |

**Index:** `@@unique([date, client_ip, traffic_type, target_domain])`

### `ict_nginx_sla` — SLA nginx per hari

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `date` | String | Unique, tanggal (YYYY-MM-DD) |
| `total_requests` | BigInt | Default 0 |
| `successful_requests` | BigInt | Default 0 |
| `client_errors` | BigInt | Default 0 |
| `server_errors` | BigInt | Default 0 |
| `attack_requests` | BigInt | Default 0 |
| `avg_response_time` | Decimal | Default 0.0000 |
| `sla_percentage` | Decimal | Default 0.00 |

### `ict_ip_whitelist` — IP whitelist

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `ip_or_cidr` | String | Unique, IP atau CIDR |
| `description` | String? | Deskripsi |
| `created_at` | DateTime | Default now() |

### `ict_ip_blacklist` — IP blacklist

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `ip` | String | Unique, IP address |
| `threat_score` | Int | Skor ancaman |
| `reason` | String | Alasan blacklist |
| `banned_at` | DateTime | Default now() |
| `expires_at` | DateTime? | Waktu kedaluwarsa |

### `ict_waf_bypass_rule` — Rule bypass WAF

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `domain` | String | Domain target |
| `url_path` | String | Path URL |
| `args_pattern` | String? | Pola query args |
| `description` | String? | Deskripsi rule |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([domain, url_path])`

---

## 18. Web Monitoring (UptimeRobot)

> File: `ict_base/prisma/schema/ict_web_monitor.prisma`

### `ict_uptimerobot_log` — Log incident UptimeRobot

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `incident_id` | String | ID incident UptimeRobot |
| `monitor_id` | BigInt | ID monitor |
| `monitor_url` | String | URL yang dimonitor |
| `monitor_name` | String | Nama monitor |
| `monitor_type` | String | Tipe monitor |
| `monitor_group` | String? | Grup monitor |
| `http_status_code` | Int | Default 0 |
| `alert_type` | Int | Tipe alert |
| `alert_type_name` | String | Nama tipe alert |
| `alert_region` | String? | Region alert |
| `alert_location` | String? | Lokasi alert |
| `alert_details` | String? | Detail alert |
| `alert_duration` | Int | Default 0 (detik) |
| `alert_datetime` | DateTime | Waktu alert |
| `incident_start_time` | DateTime? | Waktu mulai incident |
| `incident_end_time` | DateTime? | Waktu selesai incident |
| `response_time` | Int | Default 0 (ms) |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([monitor_id, alert_datetime])`, `@@index([alert_datetime])`, `@@index([incident_id])`

### `ict_uptimerobot_sum` — Ringkasan harian per monitor

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `date` | String | Tanggal (YYYY-MM-DD) |
| `monitor_id` | BigInt | ID monitor |
| `monitor_name` | String | Nama monitor |
| `monitor_url` | String | URL yang dimonitor |
| `total_checks` | BigInt | Default 0 |
| `up_count` | BigInt | Default 0 |
| `down_count` | BigInt | Default 0 |
| `downtime_seconds` | BigInt | Default 0 |
| `avg_response_time` | Decimal | Default 0.0000 |
| `uptime_percentage` | Decimal | Default 0.00 |
| `last_seen` | DateTime | Default now() |

**Index:** `@@unique([date, monitor_id])`, `@@index([monitor_id])`

### `ict_uptimerobot_sla` — SLA harian global

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `date` | String | Unique, tanggal (YYYY-MM-DD) |
| `total_checks` | BigInt | Default 0 |
| `up_count` | BigInt | Default 0 |
| `down_count` | BigInt | Default 0 |
| `downtime_seconds` | BigInt | Default 0 |
| `avg_response_time` | Decimal | Default 0.0000 |
| `overall_uptime` | Decimal | Default 0.00 |

---

## 19. Asset Ownership & Assignment

> File: `ict_base/prisma/schema/ict_asset.prisma`
> Status: **BARU**

### `ict_asset_assignment` — Riwayat peminjaman/pengembalian aset

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `asset_id` | String | FK → ict_asset |
| `assigned_to` | String | FK → dat_user (penerima aset) |
| `assigned_by` | String? | FK → dat_user (yang menyerahkan) |
| `assignment_date` | DateTime | Tanggal penugasan |
| `expected_return` | DateTime? | Tanggal pengembalian yang diharapkan |
| `actual_return` | DateTime? | Tanggal pengembalian aktual |
| `status` | enum | assigned, returned, pending_return, lost, damaged |
| `condition_out` | enum | excellent, good, fair, poor, damaged, non_functional (kondisi saat diserahkan) |
| `condition_in` | enum? | Kondisi saat dikembalikan |
| `checkout_notes` | String? | Catatan saat checkout |
| `return_notes` | String? | Catatan saat return |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([asset_id, assigned_to, assignment_date])`, `@@index([company_id])`, `@@index([asset_id])`, `@@index([assigned_to])`, `@@index([status])`, `@@index([assignment_date])`

### Enum

```prisma
enum assignment_status {
  assigned
  returned
  pending_return
  lost
  damaged
}

enum assignment_condition {
  excellent
  good
  fair
  poor
  damaged
  non_functional
}
```

### Relasi

- `ict_asset_assignment.asset_id` → `ict_asset.id`
- `ict_asset_assignment.assigned_to` → `dat_user.id`
- `ict_asset_assignment.assigned_by` → `dat_user.id`

---

## 20. Network Management (VLAN/Subnet/Route)

> File: `ict_base/prisma/schema/ict_network.prisma`
> Status: **BARU**

### `ict_network_vlan` — Definisi VLAN

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `location_id` | String? | FK → dat_location |
| `vlan_id` | Int | Nomor VLAN (1-4094) |
| `name` | String | Nama VLAN |
| `description` | String? | Deskripsi |
| `subnet_cidr` | String? | Subnet CIDR |
| `gateway` | String? | Gateway IP |
| `dns_primary` | String? | DNS primary |
| `dns_secondary` | String? | DNS secondary |
| `status` | enum | active, inactive, maintenance |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, vlan_id])`, `@@index([company_id])`, `@@index([location_id])`

### `ict_network_subnet` — Subnet allocation

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `vlan_id` | String? | FK → ict_network_vlan |
| `cidr` | String | CIDR notation (10.0.0.0/24) |
| `name` | String | Nama subnet |
| `description` | String? | Deskripsi |
| `gateway` | String | Gateway IP |
| `subnet_mask` | String | Subnet mask |
| `broadcast` | String? | Broadcast address |
| `usable_start` | String? | Range IP usable mulai |
| `usable_end` | String? | Range IP usable akhir |
| `total_addresses` | Int | Jumlah total alamat |
| `used_addresses` | Int? | Jumlah alamat terpakai |
| `status` | String | Default: "active" |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, cidr])`, `@@index([company_id])`, `@@index([vlan_id])`, `@@index([cidr])`

### `ict_network_interface` — Interface device

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `asset_id` | String? | FK → ict_asset |
| `device_name` | String | Nama device (hostname/IP) |
| `interface_name` | String | Nama interface (ether1, bridge, vlan100) |
| `interface_type` | enum | ethernet, bridge, vlan, bond, tunnel, wireless, pppoe, loopback |
| `mac_address` | String? | MAC address |
| `ip_address` | String? | IP address |
| `subnet_id` | String? | FK → ict_network_subnet |
| `vlan_id` | String? | FK → ict_network_vlan |
| `speed_mbps` | Int? | Kecepatan Mbps |
| `status` | enum | up, down, disabled |
| `mtu` | Int? | MTU value |
| `description` | String? | Deskripsi |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, device_name, interface_name])`, `@@index([company_id])`, `@@index([asset_id])`, `@@index([vlan_id])`

### `ict_network_route` — Routing table

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `destination_cidr` | String | Destination network/CIDR |
| `gateway` | String? | Gateway IP |
| `interface_id` | String? | FK → ict_network_interface (outgoing) |
| `distance` | Int | Default 1 |
| `metric` | Int? | Route metric |
| `scope` | enum | local, connected, static, dynamic, bgp |
| `status` | enum | active, inactive, disabled |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, destination_cidr, gateway, interface_id])`, `@@index([company_id])`, `@@index([destination_cidr])`, `@@index([gateway])`, `@@index([interface_id])`

### `ict_network_nat` — NAT rules

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `name` | String | Nama rule |
| `description` | String? | Deskripsi |
| `nat_type` | enum | src_nat, dst_nat, masquerade, static_nat |
| `public_ip` | String | IP public |
| `public_port` | Int? | Port public |
| `private_ip` | String | IP private (lokasi) |
| `private_port` | Int? | Port private |
| `protocol` | String? | Default: "tcp" |
| `interface` | String? | Interface terkait |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, public_ip, public_port, private_ip, private_port])`, `@@index([company_id])`, `@@index([public_ip])`, `@@index([private_ip])`, `@@index([nat_type])`

### `ict_network_dhcp_pool` — DHCP pool

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `subnet_id` | String? | FK → ict_network_subnet |
| `name` | String | Nama pool |
| `interface` | String? | Interface |
| `range_start` | String | IP mulai range |
| `range_end` | String | IP akhir range |
| `gateway` | String? | Gateway |
| `dns_primary` | String? | DNS primary |
| `dns_secondary` | String? | DNS secondary |
| `lease_time` | Int? | Lease time (detik) |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, name])`, `@@index([company_id])`, `@@index([subnet_id])`

### Enum

```prisma
enum vlan_status {
  active
  inactive
  maintenance
}

enum interface_type {
  ethernet
  bridge
  vlan
  bond
  tunnel
  wireless
  pppoe
  loopback
}

enum interface_status {
  up
  down
  disabled
}

enum route_scope {
  local
  connected
  static
  dynamic
  bgp
}

enum route_status {
  active
  inactive
  disabled
}

enum nat_type {
  src_nat
  dst_nat
  masquerade
  static_nat
}
```

---

## 21. Mikrotik Device Management

> File: `ict_base/prisma/schema/ict_mikrotik.prisma`
> Status: **BARU**

### `ict_mikrotik_device` — Data Mikrotik utama

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `location_id` | String? | FK → dat_location |
| `asset_id` | String? | FK → ict_asset |
| `name` | String | Nama device |
| `hostname` | String | Hostname/IP Mikrotik |
| `mac_address` | String? | MAC address |
| `ip_address` | String | IP address management |
| `port` | Int | Default 8728 (API port) |
| `routeros_version` | String? | Versi RouterOS |
| `board_name` | String? | Board name (hAP ac2, CCR1009, dll) |
| `identity` | String? | Identity Mikrotik |
| `uptime` | String? | Uptime |
| `status` | enum | connected, disconnected, error |
| `last_seen_at` | DateTime? | Terakhir terhubung |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, name])`, `@@index([company_id])`, `@@index([location_id])`, `@@index([asset_id])`, `@@index([status])`

### `ict_mikrotik_hotspot_profile` — Hotspot server profile

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_mikrotik_device |
| `company_id` | String | FK → dat_company |
| `name` | String | Nama profile |
| `interface` | String | Interface (wlan1, ether1) |
| `address_pool` | String? | DHCP address pool |
| `auth_method` | enum | chap, pap, mschap, none |
| `session_timeout` | Int? | Session timeout (detik) |
| `idle_timeout` | Int? | Idle timeout (detik) |
| `mac_auth` | Boolean | Default false |
| `https_redirect` | Boolean | Default false |
| `dns_name` | String? | DNS name hotspot |
| `split_dns` | Boolean | Default false |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([device_id, name])`, `@@index([device_id])`

### `ict_mikrotik_hotspot_user` — User hotspot

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_mikrotik_device |
| `profile_id` | String? | FK → ict_mikrotik_hotspot_profile |
| `company_id` | String | FK → dat_company |
| `username` | String | Username hotspot |
| `password_hash` | String? | Password hash |
| `server` | String? | Server name |
| `uptime` | String? | Total uptime |
| `bytes_in` | BigInt? | Total bytes received |
| `bytes_out` | BigInt? | Total bytes sent |
| `bytes_total` | BigInt? | Total bytes |
| `packets_in` | BigInt? | Total packets in |
| `packets_out` | BigInt? | Total packets out |
| `last_login_at` | DateTime? | Terakhir login |
| `status` | enum | active, disabled |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([device_id, username])`, `@@index([device_id])`, `@@index([profile_id])`, `@@index([status])`

### `ict_mikrotik_dhcp_lease` — DHCP lease aktif

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_mikrotik_device |
| `company_id` | String | FK → dat_company |
| `mac_address` | String | MAC address |
| `address` | String | IP address |
| `hostname` | String? | Hostname client |
| `server_name` | String? | Server name |
| `status` | String | Default: "bound" |
| `expires_at` | DateTime? | Waktu expiry lease |
| `uptime` | String? | Uptime lease |
| `bytes_in` | BigInt? | Bytes received |
| `bytes_out` | BigInt? | Bytes sent |
| `last_active_at` | DateTime? | Terakhir aktif |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([device_id, mac_address, address])`, `@@index([device_id])`, `@@index([mac_address])`, `@@index([address])`, `@@index([status])`

### `ict_mikrotik_firewall_rule` — Firewall filter/NAT rules

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_mikrotik_device |
| `company_id` | String | FK → dat_company |
| `chain` | enum | input, forward, output |
| `action` | enum | accept, drop, reject, tarpit, passthrough |
| `protocol` | enum? | tcp, udp, icmp, ip, all |
| `src_address` | String? | Source address |
| `dst_address` | String? | Destination address |
| `src_port` | String? | Source port |
| `dst_port` | String? | Destination port |
| `in_interface` | String? | In interface |
| `out_interface` | String? | Out interface |
| `src_mac_address` | String? | Source MAC |
| `connection_state` | String? | Connection state |
| `connection_nat` | String? | NAT connection |
| `comment` | String? | Rule comment |
| `bytes_counter` | BigInt? | Bytes counter |
| `packets_counter` | BigInt? | Packets counter |
| `disabled` | Boolean | Default false |
| `log` | Boolean | Default false |
| `log_prefix` | String? | Log prefix |
| `place_before` | String? | Place before rule |
| `rule_order` | Int | Default 0 |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([device_id])`, `@@index([chain])`, `@@index([action])`, `@@index([disabled])`

### `ict_mikrotik_queue` — Bandwidth queue

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_mikrotik_device |
| `company_id` | String | FK → dat_company |
| `name` | String | Nama queue |
| `target` | String | Target IP/host |
| `parent` | String? | Parent queue |
| `queue_type` | enum | simple, simple_queue, queue_tree |
| `max_limit` | String? | Max bandwidth (512k/10m) |
| `min_limit` | String? | Min bandwidth |
| `burst_limit` | String? | Burst limit |
| `burst_threshold` | String? | Burst threshold |
| `burst_time` | Int? | Burst time (detik) |
| `priority` | Int | Default 8 (1-8) |
| `bytes_in` | BigInt? | Bytes received |
| `bytes_out` | BigInt? | Bytes sent |
| `packets_in` | BigInt? | Packets in |
| `packets_out` | BigInt? | Packets out |
| `rate_in` | String? | Current rate in |
| `rate_out` | String? | Current rate out |
| `disabled` | Boolean | Default false |
| `comment` | String? | Comment |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([device_id, name])`, `@@index([device_id])`, `@@index([target])`, `@@index([disabled])`

### Enum

```prisma
enum mikrotik_status {
  connected
  disconnected
  error
}

enum hotspot_auth_type {
  chap
  pap
  mschap
  none
}

enum hotspot_status {
  active
  disabled
}

enum queue_type {
  simple
  simple_queue
  queue_tree
}

enum firewall_action {
  accept
  drop
  reject
  tarpit
  passthrough
}

enum firewall_chain {
  input
  forward
  output
}

enum firewall_protocol {
  tcp
  udp
  icmp
  ip
  all
}
```

---

## 22. Domain Management

> File: `ict_base/prisma/schema/ict_domain.prisma`
> Status: **BARU**

### `ict_domain` — Domain registry

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `location_id` | String? | FK → dat_location |
| `asset_id` | String? | FK → ict_asset (server hosting) |
| `domain_name` | String | Nama domain |
| `registrar` | String? | Registrar (GoDaddy, Cloudflare, dll) |
| `registration_date` | DateTime? | Tanggal registrasi |
| `expiry_date` | DateTime? | Tanggal kedaluwarsa |
| `auto_renew` | Boolean | Default false |
| `nameservers` | String[] | Array nameservers |
| `registrant_name` | String? | Nama registrant |
| `registrant_email` | String? | Email registrant |
| `registrant_org` | String? | Organisasi registrant |
| `admin_contact` | String? | Admin contact |
| `tech_contact` | String? | Tech contact |
| `whois_server` | String? | WHOIS server |
| `status` | enum | active, expired, suspended, pending_renewal, transferring, parked |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, domain_name])`, `@@index([company_id])`, `@@index([domain_name])`, `@@index([status])`, `@@index([expiry_date])`

### `ict_domain_dns` — DNS records

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `domain_id` | String | FK → ict_domain |
| `company_id` | String | FK → dat_company |
| `record_type` | enum | a, aaaa, cname, mx, txt, srv, ns, ptr, caa, soa, other |
| `name` | String | Record name (subdomain atau @) |
| `value` | String | Record value |
| `ttl` | Int | Default 3600 |
| `priority` | Int? | Priority (MX, SRV) |
| `status` | enum | active, inactive, pending |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([domain_id, record_type, name, value])`, `@@index([domain_id])`, `@@index([record_type])`, `@@index([name])`

### `ict_domain_nameserver` — Nameserver entries

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `domain_id` | String | FK → ict_domain |
| `company_id` | String | FK → dat_company |
| `hostname` | String | Nameserver hostname |
| `ip_address` | String? | IP address nameserver |
| `description` | String? | Deskripsi |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([domain_id, hostname])`, `@@index([domain_id])`

### Enum

```prisma
enum domain_status {
  active
  expired
  suspended
  pending_renewal
  transferring
  parked
}

enum dns_record_type {
  a
  aaaa
  cname
  mx
  txt
  srv
  ns
  ptr
  caa
  soa
  other
}

enum dns_record_status {
  active
  inactive
  pending
}
```

### Relasi

- `ict_domain` → `ict_certificate` (via domain_name)
- `ict_domain` → `ict_asset` (hosting server)
- `ict_domain_dns` → `ict_domain` (DNS records per domain)

---

## 23. IP Public & Block Management

> File: `ict_base/prisma/schema/ict_ip_public.prisma`
> Status: **BARU**

### `ict_ip_block` — IP block/CIDR yang dimiliki

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `cidr` | String | CIDR notation (203.0.113.0/24) |
| `total_ips` | Int | Jumlah total IP |
| `used_ips` | Int? | Jumlah IP terpakai |
| `asn` | String? | ASN number |
| `isp` | String? | ISP provider |
| `registry` | String? | APNIC, RIPE, ARIN, dll |
| `country_code` | String? | Kode negara |
| `announcement_date` | DateTime? | Tanggal announcement |
| `expiry_date` | DateTime? | Kedaluwarsa |
| `bgp_as_path` | String? | BGP AS path |
| `description` | String? | Deskripsi |
| `status` | enum | active, reserved, available, expired |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, cidr])`, `@@index([company_id])`, `@@index([cidr])`, `@@index([asn])`, `@@index([status])`

### `ict_ip_allocation` — Alokasi IP dari block

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `block_id` | String | FK → ict_ip_block |
| `company_id` | String | FK → dat_company |
| `ip_address` | String | IP address |
| `vlan_id` | String? | FK → ict_network_vlan |
| `subnet_id` | String? | FK → ict_network_subnet |
| `asset_id` | String? | FK → ict_asset |
| `domain_id` | String? | FK → ict_domain |
| `allocation_type` | enum | static, dynamic, nat, vpn, reserved |
| `purpose` | String? | Tujuan alokasi |
| `description` | String? | Deskripsi |
| `assigned_to` | String? | Entity yang dituju |
| `assigned_at` | DateTime? | Tanggal alokasi |
| `expires_at` | DateTime? | Kedaluwarsa alokasi |
| `status` | enum | assigned, available, reserved, deprecated, quarantined |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([block_id, ip_address])`, `@@index([block_id])`, `@@index([ip_address])`, `@@index([vlan_id])`, `@@index([asset_id])`, `@@index([domain_id])`, `@@index([status])`

### `ict_ip_nat_mapping` — NAT mapping public → private

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `name` | String | Nama mapping |
| `description` | String? | Deskripsi |
| `public_ip` | String | IP public |
| `public_port` | Int? | Port public |
| `private_ip` | String | IP private (lokasi) |
| `private_port` | Int? | Port private |
| `protocol` | String | Default: "tcp" |
| `asset_id` | String? | FK → ict_asset |
| `nat_id` | String? | FK → ict_network_nat |
| `interface` | String? | Interface |
| `is_persistent` | Boolean | Default false |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, public_ip, public_port, private_ip, private_port])`, `@@index([company_id])`, `@@index([public_ip])`, `@@index([private_ip])`, `@@index([nat_id])`, `@@index([asset_id])`

### Enum

```prisma
enum ip_block_status {
  active
  reserved
  available
  expired
}

enum ip_allocation_status {
  assigned
  available
  reserved
  deprecated
  quarantined
}

enum allocation_type {
  static
  dynamic
  nat
  vpn
  reserved
}
```

### Relasi

- `ict_ip_block` → `ict_ip_allocation` (1 block → banyak alokasi)
- `ict_ip_allocation` → `ict_network_vlan` (IP dalam VLAN)
- `ict_ip_allocation` → `ict_network_subnet` (IP dalam subnet)
- `ict_ip_allocation` → `ict_asset` (IP untuk aset tertentu)
- `ict_ip_allocation` → `ict_domain` (IP untuk domain)
- `ict_ip_nat_mapping` → `ict_network_nat` (mapping dari rule NAT)

---

## 24. Monitor Agent

> File: `ict_base/prisma/schema/ict_monitor_agent.prisma`

Mengelola deployment, status, dan komunikasi agent di seluruh endpoint. Menggantikan Wazuh Agent dan Beszel Hub/Agent.

### Shared Enum

#### `agent_platform` — Platform operating system agent

```prisma
enum agent_platform {
  linux
  windows
  macos
  docker
  network
}
```

#### `agent_status` — Status operasional agent

```prisma
enum agent_status {
  online
  offline
  error
  updating
  deprecated
}
```

#### `agent_type` — Tipe fungsi agent

```prisma
enum agent_type {
  siem
  host_monitor
  network_poller
  docker_collector
  fim
  combined
}
```

### `ict_monitor_agent_group` — Grup agent

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Kode unik grup |
| `name` | String | Nama grup |
| `description` | String? | Deskripsi grup |
| `platform` | agent_platform | Default "linux" |
| `config` | Json? | Konfigurasi grup |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`

### `ict_monitor_agent` — Agent endpoint

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `group_id` | String | FK → ict_monitor_agent_group |
| `asset_id` | String? | FK → ict_asset |
| `hostname` | String | Nama host |
| `ip_address` | String | IP address agent |
| `os_platform` | agent_platform | Default "linux" |
| `os_version` | String? | Versi OS |
| `agent_version` | String? | Versi agent |
| `status` | agent_status | Default "offline" |
| `agent_type` | agent_type | Default "combined" |
| `public_key` | String? | Public key untuk autentikasi |
| `last_heartbeat_at` | DateTime? | Waktu heartbeat terakhir |
| `last_config_sync` | DateTime? | Waktu sync config terakhir |
| `registration_token` | String? | Token registrasi |
| `config` | Json? | Konfigurasi agent |
| `tags` | String[] | Tag untuk filtering |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, hostname, ip_address])`, `@@index([company_id])`, `@@index([group_id])`, `@@index([asset_id])`, `@@index([status])`, `@@index([agent_type])`, `@@index([last_heartbeat_at])`

### `ict_monitor_poller_task` — Tugas polling agent

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `group_id` | String | FK → ict_monitor_agent_group |
| `task_type` | String | Tipe tugas polling |
| `target_entity` | String? | Tipe entitas target |
| `target_id` | String? | ID entitas target |
| `config` | Json? | Konfigurasi tugas |
| `interval_seconds` | Int | Default 300 — interval dalam detik |
| `last_run_at` | DateTime? | Waktu terakhir dijalankan |
| `next_run_at` | DateTime? | Waktu berikutnya dijalankan |
| `is_enabled` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([company_id])`, `@@index([group_id])`, `@@index([task_type])`, `@@index([next_run_at])`

### `ict_monitor_network_device` — Perangkat jaringan (SNMP)

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `location_id` | String? | FK → ict_location |
| `asset_id` | String? | FK → ict_asset |
| `device_type` | String | Tipe perangkat (router, switch, ap, dll) |
| `hostname` | String | Nama host perangkat |
| `ip_address` | String | IP address perangkat |
| `snmp_version` | String | Default "2c" |
| `snmp_community` | String? | Community string SNMP |
| `snmp_port` | Int | Default 161 |
| `snmp_sysname` | String? | sysName dari SNMP |
| `snmp_sysdesc` | String? | sysDescr dari SNMP |
| `snmp_sysobjectid` | String? | sysObjectID dari SNMP |
| `snmp_uptime` | String? | Uptime dari SNMP |
| `mac_address` | String? | MAC address |
| `vendor` | String? | Vendor perangkat |
| `model` | String? | Model perangkat |
| `serial_number` | String? | Nomor serial |
| `os_fingerprint` | String? | Deteksi OS |
| `status` | String | Default "unknown" |
| `last_polled_at` | DateTime? | Waktu polling terakhir |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, ip_address])`, `@@index([company_id])`, `@@index([location_id])`, `@@index([asset_id])`, `@@index([status])`, `@@index([device_type])`

### `ict_monitor_netif_snapshot` — Snapshot interface jaringan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_monitor_network_device |
| `company_id` | String | FK → dat_company |
| `interface_name` | String | Nama interface |
| `if_index` | Int? | SNMP ifIndex |
| `if_descr` | String? | ifDescr |
| `if_type` | Int? | ifType |
| `if_speed` | BigInt? | Kecepatan interface |
| `if_high_speed` | Int? | Kecepatan dalam Mbps |
| `if_admin_status` | String? | Status admin |
| `if_oper_status` | String? | Status operasional |
| `if_phys_address` | String? | MAC address |
| `if_alias` | String? | Alias interface |
| `in_octets` | BigInt? | Byte masuk |
| `out_octets` | BigInt? | Byte keluar |
| `in_errors` | BigInt? | Error masuk |
| `out_errors` | BigInt? | Error keluar |
| `in_discards` | BigInt? | Discard masuk |
| `out_discards` | BigInt? | Discard keluar |
| `polled_at` | DateTime | Default now() |

**Index:** `@@index([device_id, polled_at])`, `@@index([company_id])`, `@@index([polled_at])`

### `ict_monitor_snmp_trap` — SNMP trap yang diterima

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_monitor_network_device |
| `company_id` | String | FK → dat_company |
| `trap_oid` | String | OID trap |
| `trap_name` | String? | Nama trap |
| `severity` | severity_level | Default "info" |
| `source_ip` | String | IP sumber trap |
| `agent_address` | String? | Alamat agent |
| `community` | String? | Community string |
| `uptime` | String? | Uptime saat trap |
| `variables` | Json? | Variabel trap |
| `raw_message` | String? | Pesan mentah trap |
| `received_at` | DateTime | Default now() |

**Index:** `@@index([device_id, received_at])`, `@@index([company_id])`, `@@index([trap_oid])`, `@@index([received_at])`, `@@index([severity])`

---

## 25. Monitor Metric

> File: `ict_base/prisma/schema/ict_monitor_metric.prisma`

Menyimpan data metrik time-series untuk host, container, dan network device. Menggantikan LibreNMS Graphing dan Beszel Historical Data.

### Shared Enum

#### `metric_category` — Kategori metrik

```prisma
enum metric_category {
  host_cpu
  host_memory
  host_disk
  host_network
  host_load
  host_temperature
  host_gpu
  container_cpu
  container_memory
  container_network
  container_block
  container_pids
  netif_traffic
  netif_errors
  netif_utilization
  snmp_device_cpu
  snmp_device_memory
  snmp_device_temp
  snmp_device_uptime
}
```

#### `metric_source` — Sumber pengumpulan metrik

```prisma
enum metric_source {
  agent
  docker_api
  snmp_poller
  ssh_collect
  api_pull
}
```

### `ict_monitor_metric_definition` — Definisi metrik

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Kode unik metrik |
| `name` | String | Nama metrik |
| `category` | metric_category | Kategori metrik |
| `unit` | String | Satuan (%, bytes, count, dll) |
| `description` | String? | Deskripsi metrik |
| `data_type` | String | Default "gauge" — gauge, counter |
| `min_value` | Decimal? | Nilai minimum |
| `max_value` | Decimal? | Nilai maksimum |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([category])`

### `ict_monitor_metric_point` — Data point metrik

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `metric_def_id` | String | FK → ict_monitor_metric_definition |
| `source` | metric_source | Default "agent" |
| `entity_type` | String | Tipe entitas (host, container, device) |
| `entity_id` | String | ID entitas |
| `agent_id` | String? | FK → ict_monitor_agent |
| `value` | Decimal | Nilai metrik |
| `tags` | Json? | Tags tambahan |
| `collected_at` | DateTime | Default now() |

**Index:** `@@index([metric_def_id, collected_at])`, `@@index([company_id, entity_type, entity_id, collected_at])`, `@@index([agent_id])`, `@@index([collected_at])`

### `ict_monitor_metric_daily` — Agregasi metrik harian

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `metric_def_id` | String | FK → ict_monitor_metric_definition |
| `entity_type` | String | Tipe entitas |
| `entity_id` | String | ID entitas |
| `date` | DateTime | @db.Date — Tanggal agregasi |
| `min_value` | Decimal | Nilai minimum |
| `max_value` | Decimal | Nilai maksimum |
| `avg_value` | Decimal | Nilai rata-rata |
| `sum_value` | Decimal | Total nilai |
| `sample_count` | Int | Jumlah sampel |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([metric_def_id, entity_type, entity_id, date])`, `@@index([company_id, entity_type, entity_id, date])`, `@@index([date])`

### `ict_monitor_host_snapshot` — Snapshot komprehensif host

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `agent_id` | String | FK → ict_monitor_agent |
| `hostname` | String | Nama host |
| `cpu_percent` | Decimal | Persentase CPU |
| `cpu_cores` | Int? | Jumlah core CPU |
| `cpu_model` | String? | Model CPU |
| `memory_usage` | BigInt | Penggunaan memori (bytes) |
| `memory_total` | BigInt | Total memori (bytes) |
| `memory_percent` | Decimal | Persentase memori |
| `swap_usage` | BigInt? | Penggunaan swap (bytes) |
| `swap_total` | BigInt? | Total swap (bytes) |
| `disk_total` | BigInt | Total disk (bytes) |
| `disk_used` | BigInt | Terpakai disk (bytes) |
| `disk_percent` | Decimal | Persentase disk |
| `disk_io_read` | BigInt? | IO read (bytes) |
| `disk_io_write` | BigInt? | IO write (bytes) |
| `net_rx` | BigInt | Network received (bytes) |
| `net_tx` | BigInt | Network transmitted (bytes) |
| `load_avg_1` | Decimal? | Load average 1 menit |
| `load_avg_5` | Decimal? | Load average 5 menit |
| `load_avg_15` | Decimal? | Load average 15 menit |
| `uptime_seconds` | BigInt? | Uptime dalam detik |
| `temperature` | Decimal? | Suhu (Celsius) |
| `gpu_percent` | Decimal? | Penggunaan GPU |
| `gpu_memory` | Decimal? | Penggunaan memori GPU |
| `gpu_temp` | Decimal? | Suhu GPU |
| `snapshot_at` | DateTime | Default now() |

**Index:** `@@index([agent_id, snapshot_at])`, `@@index([company_id, snapshot_at])`, `@@index([snapshot_at])`

### `ict_monitor_disk_health` — Kesehatan disk (S.M.A.R.T.)

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `agent_id` | String | FK → ict_monitor_agent |
| `device_name` | String | Nama device (sda, nvme0, dll) |
| `model_name` | String? | Model disk |
| `serial_number` | String? | Nomor serial |
| `firmware` | String? | Versi firmware |
| `capacity_bytes` | BigInt? | Kapasitas (bytes) |
| `sector_size` | Int? | Ukuran sektor |
| `temperature` | Int? | Suhu (Celsius) |
| `power_on_hours` | Int? | Jam menyala |
| `power_cycle_count` | Int? | Jumlah siklus daya |
| `reallocated_sectors` | Int? | Sektor reallocated |
| `pending_sectors` | Int? | Sektor pending |
| `uncorrected_errors` | Int? | Error tidak terkoreksi |
| `wear_leveling` | Int? | Level keausan (SSD) |
| `health_status` | String | Default "unknown" — good, failing, failed |
| `smart_attributes` | Json? | Atribut S.M.A.R.T. mentah |
| `checked_at` | DateTime | Default now() |

**Index:** `@@index([agent_id, checked_at])`, `@@index([company_id])`, `@@index([device_name])`, `@@index([health_status])`

### `ict_monitor_gpu_snapshot` — Snapshot GPU

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `agent_id` | String | FK → ict_monitor_agent |
| `gpu_index` | Int | Index GPU |
| `gpu_name` | String | Nama GPU |
| `driver_version` | String? | Versi driver |
| `utilization` | Decimal? | Penggunaan GPU (%) |
| `memory_used` | BigInt? | Memori terpakai (bytes) |
| `memory_total` | BigInt? | Total memori (bytes) |
| `temperature` | Int? | Suhu (Celsius) |
| `power_draw` | Decimal? | Konsumsi daya (Watt) |
| `fan_speed` | Int? | Kecepatan fan (%) |
| `clock_core` | Int? | Clock core (MHz) |
| `clock_memory` | Int? | Clock memori (MHz) |
| `snapshot_at` | DateTime | Default now() |

**Index:** `@@index([agent_id, gpu_index, snapshot_at])`, `@@index([company_id, snapshot_at])`

---

## 26. Monitor Log

> File: `ict_base/prisma/schema/ict_monitor_log.prisma`

Sentralisasi log dari seluruh endpoint, aplikasi, dan container. Menggantikan Wazuh Log Collection dan Dockge Real-time Logs.

### Shared Enum

#### `log_source_type` — Tipe sumber log

```prisma
enum log_source_type {
  agent_syslog
  agent_journald
  agent_windows_event
  agent_file
  docker_container
  nginx_access
  nginx_error
  application
  snmp_trap
  api_webhook
}
```

#### `log_level` — Level log

```prisma
enum log_level {
  emergency
  alert
  critical
  error
  warning
  notice
  info
  debug
}
```

### `ict_monitor_log_source` — Konfigurasi sumber log

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Kode unik sumber |
| `name` | String | Nama sumber |
| `source_type` | log_source_type | Tipe sumber log |
| `config` | Json? | Konfigurasi sumber |
| `retention_days` | Int | Default 90 — Retensi data |
| `is_enabled` | Boolean | Default true |
| `last_received_at` | DateTime? | Waktu log terakhir diterima |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([source_type])`

### `ict_monitor_log_entry` — Entri log

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `source_id` | String | FK → ict_monitor_log_source |
| `agent_id` | String? | FK → ict_monitor_agent |
| `hostname` | String | Nama host sumber |
| `program` | String? | Program/app sumber |
| `pid` | Int? | Process ID |
| `level` | log_level | Default "info" |
| `message` | String | Pesan log |
| `fields` | Json? | Fields tambahan |
| `raw_log` | String? | Log mentah |
| `checksum` | String? | Checksum untuk deduplikasi |
| `decoded` | Boolean | Default false — sudah didekode |
| `rule_id` | String? | FK → ict_monitor_log_rule |
| `alert_id` | String? | ID alert terkait |
| `received_at` | DateTime | Default now() |

**Index:** `@@index([source_id, received_at])`, `@@index([company_id, received_at])`, `@@index([agent_id])`, `@@index([level])`, `@@index([hostname])`, `@@index([program])`, `@@index([rule_id])`, `@@index([received_at])`

### `ict_monitor_log_decoder` — Decoder log

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `source_id` | String | FK → ict_monitor_log_source |
| `company_id` | String | FK → dat_company |
| `name` | String | Nama decoder |
| `description` | String? | Deskripsi decoder |
| `regex_pattern` | String | Pola regex untuk dekodifikasi |
| `extract_fields` | Json? | Field yang diekstrak |
| `prematch` | String? | Pola pencocokan awal |
| `order` | Int | Default 0 — Urutan eksekusi |
| `is_enabled` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([source_id])`, `@@index([company_id])`

### `ict_monitor_log_rule` — Aturan log

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `source_id` | String? | FK → ict_monitor_log_source |
| `company_id` | String | FK → dat_company |
| `code` | String | Kode unik aturan |
| `name` | String | Nama aturan |
| `description` | String? | Deskripsi aturan |
| `condition_query` | Json | Kondisi pencocokan |
| `action` | String | Default "alert" — alert, notify, block |
| `level` | severity_level | Default "low" |
| `frequency` | Int | Default 1 — Jumlah kejadian |
| `timeframe_minutes` | Int | Default 5 — Jangka waktu |
| `mitre_tactic` | String? | Taktik MITRE ATT&CK |
| `mitre_technique` | String? | Teknik MITRE ATT&CK |
| `mitre_technique_id` | String? | ID Teknik MITRE |
| `is_enabled` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([source_id])`, `@@index([level])`, `@@index([is_enabled])`

### `ict_monitor_log_saved_query` — Query tersimpan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `source_id` | String | FK → ict_monitor_log_source |
| `company_id` | String | FK → dat_company |
| `name` | String | Nama query |
| `query` | Json | Definisi query |
| `description` | String? | Deskripsi query |
| `is_shared` | Boolean | Default false — Dibagikan ke tim |
| `created_by` | String | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([source_id])`, `@@index([company_id])`

### `ict_monitor_container_log` — Log container Docker

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `host_id` | String | ID host Docker |
| `container_id` | String | ID container |
| `stream` | String | Default "stdout" — stdout, stderr |
| `line_number` | Int? | Nomor baris |
| `message` | String | Pesan log |
| `timestamp` | DateTime | Default now() |

**Index:** `@@index([container_id, timestamp])`, `@@index([host_id, timestamp])`, `@@index([company_id, timestamp])`, `@@index([timestamp])`

---

## 27. Monitor FIM

> File: `ict_base/prisma/schema/ict_monitor_fim.prisma`

Memantau perubahan file kritis pada seluruh endpoint. Menggantikan Wazuh FIM (File Integrity Monitoring).

### Shared Enum

#### `fim_scan_mode` — Mode pemindaian FIM

```prisma
enum fim_scan_mode {
  realtime
  periodic
  scheduled
}
```

#### `fim_change_type` — Tipe perubahan file

```prisma
enum fim_change_type {
  added
  modified
  deleted
  moved
  permission_changed
  ownership_changed
}
```

#### `fim_severity` — Severity perubahan FIM

```prisma
enum fim_severity {
  critical
  high
  medium
  low
}
```

### `ict_monitor_fim_config` — Konfigurasi FIM

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Kode unik konfigurasi |
| `name` | String | Nama konfigurasi |
| `description` | String? | Deskripsi |
| `scan_mode` | fim_scan_mode | Default "periodic" |
| `scan_interval` | Int | Default 3600 — Interval detik |
| `watch_paths` | String[] | Path yang dipantau |
| `ignore_patterns` | String[] | Pola yang diabaikan |
| `ignore_dirs` | String[] | Direktori yang diabaikan |
| `check_hash` | Boolean | Default true — Cek hash file |
| `check_size` | Boolean | Default true — Cek ukuran |
| `check_mtime` | Boolean | Default true — Cek waktu modifikasi |
| `check_ctime` | Boolean | Default true — Cek waktu perubahan |
| `check_perms` | Boolean | Default true — Cek izin |
| `check_owner` | Boolean | Default true — Cek pemilik |
| `is_enabled` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([is_enabled])`

### `ict_monitor_fim_baseline` — Baseline file

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `config_id` | String | FK → ict_monitor_fim_config |
| `agent_id` | String | FK → ict_monitor_agent |
| `company_id` | String | FK → dat_company |
| `file_path` | String | Path lengkap file |
| `file_name` | String | Nama file |
| `file_dir` | String | Direktori file |
| `file_size` | BigInt | Ukuran file (bytes) |
| `file_hash_md5` | String? | Hash MD5 |
| `file_hash_sha1` | String? | Hash SHA1 |
| `file_hash_sha256` | String? | Hash SHA256 |
| `permissions` | String? | Izin file (chmod) |
| `owner_uid` | Int? | UID pemilik |
| `owner_gid` | Int? | GID pemilik |
| `inode` | BigInt? | Inode |
| `hard_links` | Int? | Jumlah hard link |
| `last_modified` | DateTime | Waktu modifikasi terakhir |
| `first_seen_at` | DateTime | Default now() — Pertama kali terlihat |
| `last_scan_at` | DateTime? | Pemindaian terakhir |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([config_id, agent_id, file_path])`, `@@index([config_id])`, `@@index([agent_id])`, `@@index([company_id])`, `@@index([file_path])`, `@@index([file_hash_sha256])`

### `ict_monitor_fim_change` — Perubahan file terdeteksi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `baseline_id` | String | FK → ict_monitor_fim_baseline |
| `agent_id` | String | FK → ict_monitor_agent |
| `company_id` | String | FK → dat_company |
| `change_type` | fim_change_type | Tipe perubahan |
| `severity` | fim_severity | Default "medium" |
| `old_hash` | String? | Hash lama |
| `new_hash` | String? | Hash baru |
| `old_size` | BigInt? | Ukuran lama |
| `new_size` | BigInt? | Ukuran baru |
| `old_permissions` | String? | Izin lama |
| `new_permissions` | String? | Izin baru |
| `old_owner` | Int? | UID lama |
| `new_owner` | Int? | UID baru |
| `description` | String? | Deskripsi perubahan |
| `is_acknowledged` | Boolean | Default false — Sudah diakui |
| `acknowledged_by` | String? | FK → dat_user |
| `acknowledged_at` | DateTime? | Waktu diakui |
| `detected_at` | DateTime | Default now() |

**Index:** `@@index([baseline_id, detected_at])`, `@@index([agent_id])`, `@@index([company_id])`, `@@index([change_type])`, `@@index([severity])`, `@@index([detected_at])`, `@@index([is_acknowledged])`

### `ict_monitor_fim_alert` — Alert FIM

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `config_id` | String | FK → ict_monitor_fim_config |
| `company_id` | String | FK → dat_company |
| `change_id` | String? | FK → ict_monitor_fim_change |
| `title` | String | Judul alert |
| `description` | String | Deskripsi alert |
| `severity` | fim_severity | Default "medium" |
| `status` | String | Default "new" — new, acknowledged, resolved |
| `assigned_to` | String? | FK → dat_user |
| `resolved_at` | DateTime? | Waktu resolved |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([config_id])`, `@@index([company_id])`, `@@index([severity])`, `@@index([status])`, `@@index([created_at])`

---

## 28. Monitor SCA

> File: `ict_base/prisma/schema/ict_monitor_sca.prisma`

Pemindaian konfigurasi keamanan terhadap benchmark (CIS, STIG, dll). Menggantikan Wazuh SCA (Security Configuration Assessment).

### Shared Enum

#### `sca_policy_status` — Status kebijakan SCA

```prisma
enum sca_policy_status {
  draft
  active
  deprecated
}
```

#### `sca_scan_status` — Status pemindaian SCA

```prisma
enum sca_scan_status {
  pending
  running
  completed
  failed
  cancelled
}
```

#### `sca_check_result` — Hasil pemeriksaan SCA

```prisma
enum sca_check_result {
  pass
  fail
  not_applicable
  error
}
```

### `ict_monitor_sca_policy` — Kebijakan SCA

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Kode unik kebijakan |
| `name` | String | Nama kebijakan |
| `description` | String? | Deskripsi |
| `benchmark_name` | String | Nama benchmark (CIS, STIG, dll) |
| `benchmark_version` | String? | Versi benchmark |
| `platform` | String[] | Platform yang didukung |
| `policy_file_url` | String? | URL file kebijakan |
| `total_checks` | Int | Default 0 — Total pemeriksaan |
| `status` | sca_policy_status | Default "draft" |
| `is_enabled` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([benchmark_name])`, `@@index([status])`

### `ict_monitor_sca_check` — Pemeriksaan SCA

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `policy_id` | String | FK → ict_monitor_sca_policy |
| `check_id` | String | ID pemeriksaan dalam kebijakan |
| `title` | String | Judul pemeriksaan |
| `description` | String? | Deskripsi |
| `rationale` | String? | Alasan/justifikasi |
| `remediation` | String? | Langkah perbaikan |
| `severity` | severity_level | Default "medium" |
| `compliance_url` | String? | URL referensi compliance |
| `condition` | Json? | Kondisi pemeriksaan |
| `is_enabled` | Boolean | Default true |
| `sort_order` | Int | Default 0 — Urutan tampil |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([policy_id, check_id])`, `@@index([policy_id])`, `@@index([severity])`

### `ict_monitor_sca_scan` — Pemindaian SCA

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `policy_id` | String | FK → ict_monitor_sca_policy |
| `company_id` | String | FK → dat_company |
| `agent_id` | String | FK → ict_monitor_agent |
| `status` | sca_scan_status | Default "pending" |
| `started_at` | DateTime? | Waktu mulai |
| `finished_at` | DateTime? | Waktu selesai |
| `total_checks` | Int | Default 0 — Total pemeriksaan |
| `passed_count` | Int | Default 0 — Lulus |
| `failed_count` | Int | Default 0 — Gagal |
| `not_applicable_count` | Int | Default 0 — Tidak berlaku |
| `error_count` | Int | Default 0 — Error |
| `compliance_score` | Decimal? | Skor compliance (%) |
| `triggered_by` | String? | FK → dat_user atau "system" |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([policy_id])`, `@@index([company_id])`, `@@index([agent_id])`, `@@index([status])`, `@@index([created_at])`

### `ict_monitor_sca_result` — Hasil pemeriksaan SCA

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `scan_id` | String | FK → ict_monitor_sca_scan |
| `check_id` | String | FK → ict_monitor_sca_check |
| `result` | sca_check_result | Hasil pemeriksaan |
| `actual_value` | String? | Nilai aktual |
| `reason` | String? | Alasan hasil |
| `remediation_applied` | Boolean | Default false — Perbaikan diterapkan |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([scan_id, check_id])`, `@@index([scan_id])`, `@@index([check_id])`, `@@index([result])`

---

## 29. Monitor Discovery

> File: `ict_base/prisma/schema/ict_monitor_discovery.prisma`

Otomatis menemukan dan menginventory perangkat jaringan. Menggantikan LibreNMS Auto-Discovery.

### Shared Enum

#### `discovery_status` — Status penemuan

```prisma
enum discovery_status {
  pending
  running
  completed
  failed
  cancelled
}
```

#### `discovery_method` — Metode penemuan

```prisma
enum discovery_method {
  snmp_scan
  ping_sweep
  cdp_neighbors
  lldp_neighbors
  arp_table
  dhcp_lease
  dns_reverse
  ospf_neighbor
  bgp_peer
}
```

### `ict_monitor_discovery_job` — Tugas penemuan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `name` | String | Nama tugas |
| `description` | String? | Deskripsi |
| `method` | discovery_method | Default "snmp_scan" |
| `target_ranges` | String[] | Range IP target (CIDR) |
| `excluded_ranges` | String[] | Range IP yang dikecualikan |
| `snmp_community` | String? | Community string SNMP |
| `snmp_version` | String | Default "2c" |
| `ping_timeout_ms` | Int | Default 1000 — Timeout ping |
| `max_concurrent` | Int | Default 50 — Maksimum konkuren |
| `config` | Json? | Konfigurasi tambahan |
| `status` | discovery_status | Default "pending" |
| `started_at` | DateTime? | Waktu mulai |
| `finished_at` | DateTime? | Waktu selesai |
| `devices_found` | Int | Default 0 — Perangkat ditemukan |
| `new_devices` | Int | Default 0 — Perangkat baru |
| `updated_devices` | Int | Default 0 — Perangkat diperbarui |
| `triggered_by` | String? | FK → dat_user atau "system" |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, name])`, `@@index([company_id])`, `@@index([status])`, `@@index([started_at])`

### `ict_monitor_discovery_result` — Hasil penemuan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `job_id` | String | FK → ict_monitor_discovery_job |
| `company_id` | String | FK → dat_company |
| `device_ip` | String | IP perangkat ditemukan |
| `hostname` | String? | Nama host |
| `mac_address` | String? | MAC address |
| `vendor` | String? | Vendor perangkat |
| `device_type` | String? | Tipe perangkat |
| `os_fingerprint` | String? | Deteksi OS |
| `snmp_sysname` | String? | sysName SNMP |
| `snmp_sysdesc` | String? | sysDescr SNMP |
| `open_ports` | Int[] | Port terbuka |
| `is_new` | Boolean | Default false — Perangkat baru |
| `network_device_id` | String? | FK → ict_monitor_network_device (jika sudah terdaftar) |
| `raw_data` | Json? | Data mentah hasil penemuan |
| `discovered_at` | DateTime | Default now() |

**Index:** `@@unique([job_id, device_ip])`, `@@index([job_id])`, `@@index([company_id])`, `@@index([device_ip])`, `@@index([is_new])`

---

## 30. Monitor Alert

> File: `ict_base/prisma/schema/ict_monitor_alert.prisma`

Unified alerting engine yang menggantikan LibreNMS Alert Rules, Wazuh Alert Engine, dan Beszel Thresholds.

### Shared Enum

#### `alert_rule_status` — Status aturan alert

```prisma
enum alert_rule_status {
  active
  disabled
  deleted
}
```

#### `alert_condition_type` — Tipe kondisi alert

```prisma
enum alert_condition_type {
  threshold_above
  threshold_below
  threshold_equals
  change_detected
  pattern_match
  absence_detect
  composite
}
```

#### `alert_severity` — Severity alert

```prisma
enum alert_severity {
  emergency
  critical
  high
  medium
  low
  info
}
```

#### `alert_state` — Status alert firing

```prisma
enum alert_state {
  pending
  firing
  acknowledged
  resolved
  silenced
}
```

#### `notification_channel` — Channel notifikasi

```prisma
enum notification_channel {
  toast
  email
  webhook
  slack
  teams
  pagerduty
}
```

### `ict_monitor_alert_rule` — Aturan alert

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `metric_def_id` | String? | FK → ict_monitor_metric_definition |
| `code` | String | Kode unik aturan |
| `name` | String | Nama aturan |
| `description` | String? | Deskripsi |
| `entity_type` | String | Tipe entitas yang dipantau |
| `condition_type` | alert_condition_type | Default "threshold_above" |
| `condition_config` | Json | Konfigurasi kondisi |
| `severity` | alert_severity | Default "medium" |
| `cooldown_minutes` | Int | Default 15 — cooldown antar alert |
| `repeat_interval` | Int? | Interval pengulangan (menit) |
| `trigger_count` | Int | Default 0 — Jumlah trigger |
| `last_triggered_at` | DateTime? | Waktu terakhir trigger |
| `status` | alert_rule_status | Default "active" |
| `created_by` | String? | FK → dat_user |
| `is_enabled` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([entity_type])`, `@@index([severity])`, `@@index([status])`, `@@index([metric_def_id])`

### `ict_monitor_alert_firing` — Alert yang sedang aktif

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `rule_id` | String | FK → ict_monitor_alert_rule |
| `company_id` | String | FK → dat_company |
| `entity_type` | String | Tipe entitas |
| `entity_id` | String | ID entitas |
| `entity_label` | String? | Label entitas |
| `state` | alert_state | Default "pending" |
| `severity` | alert_severity | Severity alert |
| `trigger_value` | Json | Nilai yang memicu alert |
| `message` | String | Pesan alert |
| `description` | String? | Deskripsi detail |
| `acknowledged_by` | String? | FK → dat_user |
| `acknowledged_at` | DateTime? | Waktu diakui |
| `resolved_at` | DateTime? | Waktu resolved |
| `silenced_until` | DateTime? | Dibisukan sampai |
| `escalation_level` | Int | Default 0 — Level eskalasi |
| `notified_users` | String[] | User yang sudah dinotifikasi |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([rule_id, created_at])`, `@@index([company_id, state])`, `@@index([entity_type, entity_id])`, `@@index([state])`, `@@index([severity])`, `@@index([created_at])`

### `ict_monitor_alert_silence` — Pembisuan alert

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `rule_id` | String | FK → ict_monitor_alert_rule |
| `company_id` | String | FK → dat_company |
| `name` | String | Nama pembisuan |
| `description` | String? | Deskripsi |
| `matchers` | Json | Kondisi pencocokan |
| `start_at` | DateTime | Waktu mulai bisu |
| `end_at` | DateTime | Waktu selesai bisu |
| `created_by` | String | FK → dat_user |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([rule_id])`, `@@index([company_id])`, `@@index([start_at, end_at])`

### `ict_monitor_alert_notification` — Riwayat notifikasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `firing_id` | String? | FK → ict_monitor_alert_firing |
| `channel` | notification_channel | Channel notifikasi |
| `recipient` | String | Penerima notifikasi |
| `subject` | String | Subjek notifikasi |
| `message` | String | Isi notifikasi |
| `status` | String | Default "pending" — pending, sent, failed |
| `sent_at` | DateTime? | Waktu terkirim |
| `error_msg` | String? | Pesan error |
| `retry_count` | Int | Default 0 — Jumlah percobaan ulang |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([company_id, created_at])`, `@@index([firing_id])`, `@@index([channel])`, `@@index([status])`

---

## 31. Monitor Docker Extensions

> File: `ict_base/prisma/schema/ict_monitor_docker_ext.prisma`

Ekstensi manajemen DockerCompose version history, web terminal sessions, dan log streaming. Menggantikan Dockge Compose Management.

### Shared Enum

#### `terminal_status` — Status terminal

```prisma
enum terminal_status {
  active
  closed
  timeout
  error
}
```

#### `log_stream_status` — Status log streaming

```prisma
enum log_stream_status {
  streaming
  paused
  stopped
  error
}
```

### `ict_docker_compose_version` — Versi Docker Compose

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `compose_id` | String | ID compose |
| `company_id` | String | FK → dat_company |
| `version` | Int | Nomor versi |
| `compose_content` | String | Isi file compose.yml |
| `env_content` | String? | Isi file .env |
| `change_note` | String? | Catatan perubahan |
| `change_type` | String | Default "manual" — manual, git, api |
| `created_by` | String | FK → dat_user |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([compose_id, version])`, `@@index([compose_id])`, `@@index([company_id])`, `@@index([created_by])`

### `ict_docker_terminal_session` — Sesi terminal web

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `host_id` | String | ID host Docker |
| `container_id` | String? | ID container (opsional) |
| `company_id` | String | FK → dat_company |
| `user_id` | String | FK → dat_user |
| `session_key` | String | Kunci sesi unik |
| `status` | terminal_status | Default "active" |
| `shell` | String | Default "/bin/sh" — Shell yang digunakan |
| `cols` | Int | Default 80 — Lebar terminal |
| `rows` | Int | Default 24 — Tinggi terminal |
| `started_at` | DateTime | Default now() |
| `ended_at` | DateTime? | Waktu berakhir |
| `duration_sec` | Int? | Durasi sesi (detik) |

**Index:** `@@unique([session_key])`, `@@index([host_id])`, `@@index([container_id])`, `@@index([company_id])`, `@@index([user_id])`, `@@index([status])`

### `ict_docker_log_stream` — Streaming log container

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `host_id` | String | ID host Docker |
| `container_id` | String | ID container |
| `company_id` | String | FK → dat_company |
| `user_id` | String | FK → dat_user |
| `stream_type` | String | Default "both" — stdout, stderr, both |
| `tail_lines` | Int | Default 100 — Jumlah baris awal |
| `follow` | Boolean | Default true — Ikuti log baru |
| `status` | log_stream_status | Default "streaming" |
| `started_at` | DateTime | Default now() |
| `stopped_at` | DateTime? | Waktu berhenti |

**Index:** `@@index([host_id, container_id])`, `@@index([company_id])`, `@@index([user_id])`, `@@index([status])`

### `ict_docker_compose_diff` — Perbandingan versi Compose

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `compose_id` | String | ID compose |
| `from_version` | Int | Versi asal |
| `to_version` | Int | Versi tujuan |
| `diff_content` | String | Isi diff |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([compose_id])`

### `ict_docker_env_override` — Override environment variable

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `compose_id` | String | ID compose |
| `key` | String | Nama variabel |
| `value` | String | Nilai variabel |
| `is_secret` | Boolean | Default false — Rahasia |
| `source` | String | Default "manual" — manual, vault, env |
| `created_by` | String? | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([compose_id, key])`, `@@index([compose_id])`

---

## 32. ISO 27001 & SOP

> File: `ict_base/prisma/schema/ict_iso.prisma`

Mengelola kebijakan, prosedur, bukti compliance, dan metrik keamanan untuk ISO 27001 dan dokumentasi SOP.

### Shared Enum

#### `policy_status` — Status kebijakan/prosedur

```prisma
enum policy_status {
  draft
  under_review
  approved
  published
  archived
  under_revision
}
```

#### `evidence_status` — Status bukti compliance

```prisma
enum evidence_status {
  pending
  collected
  verified
  expired
  rejected
}
```

#### `sop_category` — Kategori SOP

```prisma
enum sop_category {
  access_control
  incident_response
  backup_recovery
  change_management
  network_security
  physical_security
  human_resources
  vendor_management
  data_protection
  business_continuity
  monitoring_operations
  vulnerability_management
  configuration_management
  asset_management
  other
}
```

#### `metric_frequency` — Frekuensi metrik

```prisma
enum metric_frequency {
  daily
  weekly
  monthly
  quarterly
  yearly
}
```

### `ict_iso_policy` — Kebijakan Keamanan Informasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `doc_id` | String? | FK → ict_doc |
| `code` | String | Kode unik kebijakan |
| `title` | String | Judul kebijakan |
| `description` | String? | Deskripsi |
| `policy_number` | String? | Nomor kebijakan |
| `version` | Int | Default 1 — Versi saat ini |
| `status` | policy_status | Default "draft" |
| `owner_id` | String | FK → dat_user — Pemilik kebijakan |
| `approver_id` | String? | FK → dat_user — Penyetuju |
| `effective_date` | DateTime? | Tanggal berlaku |
| `review_date` | DateTime? | Tanggal review berikutnya |
| `review_cycle_months` | Int | Default 12 — Siklus review (bulan) |
| `scope` | String? | Ruang lingkup |
| `objective` | String? | Tujuan kebijakan |
| `is_mandatory` | Boolean | Default true — Wajib |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([status])`, `@@index([owner_id])`, `@@index([review_date])`

### `ict_iso_policy_version` — Versi kebijakan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `policy_id` | String | FK → ict_iso_policy (onDelete: Cascade) |
| `version` | Int | Nomor versi |
| `title` | String | Judul versi |
| `content` | String | Konten kebijakan (Markdown/HTML) |
| `change_note` | String? | Catatan perubahan |
| `author_id` | String | FK → dat_user — Penulis |
| `approved_by` | String? | FK → dat_user — Disetujui oleh |
| `approved_at` | DateTime? | Waktu persetujuan |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([policy_id, version])`, `@@index([policy_id])`

### `ict_iso_policy_control` — Mapping kebijakan ke kontrol ISO

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `policy_id` | String | FK → ict_iso_policy (onDelete: Cascade) |
| `control_id` | String | ID kontrol ISO 27001 (A.5.1.1, dll) |
| `notes` | String? | Catatan mapping |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([policy_id, control_id])`, `@@index([policy_id])`, `@@index([control_id])`

### `ict_iso_sop` — Standard Operating Procedure

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `doc_id` | String? | FK → ict_doc |
| `policy_id` | String? | FK → ict_iso_policy |
| `code` | String | Kode unik SOP |
| `title` | String | Judul SOP |
| `description` | String? | Deskripsi |
| `sop_number` | String? | Nomor SOP |
| `category` | sop_category | Default "other" |
| `version` | Int | Default 1 — Versi saat ini |
| `status` | policy_status | Default "draft" |
| `owner_id` | String | FK → dat_user — Pemilik SOP |
| `approver_id` | String? | FK → dat_user — Penyetuju |
| `effective_date` | DateTime? | Tanggal berlaku |
| `review_date` | DateTime? | Tanggal review berikutnya |
| `review_cycle_months` | Int | Default 12 — Siklus review (bulan) |
| `scope` | String? | Ruang lingkup |
| `objective` | String? | Tujuan SOP |
| `steps_count` | Int | Default 0 — Jumlah langkah |
| `is_mandatory` | Boolean | Default true — Wajib |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([category])`, `@@index([status])`, `@@index([owner_id])`, `@@index([review_date])`, `@@index([policy_id])`

### `ict_iso_sop_version` — Versi SOP

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `sop_id` | String | FK → ict_iso_sop (onDelete: Cascade) |
| `version` | Int | Nomor versi |
| `title` | String | Judul versi |
| `content` | String | Konten SOP (Markdown/HTML) |
| `change_note` | String? | Catatan perubahan |
| `author_id` | String | FK → dat_user — Penulis |
| `approved_by` | String? | FK → dat_user — Disetujui oleh |
| `approved_at` | DateTime? | Waktu persetujuan |
| `created_at` | DateTime | Default now() |

**Index:** `@@unique([sop_id, version])`, `@@index([sop_id])`

### `ict_iso_sop_step` — Langkah SOP

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `sop_id` | String | FK → ict_iso_sop (onDelete: Cascade) |
| `step_number` | Int | Nomor urut langkah |
| `title` | String | Judul langkah |
| `description` | String | Deskripsi langkah |
| `warning` | String? | Peringatan |
| `tip` | String? | Tips |
| `image_url` | String? | URL gambar ilustrasi |
| `role_required` | String? | Role yang diperlukan |
| `is_critical` | Boolean | Default false — Langkah kritis |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([sop_id, step_number])`, `@@index([sop_id])`

### `ict_iso_sop_execution` — Eksekusi SOP

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `sop_id` | String | FK → ict_iso_sop |
| `company_id` | String | FK → dat_company |
| `executed_by` | String | FK → dat_user |
| `started_at` | DateTime | Default now() |
| `finished_at` | DateTime? | Waktu selesai |
| `status` | String | Default "in_progress" — in_progress, completed, failed |
| `notes` | String? | Catatan eksekusi |
| `deviations` | String? | Penyimpangan dari SOP |
| `evidence_urls` | String[] | URL bukti eksekusi |
| `duration_min` | Int? | Durasi eksekusi (menit) |

**Index:** `@@index([sop_id])`, `@@index([company_id])`, `@@index([executed_by])`, `@@index([started_at])`, `@@index([status])`

### `ict_iso_evidence` — Bukti Compliance

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `control_id` | String? | ID kontrol ISO 27001 |
| `sop_id` | String? | FK → ict_iso_sop |
| `finding_id` | String? | ID temuan audit |
| `code` | String | Kode unik bukti |
| `title` | String | Judul bukti |
| `description` | String? | Deskripsi |
| `evidence_type` | String | Tipe bukti (screenshot, document, log, dll) |
| `source_system` | String? | Sistem sumber |
| `file_urls` | String[] | URL file bukti |
| `collected_by` | String | FK → dat_user |
| `collected_at` | DateTime | Default now() |
| `verified_by` | String? | FK → dat_user |
| `verified_at` | DateTime? | Waktu verifikasi |
| `status` | evidence_status | Default "pending" |
| `expiry_date` | DateTime? | Tanggal kedaluwarsa |
| `notes` | String? | Catatan |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([control_id])`, `@@index([sop_id])`, `@@index([status])`, `@@index([collected_at])`, `@@index([expiry_date])`

### `ict_iso_security_metric` — Definisi Metrik Keamanan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Kode unik metrik |
| `name` | String | Nama metrik |
| `description` | String? | Deskripsi |
| `category` | String | Kategori metrik |
| `unit` | String | Default "count" — Satuan |
| `frequency` | metric_frequency | Default "monthly" |
| `target_value` | Decimal? | Nilai target |
| `warning_threshold` | Decimal? | Ambang peringatan |
| `critical_threshold` | Decimal? | Ambang kritis |
| `formula` | String? | Rumus perhitungan |
| `data_source` | String? | Sumber data |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([category])`

### `ict_iso_security_metric_value` — Nilai Metrik Keamanan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `metric_id` | String | FK → ict_iso_security_metric |
| `company_id` | String | FK → dat_company |
| `period_start` | DateTime | @db.Date — Tanggal mulai periode |
| `period_end` | DateTime | @db.Date — Tanggal akhir periode |
| `actual_value` | Decimal | Nilai aktual |
| `target_value` | Decimal? | Nilai target |
| `is_met` | Boolean | Target tercapai |
| `notes` | String? | Catatan |
| `data_source` | String? | Sumber data |
| `recorded_at` | DateTime | Default now() |

**Index:** `@@unique([metric_id, period_start])`, `@@index([metric_id, period_start])`, `@@index([company_id, period_start])`

### `ict_iso_report` — Laporan ISO

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Kode unik laporan |
| `title` | String | Judul laporan |
| `report_type` | String | Tipe laporan |
| `period_start` | DateTime | @db.Date — Mulai periode |
| `period_end` | DateTime | @db.Date — Akhir periode |
| `generated_by` | String | FK → dat_user atau "system" |
| `status` | String | Default "draft" — draft, published |
| `file_path` | String? | Path file laporan |
| `file_type` | String | Default "pdf" |
| `summary` | String? | Ringkasan laporan |
| `metrics_summary` | Json? | Ringkasan metrik |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([report_type])`, `@@index([period_start])`

### `ict_iso_report_metric` — Metrik dalam laporan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `report_id` | String | FK → ict_iso_report (onDelete: Cascade) |
| `metric_id` | String | FK → ict_iso_security_metric |
| `value_id` | String? | FK → ict_iso_security_metric_value |
| `actual` | Decimal | Nilai aktual |
| `target` | Decimal? | Nilai target |
| `is_met` | Boolean | Target tercapai |
| `trend` | String? | Tren (up, down, stable) |
| `notes` | String? | Catatan |

**Index:** `@@unique([report_id, metric_id])`, `@@index([report_id])`, `@@index([metric_id])`

### `ict_iso_report_section` — Bagian laporan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `report_id` | String | FK → ict_iso_report (onDelete: Cascade) |
| `title` | String | Judul bagian |
| `content` | String | Konten bagian |
| `sort_order` | Int | Default 0 — Urutan |
| `data_source` | String? | Sumber data |
| `chart_config` | Json? | Konfigurasi chart |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([report_id])`

### `ict_iso_document_register` — Register Dokumen ISO

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `doc_number` | String | Nomor dokumen |
| `title` | String | Judul dokumen |
| `doc_type` | String | Tipe dokumen |
| `category` | String? | Kategori |
| `version` | Int | Default 1 — Versi saat ini |
| `status` | policy_status | Default "draft" |
| `owner_id` | String | FK → dat_user — Pemilik |
| `approver_id` | String? | FK → dat_user — Penyetuju |
| `effective_date` | DateTime? | Tanggal berlaku |
| `review_date` | DateTime? | Tanggal review berikutnya |
| `review_cycle_months` | Int | Default 12 — Siklus review (bulan) |
| `distribution` | String[] | Daftar distribusi |
| `confidentiality` | String | Default "internal" — public, internal, confidential, secret |
| `is_controlled` | Boolean | Default true — Dokumen terkontrol |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, doc_number])`, `@@index([company_id])`, `@@index([doc_type])`, `@@index([status])`, `@@index([owner_id])`, `@@index([review_date])`

---

## 33. Naming Convention

| Elemen | Format | Contoh |
|--------|--------|--------|
| **Tabel** | `ict_{modul}` atau `dat_{entitas}` | `ict_vuln_finding`, `dat_user` |
| **Enum** | `snake_case` sesuai Prisma | `vuln_status`, `severity_level` |
| **Field** | `snake_case` | `company_id`, `created_at` |
| **Primary Key** | `String @id @default(uuid())` | — |
| **Timestamps** | `created_at` + `updated_at` | `DateTime @default(now())` |
| **Multi-tenant** | `company_id` FK | Semua tabel domain |
| **Unique** | `@@unique([company_id, ...])` | Isolasi data per perusahaan |
| **Index** | `@@index` | Kolom yang sering di-query |

---

## 23. Quality Management System (QMS)

> File: `ict_base/prisma/schema/ict_qms.prisma`

### `ict_quality_objective` — Sasaran kualitas

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Unique per company, format: QO-XXXX |
| `title` | String | Judul sasaran |
| `description` | String | Deskripsi |
| `department` | String? | Departemen terkait |
| `owner_id` | String? | FK → dat_user |
| `target_value` | Decimal? | Target numerik |
| `current_value` | Decimal? | Nilai aktual |
| `unit` | String? | Satuan pengukuran |
| `target_date` | DateTime? | Tanggal target tercapai |
| `status` | enum | on_track, at_risk, behind, achieved, cancelled |
| `review_frequency` | String? | monthly, quarterly, annually |
| `last_review_date` | DateTime? | Terakhir di-review |
| `next_review_date` | DateTime? | Review berikutnya |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([status])`, `@@index([department])`

### `ict_quality_kpi` — Pengukuran KPI kualitas

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `objective_id` | String | FK → ict_quality_objective |
| `name` | String | Nama KPI |
| `target_value` | Decimal | Target |
| `current_value` | Decimal? | Aktual |
| `unit` | String? | Satuan |
| `measurement_period` | String? | Periode pengukuran |
| `last_measured_at` | DateTime? | Terakhir diukur |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([objective_id])`

### `ict_supplier` — Data supplier/vendor

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Unique per company, format: SUP-XXXX |
| `name` | String | Nama supplier |
| `contact_name` | String? | PIC |
| `contact_email` | String? | Email |
| `contact_phone` | String? | Telepon |
| `address` | String? | Alamat |
| `website` | String? | Website |
| `category` | String | Kategori: hardware, software, cloud, maintenance, consulting |
| `status` | enum | active, inactive, suspended, blacklisted |
| `rating` | enum? | excellent, good, satisfactory, poor, unacceptable |
| `last_evaluation_date` | DateTime? | Terakhir dievaluasi |
| `next_evaluation_date` | DateTime? | Evaluasi berikutnya |
| `is_critical` | Boolean | Default false |
| `contract_start` | DateTime? | Mulai kontrak |
| `contract_end` | DateTime? | Akhir kontrak |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([status])`, `@@index([category])`

### `ict_supplier_evaluation` — Evaluasi supplier periodik

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `supplier_id` | String | FK → ict_supplier |
| `company_id` | String | FK → dat_company |
| `evaluation_date` | DateTime | Tanggal evaluasi |
| `evaluator_id` | String | FK → dat_user |
| `quality_score` | Decimal? | Skor kualitas |
| `delivery_score` | Decimal? | Skor pengiriman |
| `price_score` | Decimal? | Skor harga |
| `service_score` | Decimal? | Skor layanan |
| `overall_score` | Decimal? | Rata-rata/bobot |
| `rating` | enum | excellent, good, satisfactory, poor, unacceptable |
| `strengths` | String? | Kelebihan |
| `weaknesses` | String? | Kekurangan |
| `improvement_plan` | String? | Rencana perbaikan |
| `next_evaluation_date` | DateTime? | Evaluasi berikutnya |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([supplier_id])`, `@@index([company_id])`, `@@index([evaluation_date])`

### `ict_training_program` — Program pelatihan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `code` | String | Unique per company, format: TRN-XXXX |
| `title` | String | Judul program |
| `description` | String? | |
| `training_type` | enum | onboarding, refresher, compliance, technical, safety, awareness, certification |
| `duration_hours` | Decimal? | Durasi (jam) |
| `provider` | String? | Penyelenggara/internal/external |
| `is_mandatory` | Boolean | Default false |
| `validity_months` | Int? | Masa berlaku sertifikat (bulan) |
| `is_active` | Boolean | Default true |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, code])`, `@@index([company_id])`, `@@index([training_type])`, `@@index([is_mandatory])`

### `ict_training_record` — Record partisipasi & sertifikasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `program_id` | String | FK → ict_training_program |
| `user_id` | String | FK → dat_user |
| `company_id` | String | FK → dat_company |
| `training_date` | DateTime | Tanggal pelatihan |
| `completion_date` | DateTime? | Tanggal selesai |
| `status` | enum | planned, in_progress, completed, cancelled |
| `score` | Decimal? | Nilai (jika ada) |
| `passed` | Boolean? | Lulus/tidak |
| `certificate_url` | String? | URL sertifikat |
| `expiry_date` | DateTime? | Berlaku sampai |
| `notes` | String? | Catatan |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([program_id, user_id, training_date])`, `@@index([program_id])`, `@@index([user_id])`, `@@index([company_id])`, `@@index([status])`, `@@index([expiry_date])`

### `ict_calibration_asset` — Aset yang butuh kalibrasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `asset_id` | String | FK → ict_asset (unique) |
| `calibration_date` | DateTime | Tanggal kalibrasi terakhir |
| `next_calibration_date` | DateTime | Tanggal kalibrasi berikutnya |
| `calibration_interval_months` | Int | Default 12 |
| `status` | enum | due, in_progress, completed, overdue, out_of_tolerance |
| `certificate_number` | String? | Nomor sertifikat |
| `calibrated_by` | String? | Yang melakukan kalibrasi |
| `certificate_url` | String? | URL sertifikat |
| `last_result` | String? | Hasil terakhir |
| `tolerance` | String? | Toleransi |
| `notes` | String? | Catatan |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([asset_id])`, `@@index([company_id])`, `@@index([status])`, `@@index([next_calibration_date])`

### `ict_calibration_record` — Record kalibrasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `calibration_asset_id` | String | FK → ict_calibration_asset |
| `calibration_date` | DateTime | Tanggal kalibrasi |
| `performed_by` | String? | Yang melakukan |
| `result` | String | pass, fail, conditional |
| `is_within_tolerance` | Boolean | Apakah dalam toleransi |
| `certificate_number` | String? | Nomor sertifikat |
| `certificate_url` | String? | URL sertifikat |
| `notes` | String? | Catatan |
| `created_at` | DateTime | Default now() |

**Index:** `@@index([calibration_asset_id])`, `@@index([calibration_date])`

### `ict_customer_feedback` — Umpan balik pelanggan

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `source_type` | String | Jenis sumber: email, phone, form, meeting, other |
| `customer_name` | String? | Nama pelanggan |
| `customer_email` | String? | Email pelanggan |
| `subject` | String | Judul feedback |
| `description` | String | Deskripsi detail |
| `rating` | Int? | Rating 1-5 |
| `category` | String? | Kategori feedback |
| `status` | String | Default "new" |
| `assigned_to` | String? | FK → dat_user |
| `response` | String? | Respon |
| `responded_at` | DateTime? | Waktu respon |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([company_id])`, `@@index([status])`, `@@index([source_type])`, `@@index([created_at])`

### Enum

```prisma
enum objective_status {
  on_track
  at_risk
  behind
  achieved
  cancelled
}

enum supplier_status {
  active
  inactive
  suspended
  blacklisted
}

enum supplier_rating {
  excellent
  good
  satisfactory
  poor
  unacceptable
}

enum training_status {
  planned
  in_progress
  completed
  cancelled
}

enum training_type {
  onboarding
  refresher
  compliance
  technical
  safety
  awareness
  certification
}

enum calibration_status {
  due
  in_progress
  completed
  overdue
  out_of_tolerance
}
```

---

## 24. Store Management

> File: `ict_base/prisma/schema/ict_store.prisma`

### `ict_store_info` — Informasi store (one-to-one dengan location)

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `location_id` | String | FK → dat_location (unique, one-to-one) |
| `company_id` | String | FK → dat_company |
| `ip_address` | String? | IP address store |
| `phone_type` | String? | Tipe phone asset |
| `phone_number` | String? | Nomor telepon |
| `phone_imei` | String? | IMEI phone |
| `notes` | String? | Catatan |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([location_id])`, `@@index([company_id])`

### `ict_store_edc` — Daftar EDC

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `store_info_id` | String | FK → ict_store_info |
| `company_id` | String | FK → dat_company |
| `bank` | String | Nama bank |
| `type` | String | Tipe EDC |
| `connect_with` | String | LAN, COM, H2H |
| `connect_to` | String | SS, Cashier 1, Cashier 2 |
| `mid` | String? | Merchant ID |
| `tid` | String? | Terminal ID |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([store_info_id])`, `@@index([company_id])`

### `ict_store_account` — Daftar akun

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `store_info_id` | String | FK → ict_store_info |
| `company_id` | String | FK → dat_company |
| `type` | String | Jenis akun |
| `username` | String | Username |
| `password` | String? | Password |
| `phone_number` | String? | Nomor telepon |
| `sid` | String? | Settlement ID |
| `mid` | String? | Merchant ID |
| `uuid` | String? | Unique identifier |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([store_info_id])`, `@@index([company_id])`

### `ict_store_asset` — Aset store (SS, C1, C2)

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `store_info_id` | String | FK → ict_store_info |
| `company_id` | String | FK → dat_company |
| `asset_type` | String | SS (Self Service), C1 (Cashier 1), C2 (Cashier 2) |
| `pos` | String? | POS |
| `ram` | String? | RAM |
| `printer` | String? | Printer |
| `tahun` | Int? | Tahun |
| `detail` | String? | Detail tambahan |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([store_info_id, asset_type])`, `@@index([store_info_id])`, `@@index([company_id])`

### `ict_store_digisign` — Daftar DigiSign

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `store_info_id` | String | FK → ict_store_info |
| `company_id` | String | FK → dat_company |
| `mini_pc` | String? | Mini PC |
| `status_sewa` | String? | Status sewa |
| `sn_pc` | String? | Serial number PC |
| `tv_1` | String? | TV 1 |
| `sn_tv_1` | String? | Serial number TV 1 |
| `tv_2` | String? | TV 2 |
| `sn_tv_2` | String? | Serial number TV 2 |
| `tv_3` | String? | TV 3 |
| `sn_tv_3` | String? | Serial number TV 3 |
| `detail` | String? | Detail tambahan |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@index([store_info_id])`, `@@index([company_id])`

---

## 25. Preventive Maintenance — WiFi & Vicon Check

> File: `ict_base/prisma/schema/ict_preventive.prisma`

Model tambahan untuk preventive maintenance harian pada lokasi HO/Regional.

### `ict_pm_wifi_check` — Pemeriksaan WiFi harian

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `schedule_id` | String? | FK → ict_pm_schedule |
| `room_id` | String | FK → ict_pm_room |
| `check_date` | DateTime | Tanggal pemeriksaan |
| `rssi` | Decimal? | RSSI signal strength |
| `snr` | Decimal? | Signal-to-Noise Ratio |
| `latency` | Decimal? | Latency (ms) |
| `overlap` | Decimal? | Channel overlap (%) |
| `notes` | String? | Catatan |
| `checked_by` | String? | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, room_id, check_date])`, `@@index([company_id])`, `@@index([room_id])`, `@@index([check_date])`

### `ict_pm_vicon_check` — Pemeriksaan ruang meeting/vicon harian

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `schedule_id` | String? | FK → ict_pm_schedule |
| `room_id` | String | FK → ict_pm_room |
| `check_date` | DateTime | Tanggal pemeriksaan |
| `pc_vicon` | String? | Status PC Vicon |
| `infocus` | String? | Status Infocus |
| `camera` | String? | Status Camera |
| `sound` | String? | Status Sound |
| `notes` | String? | Catatan |
| `checked_by` | String? | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, room_id, check_date])`, `@@index([company_id])`, `@@index([room_id])`, `@@index([check_date])`

### `ict_pm_store_check` — Pemeriksaan store harian

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `schedule_id` | String? | FK → ict_pm_schedule |
| `location_id` | String | FK → dat_location |
| `check_date` | DateTime | Tanggal pemeriksaan |
| `opening_cash` | Decimal? | Saldo awal |
| `closing_cash` | Decimal? | Saldo akhir |
| `actual_cash` | Decimal? | Saldo aktual |
| `diff_cash` | Decimal? | Selisih |
| `edc_total` | Decimal? | Total EDC |
| `qris_total` | Decimal? | Total QRIS |
| `notes` | String? | Catatan |
| `checked_by` | String? | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, location_id, check_date])`, `@@index([company_id])`, `@@index([location_id])`, `@@index([check_date])`

### `ict_pm_warehouse_check` — Pemeriksaan gudang harian

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `schedule_id` | String? | FK → ict_pm_schedule |
| `location_id` | String | FK → dat_location |
| `check_date` | DateTime | Tanggal pemeriksaan |
| `cleanliness` | String? | Kebersihan |
| `stock_accuracy` | String? | Akurasi stok |
| `fire_safety` | String? | Keselamatan api |
| `access_control` | String? | Kontrol akses |
| `notes` | String? | Catatan |
| `checked_by` | String? | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, location_id, check_date])`, `@@index([company_id])`, `@@index([location_id])`, `@@index([check_date])`

### `ict_pm_production_check` — Pemeriksaan produksi harian

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `company_id` | String | FK → dat_company |
| `schedule_id` | String? | FK → ict_pm_schedule |
| `location_id` | String | FK → dat_location |
| `check_date` | DateTime | Tanggal pemeriksaan |
| `machine_status` | String? | Status mesin |
| `output_quality` | String? | Kualitas output |
| `safety_check` | String? | Pemeriksaan keselamatan |
| `maintenance_due` | String? | Jadwal maintenance |
| `notes` | String? | Catatan |
| `checked_by` | String? | FK → dat_user |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([company_id, location_id, check_date])`, `@@index([company_id])`, `@@index([location_id])`, `@@index([check_date])`

---

## 26. Mikrotik — Port Group, VLAN Group, Client & Access List

> File: `ict_base/prisma/schema/ict_mikrotik.prisma`

### `ict_mikrotik_port_group` — Group port firewall

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_mikrotik_device |
| `company_id` | String | FK → dat_company |
| `ros_id` | Int | ID dari Mikrotik ROS (input user) |
| `name` | String | Nama group |
| `port` | String | Nomor port |
| `description` | String? | Deskripsi |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([device_id, name])`, `@@index([device_id])`, `@@index([company_id])`

### `ict_mikrotik_vlan_group` — Group VLAN

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_mikrotik_device |
| `company_id` | String | FK → dat_company |
| `ros_id` | Int | ID dari Mikrotik ROS |
| `vlan_id` | Int | VLAN ID (numeric) |
| `name` | String | Nama group |
| `description` | String? | Deskripsi |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([device_id, vlan_id])`, `@@index([device_id])`, `@@index([company_id])`

### `ict_mikrotik_client` — Client IP

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_mikrotik_device |
| `company_id` | String | FK → dat_company |
| `client_ip` | String | IP address client |
| `device_type` | String | Laptop, Computer, Ponsel, Server, Bastion |
| `name` | String | Nama client |
| `status` | String | Default "active" |
| `ros_id` | Int? | ID dari Mikrotik ROS |
| `ros_name` | String? | Nama dari Mikrotik |
| `ros_comment` | String? | Comment dari Mikrotik |
| `ros_status` | String? | Status dari Mikrotik |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([device_id, client_ip])`, `@@index([device_id])`, `@@index([company_id])`, `@@index([client_ip])`, `@@index([device_type])`

### `ict_mikrotik_access_normalization` — Normalisasi VLAN Group + Port Group

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_mikrotik_device |
| `company_id` | String | FK → dat_company |
| `vlan_group_id` | String | FK → ict_mikrotik_vlan_group |
| `port_group_id` | String? | FK → ict_mikrotik_port_group |
| `status` | String | Default "active" |
| `ros_id` | Int? | ID dari Mikrotik ROS |
| `ros_name` | String? | Nama dari Mikrotik |
| `ros_comment` | String? | Comment dari Mikrotik |
| `ros_status` | String? | Status dari Mikrotik |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([device_id, vlan_group_id, port_group_id])`, `@@index([device_id])`, `@@index([company_id])`, `@@index([vlan_group_id])`, `@@index([port_group_id])`

### `ict_mikrotik_access_list` — Daftar akses client → normalisasi

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | String (UUID) | PK |
| `device_id` | String | FK → ict_mikrotik_device |
| `company_id` | String | FK → dat_company |
| `client_id` | String | FK → ict_mikrotik_client |
| `access_normalization_id` | String | FK → ict_mikrotik_access_normalization |
| `status` | String | Default "active" |
| `ros_id` | Int? | ID dari Mikrotik ROS |
| `ros_name` | String? | Nama dari Mikrotik |
| `ros_comment` | String? | Comment dari Mikrotik |
| `ros_status` | String? | Status dari Mikrotik |
| `created_at` | DateTime | Default now() |
| `updated_at` | DateTime | Default now() |

**Index:** `@@unique([device_id, client_id, access_normalization_id])`, `@@index([device_id])`, `@@index([company_id])`, `@@index([client_id])`, `@@index([access_normalization_id])`
