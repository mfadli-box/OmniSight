# Checklist Eksekusi - Platform Replacement OmniSight

Dokumen ini menerjemahkan gap analysis menjadi backlog implementasi yang bisa dieksekusi bertahap untuk domain:
- Docker (Dockge/Beszel)
- Nginx (nginx-ui)
- SIEM (Wazuh)
- MikroTik Monitoring (LibreNMS)
- Bastion (JumpServer: SSH, RDP, FTP, WebAppProxy)
- VM Management/Monitoring
- SOP/ISO Governance

## Cara Baca Status

- `Not Started`: belum ada implementasi aktif.
- `In Progress`: sudah ada rancangan/komponen awal.
- `Ready for Pilot`: siap uji terbatas.
- `Production Ready`: layak produksi terkontrol.

## Prioritas Eksekusi

1. SOP/ISO Governance
2. Bastion (SSH, RDP, FTP, WebAppProxy)
3. VM + Docker + Nginx control plane
4. MikroTik monitoring
5. SIEM correlation pipeline

## A. SOP/ISO Governance

### Status Saat Ini
- In Progress (fondasi data, template, dan sample sudah tersedia).

### Deliverable Minimum
1. Workflow dokumen terkendali berbasis signature type aktif.
2. Approval matrix terpasang pada dokumen produksi.
3. Evidence checklist dan revision history berjalan konsisten.

### Backlog Teknis
1. Aktifkan kode signature target di SM04:
   - SIG_POLICY_V1
   - SIG_SOP_V1
   - SIG_WI_V1
   - SIG_FORM_V1
2. Tautkan tiap kategori dokumen ke signature type yang sesuai.
3. Pastikan audit action untuk create/review/approve/reject tercatat.
4. Buat paket uji QA dokumen (happy path + reject path).

### Exit Criteria
- 4 signature type aktif dan tervalidasi.
- 4 kategori dokumen lulus uji approval end-to-end.
- Evidence audit tersedia dan dapat ditelusuri.

## B. Bastion (JumpServer Replacement)

### Status Saat Ini
- In Progress (level blueprint + roadmap + schema draft JMS).

### Deliverable Minimum
1. Web SSH terminal dengan session audit.
2. Web RDP browser access dengan policy approval.
3. FTP/file transfer browser-based dengan kontrol upload/download.
4. WebAppProxy untuk aplikasi internal melalui board.

### Backlog Teknis
1. Backend JMS runtime:
   - endpoint connect token
   - websocket bridge session
   - revoke session
2. File transfer runtime:
   - file token endpoint
   - upload/download stream
   - audit transfer event
3. WebAppProxy runtime:
   - route proxy internal app
   - policy allowlist host/path
   - access log per request
4. Security controls:
   - short-lived one-time token
   - masked secret handling
   - session timeout + forced disconnect

### Exit Criteria
- SSH, RDP, FTP, dan WebAppProxy bisa dipakai user pilot.
- Seluruh aksi session terekam audit trail.
- Approval access minimal untuk asset kritikal berjalan.

## C. Docker (Dockge/Beszel Replacement)

### Status Saat Ini
- Not Started (belum ada modul runtime aktif).

### Deliverable Minimum
1. Inventory host dan compose project.
2. Deploy action start/stop/redeploy.
3. Monitoring host/container dasar + alert threshold.

### Backlog Teknis
1. Buat domain data host, project, service, deploy history.
2. Implement API deploy command terkontrol dengan approval opsional.
3. Tambah collector metrics CPU/Mem/Disk/Restart count.
4. Buat dashboard list stack + status service + log ringkas.

### Exit Criteria
- Minimal 1 host pilot dan 3 compose project terkelola.
- Deploy history dan rollback dasar tersedia.
- Alert dasar resource aktif.

## D. Nginx (nginx-ui Replacement)

### Status Saat Ini
- Not Started.

### Deliverable Minimum
1. CRUD site dan upstream.
2. SSL inventory + binding.
3. Render config + test + reload + rollback.

