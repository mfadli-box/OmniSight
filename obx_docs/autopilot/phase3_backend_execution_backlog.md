# Backlog Backend Fase 3 - Infra Monitoring dan SIEM Fondasi

Dokumen ini memecah Fase 3 menjadi urutan kerja backend yang paling siap dieksekusi setelah Bastion stabil.

## Tujuan

1. Menyediakan inventory backend untuk host, stack, VM, site, upstream, device, dan network sample.
2. Menyediakan alerting dasar dan backup job untuk infra yang didukung.
3. Menyediakan fondasi event store dan correlation untuk SIEM bertahap.
4. Menjaga perubahan tetap bertahap dan mudah diverifikasi.

## Urutan Backend

### Paket 1 - Inventory Host, Stack, dan VM

#### Modul
- infra_host
- infra_stack
- vm_host

#### Aksi
1. Finalisasi model host dan stack/container inventory.
2. Finalisasi model VM host dan cluster inventory.
3. Tentukan company scope, status aktif, dan metadata minimum.
4. Siapkan data pilot untuk dashboard dan collector.

#### Output
- Inventory infra inti tersedia.

### Paket 2 - Inventory Nginx

#### Modul
- web_site
- web_upstream
- web_certificate
- web_config_version
- web_reload_history

#### Aksi
1. Definisikan data site, upstream, dan certificate inventory.
2. Siapkan versioning config dan histori reload.
3. Pastikan syntax test dan rollback dapat ditelusuri.
4. Kaitkan data dengan company dan actor audit.

#### Output
- Nginx inventory dan versioning siap digunakan.

### Paket 3 - Inventory MikroTik dan Polling

#### Modul
- net_device
- net_interface
- net_poll_sample

#### Aksi
1. Finalisasi model perangkat dan interface.
2. Siapkan collector polling berkala.
3. Definisikan format sample untuk health dan traffic.
4. Pastikan hasil polling terhubung ke dashboard dan alert.

#### Output
- Monitoring MikroTik dasar berjalan.

### Paket 4 - Alert, Backup, dan Resource Metric Infra

#### Modul
- infra_metric_sample
- infra_alert_rule
- net_alert_rule
- net_backup_job
- vm_resource_sample
- vm_permission
- vm_action_log

#### Aksi
1. Definisikan rule threshold resource dan network.
2. Tambahkan backup config terjadwal untuk device yang didukung.
3. Simpan histori metric sample dan alert.
4. Catat aksi perubahan penting pada VM dan resource infra.

#### Output
- Alerting dasar dan backup infra tersedia.

### Paket 5 - Event Store dan SIEM Fondasi

#### Modul
- sec_event
- sec_event_source
- sec_event_parser
- sec_rule
- sec_alert
- sec_incident

#### Aksi
1. Definisikan schema event inti lintas sumber.
2. Bangun pipeline ingest dan normalisasi.
3. Tambahkan correlation rule tahap awal.
4. Buat alur incident dan timeline.

#### Output
- Fondasi SIEM bertahap tersedia.

## Dependensi

| Area | Dependensi Minimum |
|---|---|
| Inventory infra | Fase 1 dan Fase 2 stabil |
| Nginx management | inventory host dan audit dasar tersedia |
| MikroTik monitoring | collector dan schema sample siap |
| Alerting | metric sample dan rule model tersedia |
| SIEM fondasi | event schema dan ingest pipeline siap |

## DoD Backend Fase 3

- Host, stack, dan VM inventory aktif.
- Nginx site, upstream, certificate, dan reload history aktif.
- MikroTik device, interface, dan polling sample aktif.
- Alert dasar dan backup config berjalan.
- Event store dan correlation rule dasar tersedia.

## Risiko

1. Collector runtime tidak stabil.
2. Metrics dan log tumbuh lebih cepat daripada kapasitas penyimpanan.
3. Alert noise terlalu tinggi.
4. SIEM dipaksa terlalu cepat sebelum inventory infra stabil.

## Mitigasi

1. Mulai dari pilot kecil dan sumber data terbatas.
2. Tetapkan retensi data dan batas sampling.
3. Kalibrasi rule alert bertahap.
4. Pertahankan SIEM sebagai tahap terakhir.

## Referensi

- roadmap_platform_3_phase.md
- phase3_execution_backlog.md
- backend_module_matrix_platform_replacement.md
- ai_runbook.md
