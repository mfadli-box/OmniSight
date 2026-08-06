# Rencana Task File-by-File Paket 4 Fase 3

Dokumen ini memecah Paket 4 Fase 3 menjadi task file-by-file untuk klaster Alert, Backup, dan Resource Metric.

## Scope Paket 4

- Modul domain: infra_metric_sample, infra_alert_rule, net_alert_rule, net_backup_job, vm_resource_sample, vm_permission, vm_action_log.
- Tujuan: alerting dasar dan backup infra aktif dengan jejak audit yang jelas.

## Batch Task A - Finalisasi Schema Alert dan Metric

### File Target
- obx_base/prisma/schema/ict_monitoring.prisma
- obx_base/prisma/schema/ict_machine.prisma
- obx_base/prisma/schema/ict_mikrotik.prisma
- obx_base/prisma/schema/schema.prisma

### Verifikasi
- npx prisma generate di obx_base lulus.

## Batch Task B - API Rule, Metric, Backup, Action Log

### File Target
- obx_rest/backbone/routes.go
- obx_rest/skeleton/{KODE_INFRA_ALERT}/template.go
- obx_rest/skeleton/{KODE_INFRA_ALERT}/repository.go
- obx_rest/skeleton/{KODE_INFRA_ALERT}/usecase.go
- obx_rest/skeleton/{KODE_INFRA_ALERT}/handler.go

### Verifikasi
- go build ./... di obx_rest lulus.
- Rule alert dan backup job tercatat.

## Batch Task C - Frontend Board Alert dan Backup

### File Target
- obx_site/src/app/board/pages/{KODE_INFRA_ALERT}/page.tsx
- obx_site/src/app/board/model/module.ts

### Verifikasi
- npm run lint di obx_site lulus.
- UI menampilkan metric sample, alert rule, dan backup job.

## Batch Task D - Sinkronisasi Dokumentasi

### File Target
- obx_docs/blueprint/BASE/ict_machine.md
- obx_docs/guide/BASE/ict_machine.md
- obx_docs/blueprint/BASE/ict_mikrotik.md
- obx_docs/guide/BASE/ict_mikrotik.md
- obx_docs/blueprint/PLAN/README.md

### Verifikasi
- Dokumen konsisten dengan route dan payload API.

## Referensi

- phase3_backend_coding_plan.md
- phase3_backend_execution_backlog.md
- phase3_execution_backlog.md
- ai_runbook.md
