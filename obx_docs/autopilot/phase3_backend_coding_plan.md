# Rencana Coding Backend Fase 3 - Infra Monitoring dan SIEM Fondasi

Dokumen ini menurunkan backlog backend Fase 3 menjadi paket coding yang lebih konkret untuk dieksekusi di repo.

## Tujuan

1. Mengaktifkan inventory infra host, stack, VM, network, dan web management.
2. Menyediakan alerting dasar, backup job, dan metric sample.
3. Menyediakan fondasi event store serta correlation rule untuk SIEM bertahap.
4. Menjaga setiap paket kecil, terverifikasi, dan mudah di-rollback.

## Prinsip Coding

- Satu paket fokus pada satu klaster domain infra.
- Schema, backend API, frontend board, dan dokumentasi terkait diselesaikan dalam batch yang sama bila memungkinkan.
- Jangan lanjut ke paket berikutnya sebelum verifikasi paket berjalan.

## Paket 1 - Inventory Host, Stack, dan VM

### Target Modul
- infra_host
- infra_stack
- vm_host

### File yang Ditargetkan
- obx_base/prisma/schema/ict_machine.prisma
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_rest/skeleton/{KODE_INFRA_CORE}/template.go
- obx_rest/skeleton/{KODE_INFRA_CORE}/repository.go
- obx_rest/skeleton/{KODE_INFRA_CORE}/usecase.go
- obx_rest/skeleton/{KODE_INFRA_CORE}/handler.go
- obx_rest/backbone/routes.go
- obx_site/src/app/board/pages/{KODE_INFRA_CORE}/page.tsx
- obx_site/src/app/board/model/module.ts
- obx_docs/blueprint/BASE/ict_machine.md
- obx_docs/guide/BASE/ict_machine.md

### Verifikasi
- CRUD inventory host dan VM berjalan.
- Company scope konsisten.
- Build backend lulus.

## Paket 2 - Inventory Nginx

### Target Modul
- web_site
- web_upstream
- web_certificate
- web_config_version
- web_reload_history

### File yang Ditargetkan
- obx_base/prisma/schema/ict_website.prisma
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_rest/skeleton/{KODE_WEB_MGMT}/template.go
- obx_rest/skeleton/{KODE_WEB_MGMT}/repository.go
- obx_rest/skeleton/{KODE_WEB_MGMT}/usecase.go
- obx_rest/skeleton/{KODE_WEB_MGMT}/handler.go
- obx_rest/backbone/routes.go
- obx_site/src/app/board/pages/{KODE_WEB_MGMT}/page.tsx
- obx_site/src/app/board/model/module.ts
- obx_docs/blueprint/BASE/ict_website.md
- obx_docs/guide/BASE/ict_website.md

### Verifikasi
- Site, upstream, certificate, dan reload history tercatat.
- Build backend lulus.

## Paket 3 - Inventory MikroTik dan Polling

### Target Modul
- net_device
- net_interface
- net_poll_sample

### File yang Ditargetkan
- obx_base/prisma/schema/ict_mikrotik.prisma
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_rest/skeleton/{KODE_NET_MON}/template.go
- obx_rest/skeleton/{KODE_NET_MON}/repository.go
- obx_rest/skeleton/{KODE_NET_MON}/usecase.go
- obx_rest/skeleton/{KODE_NET_MON}/handler.go
- obx_rest/backbone/routes.go
- obx_site/src/app/board/pages/{KODE_NET_MON}/page.tsx
- obx_site/src/app/board/model/module.ts
- obx_docs/blueprint/BASE/ict_mikrotik.md
- obx_docs/guide/BASE/ict_mikrotik.md

### Verifikasi
- Device, interface, dan poll sample tersimpan.
- Build backend lulus.

## Paket 4 - Alert, Backup, dan Resource Metric

### Target Modul
- infra_metric_sample
- infra_alert_rule
- net_alert_rule
- net_backup_job
- vm_resource_sample
- vm_permission
- vm_action_log

### File yang Ditargetkan
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_base/prisma/schema/ict_machine.prisma
- obx_base/prisma/schema/ict_mikrotik.prisma
- obx_rest/skeleton/{KODE_INFRA_ALERT}/template.go
- obx_rest/skeleton/{KODE_INFRA_ALERT}/repository.go
- obx_rest/skeleton/{KODE_INFRA_ALERT}/usecase.go
- obx_rest/skeleton/{KODE_INFRA_ALERT}/handler.go
- obx_rest/backbone/routes.go
- obx_site/src/app/board/pages/{KODE_INFRA_ALERT}/page.tsx
- obx_site/src/app/board/model/module.ts

### Verifikasi
- Rule alert dapat dibuat dan dipicu.
- Backup job tercatat.
- Build backend lulus.

## Paket 5 - Event Store dan SIEM Fondasi

### Target Modul
- sec_event
- sec_event_source
- sec_event_parser
- sec_rule
- sec_alert
- sec_incident

### File yang Ditargetkan
- obx_base/prisma/schema/ict_security.prisma
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_rest/skeleton/{KODE_SIEM_CORE}/template.go
- obx_rest/skeleton/{KODE_SIEM_CORE}/repository.go
- obx_rest/skeleton/{KODE_SIEM_CORE}/usecase.go
- obx_rest/skeleton/{KODE_SIEM_CORE}/handler.go
- obx_rest/backbone/routes.go
- obx_site/src/app/board/pages/{KODE_SIEM_CORE}/page.tsx
- obx_site/src/app/board/model/module.ts
- obx_docs/blueprint/BASE/ict_security.md
- obx_docs/guide/BASE/ict_security.md

### Verifikasi
- Event dari satu sumber bisa ingest.
- Rule dasar menghasilkan alert.
- Incident timeline terbaca.

## Urutan Coding yang Disarankan

1. Paket 1 - Inventory Host, Stack, dan VM
2. Paket 2 - Inventory Nginx
3. Paket 3 - Inventory MikroTik dan Polling
4. Paket 4 - Alert, Backup, dan Resource Metric
5. Paket 5 - Event Store dan SIEM Fondasi

## DoD Coding Fase 3

- Inventory host, stack, VM, site, upstream, certificate, device, dan interface aktif.
- Poll sample, metric sample, alert rule, dan backup job aktif.
- Event store, parser, rule, alert, dan incident dasar tersedia.
- Dokumentasi BASE dan PLAN sinkron.

## Referensi

- phase3_backend_execution_backlog.md
- phase3_execution_backlog.md
- roadmap_platform_3_phase.md
- backend_module_matrix_platform_replacement.md
- ai_runbook.md
