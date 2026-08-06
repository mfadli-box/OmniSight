# agent_vm - Panduan VM & Infrastructure Agent

## Ringkasan Fitur

Agent ini mengelola Docker, Nginx, SSL, firewall, update, dan scan keamanan di VM host.

## Prasyarat Akses

- Docker daemon aktif
- lego, Trivy, git, nftables tersedia
- KEY_ROBOT dan IDA_HOST valid

## Langkah Penggunaan

1. Isi `.env` agent_vm.
2. Build dan jalankan container agent.
3. Cek log untuk heartbeat dan job.

## Validasi Hasil

- Host muncul online di backend.
- Job status berubah dari pending ke success/failed.

## Troubleshooting

- Jika heartbeat gagal, cek API_URL dan KEY_ROBOT.
- Jika Docker tidak terbaca, cek socket mount.
