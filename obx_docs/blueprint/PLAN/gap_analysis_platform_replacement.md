# Gap Analysis - Pengganti Platform Eksternal di OmniSight

Dokumen ini memetakan gap antara kebutuhan platform target dan kondisi OmniSight saat ini berdasarkan source code aktif, blueprint, dan roadmap rancangan.

## Ringkasan Status

| Domain | Produk Acuan | Target OmniSight | Status Saat Ini | Kesiapan Pengganti |
|---|---|---|---|---|
| Docker orchestration | Dockge | Kelola compose, container, deploy, history | Belum aktif di source code saat ini | Rendah |
| Host/container monitoring | Beszel | Monitoring host, container, resource, alert | Belum aktif di source code saat ini | Rendah |
| Nginx management | nginx-ui | Site, upstream, SSL, config versioning | Belum aktif di source code saat ini | Rendah |
| SIEM | Wazuh | Event ingest, correlation, alert, case audit | Baru level domain dan konsep | Rendah |
| Network monitoring | LibreNMS | Inventori MikroTik, polling, alert, backup | Baru level domain dan konsep | Rendah |
| Bastion | JumpServer | SSH, RDP, FTP, WebAppProxy, audit session | Sudah ada blueprint PLAN jsm_stack | Menengah untuk rancangan, rendah untuk implementasi |
| VM management | Proxmox/monitoring stack | Inventory, monitoring, group, permission, action | Belum aktif di source code saat ini | Rendah |
| SOP/ISO docs | DMS / QMS | Controlled docs, approval, revision history, audit | Paling dekat secara arsitektur | Menengah |

## 1. Docker - Dockge + Beszel

### Kapabilitas produk acuan
- Manajemen compose file
- Start/stop/redeploy stack
- Monitoring container dan host
- Log, metrics, deploy history
- Multi-host visibility

### Kebutuhan minimum di OmniSight
- Inventory host dan compose project
- CRUD compose project
- Deploy action + riwayat deploy
- Metrics host/container berkala
- Alert untuk CPU, memory, disk, restart loop

### Status OmniSight saat ini
- Dokumentasi domain ada di blueprint BASE dan histori project.
- Source code aktif workspace belum menampilkan module DK atau VM terdaftar pada backend/frontend saat ini.

### Gap utama
- Belum ada route aktif dan page operasional untuk compose/container.
- Belum ada agent runtime terbukti untuk metrik real-time host/container.
- Belum ada UI deploy workflow setara Dockge.

### Kesimpulan
- Belum bisa menggantikan Dockge/Beszel saat ini.
- Layak dijadikan fase setelah fondasi asset, auth, dan inventory matang.

## 2. Nginx - nginx-ui

### Kapabilitas produk acuan
- CRUD virtual host
- CRUD upstream
- Manage SSL/certificate
- Edit config dan reload service
- Versioning konfigurasi

### Kebutuhan minimum di OmniSight
- Site inventory
- Upstream inventory
- SSL inventory dan binding
- Render config dan push ke target host
- Riwayat perubahan dan rollback dasar

### Status OmniSight saat ini
- Domain `ict_website` sudah terdokumentasi.
- Source code aktif saat ini belum menampilkan modul NX aktif pada backend/frontend.

### Gap utama
- Tidak ada UI runtime aktif untuk site/upstream/certificate.
- Tidak ada flow deploy config dan verifikasi reload Nginx yang aktif.
- Tidak ada bukti rollback config end-to-end.

### Kesimpulan
- Belum bisa menggantikan nginx-ui saat ini.

## 3. SIEM - Wazuh

### Kapabilitas produk acuan
- Agent log/security event collection
- Rule engine dan correlation
- Alerting
- Incident trail dan case workflow
- Dashboard dan forensic data

### Kebutuhan minimum di OmniSight
- Central event store
- Parsing dan normalisasi log
- Rule engine severity
- Alert workflow
- Incident/timeline view

### Status OmniSight saat ini
- Domain security ada di schema dan blueprint.
- Agent/log domain ada secara konsep pada AUTO.
- Belum ada bukti implementasi SIEM aktif pada source code workspace saat ini.

### Gap utama
- Belum ada rule correlation engine.
- Belum ada active alert workflow yang terbukti berjalan.
- Belum ada UI analitik SIEM penuh.

### Kesimpulan
- Belum bisa menggantikan Wazuh.
- Ini domain yang paling berat dan sebaiknya menjadi fase akhir.

