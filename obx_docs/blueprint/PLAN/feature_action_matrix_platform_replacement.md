# Matriks Aksi Implementasi - Platform Replacement OmniSight

Dokumen ini memecah rekomendasi menjadi aksi implementasi per fitur agar tim bisa mengeksekusi secara bertahap tanpa mencampur semua domain sekaligus.

## Legenda Status

- `Belum`: belum ada implementasi aktif.
- `Rancang`: sudah ada blueprint, roadmap, atau skema data.
- `Pilot`: cukup untuk uji terbatas.
- `Produksi`: layak dipakai operasional penuh secara terkontrol.

## Ringkasan Eksekusi

| Domain | Status | Fokus Aksi Berikutnya |
|---|---|---|
| Docker (Dockge/Beszel) | Belum | inventory host, deploy history, metrics collector |
| Nginx (nginx-ui) | Belum | site/upstream model, config renderer, reload flow |
| SIEM (Wazuh) | Belum | event schema, ingest pipeline, correlation rule |
| MikroTik (LibreNMS) | Belum | device inventory, polling collector, alert threshold |
| Bastion (JumpServer) | Rancang | SSH, RDP, FTP, WebAppProxy, session audit |
| VM Monitoring | Belum | VM inventory, health collector, action log |
| SOP / ISO | Rancang menuju Pilot | approval workflow, versioning, evidence trail |

## 1. Docker - Dockge / Beszel

| Fitur | Prioritas | Status Saat Ini | Aksi Minimal | Output yang Diharapkan |
|---|---|---|---|---|
| Host inventory | Tinggi | Belum | buat model host dan koneksi | daftar host terkelola |
| Compose project inventory | Tinggi | Belum | buat model project/stack | daftar stack per host |
| Deploy history | Tinggi | Belum | simpan aksi deploy/rollback | riwayat deploy bisa ditelusuri |
| Resource metrics | Tinggi | Belum | collector CPU, memory, disk | metrik host/container |
| Alerting | Sedang | Belum | threshold dasar | notifikasi resource kritis |
| Multi-host dashboard | Sedang | Belum | agregasi view host | dashboard operasi |

### Urutan Eksekusi

1. Inventory host dan project.
2. Deploy history dan action log.
3. Collector metrics.
4. Alert dasar.

## 2. Nginx - nginx-ui

| Fitur | Prioritas | Status Saat Ini | Aksi Minimal | Output yang Diharapkan |
|---|---|---|---|---|
| Site inventory | Tinggi | Belum | model virtual host | daftar site |
| Upstream inventory | Tinggi | Belum | model upstream | daftar upstream |
| SSL inventory | Tinggi | Belum | model certificate dan binding | status certificate |
| Config renderer | Tinggi | Belum | template config nginx | file config ter-render |
| Syntax test | Tinggi | Belum | validasi sebelum reload | config valid |
| Reload flow | Tinggi | Belum | aksi reload terkontrol | reload tercatat |
| Rollback config | Sedang | Belum | simpan versi config | rollback dasar |

### Urutan Eksekusi

1. Inventory site, upstream, SSL.
2. Renderer config.
3. Syntax test dan reload.
4. Rollback versi.

## 3. SIEM - Wazuh

| Fitur | Prioritas | Status Saat Ini | Aksi Minimal | Output yang Diharapkan |
|---|---|---|---|---|
| Event schema | Tinggi | Belum | definisikan format event inti | event seragam |
| Ingest pipeline | Tinggi | Belum | endpoint/collector event | event masuk terpusat |
| Normalisasi | Tinggi | Belum | parser per sumber | data siap correlation |
| Correlation rule | Tinggi | Belum | rule engine sederhana | alert berbasis rule |
| Severity mapping | Sedang | Belum | level severity | prioritas insiden |
| Incident timeline | Sedang | Belum | view timeline | alur insiden terlihat |
| Forensic search | Sedang | Belum | pencarian event | investigasi dasar |

### Urutan Eksekusi

1. Event schema dan ingest.
2. Normalisasi.
3. Correlation rule sederhana.
4. Timeline dan forensic search.

## 4. MikroTik - LibreNMS

