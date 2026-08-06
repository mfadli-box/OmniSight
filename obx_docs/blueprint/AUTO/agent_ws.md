# agent_ws - Website/Log Ingest Agent

## Ringkasan

- Scope: agent Go untuk ingest log Nginx, deteksi serangan, rotasi arsip, dan sync Elasticsearch
- Entry: `obx_auto/agent_ws/main.go`
- Backend DB: PostgreSQL langsung + Elasticsearch

## Technical Notes

- Membaca env PG_HOST, PG_PORT, PG_USER, PG_PASS, PG_DATA, ES_LINK, FT_HTTP, RE_PATH, RE_NORMAL, RE_ATTACK
- Fungsi utama: LoadConfig dan InitDB
- Memproses log Nginx ke PostgreSQL
- Melakukan flood detection, blacklist update, dan sync ke Elasticsearch

## Checklist

- [ ] DB dan ES aktif
- [ ] Log ingest aktif
- [ ] Flood detection aktif
- [ ] Rotation berjalan
