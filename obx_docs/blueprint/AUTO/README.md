# Rancangan Submodul: AUTO - obx_auto Agent & Microservices

## Ringkasan

- Scope: tiga agent Go yang berjalan sebagai microservice terpisah
- Tujuan: otomasi VM, MikroTik, dan log ingest website
- Komunikasi ke backend: HTTP REST via obx_rest robot endpoint

## Struktur Dokumen

- Summary teknis: file ini
- Detail page teknis: `obx_docs/blueprint/AUTO/*.md`
- Summary user guide: `obx_docs/guide/AUTO.md`
- Detail page user guide: `obx_docs/guide/AUTO/*.md`

## Agent Domain

- `agent_vm`: VM host, Docker, Nginx, SSL, firewall, GitOps, vuln scan, update
- `agent_mx`: MikroTik stats, firewall, address list, backup, syslog
- `agent_ws`: Nginx log ingest, attack detection, rotation, Elasticsearch sync

## Pattern Summary

- Setiap agent punya `.env` sendiri
- agent_vm dan agent_mx autentikasi ke backend via `KEY_ROBOT`
- agent_ws memakai PostgreSQL dan Elasticsearch langsung
- Job dispatch, stats, heartbeat, dan ingest dijalankan lewat loop/ticker

## Page Index

| Page | Detail Teknis | Detail Guide |
|---|---|---|
| agent_vm | agent_vm.md | agent_vm.md |
| agent_mx | agent_mx.md | agent_mx.md |
| agent_ws | agent_ws.md | agent_ws.md |

## Checklist

- [ ] Env per agent terisi
- [ ] Token robot aktif
- [ ] Docker image build berhasil
- [ ] Heartbeat dan job dispatch berjalan
- [ ] Log ingest dan rotation berjalan
- [ ] Blueprint DB terkait sinkron

## Referensi

- `obx_auto/agent_vm/*`
- `obx_auto/agent_mx/*`
- `obx_auto/agent_ws/*`
- `obx_docs/blueprint/BASE/README.md`