| Fitur | Prioritas | Status Saat Ini | Aksi Minimal | Output yang Diharapkan |
|---|---|---|---|---|
| Device inventory | Tinggi | Belum | model perangkat dan group | daftar device |
| Polling collector | Tinggi | Belum | SNMP/API polling berkala | data health device |
| Interface monitoring | Tinggi | Belum | status port/interface | perubahan link terlihat |
| Traffic graph | Sedang | Belum | penyimpanan metrik historis | grafik trafik |
| Alert threshold | Sedang | Belum | rule resource/interface | notifikasi kondisi buruk |
| Config backup | Sedang | Belum | backup konfigurasi | arsip backup perangkat |

### Urutan Eksekusi

1. Inventory device.
2. Collector polling.
3. Interface monitoring.
4. Alert dan backup config.

## 5. Bastion - JumpServer

| Fitur | Prioritas | Status Saat Ini | Aksi Minimal | Output yang Diharapkan |
|---|---|---|---|---|
| Asset inventory | Tinggi | Rancang | finalisasi model asset | daftar asset bastion |
| Account vault | Tinggi | Rancang | model akun dan secret ref | akun terkelola |
| Approval policy | Tinggi | Rancang | kaitkan request dengan policy | akses terotorisasi |
| SSH web terminal | Tinggi | Rancang | connect token + PTY bridge | SSH via browser |
| RDP browser access | Tinggi | Rancang | browser client + gateway | RDP via browser |
| FTP / file transfer | Tinggi | Rancang | token file + stream transfer | upload/download terkontrol |
| WebAppProxy | Tinggi | Rancang | proxy path-based | akses aplikasi internal |
| Session audit | Tinggi | Rancang | event connect/disconnect | audit trail sesi |
| Session recording | Sedang | Belum | simpan artefak sesi | rekaman sesi |
| Credential rotation | Sedang | Belum | sinkron secret / rotate | kredensial lebih aman |

### Urutan Eksekusi

1. Asset dan account inventory.
2. Connect token dan session audit.
3. SSH web terminal.
4. FTP dan WebAppProxy.
5. RDP dan recording.

## 6. VM Monitoring / Management

| Fitur | Prioritas | Status Saat Ini | Aksi Minimal | Output yang Diharapkan |
|---|---|---|---|---|
| VM inventory | Tinggi | Belum | model host, cluster, owner | daftar VM |
| Health monitoring | Tinggi | Belum | collector status hidup/mati | status VM aktif |
| Resource metrics | Sedang | Belum | CPU, memory, disk | grafik kesehatan |
| Permission model | Sedang | Belum | grup dan akses | akses terkontrol |
| Action log | Sedang | Belum | catat perubahan penting | histori tindakan |

### Urutan Eksekusi

1. Inventory VM.
2. Health monitoring.
3. Resource metrics.
4. Permission dan action log.

## 7. SOP / ISO Documentation

| Fitur | Prioritas | Status Saat Ini | Aksi Minimal | Output yang Diharapkan |
|---|---|---|---|---|
| Controlled document | Tinggi | Rancang menuju Pilot | versioning dan status | dokumen terkendali |
| Approval workflow | Tinggi | Rancang menuju Pilot | reviewer/approver/signature | alur approval hidup |
| Evidence checklist | Tinggi | Rancang menuju Pilot | checklist bukti | audit lebih rapi |
| Revision history | Tinggi | Rancang | simpan perubahan | histori revisi |
| Owner mapping | Sedang | Rancang | peran owner/reviewer/approver | akuntabilitas jelas |
| Archive policy | Sedang | Rancang | masa simpan dan arsip | dokumen tertata |

### Urutan Eksekusi

1. Finalisasi approval workflow.
2. Aktifkan signature type.
3. Terapkan evidence checklist.
4. Tambahkan revision history dan archive policy.

## Rekomendasi Phase-Gate

### Gate 1 - Quick Win
- SOP / ISO pilot.

### Gate 2 - Access Control
- Bastion SSH dan session audit.

### Gate 3 - Infra Inventory
- Docker, Nginx, dan VM inventory.

### Gate 4 - Network Monitoring
- MikroTik polling dan alert.

### Gate 5 - Security Platform
- SIEM pipeline dan correlation.

## Kesimpulan

OmniSight paling masuk akal dimulai dari SOP/ISO dan Bastion, lalu bergerak ke inventory/monitoring infra. Dockge, Beszel, nginx-ui, LibreNMS, dan Wazuh belum bisa diganti penuh tanpa membangun runtime, collector, dan workflow operasional yang cukup besar.