### Backlog Teknis
1. Model data untuk server block, upstream, certificate.
2. Template renderer untuk nginx config.
3. Pre-flight validation (syntax test) sebelum reload.
4. Audit log perubahan config dan actor.

### Exit Criteria
- 5 site pilot bisa dikelola dari OmniSight.
- Deploy config gagal dapat rollback otomatis.
- Sertifikat dan expiry terpantau.

## E. MikroTik Monitoring (LibreNMS Replacement)

### Status Saat Ini
- Not Started (baru domain konseptual).

### Deliverable Minimum
1. Inventory perangkat MikroTik.
2. Polling berkala resource + interface status.
3. Alert threshold dasar.

### Backlog Teknis
1. Collector SNMP/API untuk router pilot.
2. Penyimpanan timeseries metrik inti.
3. Rule alert untuk CPU, memory, link down.
4. Backup config terjadwal + status terakhir backup.

### Exit Criteria
- Minimal 10 perangkat pilot termonitor.
- Alert dasar berjalan dan tervalidasi.
- Backup config otomatis berhasil.

## F. VM Management/Monitoring

### Status Saat Ini
- Not Started.

### Deliverable Minimum
1. VM inventory (host, owner, environment, criticality).
2. Monitoring performa dasar dan status hidup/mati.
3. Integrasi approval untuk tindakan berisiko.

### Backlog Teknis
1. Sinkronisasi metadata VM dari sumber infrastruktur.
2. Daftar tindakan operasional yang diizinkan (read-only dulu).
3. Dashboard health ringkas per cluster/host.
4. Audit perubahan status dan assignment owner.

### Exit Criteria
- 100% VM kritikal terinventaris.
- Health dashboard akurat untuk pilot environment.
- Jejak audit perubahan owner/status tersedia.

## G. SIEM (Wazuh Replacement)

### Status Saat Ini
- Not Started (kompleksitas tertinggi).

### Deliverable Minimum
1. Event ingestion terpusat.
2. Normalisasi event lintas sumber.
3. Rule correlation dasar dan severity mapping.
4. Alert queue + incident timeline.

### Backlog Teknis
1. Definisikan event schema inti (auth, policy, infra, network).
2. Bangun pipeline ingest + parser standar.
3. Terapkan rules tahap awal (failed login burst, privilege abuse, service anomaly).
4. Buat dashboard insiden + drill-down log.

### Exit Criteria
- Sumber event utama berhasil di-ingest.
- Correlation rule minimum berjalan stabil.
- Incident timeline dapat dipakai audit internal.

## Rencana 90 Hari (Target Realistis)

### Hari 1-30
1. Finalisasi SOP/ISO workflow production pilot.
2. Rollout SM04 signature type target.
3. Mulai implementasi backend JMS connect token + session event.

### Hari 31-60
1. Uji pilot Web SSH + session audit.
2. Implement FTP browser transfer dan audit.
3. Implement WebAppProxy dasar untuk aplikasi internal prioritas.

### Hari 61-90
1. Tambah pilot Web RDP.
2. Mulai fondasi inventory Docker/VM/Nginx.
3. Siapkan design event schema untuk SIEM fase berikutnya.

## Risiko Utama Lintas Domain

1. Scope creep lintas domain terlalu cepat.
2. Kekurangan resource engineer untuk domain observability.
3. Kualitas audit trail tertinggal dari fitur runtime.
4. Kebutuhan keamanan token/session lebih tinggi dari estimasi awal.

## Mitigasi

1. Gunakan pendekatan phase-gate, tidak paralel penuh semua domain.
2. Wajibkan Definition of Done dan evidence checklist di setiap milestone.
3. Jadikan security review sebagai syarat masuk pilot untuk bastion.
4. Simpan roadmap SIEM sebagai fase akhir setelah fondasi operasi stabil.

## Referensi

- gap_analysis_platform_replacement.md
- capability_matrix_platform_replacement.md
- roadmap_platform_3_phase.md
- checklist_sm04_signature_rollout.md
