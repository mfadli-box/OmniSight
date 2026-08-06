# agent_mx - Panduan MikroTik Agent

## Ringkasan Fitur

Agent ini memantau dan menyinkronkan perangkat MikroTik.

## Prasyarat Akses

- Perangkat MikroTik reachable
- KEY_ROBOT dan IDA_HOST valid
- Port syslog terbuka

## Langkah Penggunaan

1. Isi `.env` agent_mx.
2. Jalankan container agent.
3. Arahkan syslog MikroTik ke host agent.

## Validasi Hasil

- Perangkat tampil online.
- Log syslog masuk ke backend.

## Troubleshooting

- Jika syslog tidak masuk, cek port 514 UDP dan firewall host.