## 4. MikroTik / Network Monitoring - LibreNMS

### Kapabilitas produk acuan
- Discovery perangkat
- Polling SNMP/API
- Resource graph
- Alerting
- Inventory interface dan config backup

### Kebutuhan minimum di OmniSight
- Inventori perangkat MikroTik
- Polling berkala CPU, memory, traffic, uptime
- Daftar interface dan perubahan status
- Backup config terjadwal
- Alert threshold dasar

### Status OmniSight saat ini
- Domain `ict_mikrotik` sudah ada di schema/blueprint.
- Agent/collector belum terbukti aktif sebagai pengganti LibreNMS penuh.

### Gap utama
- Discovery dan polling massal belum terbukti.
- Visualisasi historis belum terbukti.
- Alert threshold dan eskalasi belum ada bukti implementasi aktif.

### Kesimpulan
- Belum bisa menggantikan LibreNMS saat ini.

## 5. Bastion - JumpServer

### Kapabilitas produk acuan
- Asset inventory
- Account vault
- Approval / access policy
- Session audit
- SSH web terminal
- RDP browser access
- FTP / file transfer
- WebAppProxy

### Status OmniSight saat ini
- Sudah ada blueprint dan roadmap JMS pada PLAN.
- SSH, RDP, WebAppProxy sudah masuk blueprint/roadmap.
- FTP dan file transfer browser-based perlu ditambahkan eksplisit ke rancangan.

### Gap utama
- Belum ada implementasi route, gateway, dan UI nyata.
- Belum ada session recording runtime.
- Belum ada credential rotation dan approval policy yang aktif.

### Kesimpulan
- Secara arsitektur paling dekat untuk dirancang.
- Secara implementasi masih belum siap menggantikan JumpServer saat ini.

## 6. VM Management / Monitoring

### Kapabilitas target
- Inventory VM host
- Resource monitoring
- Host group dan permission
- Action operasional dasar
- Riwayat perubahan / deploy / patch

### Status OmniSight saat ini
- Domain VM dan machine sudah ada di dokumentasi BASE.
- Tidak ada modul aktif di backend/frontend saat ini pada source code workspace.

### Gap utama
- UI dan route runtime belum aktif.
- Monitoring pipeline host belum terbukti.
- Integrasi approval/action belum aktif.

### Kesimpulan
- Belum bisa menggantikan platform VM management/monitoring saat ini.

## 7. SOP / ISO Documentation

### Kapabilitas target
- Controlled documents
- Approval workflow
- Versioning
- Evidence dan audit trail
- Mapping policy/procedure/work instruction

### Status OmniSight saat ini
- Domain document dan signature paling cocok untuk use case ini.
- Struktur obx_docs, dat_document, dat_signature, dan dat_request mendukung arah ini.

### Gap utama
- Belum ada template ISO/SOP formal yang lengkap per proses.
- Belum ada matriks dokumen terkendali, owner, masa berlaku, dan evidence checklist yang eksplisit.

### Kesimpulan
- Ini area paling realistis untuk dijadikan kemenangan awal.
- OmniSight paling cepat memberi nilai pada dokumentasi SOP/ISO dibanding replacement tool infra berat.

## Prioritas Kelayakan

| Prioritas | Domain | Alasan |
|---|---|---|
| 1 | SOP / ISO Documentation | Fondasi data dan alur approval paling dekat |
| 2 | Bastion (JMS) | Sudah ada blueprint dan roadmap aktif |
| 3 | VM / Docker / Nginx Inventory | Masuk akal setelah asset dan bastion dasar matang |
| 4 | MikroTik monitoring | Butuh collector dan visualisasi historis |
| 5 | SIEM | Paling kompleks dan mahal secara operasional |

## Rekomendasi Keputusan

- Jangan memosisikan OmniSight sebagai pengganti total semua tool dalam jangka pendek.
- Gunakan strategi bertahap:
  1. menangkan domain SOP/ISO dan controlled documentation,
  2. bangun bastion dasar (SSH, RDP, FTP, WebAppProxy),
  3. lanjutkan inventory dan monitoring infra,
  4. baru masuk ke SIEM penuh.

## Referensi Detail

- [Matriks detail kelayakan pengganti platform](detailed_capability_matrix_platform_replacement.md)
- [Matriks aksi implementasi per fitur](feature_action_matrix_platform_replacement.md)
- [Matriks modul backend per domain](backend_module_matrix_platform_replacement.md)
