# ict_monitoring - Monitoring Infra Domain

## Ringkasan

- File schema: obx_base/prisma/schema/ict_monitoring.prisma
- Isi utama: metric sample, alert rule, backup job, dan histori operasional monitoring

## Model Inti

- infra_metric_sample
- infra_alert_rule
- net_alert_rule
- net_backup_job
- vm_resource_sample
- vm_permission
- vm_action_log

## Catatan Teknis

- Dipakai sebagai fondasi monitoring lintas host, VM, dan network.
- Menjadi sumber data untuk alert threshold dan action log operasional.
- Company scope wajib diterapkan pada query list, detail, dan histori.

## Checklist

- [ ] Metric sample tersimpan konsisten
- [ ] Alert rule dan backup job terdokumentasi
- [ ] Action log VM dapat ditelusuri
