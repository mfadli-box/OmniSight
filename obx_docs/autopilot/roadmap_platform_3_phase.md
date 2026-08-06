# Roadmap Prioritas 3 Fase - OmniSight sebagai Platform Pengganti

Dokumen ini menyusun urutan prioritas realistis agar OmniSight berkembang menjadi platform terpadu pengganti tool terpisah.

## Tujuan Besar

Mengarahkan OmniSight menjadi platform internal untuk:
- dokumentasi SOP dan ISO,
- bastion access,
- inventory dan monitoring infra,
- automation agent,
- dan secara bertahap event/security management.

## Pemetaan Implementasi

Roadmap ini diturunkan dari tiga dokumen kerja yang harus dipakai bersama:

- [Matriks detail kelayakan pengganti platform](../blueprint/PLAN/detailed_capability_matrix_platform_replacement.md)
- [Matriks aksi implementasi per fitur](../blueprint/PLAN/feature_action_matrix_platform_replacement.md)
- [Matriks modul backend per domain](../blueprint/PLAN/backend_module_matrix_platform_replacement.md)

Urutan eksekusi operasional harus mengikuti prioritas berikut:

1. SOP/ISO governance sebagai quick win.
2. Bastion core untuk SSH, RDP, FTP, dan WebAppProxy.
3. Inventory infra untuk Docker, Nginx, dan VM.
4. Monitoring MikroTik.
5. SIEM terakhir.

## Fase 1 - Governance, SOP/ISO, dan Fondasi Asset

### Fokus
- Controlled document
- Approval workflow
- Master data company, module, user, privilege
- Asset inventory dasar

### Deliverable
- Dokumen SOP/ISO dengan versioning dan approval
- Request workflow memakai `dat_request`
- Audit trail perubahan dokumen
- Asset inventory dasar untuk host, app, dan kredensial referensi
- Signature type aktif untuk approval workflow
- Evidence checklist dan revision history untuk dokumen terkendali

### Outcome bisnis
- Organisasi punya single source of truth untuk prosedur dan bukti audit.
- Fondasi identitas, akses, dan company scope stabil.

### Risiko utama
- User menganggap fitur infra kompleks harus jadi dulu.
- Scope dokumentasi terlalu longgar tanpa template baku.

## Fase 2 - Bastion Access dan Operasi Web-Based

### Fokus
- JMS / bastion
- Web SSH
- Web RDP
- FTP / file transfer berbasis browser
- WebAppProxy
- Session audit dan approval access

### Deliverable
- Asset + account + session audit
- Connect token, gateway bridge, session event
- File transfer browser-based untuk upload/download terkontrol
- Proxy aplikasi internal lewat board
- Audit connect/disconnect, revoke, dan file transfer yang konsisten

### Outcome bisnis
- Operator tidak lagi bergantung ke banyak tool akses terpisah.
- Aktivitas remote access lebih mudah diaudit.

### Risiko utama
- Kompleksitas gateway stream dan keamanan token
- Harapan menggantikan JumpServer penuh terlalu cepat

## Fase 3 - Monitoring Infra, Network, Docker, dan Security Event

### Fokus
- Monitoring VM/host/container
- Manajemen Nginx
- Monitoring MikroTik
- Alerting terpusat
- Security event pipeline bertahap

### Deliverable
- Inventory dan metrik host/container
- Nginx site/upstream/certificate management
- Polling device MikroTik dan backup
- Alert threshold dasar
- Event store dan korelasi ringan sebelum SIEM penuh
- Model backend inventory, collector, dan history untuk tiap domain infra

### Outcome bisnis
- OmniSight berkembang dari sistem dokumentasi + akses menjadi control plane operasional.

### Risiko utama
- Kompleksitas observability dan agent runtime
- Biaya storage metrics/log meningkat
- Scope melebar sebelum fase 2 stabil

## Exit Criteria per Fase

### Fase 1
- Dokumen SOP/ISO berjalan dengan approval dan histori.
- Asset inventory dasar tersedia.
- User/company/privilege stabil.

### Fase 2
- Web SSH, Web RDP, FTP browser-based, dan WebAppProxy aktif.
- Session audit dan revoke berjalan.
- Approval access minimum tersedia.

### Fase 3
- Monitoring host/container/device stabil.
- Alerting dasar berjalan.
- Nginx/network/security event punya visibilitas operasional.

## Prinsip Eksekusi

- Jangan kejar semua domain sekaligus.
- Setiap fase harus memberi nilai mandiri.
- Domain paling berat seperti SIEM dikerjakan terakhir.
- Dokumentasi dan audit trail tidak boleh tertinggal dari implementasi.
- Gunakan matriks backend sebagai daftar urutan modul saat mulai coding.
