# agent_ws - Panduan Website/Log Ingest Agent

## Ringkasan Fitur

Agent ini dipakai untuk ingest log website, deteksi flood, dan sync analitik.

## Prasyarat Akses

- PostgreSQL dan Elasticsearch tersedia
- Direktori arsip writable

## Langkah Penggunaan

1. Isi `.env` agent_ws.
2. Jalankan container agent.
3. Pastikan log Nginx ter-mount ke container.

## Validasi Hasil

- Log masuk ke database.
- Elasticsearch terisi.

## Troubleshooting

- Jika DB error, cek PG env.
- Jika log tidak masuk, cek mount path.
