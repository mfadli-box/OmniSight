# User Guide AUTO - obx_auto Agent & Microservices

## Ringkasan

obx_auto terdiri dari tiga agent Go yang berjalan sebagai container terpisah untuk manajemen infrastruktur, jaringan, dan log secara otomatis.

| Agent | Tugas Utama |
|---|---|
| agent_vm | Kelola Docker, Nginx, SSL, firewall, update, dan vuln scan di VM host |
| agent_mx | Monitor dan sinkronisasi perangkat MikroTik |
| agent_ws | Ingest dan analisis log Nginx, deteksi serangan, sinkronisasi ke Elasticsearch |

---

## agent_vm

### Prasyarat

- Docker daemon berjalan di host target
- Binary: lego (ACME), Trivy, git, nftables sudah tersedia di host
- Key robot sudah didaftarkan di backend (ict_vm_host aktif)
- .env sudah diisi: API_URL, KEY_ROBOT, IDA_HOST

### Menjalankan Agent

```sh
docker build -t agent_vm obx_auto/agent_vm/
docker run -d --env-file obx_auto/agent_vm/.env \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /etc/nginx:/etc/nginx \
  -v /etc/lego:/etc/lego \
  --name agent_vm agent_vm
```

### Memverifikasi Agent Aktif

1. Cek log container: `docker logs agent_vm`
2. Pastikan log menampilkan `[agent] VM Agent starting...` tanpa error.
3. Di dashboard backend, host target harus menunjukkan status online dan last_heartbeat yang baru.

### Cara Kerja Harian

- Agent **otomatis** memproses job yang dibuat dari dashboard backend.
- Tidak perlu intervensi manual selama heartbeat normal dan job tidak stuck.
- Jika job gagal, backend akan menampilkan status `failed` beserta log error.

### Troubleshooting

| Gejala | Kemungkinan Penyebab | Langkah |
|---|---|---|
| Host offline di dashboard | Heartbeat tidak sampai | Cek API_URL dan KEY_ROBOT, cek jaringan container ke backend |
| Job stuck di pending | Agent tidak poll | Restart container, cek log IVL_POOL |
| SSL gagal issue | lego tidak terinstall atau DNS salah | Cek instalasi lego dan konfigurasi domain |
| Trivy scan error | Cache corrupt atau trivy tidak terinstall | Hapus DIR_TRIVY, cek instalasi trivy |
| Docker command gagal | Socket tidak ter-mount | Pastikan volume `/var/run/docker.sock` di-mount |

---

## agent_mx

### Prasyarat

- Perangkat MikroTik dapat diakses dari host agent (TCP API port 8728/8729)
- Key robot sudah didaftarkan di backend (ict_mikrotik_device aktif)
- .env sudah diisi: API_URL, KEY_ROBOT, IDA_HOST, IPP_SYSLOG

### Menjalankan Agent

```sh
docker build -t agent_mx obx_auto/agent_mx/
docker run -d --env-file obx_auto/agent_mx/.env \
  -p 514:514/udp \
  --name agent_mx agent_mx
```

### Konfigurasi Syslog di MikroTik

Di RouterOS, arahkan syslog ke IP host agent_mx port 514 UDP:
```
/system logging action add name=remote target=remote remote=<IP_AGENT>
/system logging add action=remote topics=firewall,info,warning,error
```

### Memverifikasi Agent Aktif

1. Cek log container: `docker logs agent_mx`
2. Di dashboard backend, perangkat harus menunjukkan status online.
3. Log MikroTik harus mulai masuk ke ict_mikrotik_log.

### Troubleshooting

| Gejala | Kemungkinan Penyebab | Langkah |
|---|---|---|
| Device offline di dashboard | Heartbeat gagal | Cek API_URL, KEY_ROBOT, dan koneksi ke backend |
| Stats tidak update | statsLoop error | Cek akses ke MikroTik API port, cek kredensial |
| Syslog tidak masuk | Port 514 tidak terbuka | Cek port mapping container dan firewall host |
| Firewall sync gagal | Hak akses RouterOS kurang | Pastikan user RouterOS punya izin read untuk firewall |

---

## agent_ws

### Prasyarat

- PostgreSQL dan Elasticsearch tersedia dan dapat diakses dari container
- .env sudah diisi: PG_HOST, PG_PORT, PG_USER, PG_PASS, PG_DATA, ES_LINK
- Direktori arsip (RE_PATH) tersedia dan writable

### Menjalankan Agent

```sh
docker build -t agent_ws obx_auto/agent_ws/
docker run -d --env-file obx_auto/agent_ws/.env \
  -v /var/log/nginx:/var/log/nginx:ro \
  -v /data/archive:/archive \
  --name agent_ws agent_ws
```

### Memverifikasi Agent Aktif

1. Cek log container: `docker logs agent_ws`
2. Data log Nginx harus mulai muncul di tabel ict_nginx_log / ict_nginx_atc.
3. Di Elasticsearch, indeks log harus terupdate.

### Flood Detection dan Blacklist

- Agent otomatis mendeteksi IP yang melebihi FT_HTTP request dalam satu window waktu.
- IP yang terdeteksi akan ditambahkan ke ict_ip_blacklist.
- Nginx/WAF di host akan membaca blacklist ini untuk memblokir IP.

### Rotasi Log

- Log normal dihapus dari DB setelah RE_NORMAL hari (default 7).
- Log serangan dipertahankan RE_ATTACK hari (default 90).
- File arsip disimpan di RE_PATH.

### Troubleshooting

| Gejala | Kemungkinan Penyebab | Langkah |
|---|---|---|
| DB connection error | PG env salah atau DB tidak jalan | Verifikasi PG_HOST/PORT/USER/PASS/DATA |
| Elasticsearch error | ES_LINK tidak valid atau ES down | Cek URL dan status Elasticsearch |
| Log tidak masuk | Path file log tidak ter-mount | Cek volume mount log Nginx |
| Arsip gagal | RE_PATH tidak writable | Cek permission direktori arsip |

---

## Monitoring Status Semua Agent

Di dashboard backend (obx_site):
- **VM Host**: menu ict_machine → status host, last_heartbeat, total container/image
- **MikroTik**: menu ict_mikrotik → status device, last_seen
- **Website Log**: menu ict_website → statistik log, SLA, blacklist IP
