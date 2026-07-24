# OmniSight — Referensi Utama

Dokumen ini adalah **referensi utama** untuk seluruh pengembangan, review, dan debug OmniSight.

---

## Daftar Isi

1. [Status Proyek](#status-proyek)
2. [Kode Halaman (Module Code)](#kode-halaman-module-code)
3. [Struktur Backend (ict_rest)](#struktur-backend-ict_rest)
4. [Struktur Frontend (ict_site)](#struktur-frontend-ict_site)
5. [Database Schema (ict_base)](#database-schema-ict_base)
6. [Sistem Akses (Privilege)](#sistem-akses-privilege)
7. [Tech Stack](#tech-stack)
8. [Catatan Teknis](#catatan-teknis)
9. [Update Terakhir](#update-terakhir)

---

## Status Proyek

| Aspek | Status |
|-------|--------|
| **Fase Saat Ini** | Phase 1 — Fondasi (MVP): Location Management + IT Infrastructure |
| **Total Task** | 380 |
| **Selesai** | 176 |
| **Progress** | 46.3% |

---

## Kode Halaman (Module Code)

Setiap halaman WAJIB memiliki kode unik. Kode ini digunakan di:
- **Backend**: nama folder skeleton (`ict_rest/skeleton/{KODE}/`), route path (`/rest/pages/{KODE}`)
- **Frontend**: path halaman (`/board/pages/{KODE}`), sidebar navigation
- **Database**: `dat_module.code` dan `dat_user_privilege.module_id`

### Konvensi Penamaan

| Elemen | Format | Contoh |
|--------|--------|--------|
| **Group/Modul** | 2 huruf kapital | `SP`, `SM`, `AM` |
| **Halaman Turunan** | Group + 2 digit | `SP01`, `SM01`, `AM01` |
| **Route Backend** | `/rest/pages/{KODE}` | `/rest/pages/AM01` |
| **Route Frontend** | `/board/pages/{KODE}` | `/board/pages/AM01` |
| **Skeleton Folder** | `ict_rest/skeleton/{KODE}/` | `ict_rest/skeleton/AM01/` |

### Level Privilege (`dat_action_type`)

| Level | HTTP Method | Keterangan |
|-------|-------------|------------|
| `hide` | — | Halaman tidak terlihat di sidebar |
| `view` | `GET` | Hanya bisa melihat data |
| `book` | `GET` + `POST` + `PUT` reservasi | Melihat + menandai/reservasi |
| `post` | `GET` + `POST` + `PUT` + `DELETE` posting | Full CRUD + Posting |

### Aturan Middleware

| Tipe Halaman | Privilege Minimum | Contoh |
|--------------|-------------------|--------|
| All User | `have session` (USLoad) | SP00–SP03 |
| Admin Only | `isAdmin = true` (USLock) | SM01–SM04, XX01 |
| CRUD Page | `post` | AM01, VL02, IN01 |
| View Only | `view` | AM03, VL05, DB01 |
| Tree/Hierarchy | `view` | AM02 |

### 1. SP — Session Profile

> **Group**: Session Profile | **Icon**: Users | **Akses**: Semua user login

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **SP** | Session Profile *(group)* | — | — | — | — | ✅ |
| SP01 | Profile | `/rest/pages/SP01` | `/board/pages/SP01` | view | GET | ✅ |
| SP02 | Change Password | `/rest/pages/SP02` | `/board/pages/SP02` | post | PUT | ✅ |
| SP03 | Action History | `/rest/pages/SP03` | `/board/pages/SP03` | view | GET | ✅ |

### 2. SM — System Manager

> **Group**: System Manager | **Icon**: Settings | **Akses**: Admin only (`USLock`)

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **SM** | System Manager *(group)* | — | — | — | — | ✅ |
| SM01 | User Management | `/rest/pages/SM01` | `/board/pages/SM01` | admin | GET/POST/PUT/DELETE | ✅ |
| SM02 | Module Management | `/rest/pages/SM02` | `/board/pages/SM02` | admin | GET/POST/PUT | ✅ |
| SM03 | Company Management | `/rest/pages/SM03` | `/board/pages/SM03` | admin | GET/POST/PUT | ✅ |
| SM04 | Signature/Doc Type | `/rest/pages/SM04` | `/board/pages/SM04` | admin | GET/POST/PUT | ✅ |
| SM05 | Session Management | `/rest/pages/SM05` | `/board/pages/SM05` | admin | GET/DELETE | ✅ |
| SM06 | Location Management | `/rest/pages/SM06` | `/board/pages/SM06` | post | GET/POST/PUT/DELETE | ✅ |

### 3. AM — Asset Management

> **Group**: Asset Management | **Icon**: MapPin | **Akses**: User privilege (`dat_user_privilege`)

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **AM** | Asset Management *(group)* | — | — | — | — | ⬜ |
| AM01 | Asset Category | `/rest/pages/AM01` | `/board/pages/AM01` | post | GET/POST/PUT/DELETE | ⬜ |
| AM02 | Asset List | `/rest/pages/AM02` | `/board/pages/AM02` | post | GET/POST/PUT/DELETE | ⬜ |
| AM03 | Asset Detail | `/rest/pages/AM03` | `/board/pages/AM03` | view | GET | ⬜ |
| AM04 | IP Management | `/rest/pages/AM04` | `/board/pages/AM04` | post | GET/POST/PUT/DELETE | ⬜ |
| AM05 | Asset Relations | `/rest/pages/AM05` | `/board/pages/AM05` | post | GET/POST/PUT/DELETE | ⬜ |

### 4. VL — Vulnerability Management

> **Group**: Vulnerability Management | **Icon**: Shield | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **VL** | Vulnerability *(group)* | — | — | — | — | ⬜ |
| VL01 | Scanner Types | `/rest/pages/VL01` | `/board/pages/VL01` | post | GET/POST/PUT/DELETE | ⬜ |
| VL02 | Scan Jobs | `/rest/pages/VL02` | `/board/pages/VL02` | post | GET/POST/PUT/DELETE | ⬜ |
| VL03 | Findings | `/rest/pages/VL03` | `/board/pages/VL03` | post | GET/POST/PUT/DELETE | ⬜ |
| VL04 | Remediation Log | `/rest/pages/VL04` | `/board/pages/VL04` | post | GET/POST/PUT/DELETE | ⬜ |
| VL05 | Vuln Dashboard | `/rest/pages/VL05` | `/board/pages/VL05` | view | GET | ⬜ |

### 5. IN — Incident Management

> **Group**: Incident Management | **Icon**: AlertTriangle | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **IN** | Incident *(group)* | — | — | — | — | ⬜ |
| IN01 | Incident Types | `/rest/pages/IN01` | `/board/pages/IN01` | post | GET/POST/PUT/DELETE | ⬜ |
| IN02 | Incident List | `/rest/pages/IN02` | `/board/pages/IN02` | post | GET/POST/PUT/DELETE | ⬜ |
| IN03 | Incident Detail | `/rest/pages/IN03` | `/board/pages/IN03` | view | GET | ⬜ |
| IN04 | Timeline | `/rest/pages/IN04` | `/board/pages/IN04` | view | GET | ⬜ |
| IN05 | Incident Dashboard | `/rest/pages/IN05` | `/board/pages/IN05` | view | GET | ⬜ |

### 6. SI — SIEM / Log Monitoring

> **Group**: SIEM | **Icon**: Activity | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **SI** | SIEM *(group)* | — | — | — | — | ⬜ |
| SI01 | Log Sources | `/rest/pages/SI01` | `/board/pages/SI01` | post | GET/POST/PUT/DELETE | ⬜ |
| SI02 | Detection Rules | `/rest/pages/SI02` | `/board/pages/SI02` | post | GET/POST/PUT/DELETE | ⬜ |
| SI03 | Alert Management | `/rest/pages/SI03` | `/board/pages/SI03` | post | GET/POST/PUT/DELETE | ⬜ |
| SI04 | MITRE ATT&CK | `/rest/pages/SI04` | `/board/pages/SI04` | view | GET | ⬜ |
| SI05 | SIEM Dashboard | `/rest/pages/SI05` | `/board/pages/SI05` | view | GET | ⬜ |

### 7. RK — Risk Assessment

> **Group**: Risk Assessment | **Icon**: BarChart3 | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **RK** | Risk *(group)* | — | — | — | — | ⬜ |
| RK01 | Risk Categories | `/rest/pages/RK01` | `/board/pages/RK01` | post | GET/POST/PUT/DELETE | ⬜ |
| RK02 | Risk Register | `/rest/pages/RK02` | `/board/pages/RK02` | post | GET/POST/PUT/DELETE | ⬜ |
| RK03 | Risk Assessment | `/rest/pages/RK03` | `/board/pages/RK03` | post | GET/POST/PUT/DELETE | ⬜ |
| RK04 | Treatment Plans | `/rest/pages/RK04` | `/board/pages/RK04` | post | GET/POST/PUT/DELETE | ⬜ |
| RK05 | Risk Matrix | `/rest/pages/RK05` | `/board/pages/RK05` | view | GET | ⬜ |
| RK06 | Risk Dashboard | `/rest/pages/RK06` | `/board/pages/RK06` | view | GET | ⬜ |

### 8. CP — Compliance Monitoring (ISO 27001/Standar)

> **Group**: Compliance | **Icon**: CheckCircle | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **CP** | Compliance *(group)* | — | — | — | — | ⬜ |
| CP01 | Standards | `/rest/pages/CP01` | `/board/pages/CP01` | post | GET/POST/PUT/DELETE | ⬜ |
| CP02 | Controls | `/rest/pages/CP02` | `/board/pages/CP02` | post | GET/POST/PUT/DELETE | ⬜ |
| CP03 | Assessments | `/rest/pages/CP03` | `/board/pages/CP03` | post | GET/POST/PUT/DELETE | ⬜ |
| CP04 | Findings | `/rest/pages/CP04` | `/board/pages/CP04` | post | GET/POST/PUT/DELETE | ⬜ |
| CP05 | Compliance Dashboard | `/rest/pages/CP05` | `/board/pages/CP05` | view | GET | ⬜ |
| CP06 | Statement of Applicability | `/rest/pages/CP06` | `/board/pages/CP06` | post | GET/POST/PUT | ⬜ |

### 9. CA — CAPA (Corrective and Preventive Action)

> **Group**: CAPA | **Icon**: AlertTriangle | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **CA** | CAPA *(group)* | — | — | — | — | ⬜ |
| CA01 | CAPA List | `/rest/pages/CA01` | `/board/pages/CA01` | post | GET/POST/PUT/DELETE | ⬜ |
| CA02 | CAPA Dashboard | `/rest/pages/CA02` | `/board/pages/CA02` | view | GET | ⬜ |

### 10. IA — Internal Audit

> **Group**: Internal Audit | **Icon**: ClipboardCheck | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **IA** | Internal Audit *(group)* | — | — | — | — | ⬜ |
| IA01 | Audit Programs | `/rest/pages/IA01` | `/board/pages/IA01` | post | GET/POST/PUT/DELETE | ⬜ |
| IA02 | Audit Instances | `/rest/pages/IA02` | `/board/pages/IA02` | post | GET/POST/PUT/DELETE | ⬜ |
| IA03 | Audit Dashboard | `/rest/pages/IA03` | `/board/pages/IA03` | view | GET | ⬜ |

### 11. MR — Management Review

> **Group**: Management Review | **Icon**: Users | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **MR** | Management Review *(group)* | — | — | — | — | ⬜ |
| MR01 | Review Meetings | `/rest/pages/MR01` | `/board/pages/MR01` | post | GET/POST/PUT/DELETE | ⬜ |
| MR02 | Review Dashboard | `/rest/pages/MR02` | `/board/pages/MR02` | view | GET | ⬜ |

### 12. QM — Quality Management System (QMS)

> **Group**: Quality Management | **Icon**: Award | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **QM** | Quality Management *(group)* | — | — | — | — | ⬜ |
| QM01 | Quality Objectives | `/rest/pages/QM01` | `/board/pages/QM01` | post | GET/POST/PUT/DELETE | ⬜ |
| QM02 | Quality KPIs | `/rest/pages/QM02` | `/board/pages/QM02` | post | GET/POST/PUT/DELETE | ⬜ |
| QM03 | Supplier Management | `/rest/pages/QM03` | `/board/pages/QM03` | post | GET/POST/PUT/DELETE | ⬜ |
| QM04 | Training Programs | `/rest/pages/QM04` | `/board/pages/QM04` | post | GET/POST/PUT/DELETE | ⬜ |
| QM05 | Training Records | `/rest/pages/QM05` | `/board/pages/QM05` | post | GET/POST/PUT/DELETE | ⬜ |
| QM06 | Calibration Assets | `/rest/pages/QM06` | `/board/pages/QM06` | post | GET/POST/PUT/DELETE | ⬜ |
| QM07 | Customer Feedback | `/rest/pages/QM07` | `/board/pages/QM07` | post | GET/POST/PUT/DELETE | ⬜ |
| QM08 | QMS Dashboard | `/rest/pages/QM08` | `/board/pages/QM08` | view | GET | ⬜ |

### 13. CT — Certificate Management

> **Group**: Certificate | **Icon**: Key | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **CT** | Certificate *(group)* | — | — | — | — | ⬜ |
| CT01 | Certificate Types | `/rest/pages/CT01` | `/board/pages/CT01` | post | GET/POST/PUT/DELETE | ⬜ |
| CT02 | Certificate List | `/rest/pages/CT02` | `/board/pages/CT02` | post | GET/POST/PUT/DELETE | ⬜ |
| CT03 | Expiring Soon | `/rest/pages/CT03` | `/board/pages/CT03` | view | GET | ⬜ |
| CT04 | Certificate Dashboard | `/rest/pages/CT04` | `/board/pages/CT04` | view | GET | ⬜ |

### 14. DK — Docker Management

> **Group**: Docker | **Icon**: Container | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **DK** | Docker *(group)* | — | — | — | — | ⬜ |
| DK01 | Docker Hosts | `/rest/pages/DK01` | `/board/pages/DK01` | post | GET/POST/PUT/DELETE | ⬜ |
| DK02 | Containers | `/rest/pages/DK02` | `/board/pages/DK02` | post | GET/POST/PUT/DELETE | ⬜ |
| DK03 | Compose Stacks | `/rest/pages/DK03` | `/board/pages/DK03` | post | GET/POST/PUT/DELETE | ⬜ |
| DK04 | Images | `/rest/pages/DK04` | `/board/pages/DK04` | view | GET | ⬜ |
| DK05 | Networks | `/rest/pages/DK05` | `/board/pages/DK05` | view | GET | ⬜ |
| DK06 | Volumes | `/rest/pages/DK06` | `/board/pages/DK06` | view | GET | ⬜ |
| DK07 | Monitoring | `/rest/pages/DK07` | `/board/pages/DK07` | view | GET | ⬜ |
| DK08 | Alerts | `/rest/pages/DK08` | `/board/pages/DK08` | post | GET/POST/PUT/DELETE | ⬜ |
| DK09 | Docker Dashboard | `/rest/pages/DK09` | `/board/pages/DK09` | view | GET | ⬜ |

### 15. DM — Document Management

> **Group**: Document Management | **Icon**: FileText | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **DM** | Document *(group)* | — | — | — | — | ⬜ |
| DM01 | Document List | `/rest/pages/DM01` | `/board/pages/DM01` | post | GET/POST/PUT/DELETE | ⬜ |
| DM02 | Editor | `/rest/pages/DM02` | `/board/pages/DM02` | post | GET/POST/PUT | ⬜ |
| DM03 | Todo Lists | `/rest/pages/DM03` | `/board/pages/DM03` | post | GET/POST/PUT/DELETE | ⬜ |
| DM04 | Templates | `/rest/pages/DM04` | `/board/pages/DM04` | post | GET/POST/PUT/DELETE | ⬜ |
| DM05 | Approval Workflow | `/rest/pages/DM05` | `/board/pages/DM05` | post | GET/POST/PUT | ⬜ |
| DM06 | Sharing | `/rest/pages/DM06` | `/board/pages/DM06` | post | GET/POST/PUT/DELETE | ⬜ |
| DM07 | Publishing | `/rest/pages/DM07` | `/board/pages/DM07` | post | GET/POST/PUT | ⬜ |
| DM08 | Docman Dashboard | `/rest/pages/DM08` | `/board/pages/DM08` | view | GET | ⬜ |

### 16. SR — Service Request

> **Group**: Service Request | **Icon**: ClipboardList | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **SR** | Service Request *(group)* | — | — | — | — | ⬜ |
| SR01 | Service Catalog | `/rest/pages/SR01` | `/board/pages/SR01` | view | GET | ⬜ |
| SR02 | My Requests | `/rest/pages/SR02` | `/board/pages/SR02` | view | GET | ⬜ |
| SR03 | Submit Request | `/rest/pages/SR03` | `/board/pages/SR03` | post | GET/POST | ⬜ |
| SR04 | Admin Panel | `/rest/pages/SR04` | `/board/pages/SR04` | post | GET/POST/PUT/DELETE | ⬜ |
| SR05 | Approvals | `/rest/pages/SR05` | `/board/pages/SR05` | post | GET/POST/PUT | ⬜ |
| SR06 | Fulfillment | `/rest/pages/SR06` | `/board/pages/SR06` | post | GET/POST/PUT | ⬜ |
| SR07 | SR Dashboard | `/rest/pages/SR07` | `/board/pages/SR07` | view | GET | ⬜ |

### 17. BL — Billing & Invoice

> **Group**: Billing | **Icon**: CreditCard | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **BL** | Billing *(group)* | — | — | — | — | ⬜ |
| BL01 | Invoice List | `/rest/pages/BL01` | `/board/pages/BL01` | post | GET/POST/PUT/DELETE | ⬜ |
| BL02 | Create Invoice | `/rest/pages/BL02` | `/board/pages/BL02` | post | GET/POST | ⬜ |
| BL03 | Overdue Invoices | `/rest/pages/BL03` | `/board/pages/BL03` | view | GET | ⬜ |
| BL04 | Vendor Management | `/rest/pages/BL04` | `/board/pages/BL04` | post | GET/POST/PUT/DELETE | ⬜ |
| BL05 | Billing Categories | `/rest/pages/BL05` | `/board/pages/BL05` | post | GET/POST/PUT/DELETE | ⬜ |
| BL06 | Billing Dashboard | `/rest/pages/BL06` | `/board/pages/BL06` | view | GET | ⬜ |

### 18. PM — Preventive Maintenance

> **Group**: Preventive Maintenance | **Icon**: Wrench | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **PM** | Preventive Maintenance *(group)* | — | — | — | — | ⬜ |
| PM01 | Rooms | `/rest/pages/PM01` | `/board/pages/PM01` | post | GET/POST/PUT/DELETE | ⬜ |
| PM02 | Equipment | `/rest/pages/PM02` | `/board/pages/PM02` | post | GET/POST/PUT/DELETE | ⬜ |
| PM03 | Checklists | `/rest/pages/PM03` | `/board/pages/PM03` | post | GET/POST/PUT/DELETE | ⬜ |
| PM04 | Check Calendar | `/rest/pages/PM04` | `/board/pages/PM04` | view | GET | ⬜ |
| PM05 | Issues | `/rest/pages/PM05` | `/board/pages/PM05` | post | GET/POST/PUT/DELETE | ⬜ |
| PM06 | PM Dashboard | `/rest/pages/PM06` | `/board/pages/PM06` | view | GET | ⬜ |

### 19. DB — Dashboard & Reporting

> **Group**: Dashboard | **Icon**: LayoutDashboard | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **DB** | Dashboard *(group)* | — | — | — | — | ⬜ |
| DB01 | Overview Dashboard | `/rest/pages/DB01` | `/board/pages/DB01` | view | GET | ⬜ |
| DB02 | Widget Management | `/rest/pages/DB02` | `/board/pages/DB02` | post | GET/POST/PUT/DELETE | ⬜ |
| DB03 | Report Templates | `/rest/pages/DB03` | `/board/pages/DB03` | post | GET/POST/PUT/DELETE | ⬜ |
| DB04 | Generate Reports | `/rest/pages/DB04` | `/board/pages/DB04` | post | GET/POST | ⬜ |
| DB05 | Report History | `/rest/pages/DB05` | `/board/pages/DB05` | view | GET | ⬜ |

### 16. XX — Utility

> **Group**: Utility | **Icon**: Tool | **Akses**: Admin only

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **XX** | Utility *(group)* | — | — | — | — | ⬜ |
| XX01 | Cleanup/Purge | `/rest/pages/XX01` | — | admin | POST | ✅ |

### 17. NK — Network Management (VLAN/Subnet/Route)

> **Group**: Network | **Icon**: Network | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **NK** | Network *(group)* | — | — | — | — | ⬜ |
| NK01 | VLAN Management | `/rest/pages/NK01` | `/board/pages/NK01` | post | GET/POST/PUT/DELETE | ⬜ |
| NK02 | Subnet Allocation | `/rest/pages/NK02` | `/board/pages/NK02` | post | GET/POST/PUT/DELETE | ⬜ |
| NK03 | Routing Table | `/rest/pages/NK03` | `/board/pages/NK03` | post | GET/POST/PUT/DELETE | ⬜ |
| NK04 | NAT Management | `/rest/pages/NK04` | `/board/pages/NK04` | post | GET/POST/PUT/DELETE | ⬜ |
| NK05 | Network Dashboard | `/rest/pages/NK05` | `/board/pages/NK05` | view | GET | ⬜ |

### 18. MK — Mikrotik Device Management

> **Group**: Mikrotik | **Icon**: Router | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **MK** | Mikrotik *(group)* | — | — | — | — | ⬜ |
| MK01 | Mikrotik Devices | `/rest/pages/MK01` | `/board/pages/MK01` | post | GET/POST/PUT/DELETE | ⬜ |
| MK02 | Hotspot Management | `/rest/pages/MK02` | `/board/pages/MK02` | post | GET/POST/PUT/DELETE | ⬜ |
| MK03 | DHCP Leases | `/rest/pages/MK03` | `/board/pages/MK03` | post | GET/POST/PUT/DELETE | ⬜ |
| MK04 | Firewall Rules | `/rest/pages/MK04` | `/board/pages/MK04` | post | GET/POST/PUT/DELETE | ⬜ |
| MK05 | Bandwidth Queue | `/rest/pages/MK05` | `/board/pages/MK05` | post | GET/POST/PUT/DELETE | ⬜ |
| MK06 | Port Group | `/rest/pages/MK06` | `/board/pages/MK06` | post | GET/POST/PUT/DELETE | ⬜ |
| MK07 | VLAN Group | `/rest/pages/MK07` | `/board/pages/MK07` | post | GET/POST/PUT/DELETE | ⬜ |
| MK08 | Mikrotik Client | `/rest/pages/MK08` | `/board/pages/MK08` | post | GET/POST/PUT/DELETE | ⬜ |
| MK09 | Access Normalization | `/rest/pages/MK09` | `/board/pages/MK09` | post | GET/POST/PUT/DELETE | ⬜ |
| MK10 | Access List | `/rest/pages/MK10` | `/board/pages/MK10` | post | GET/POST/PUT/DELETE | ⬜ |

### 19. DN — Domain Management

> **Group**: Domain | **Icon**: Globe | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **DN** | Domain *(group)* | — | — | — | — | ⬜ |
| DN01 | Domain Registry | `/rest/pages/DN01` | `/board/pages/DN01` | post | GET/POST/PUT/DELETE | ⬜ |
| DN02 | DNS Records | `/rest/pages/DN02` | `/board/pages/DN02` | post | GET/POST/PUT/DELETE | ⬜ |
| DN03 | Domain Dashboard | `/rest/pages/DN03` | `/board/pages/DN03` | view | GET | ⬜ |

### 20. IP — IP Public & Block Management

> **Group**: IP Public | **Icon**: GlobeLock | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **IP** | IP Public *(group)* | — | — | — | — | ⬜ |
| IP01 | IP Blocks | `/rest/pages/IP01` | `/board/pages/IP01` | post | GET/POST/PUT/DELETE | ⬜ |
| IP02 | IP Allocation | `/rest/pages/IP02` | `/board/pages/IP02` | post | GET/POST/PUT/DELETE | ⬜ |
| IP03 | NAT Mapping | `/rest/pages/IP03` | `/board/pages/IP03` | post | GET/POST/PUT/DELETE | ⬜ |
| IP04 | IP Dashboard | `/rest/pages/IP04` | `/board/pages/IP04` | view | GET | ⬜ |

### 21. AM — Asset Management *(update)*

> **Group**: Asset Management | **Icon**: MapPin | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **AM** | Asset Management *(group)* | — | — | — | — | ⬜ |
| AM01 | Asset Category | `/rest/pages/AM01` | `/board/pages/AM01` | post | GET/POST/PUT/DELETE | ⬜ |
| AM02 | Asset List | `/rest/pages/AM02` | `/board/pages/AM02` | post | GET/POST/PUT/DELETE | ⬜ |
| AM03 | Asset Detail | `/rest/pages/AM03` | `/board/pages/AM03` | view | GET | ⬜ |
| AM04 | IP Management | `/rest/pages/AM04` | `/board/pages/AM04` | post | GET/POST/PUT/DELETE | ⬜ |
| AM05 | Asset Relations | `/rest/pages/AM05` | `/board/pages/AM05` | post | GET/POST/PUT/DELETE | ⬜ |
| AM06 | Asset Assignment | `/rest/pages/AM06` | `/board/pages/AM06` | post | GET/POST/PUT/DELETE | ⬜ |

### 22. MO — Monitor Agent & Discovery

> **Group**: Monitor Agent | **Icon**: Radio | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **MO** | Monitor Agent *(group)* | — | — | — | — | ⬜ |
| MO01 | Agent Groups | `/rest/pages/MO01` | `/board/pages/MO01` | post | GET/POST/PUT/DELETE | ⬜ |
| MO02 | Agent List | `/rest/pages/MO02` | `/board/pages/MO02` | post | GET/POST/PUT/DELETE | ⬜ |
| MO03 | Network Devices | `/rest/pages/MO03` | `/board/pages/MO03` | post | GET/POST/PUT/DELETE | ⬜ |
| MO04 | Discovery Jobs | `/rest/pages/MO04` | `/board/pages/MO04` | post | GET/POST/PUT/DELETE | ⬜ |
| MO05 | Discovery Results | `/rest/pages/MO05` | `/board/pages/MO05` | view | GET | ⬜ |
| MO06 | SNMP Traps | `/rest/pages/MO06` | `/board/pages/MO06` | view | GET | ⬜ |
| MO07 | Agent Dashboard | `/rest/pages/MO07` | `/board/pages/MO07` | view | GET | ⬜ |

### 23. MT — Metrics & Alerting

> **Group**: Metrics | **Icon**: Activity | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **MT** | Metrics *(group)* | — | — | — | — | ⬜ |
| MT01 | Metric Definitions | `/rest/pages/MT01` | `/board/pages/MT01` | post | GET/POST/PUT/DELETE | ⬜ |
| MT02 | Host Snapshots | `/rest/pages/MT02` | `/board/pages/MT02` | view | GET | ⬜ |
| MT03 | Disk Health | `/rest/pages/MT03` | `/board/pages/MT03` | view | GET | ⬜ |
| MT04 | GPU Monitoring | `/rest/pages/MT04` | `/board/pages/MT04` | view | GET | ⬜ |
| MT05 | Alert Rules | `/rest/pages/MT05` | `/board/pages/MT05` | post | GET/POST/PUT/DELETE | ⬜ |
| MT06 | Active Alerts | `/rest/pages/MT06` | `/board/pages/MT06` | post | GET/POST/PUT | ⬜ |
| MT07 | Alert Silence | `/rest/pages/MT07` | `/board/pages/MT07` | post | GET/POST/PUT/DELETE | ⬜ |
| MT08 | Notification Config | `/rest/pages/MT08` | `/board/pages/MT08` | post | GET/POST/PUT | ⬜ |
| MT09 | Metrics Dashboard | `/rest/pages/MT09` | `/board/pages/MT09` | view | GET | ⬜ |

### 24. ML — Log Management

> **Group**: Log Management | **Icon**: ScrollText | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **ML** | Log *(group)* | — | — | — | — | ⬜ |
| ML01 | Log Sources | `/rest/pages/ML01` | `/board/pages/ML01` | post | GET/POST/PUT/DELETE | ⬜ |
| ML02 | Log Search | `/rest/pages/ML02` | `/board/pages/ML02` | view | GET | ⬜ |
| ML03 | Decoders | `/rest/pages/ML03` | `/board/pages/ML03` | post | GET/POST/PUT/DELETE | ⬜ |
| ML04 | Log Rules | `/rest/pages/ML04` | `/board/pages/ML04` | post | GET/POST/PUT/DELETE | ⬜ |
| ML05 | Saved Queries | `/rest/pages/ML05` | `/board/pages/ML05` | post | GET/POST/PUT/DELETE | ⬜ |
| ML06 | Container Logs | `/rest/pages/ML06` | `/board/pages/ML06` | view | GET | ⬜ |
| ML07 | Log Dashboard | `/rest/pages/ML07` | `/board/pages/ML07` | view | GET | ⬜ |

### 25. MF — File Integrity Monitoring

> **Group**: FIM | **Icon**: FileCheck | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **MF** | FIM *(group)* | — | — | — | — | ⬜ |
| MF01 | FIM Configurations | `/rest/pages/MF01` | `/board/pages/MF01` | post | GET/POST/PUT/DELETE | ⬜ |
| MF02 | Baselines | `/rest/pages/MF02` | `/board/pages/MF02` | view | GET | ⬜ |
| MF03 | Changes Detected | `/rest/pages/MF03` | `/board/pages/MF03` | post | GET/POST/PUT | ⬜ |
| MF04 | FIM Alerts | `/rest/pages/MF04` | `/board/pages/MF04` | post | GET/POST/PUT | ⬜ |
| MF05 | FIM Dashboard | `/rest/pages/MF05` | `/board/pages/MF05` | view | GET | ⬜ |

### 26. MS — Security Configuration Assessment

> **Group**: SCA | **Icon**: ShieldCheck | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **MS** | SCA *(group)* | — | — | — | — | ⬜ |
| MS01 | SCA Policies | `/rest/pages/MS01` | `/board/pages/MS01` | post | GET/POST/PUT/DELETE | ⬜ |
| MS02 | Policy Checks | `/rest/pages/MS02` | `/board/pages/MS02` | post | GET/POST/PUT/DELETE | ⬜ |
| MS03 | Scan Results | `/rest/pages/MS03` | `/board/pages/MS03` | view | GET | ⬜ |
| MS04 | Compliance Score | `/rest/pages/MS04` | `/board/pages/MS04` | view | GET | ⬜ |
| MS05 | SCA Dashboard | `/rest/pages/MS05` | `/board/pages/MS05` | view | GET | ⬜ |

### 27. MD — Docker Extensions (Dockge-like)

> **Group**: Docker Extensions | **Icon**: Container | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **MD** | Docker Ext *(group)* | — | — | — | — | ⬜ |
| MD01 | Compose Versions | `/rest/pages/MD01` | `/board/pages/MD01` | post | GET/POST/PUT/DELETE | ⬜ |
| MD02 | Compose Diff | `/rest/pages/MD02` | `/board/pages/MD02` | view | GET | ⬜ |
| MD03 | Web Terminal | `/rest/pages/MD03` | `/board/pages/MD03` | post | GET/POST/DELETE | ⬜ |
| MD04 | Env Overrides | `/rest/pages/MD04` | `/board/pages/MD04` | post | GET/POST/PUT/DELETE | ⬜ |
| MD05 | Docker Ext Dashboard | `/rest/pages/MD05` | `/board/pages/MD05` | view | GET | ⬜ |

### 28. IS — ISO 27001 & SOP

> **Group**: ISO Compliance | **Icon**: Award | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **IS** | ISO *(group)* | — | — | — | — | ⬜ |
| IS01 | Policies | `/rest/pages/IS01` | `/board/pages/IS01` | post | GET/POST/PUT/DELETE | ⬜ |
| IS02 | SOP List | `/rest/pages/IS02` | `/board/pages/IS02` | post | GET/POST/PUT/DELETE | ⬜ |
| IS03 | SOP Detail | `/rest/pages/IS03` | `/board/pages/IS03` | view | GET | ⬜ |
| IS04 | SOP Execution | `/rest/pages/IS04` | `/board/pages/IS04` | post | GET/POST/PUT | ⬜ |
| IS05 | Evidence Collection | `/rest/pages/IS05` | `/board/pages/IS05` | post | GET/POST/PUT/DELETE | ⬜ |
| IS06 | Security Metrics | `/rest/pages/IS06` | `/board/pages/IS06` | post | GET/POST/PUT/DELETE | ⬜ |
| IS07 | ISO Reports | `/rest/pages/IS07` | `/board/pages/IS07` | post | GET/POST | ⬜ |
| IS08 | Document Register | `/rest/pages/IS08` | `/board/pages/IS08` | post | GET/POST/PUT/DELETE | ⬜ |
| IS09 | ISO Dashboard | `/rest/pages/IS09` | `/board/pages/IS09` | view | GET | ⬜ |

### 29. WM — Web Monitoring (UptimeRobot)

> **Group**: Web Monitor | **Icon**: Globe | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **WM** | Web Monitor *(group)* | — | — | — | — | ⬜ |
| WM01 | UptimeRobot Monitors | `/rest/pages/WM01` | `/board/pages/WM01` | post | GET/POST/PUT/DELETE | ⬜ |
| WM02 | Incident Logs | `/rest/pages/WM02` | `/board/pages/WM02` | view | GET | ⬜ |
| WM03 | Uptime Dashboard | `/rest/pages/WM03` | `/board/pages/WM03` | view | GET | ⬜ |

### 30. WS — Web Security (Nginx/WAF)

> **Group**: Web Security | **Icon**: Shield | **Akses**: User privilege

| Kode | Nama | Route Backend | Route Frontend | Privilege | HTTP | Status |
|------|------|---------------|----------------|-----------|------|--------|
| **WS** | Web Security *(group)* | — | — | — | — | ⬜ |
| WS01 | Nginx Logs | `/rest/pages/WS01` | `/board/pages/WS01` | view | GET | ⬜ |
| WS02 | Attack Logs | `/rest/pages/WS02` | `/board/pages/WS02` | view | GET | ⬜ |
| WS03 | IP Whitelist | `/rest/pages/WS03` | `/board/pages/WS03` | post | GET/POST/PUT/DELETE | ⬜ |
| WS04 | IP Blacklist | `/rest/pages/WS04` | `/board/pages/WS04` | post | GET/POST/PUT/DELETE | ⬜ |
| WS05 | WAF Bypass Rules | `/rest/pages/WS05` | `/board/pages/WS05` | post | GET/POST/PUT/DELETE | ⬜ |
| WS06 | SLA Reports | `/rest/pages/WS06` | `/board/pages/WS06` | view | GET | ⬜ |
| WS07 | Security Dashboard | `/rest/pages/WS07` | `/board/pages/WS07` | view | GET | ⬜ |

### Ringkasan Module Code

| Group | Kode | Jumlah Halaman | Selesai |
|-------|------|----------------|---------|
| Session Profile | SP | 3 | 3/3 ✅ |
| System Manager | SM | 6 | 6/6 ✅ |
| Asset Management | AM | 6 | 0/6 |
| Network | NK | 5 | 0/5 |
| Mikrotik | MK | 10 | 0/10 |
| Domain | DN | 3 | 0/3 |
| IP Public | IP | 4 | 0/4 |
| Vulnerability | VL | 5 | 0/5 |
| Incident | IN | 5 | 0/5 |
| SIEM | SI | 5 | 0/5 |
| Risk | RK | 6 | 0/6 |
| Compliance | CP | 6 | 0/6 |
| Certificate | CT | 4 | 0/4 |
| Docker | DK | 9 | 0/9 |
| Document | DM | 8 | 0/8 |
| Service Request | SR | 7 | 0/7 |
| Billing | BL | 6 | 0/6 |
| Preventive Maintenance | PM | 6 | 0/6 |
| Dashboard | DB | 5 | 0/5 |
| Monitor Agent | MO | 7 | 0/7 |
| Metrics & Alerting | MT | 9 | 0/9 |
| Log Management | ML | 7 | 0/7 |
| FIM | MF | 5 | 0/5 |
| SCA | MS | 5 | 0/5 |
| Docker Extensions | MD | 5 | 0/5 |
| ISO & SOP | IS | 9 | 0/9 |
| Web Monitor | WM | 3 | 0/3 |
| Web Security | WS | 7 | 0/7 |
| Utility | XX | 1 | 1/1 ✅ |
| **Total** | | **164** | **9/164** |

---

## Struktur Backend (ict_rest)

### Arsitektur

```
ict_rest/
├── main.go                    Entry point (port 36665, graceful shutdown)
├── backbone/                  Infrastructure layer
│   ├── database.go            PostgreSQL connection
│   ├── routes.go              Gin router, CORS, DI wiring
│   ├── memory.go              Session middleware (USLoad, USLock, USRole, USLogs)
│   ├── upload.go              File upload handler (20MB)
│   ├── logger.go              zerolog structured logging + RequestID
│   ├── recovery.go            CustomRecovery with stack trace
│   └── cleanup.go             CleanupExpiredSessions
├── mechanic/                  Helper layer
│   ├── helper.go              AppError + Error() handler
│   ├── crypto.go              AES-GCM + NullableString
│   └── typography.go          Text utilities (CheckMeta)
├── skeleton/                  Linear Layer Pattern (per halaman)
│   ├── SP00/                  Login/Auth
│   ├── SP01/                  Session Profile / Board Data
│   ├── SP02/                  Change Password
│   ├── SP03/                  User Action History
│   ├── SM01/                  User Management
│   ├── SM02/                  Module Management
│   ├── SM03/                  Company Management
│   ├── SM04/                  Signature/Approval Type
│   ├── SM05/                  Session Management
│   └── XX99/                  Utility/Cleanup
└── Dockerfile                 Multi-stage build
```

### Architecture Notes

- **Port**: 36665
- **Auth**: Bearer token via `dat_user_session`
- **Groups**: `/rest/guest` (public), `/rest/pages` (auth via USLoad), `/rest/pages` admin sub-group (auth + USLock)
- **Pattern**: Linear Layer Execution (template → repository → usecase → handler)
- **DI**: `SM0X.NRepo(PgSQL)` → `SM0X.NCase(repo)` → `SM0X.NHand(uc)`
- **Middleware Order**: RequestID → CustomRecovery → Logger → CORS → [USLoad → USRole → USLock → USLogs]
- **USRole**: Sudah diimplementasi di `backbone/memory.go`, belum di-wiring ke route manapun
- **USLogs**: Middleware audit trail untuk mutasi (POST/PUT/DELETE)

### Endpoint Mapping

| Route | Method | Description |
|-------|--------|-------------|
| `/rest/guest/SP00` | GET | List HRIS companies |
| `/rest/guest/SP00` | POST | Login |
| `/rest/pages/SP00` | DELETE | Logout |
| `/rest/pages/SP01/company` | GET | User companies |
| `/rest/pages/SP01/module` | GET | User modules |
| `/rest/pages/SP02` | PUT | Change password |
| `/rest/pages/SP03` | GET | User action history |
| `/rest/pages/admin/SM01` | GET/POST/PUT/DELETE | User CRUD |
| `/rest/pages/admin/SM01/:id/company` | GET/POST/DELETE | User company assignment |
| `/rest/pages/admin/SM01/:id/privilege` | GET/POST/PUT/DELETE | User privilege assignment |
| `/rest/pages/admin/SM01/:id/location` | GET/POST/PUT/DELETE | User location assignment |
| `/rest/pages/admin/SM01/area` | GET | Location areas by company |
| `/rest/pages/admin/SM01/hris-companies` | GET | HRIS companies list |
| `/rest/pages/admin/SM02` | GET/POST/PUT | Module CRUD |
| `/rest/pages/admin/SM03` | GET/POST/PUT | Company CRUD |
| `/rest/pages/admin/SM03/:id/module` | GET/POST/DELETE | Company module assignment |
| `/rest/pages/admin/SM03/:id/location` | GET/POST/PUT/DELETE | Company location area CRUD |
| `/rest/pages/admin/SM04` | GET/POST/PUT | Signature type CRUD |
| `/rest/pages/admin/SM04/users` | GET | Search users for signer assignment |
| `/rest/pages/admin/SM05` | GET | List active sessions |
| `/rest/pages/admin/SM05/:id` | GET | Find session by ID |
| `/rest/pages/admin/SM05/:id` | DELETE | Revoke/invalidate session |
| `/rest/pages/admin/XX99` | POST | Cleanup expired sessions |

### Middleware Flow

```
Request → RequestID → CustomRecovery → Logger → CORS
         → USLoad()          (authenticated routes)
         → USRole()          (per-module privilege check)
         → USLock()          (admin-only routes)
         → USLogs()          (audit trail for mutations)
         → Handler
```

---

## Struktur Frontend (ict_site)

### Arsitektur

```
ict_site/src/
├── app/
│   ├── layout.tsx             Root layout (ErrorBoundary)
│   ├── page.tsx               Home page
│   ├── login/page.tsx         Login (react-hook-form + zod)
│   ├── proxy/guest/           Guest proxy route
│   └── board/
│       ├── layout.tsx         Board layout (SidebarProvider)
│       ├── page.tsx           Dashboard page
│       ├── model/module.ts    Nav types + useModuleItem()
│       ├── widget/            Sidebar, breadcrumb, company, dashboard
│       └── pages/
│           ├── SP01/          Profile
│           ├── SP02/          Change Password
│           ├── SP03/          Action History
│           ├── SM01/          User Management
│           ├── SM02/          Module Management
│           ├── SM03/          Company Management
│           ├── SM04/          Signature Type
│           └── SM05/          Session Management
├── lib/                       17 files (API client, hooks, utilities)
│   ├── client-api.ts          API client (code + requestId)
│   ├── server-api.ts          Server-side API client
│   ├── error-message.ts       Error code → user message
│   ├── storage.ts             LocalStorage + preference helpers
│   ├── view-mobile.ts         useIsMobile() hook
│   └── use-*.ts               9 custom React hooks (debounce, media-query, etc.)
└── uix/                       81 components
```

### Tech Stack

| Package | Version | Status |
|---------|---------|--------|
| Next.js | 16.2.10 | ✅ Digunakan |
| React | 19.2.4 | ✅ Digunakan |
| TypeScript | 5.x | ✅ Digunakan |
| Tailwind CSS | 4.x | ✅ Digunakan |
| shadcn/ui | @shadcn/react 0.2.0 | ✅ Digunakan |
| Zustand | 5.x | ✅ Digunakan |
| Zod | 4.x | ✅ Digunakan |
| React Hook Form | 7.80.0 | ✅ Digunakan |
| Sonner | 2.x | ✅ Digunakan |

### Backend Connection

- **API Base**: `/rest/pages` (via `BE_POOL`)
- **Port**: 36666 (dev: `next dev -p 36666`)
- **Backend**: `http://ict_rest:36665`

---

## Database Schema (ict_base)

### Schema Yang Sudah Ada

| File | Model | Status |
|------|-------|--------|
| `schema.prisma` | Main config | ✅ Existing |
| `dat_company.prisma` | `dat_company`, `dat_company_module` | ✅ Existing |
| `dat_user.prisma` | `dat_action_type`, `dat_user`, `dat_user_company`, `dat_user_location`, `dat_user_privilege`, `dat_user_session`, `dat_user_action` | ✅ Existing |
| `dat_module.prisma` | `dat_module` | ✅ Existing |
| `dat_signature.prisma` | Approval workflow | ✅ Existing |
| `dat_location.prisma` | `dat_location`, `dat_location_type`, `dat_location_detail`, `dat_location_contact`, `dat_location_asset` | ✅ Existing |
| `ict_web_security.prisma` | Security logs | ✅ Existing |
| `ict_web_monitor.prisma` | Monitoring | ✅ Existing |

### Schema Yang Direncanakan

| File | Modul | Jumlah Tabel | Status |
|------|-------|--------------|--------|
| `ict_asset.prisma` | Asset/Inventory | 5 | ✅ (+assignment) |
| `ict_network.prisma` | Network (VLAN/Subnet/Route) | 6 | ✅ **BARU** |
| `ict_mikrotik.prisma` | Mikrotik Device | 11 | ✅ **BARU** |
| `ict_domain.prisma` | Domain Management | 3 | ✅ **BARU** |
| `ict_ip_public.prisma` | IP Public & Block | 3 | ✅ **BARU** |
| `ict_vulnerability.prisma` | Vulnerability | 4 | ⬜ |
| `ict_incident.prisma` | Incident | 4 | ⬜ |
| `ict_siem.prisma` | SIEM | 4 | ⬜ |
| `ict_risk.prisma` | Risk | 4 | ⬜ |
| `ict_compliance.prisma` | Compliance + SoA + CAPA + Audit + Mgmt Review | 14 | ✅ |
| `ict_certificate.prisma` | Certificate | 3 | ⬜ |
| `ict_dashboard.prisma` | Dashboard + Report Templates | 7 | ✅ (updated) |
| `ict_document.prisma` | Document Management | 11 | ✅ |
| `ict_docker.prisma` | Docker Management | 14 | ⬜ |
| `ict_service_request.prisma` | Service Request | 8 | ⬜ |
| `ict_billing.prisma` | Billing & Invoice | 4 | ⬜ |
| `ict_preventive.prisma` | Preventive Maintenance | 13 | ✅ |
| `ict_qms.prisma` | Quality Management (QMS) | 9 | ✅ |
| `ict_generic.prisma` | Generic (shared tables) | 4 | ✅ |
| `ict_web_security.prisma` | Nginx/WAF Security (WS) | 8 | ⬜ |
| `ict_web_monitor.prisma` | UptimeRobot Monitoring (WM) | 3 | ⬜ |
| `ict_store.prisma` | Store Operations | 5 | ⬜ |
| `ict_monitor_agent.prisma` | Agent Management (Wazuh+Beszel) | 5 | ✅ **BARU** |
| `ict_monitor_metric.prisma` | Time-Series Metrics (LibreNMS+Beszel) | 6 | ✅ **BARU** |
| `ict_monitor_log.prisma` | Log Aggregation (Wazuh+Dockge) | 6 | ✅ **BARU** |
| `ict_monitor_fim.prisma` | File Integrity Monitoring (Wazuh) | 4 | ✅ **BARU** |
| `ict_monitor_sca.prisma` | Security Config Assessment (Wazuh) | 4 | ✅ **BARU** |
| `ict_monitor_discovery.prisma` | Network Discovery (LibreNMS) | 2 | ✅ **BARU** |
| `ict_monitor_alert.prisma` | Alert Rules & Thresholds | 4 | ✅ **BARU** |
| `ict_monitor_docker_ext.prisma` | Docker Extensions (Dockge) | 5 | ✅ **BARU** |
| `ict_iso.prisma` | ISO 27001 & SOP Documentation | 14 | ✅ **BARU** |
| **Total** | | **~210** | |

### Naming Convention

- **Tabel**: `ict_{modul}` atau `dat_{entitas}`
- **Enum**: `snake_case` sesuai Prisma convention
- **Field**: `snake_case` (konsisten dengan existing)

### Pola Konsisten

- Primary key: `String @id @default(uuid())`
- Timestamps: `created_at DateTime @default(now())` + `updated_at DateTime @default(now())`
- Soft delete: gunakan `is_active` atau `status` (tidak hard delete)
- Multi-tenant: semua tabel domain menggunakan `company_id`
- Unique constraint: `@@unique([company_id, ...])` untuk isolasi data per perusahaan
- Index: `@@index` untuk kolom yang sering di-query

---

## Sistem Akses (Privilege)

### Level Akses

| Level | Keterangan |
|-------|------------|
| `hide` | Sembunyikan modul |
| `view` | Lihat saja (read-only) |
| `book` | Kelola data (CRUD) |
| `post` | Publish/Escalate |

### Role Templates (IT)

| Role | Default Level |
|------|---------------|
| `admin` | `post` semua modul |
| `manager` | `book`~`post` modul inti |
| `analyst` | `view`~`book` modul security |
| `operator` | `view`~`book` modul operasional |
| `staff` | `view` modul dasar |
| `viewer` | `view` semua modul |

### Role Templates (non-IT)

| Role | Service Access |
|------|---------------|
| `requester` | Submit request, lihat status sendiri |
| `department_head` | Submit + approve department |
| `finance` | Lihat request biaya, approve budget |
| `hrd` | Request aset untuk karyawan baru |
| `project_manager` | Request server/aplikasi project |

---

## Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| **Backend** | Go 1.26.4, Gin Gonic, `database/sql` + `lib/pq` (raw SQL, tanpa ORM) |
| **Frontend** | Next.js 16.2.10, React 19.2.4, TypeScript 5 (Strict Mode), Tailwind CSS 4, shadcn/ui (@base-ui) |
| **Database** | PostgreSQL 18, Prisma ORM 7.8.0 (hanya untuk schema/migration) |
| **Architecture** | Linear Layer Execution (Template → Repository → Usecase → Handler) |
| **State Management** | Zustand 5.x (client preferences/session) |
| **Form** | React Hook Form 7.80.0 + Zod 4.x |
| **Data Fetching** | TanStack Query (client), `serverApi` (server components) |
| **Logging** | Zerolog (backend structured logging) |
| **Encryption** | AES-GCM (mechanic/crypto.go), bcrypt (password hashing) |

---

## Catatan Teknis

### Mapping Kode ke Skeleton Folder (Backend)

| Kode Halaman | Folder Skeleton | Keterangan |
|--------------|-----------------|------------|
| SP00 | `SP00/` | Login (guest route, bukan pages) |
| SP01–SP03 | `SP01/`, `SP02/`, `SP03/` | Per halaman |
| SM01–SM05 | `SM01/`, `SM02/`, `SM03/`, `SM04/`, `SM05/` | Per halaman |
| XX99 | `XX99/` | Utility/Cleanup (admin route) |

### Frontend Route Mapping

```
/board/pages/SP01  → proxy → /rest/pages/SP01
/board/pages/SM01  → proxy → /rest/pages/SM01
```

---

## Dokumentasi Terkait

| File | Deskripsi |
|------|-----------|
| `ict_auto.md` | Rencana otomasi/agent (Docker agent, SIEM agents, certificate renewal, SR notification, PM reminder) |
| `ict_base.md` | Rencana database schema (30+ modul, ~210 tabel, enum, naming convention) |
| `ict_rest.md` | Rencana backend API (error handling, privilege, semua endpoint) |
| `ict_site.md` | Rencana frontend website (error handling, semua halaman & komponen) |

---

## Update Terakhir

- **Tanggal**: 22 Juli 2026
- **Oleh**: System
- **Keterangan**: Sinkronisasi dokumentasi dengan source code aktual:
  - Perbaiki Tech Stack: Backend menggunakan raw SQL (bukan Prisma ORM), Next.js 16.2.10
  - Tambah skeleton SM05 (Session Management) dan XX99 (Utility/Cleanup) ke struktur backend
  - Update component count: 81 UIX components, 17 lib files
  - Perbaiki Endpoint Mapping: tambah SM05 dan XX99 routes
  - Catatan: USRole sudah ada di codebase tapi belum di-wiring ke route manapun
  - SM06 (Location Management) sudah diimplementasi
