# Matriks Detail Kelayakan Pengganti Platform OmniSight

Dokumen ini memecah penilaian pada level fitur agar terlihat jelas bagian mana yang sudah bisa dikejar, mana yang masih rancangan, dan mana yang masih belum siap.

## Legenda Status

- `Belum`: belum ada bukti implementasi aktif atau rancangan detail yang cukup.
- `Rancang`: sudah ada blueprint/roadmap/skeleton data, tetapi belum runtime.
- `Siap Pilot`: cukup untuk uji terbatas pada scope kecil.
- `Produksi`: layak dipakai operasional penuh secara terkontrol.

## Ringkasan Cepat

| Domain | Status Umum | Kesimpulan |
|---|---|---|
| Docker (Dockge, Beszel) | Belum | Belum bisa menggantikan |
| Nginx (nginx-ui) | Belum | Belum bisa menggantikan |
| SIEM (Wazuh) | Belum | Belum bisa menggantikan |
| MikroTik / Network Monitoring (LibreNMS) | Belum | Belum bisa menggantikan |
| Bastion (JumpServer) | Rancang | Paling dekat, tetapi belum runtime |
| VM Monitoring / Management | Belum | Belum bisa menggantikan |
| SOP / ISO Documentation | Rancang menuju Siap Pilot | Paling realistis untuk dimulai |

## 1. Docker - Dockge / Beszel

| Fitur | Dockge / Beszel | OmniSight Saat Ini | Status |
|---|---|---|---|
| Inventory host | Ya | Belum aktif | Belum |
| Inventory compose project | Ya | Belum aktif | Belum |
| Deploy start/stop/redeploy | Ya | Belum aktif | Belum |
| Log deploy history | Ya | Belum aktif | Belum |
| Metrics CPU/memory/disk | Ya | Belum aktif | Belum |
| Alert resource | Ya | Belum aktif | Belum |
| Multi-host dashboard | Ya | Belum aktif | Belum |

### Kesimpulan Docker

Belum layak menggantikan Dockge atau Beszel. Hambatan utamanya adalah belum ada domain runtime dan collector metrik aktif.

## 2. Nginx - nginx-ui

| Fitur | nginx-ui | OmniSight Saat Ini | Status |
|---|---|---|---|
| CRUD virtual host | Ya | Belum aktif | Belum |
| CRUD upstream | Ya | Belum aktif | Belum |
| SSL / certificate binding | Ya | Belum aktif | Belum |
| Render config | Ya | Belum aktif | Belum |
| Syntax test sebelum reload | Ya | Belum aktif | Belum |
| Reload service terkontrol | Ya | Belum aktif | Belum |
| Rollback konfigurasi | Ya | Belum aktif | Belum |

### Kesimpulan Nginx

Belum layak menggantikan nginx-ui.

## 3. SIEM - Wazuh

| Fitur | Wazuh | OmniSight Saat Ini | Status |
|---|---|---|---|
| Event ingest | Ya | Baru konsep domain | Belum |
| Parsing dan normalisasi | Ya | Belum aktif | Belum |
| Rule correlation | Ya | Belum aktif | Belum |
| Severity mapping | Ya | Belum aktif | Belum |
| Incident timeline | Ya | Belum aktif | Belum |
| Alert workflow | Ya | Belum aktif | Belum |
| Forensic search | Ya | Belum aktif | Belum |

### Kesimpulan SIEM

Belum layak menggantikan Wazuh. Ini adalah domain paling kompleks dan paling akhir untuk dikerjakan.

## 4. MikroTik / Network Monitoring - LibreNMS

| Fitur | LibreNMS | OmniSight Saat Ini | Status |
|---|---|---|---|
| Device discovery | Ya | Belum aktif | Belum |
| SNMP/API polling | Ya | Belum aktif | Belum |
| Interface status monitoring | Ya | Belum aktif | Belum |
| CPU/memory/traffic graph | Ya | Belum aktif | Belum |
| Alert threshold | Ya | Belum aktif | Belum |
| Config backup | Ya | Belum aktif | Belum |

### Kesimpulan MikroTik

Belum layak menggantikan LibreNMS.

## 5. Bastion - JumpServer

| Fitur | JumpServer | OmniSight Saat Ini | Status |
|---|---|---|---|
| Asset inventory | Ya | Blueprint PLAN ada | Rancang |
| Account vault | Ya | Blueprint PLAN ada | Rancang |
| Access approval policy | Ya | Blueprint + checklist ada | Rancang |
| SSH web terminal | Ya | Blueprint ada | Rancang |
| RDP browser access | Ya | Blueprint ada | Rancang |
| FTP / file transfer | Ya | Sudah masuk rancangan eksplisit | Rancang |
| WebAppProxy | Ya | Sudah masuk rancangan eksplisit | Rancang |
| Session audit | Ya | Blueprint ada | Rancang |
| Session recording | Ya | Belum runtime | Belum |
| Credential rotation | Ya | Belum runtime | Belum |

### Kesimpulan Bastion

Paling dekat untuk dijadikan pengganti bertahap, tetapi belum bisa menggantikan JumpServer penuh saat ini.

## 6. VM Management / Monitoring

| Fitur | Target Umum | OmniSight Saat Ini | Status |
|---|---|---|---|
| VM inventory | Ya | Belum aktif | Belum |
| Host/cluster resource monitoring | Ya | Belum aktif | Belum |
| Group / permission | Ya | Belum aktif | Belum |
| Action operasional | Ya | Belum aktif | Belum |
| Change history | Ya | Belum aktif | Belum |

### Kesimpulan VM

Belum layak menggantikan platform VM management/monitoring.

## 7. SOP / ISO Documentation

| Fitur | Target DMS/QMS | OmniSight Saat Ini | Status |
|---|---|---|---|
| Controlled document | Ya | Template dan sample ada | Rancang menuju Siap Pilot |
| Approval workflow | Ya | Checklist dan matrix ada | Rancang menuju Siap Pilot |
| Versioning | Ya | Direncanakan di governance docs | Rancang |
| Evidence checklist | Ya | Ada template | Rancang menuju Siap Pilot |
| Audit trail | Ya | Selaras dengan domain request/signature | Rancang menuju Siap Pilot |
| Owner/reviewer/approver | Ya | Sudah ada sample dokumen | Rancang |

### Kesimpulan SOP / ISO

Ini domain yang paling realistis untuk dibawa ke pilot terlebih dahulu, karena fondasi dokumen, matrix approval, dan checklist sudah tersedia.

## Penilaian Akhir

| Urutan | Domain | Verdict |
|---|---|---|
| 1 | SOP / ISO Documentation | Siap diarahkan ke pilot |
| 2 | Bastion | Rancang, paling dekat setelah SOP/ISO |
| 3 | VM / Docker / Nginx | Masih perlu fondasi runtime dan collector |
| 4 | MikroTik monitoring | Masih butuh collector dan alerting |
| 5 | SIEM | Paling jauh dari siap |

## Rekomendasi Praktis

1. Jadikan SOP/ISO sebagai quick win untuk membangun governance dan audit trail.
2. Lanjutkan Bastion secara bertahap, mulai dari SSH dan session audit.
3. Baru setelah itu bangun inventory dan monitoring infra.
4. Simpan SIEM sebagai fase terakhir.
