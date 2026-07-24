# Plan: Automation / CLI / Scheduler (ict_auto)

Dokumen ini berisi rencana pengembangan komponen **ict_auto** — Go CLI applications yang berfungsi sebagai agent atau isolated automate microservice.

---

## Struktur Folder

```
ict_auto/
├── ict_log_nginx/          ✅ Sudah ada
│   ├── main.go
│   ├── go.mod / go.sum
│   ├── Dockerfile
│   └── .env / .env.dev / .env.pub
├── ict_log_rotate/         ✅ Sudah ada
│   ├── main.go
│   ├── go.mod / go.sum
│   ├── Dockerfile
│   └── .env / .env.dev / .env.pub
├── ict_mikrotik_agent/     ⬜ Direncanakan
│   ├── main.go
│   ├── config.go
│   ├── collector/ (snmp.go, api.go)
│   ├── syncer/ (device.go, hotspot.go, dhcp.go, firewall.go, queue.go, port.go, vlan.go, client.go, access.go)
│   ├── go.mod / go.sum
│   ├── Dockerfile
│   └── .env / .env.dev / .env.pub
├── ict_docker_agent/       ⬜ Direncanakan
│   ├── main.go
│   ├── config.go
│   ├── api.go
│   ├── collector/ (stats.go, host.go, list.go)
│   ├── manager/ (container.go, image.go, network.go, volume.go)
│   ├── compose/ (scanner.go, parser.go, deployer.go)
│   ├── hub/ (client.go, sync.go)
│   ├── go.mod / go.sum
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── .env / .env.dev / .env.pub
├── ict_host_monitor/       ⬜ Direncanakan
│   ├── main.go
│   ├── config.go
│   ├── collector/ (cpu.go, memory.go, disk.go, network.go, load.go)
│   ├── hub/ (client.go, sync.go)
│   ├── go.mod / go.sum
│   ├── Dockerfile
│   └── .env / .env.dev / .env.pub
├── ict_network_poller/     ⬜ Direncanakan
│   ├── main.go
│   ├── config.go
│   ├── poller/ (snmp.go, interface.go, system.go)
│   ├── trap/ (receiver.go)
│   ├── hub/ (client.go, sync.go)
│   ├── go.mod / go.sum
│   ├── Dockerfile
│   └── .env / .env.dev / .env.pub
├── ict_fim_agent/          ⬜ Direncanakan
│   ├── main.go
│   ├── config.go
│   ├── scanner/ (baseline.go, hasher.go, comparator.go)
│   ├── watcher/ (fsnotify.go)
│   ├── hub/ (client.go, sync.go)
│   ├── go.mod / go.sum
│   ├── Dockerfile
│   └── .env / .env.dev / .env.pub
├── ict_uptimerobot_agent/  ⬜ Direncanakan
│   ├── main.go
│   ├── config.go
│   ├── syncer/ (monitor.go, incident.go, sla.go)
│   ├── hub/ (client.go, sync.go)
│   ├── go.mod / go.sum
│   ├── Dockerfile
│   └── .env / .env.dev / .env.pub
└── ict_sca_agent/          ⬜ Direncanakan
    ├── main.go
    ├── config.go
    ├── scanner/ (policy.go, ssh.go, firewall.go, users.go, permissions.go)
    ├── hub/ (client.go, sync.go)
    ├── go.mod / go.sum
    ├── Dockerfile
    └── .env / .env.dev / .env.pub
```

Semua agent berdiri sendiri (*standalone*) tanpa shared library. Masing-masing adalah single `main.go` dengan `package main` yang berjalan sebagai long-running daemon (bukan CLI dengan subcommand).

---

## 1. Nginx Log Sync & WAF Agent (`ict_log_nginx`)

Agent Go yang berjalan sebagai daemon untuk menarik log nginx dari Elasticsearch, mengklasifikasi traffic untuk WAF (Web Application Firewall), menulis log terstruktur ke PostgreSQL, menghitung metrik SLA, melacak threat score, dan auto-ban IP berbahaya. Juga melakukan cleanup indeks Elasticsearch yang sudah usang.

### Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Agent | Go 1.26.4 (single binary, `package main`) |
| Database | PostgreSQL (via `database/sql` + `lib/pq`) |
| Search Engine | Elasticsearch 8.x (via `go-elasticsearch/v8`) |
| Config | `godotenv` (`.env` file) |

### Dependencies

| Library | Versi | Fungsi |
|---------|-------|--------|
| `github.com/elastic/go-elasticsearch/v8` | v8.19.6 | Klien Elasticsearch 8.x |
| `github.com/google/uuid` | v1.6.0 | Pembuatan UUID untuk row ID |
| `github.com/joho/godotenv` | v1.5.1 | Loading file `.env` |
| `github.com/lib/pq` | v1.12.3 | Driver PostgreSQL |

### Arsitektur

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│  Elasticsearch   │────>│  ict_log_nginx   │────>│    PostgreSQL    │
│  (Logstash)      │     │  (WAF + Sync)    │     │  (ict database)  │
└──────────────────┘     └──────────────────┘     └──────────────────┘
                              │
                         ┌────┴────┐
                         │ Threat  │
                         │ Classify│
                         │ & Auto  │
                         │  Ban    │
                         └─────────┘
```

### Alur Kerja (setiap 1 menit)

1. **Hapus dokumen ES yang rusak** — dokumen JSON tidak valid, `status=0`, `client_ip` kosong
2. **Bersihkan baris DB** — hapus baris dengan `client_ip` kosong, hitung ulang SLA
3. **Ambil hingga 2500 dokumen** dari indeks Logstash (kemarin + hari ini)
4. **Hitung jumlah hit per IP** untuk deteksi flood
5. **Untuk setiap log entry:**
   - Klasifikasi traffic (WAF)
   - Cek whitelist/bypass
   - Route ke tabel DB yang sesuai
   - Lacak serangan, update threat score
6. **Tulis ringkasan serangan** ke `ict_nginx_atc_sum`
7. **Upsert metrik SLA** ke `ict_nginx_sla`
8. **Hapus dokumen yang sudah diproses** dari ES

### Threat Classification (WAF)

Agent mengklasifikasi traffic menggunakan regex pattern matching:

| Threat Type | Pattern | Bobot Threat |
|-------------|---------|-------------|
| `RCE_COMMAND_INJECTION` | Command injection patterns | 50 |
| `SQL_INJECTION` | SQL injection patterns | 40 |
| `PATH_TRAVERSAL` | Directory traversal (`../`, `etc/passwd`) | 35 |
| `XSS` | Cross-site scripting patterns | 20 |
| `SENSITIVE_FILE_PROBING` | Sensitive file access attempts | 15 |
| `BOT_SCANNER` | Bot/scanner user-agent | 10 |
| `NORMAL` | Tidak terdeteksi | 0 |

### Auto-Ban Mechanism

- Setiap threat yang terdeteksi menambahkan *threat score* ke IP
- Jika IP sudah di-whitelist: skip
- Jika score ≥ 100: IP otomatis di-ban (INSERT/UPSERT ke `ict_ip_blacklist` dengan expiry 24 jam)
- Score maksimum: 150

### PostgreSQL Tables

| Tabel | Operasi | Keterangan |
|-------|---------|------------|
| `ict_nginx_log` | INSERT, DELETE | Log nginx normal |
| `ict_nginx_app` | INSERT, DELETE | Log traffic aplikasi/whitelisted |
| `ict_nginx_atc` | INSERT, DELETE | Log traffic serangan |
| `ict_nginx_atc_sum` | INSERT UPSERT | Ringkasan serangan per hari/IP/tipe/domain |
| `ict_nginx_sla` | INSERT UPSERT, UPDATE | Metrik SLA harian (total, success, error, attack, avg response time, SLA %) |
| `ict_ip_whitelist` | SELECT | IP dan CIDR yang dipercaya |
| `ict_ip_blacklist` | SELECT, INSERT UPSERT | IP berbahaya yang di-ban otomatis |
| `ict_waf_bypass_rule` | SELECT | Aturan bypass WAF untuk domain/path tertentu |

### In-Memory Cache

Cache di-refresh setiap 5 menit dari database:

| Data | Sumber | Fungsi |
|------|--------|--------|
| Whitelist IP/CIDR | `ict_ip_whitelist` | Mengecualikan IP dari WAF classification |
| Banned IP | `ict_ip_blacklist` | Mencegah akses dari IP berbahaya |
| Bypass Rules | `ict_waf_bypass_rule` | Melewati klasifikasi WAF untuk domain/path tertentu |

### Konfigurasi (Environment Variable)

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `PG_HOST` | `localhost` | Host PostgreSQL |
| `PG_PORT` | `5432` | Port PostgreSQL |
| `PG_USER` | `dbe` | User PostgreSQL |
| `PG_PASS` | `rahasia` | Password PostgreSQL |
| `PG_DATA` | `erp` | Nama database |
| `IS_POOL` | `false` | Gunakan connection pool |
| `ES_LINK` | `http://localhost:9200` | URL Elasticsearch |
| `FT_HTTP` | `150` | Threshold flood detection (request per IP per batch) |

### Self-Healing

Jika komunikasi dengan Elasticsearch gagal, agent akan me-restart diri sendiri dengan me-re-execute binary (`restartSelf()`).

### Dockerfile

Multi-stage build:
- **Stage 1 (builder):** `golang:1.26.4-alpine` — kompilasi statik dengan `CGO_ENABLED=0 GOOS=linux`
- **Stage 2 (runtime):** `alpine:3.20` — minimal image dengan `ca-certificates` dan `tzdata`

### Deployment

```bash
# Docker
docker run -d \
  --name ict-log-nginx \
  -e PG_HOST=erp_pool \
  -e PG_PORT=5432 \
  -e PG_USER=dbe \
  -e PG_PASS=rahasia \
  -e PG_DATA=erp \
  -e ES_LINK=http://elasticsearch:9200 \
  -e FT_HTTP=150 \
  ict_auto/ict_log_nginx
```

---

## 2. Log Archive & Retention Worker (`ict_log_rotate`)

Agent Go yang berjalan sebagai daemon untuk melakukan arsip dan retensi data harian. Menghapus baris lama dari tabel nginx log di PostgreSQL, menyimpannya sebagai file newline-delimited JSON (JSONL) sebelum dihapus. Tabel berbeda memiliki periode retensi berbeda.

### Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Agent | Go 1.26.4 (single binary, `package main`) |
| Database | PostgreSQL (via `database/sql` + `lib/pq`) |
| Config | `godotenv` (`.env` file) |

### Dependencies

| Library | Versi | Fungsi |
|---------|-------|--------|
| `github.com/joho/godotenv` | v1.5.1 | Loading file `.env` |
| `github.com/lib/pq` | v1.12.3 | Driver PostgreSQL |

### Alur Kerja (setiap 24 jam)

1. **Arsip dan hapus** baris dari `ict_nginx_log` (retensi normal)
2. **Arsip dan hapus** baris dari `ict_nginx_app` (retensi serangan)
3. **Arsip dan hapus** baris dari `ict_nginx_atc` (retensi serangan)

### Mekanisme Arsip

- Menggunakan `DELETE ... WHERE ctid IN (SELECT ctid ... LIMIT 5000) RETURNING row_to_json(...)` untuk menghapus dan mengekstrak baris secara atomik
- Batch processing: 5000 baris per iterasi sampai sisa kurang dari 5000
- File kosong otomatis dihapus
- Menggunakan buffered writer (64KB) untuk I/O efisien

### Naming File Arsip

Format: `{nama_tabel}_archive_{YYYY-MM-DD}.log`

Contoh: `ict_nginx_log_archive_2026-07-22.log`

Setiap baris adalah satu JSON object lengkap (`row_to_json`) yang mewakili satu baris yang dihapus.

### PostgreSQL Tables

| Tabel | Operasi | Retensi Default |
|-------|---------|-----------------|
| `ict_nginx_log` | DELETE + RETURNING | 30 hari (`.env`) / 7 hari (Dockerfile) |
| `ict_nginx_app` | DELETE + RETURNING | 90 hari |
| `ict_nginx_atc` | DELETE + RETURNING | 90 hari |

### Konfigurasi (Environment Variable)

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `PG_HOST` | `localhost` | Host PostgreSQL |
| `PG_PORT` | `5432` | Port PostgreSQL |
| `PG_USER` | `dbe` | User PostgreSQL |
| `PG_PASS` | `rahasia` | Password PostgreSQL |
| `PG_DATA` | `erp` | Nama database |
| `IS_POOL` | `false` | Gunakan connection pool |
| `RE_PATH` | `archive` | Direktori penyimpanan arsip |
| `RE_NORMAL` | `30` | Hari retensi log normal |
| `RE_ATTACK` | `90` | Hari retensi log serangan |

### Dockerfile

Multi-stage build:
- **Stage 1 (builder):** `golang:1.26.4-alpine` — kompilasi statik dengan `CGO_ENABLED=0 GOOS=linux`
- **Stage 2 (runtime):** `alpine:3.20` — minimal image dengan `ca-certificates` dan `tzdata`

### Deployment

```bash
# Docker
docker run -d \
  --name ict-log-rotate \
  -e PG_HOST=erp_pool \
  -e PG_PORT=5432 \
  -e PG_USER=dbe \
  -e PG_PASS=rahasia \
  -e PG_DATA=erp \
  -e IS_POOL=true \
  -e RE_PATH=/var/log/ict_log_rotate \
  -e RE_NORMAL=30 \
  -e RE_ATTACK=90 \
  -v /var/log/ict_log_rotate:/app/archive \
  ict_auto/ict_log_rotate
```

---

## 3. Pipeline Integrasi: `ict_log_nginx` → `ict_log_rotate`

Kedua agent beroperasi pada database PostgreSQL yang sama (`ict`) dan berbagi tabel nginx log:

| Tabel | Ditulis oleh `ict_log_nginx` | Diarsip/dihapus oleh `ict_log_rotate` |
|-------|------------------------------|----------------------------------------|
| `ict_nginx_log` | INSERT | DELETE + archive |
| `ict_nginx_app` | INSERT | DELETE + archive |
| `ict_nginx_atc` | INSERT | DELETE + archive |
| `ict_nginx_atc_sum` | INSERT/UPSERT | — |
| `ict_nginx_sla` | INSERT/UPSERT | — |
| `ict_ip_whitelist` | SELECT (read-only) | — |
| `ict_ip_blacklist` | SELECT, INSERT/UPSERT | — |
| `ict_waf_bypass_rule` | SELECT (read-only) | — |

Alur pipeline:
1. **`ict_log_nginx`** menarik log dari Elasticsearch, mengklasifikasinya (WAF), dan menulis ke PostgreSQL. Berjalan setiap 1 menit.
2. **`ict_log_rotate`** menghapus baris lama dari tabel-tabel tersebut (setelah diarsip ke file JSONL di disk). Berjalan setiap 24 jam.

---

## 4. Agent yang Direncanakan

### 4.1 Docker Agent (`ict_docker_agent`)

Agent Go yang berjalan di setiap Docker host (**per-host deployment**) untuk monitoring, manajemen container, dan deployment Docker Compose stacks. Berjalan sebagai **non-root user** dengan akses ke Docker socket via group permission.

#### Arsitektur Deployment

```
Host A (Docker)                      Host B (Docker)
┌─────────────────────────┐         ┌─────────────────────────┐
│ Docker Engine           │         │ Docker Engine           │
│ /var/run/docker.sock    │         │ /var/run/docker.sock    │
│   (group: docker, 660)  │         │   (group: docker, 660)  │
└─────────┬───────────────┘         └─────────┬───────────────┘
          │                                   │
┌─────────┴───────────────┐         ┌─────────┴───────────────┐
│ ict_docker_agent        │         │ ict_docker_agent        │
│ user: omnisight (gid:   │         │ user: omnisight (gid:   │
│        docker)          │         │        docker)          │
│                         │         │                         │
│ /opt/omnisight/stacks/  │         │ /opt/omnisight/stacks/  │
│ └── docker-compose.yml  │         │ └── docker-compose.yml  │
└─────────┬───────────────┘         └─────────┬───────────────┘
          │                                   │
          │  REST API (ict_rest)              │
          └───────────────┬───────────────────┘
                          ▼
              ┌─────────────────────┐
              │  PostgreSQL Server  │
              │  (OmniSight DB)     │
              └─────────────────────┘
```

**Prinsip desain:**
- **Per-host** — 1 agent instance per Docker host, Docker socket di-mount
- **Non-root** — user `omnisight` (atau custom) dalam group `docker`
- **Compose-local** — compose files di folder lokal host (`/opt/omnisight/stacks/`)
- **State sync** — semua data dikirim ke PostgreSQL via `ict_rest` API
- **Audit trail** — semua operasi management di-log

#### Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Agent | Go 1.26.4 (single binary, `package main`) |
| Docker SDK | `github.com/docker/docker/client` (official Docker SDK) |
| Docker Compose | `github.com/compose-spec/compose-go` (parse compose files) |
| HTTP API | `net/http` (agent HTTP endpoint untuk remote commands) |
| Database | PostgreSQL via `ict_rest` API (REST, bukan direct DB) |
| Config | `godotenv` (`.env` file) |

#### Dependencies

| Library | Versi | Fungsi |
|---------|-------|--------|
| `github.com/docker/docker` | v28.x | Official Docker Engine SDK |
| `github.com/docker/go-connections` | v0.5.x | Docker connection helpers |
| `github.com/compose-spec/compose-go` | v2.x | Parse Docker Compose files |
| `github.com/joho/godotenv` | v1.5.1 | Loading file `.env` |
| `gopkg.in/yaml.v3` | v3.x | YAML parsing (compose files) |

#### Non-Root Setup

Agent berjalan sebagai user non-root dengan akses Docker socket:

```bash
# 1. Buat user omnisight
useradd -r -s /bin/false -g docker omnisight

# 2. Set permissions pada Docker socket
chmod 660 /var/run/docker.sock
chown root:docker /var/run/docker.sock

# 3. Buat direktori stacks
mkdir -p /opt/omnisight/stacks
chown -R omnisight:docker /opt/omnisight/stacks

# 4. Jalankan agent
sudo -u omnisight ./ict_docker_agent
```

**Dockerfile:**
```dockerfile
FROM golang:1.26.4-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ict_docker_agent .

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata docker-cli docker-cli-compose
RUN addgroup -g 1001 docker && \
    adduser -u 1001 -G docker -s /bin/false -D omnisight

WORKDIR /root/
COPY --from=builder /app/ict_docker_agent .

ENV DOCKER_HOST=unix:///var/run/docker.sock
ENV COMPOSE_DIR=/opt/omnisight/stacks
ENV AGENT_PORT=9100

CMD ["su", "-s", "/bin/sh", "-c", "./ict_docker_agent", "omnisight"]
```

**docker-compose.yml (deploy agent):**
```yaml
version: "3.8"
services:
  ict_docker_agent:
    build: .
    container_name: ict_docker_agent
    restart: unless-stopped
    user: "1001:1001"  # omnisight:docker
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /opt/omnisight/stacks:/opt/omnisight/stacks:rw
    ports:
      - "127.0.0.1:9100:9100"  # Localhost only
    environment:
      - HUB_URL=http://ict_rest:8080
      - HUB_TOKEN=${HUB_TOKEN}
      - AGENT_PORT=9100
      - COMPOSE_DIR=/opt/omnisight/stacks
      - COLLECT_STATS_INTERVAL=10
      - COLLECT_HOST_INTERVAL=30
      - SYNC_INTERVAL=60
```

#### Alur Kerja

Agent menjalankan **3 worker loop** secara paralel:

**Worker 1: Metrics Collector** (setiap 10-60 detik)
1. Container stats (CPU, memory, network, block I/O) → setiap 10 detik
2. Host stats (CPU, memory, disk, load) → setiap 30 detik
3. Container/image/network/volume list → setiap 60 detik
4. Kirim metrics ke Hub via REST API

**Worker 2: Compose Scanner** (setiap 5 menit)
1. Scan `COMPOSE_DIR` untuk `docker-compose*.yml`
2. Parse setiap compose file → identifikasi services
3. Compare dengan container yang running di Docker
4. Update state stack di Hub (running, stopped, missing)
5. Auto-detect perubahan compose file (hash comparison)

**Worker 3: HTTP API Server** (port 9100, localhost only)
1. Terima command dari Hub/Frontend via REST API
2. Eksekusi operasi management
3. Return result + audit log

#### Container Management

| Operasi | HTTP Endpoint | SDK Method |
|---------|--------------|------------|
| List containers | `GET /containers` | `ContainerList()` |
| Inspect container | `GET /containers/:id` | `ContainerInspect()` |
| Create container | `POST /containers` | `ContainerCreate()` |
| Start container | `POST /containers/:id/start` | `ContainerStart()` |
| Stop container | `POST /containers/:id/stop` | `ContainerStop()` |
| Restart container | `POST /containers/:id/restart` | `ContainerRestart()` |
| Remove container | `DELETE /containers/:id` | `ContainerRemove()` |
| Container logs | `GET /containers/:id/logs` | `ContainerLogs()` |
| Exec in container | `POST /containers/:id/exec` | `ContainerExecCreate()` |
| Container stats | `GET /containers/:id/stats` | `ContainerStats()` |

#### Image Management

| Operasi | HTTP Endpoint | SDK Method |
|---------|--------------|------------|
| List images | `GET /images` | `ImageList()` |
| Pull image | `POST /images/pull` | `ImagePull()` |
| Remove image | `DELETE /images/:id` | `ImageRemove()` |
| Image inspect | `GET /images/:id` | `ImageInspect()` |

#### Network Management

| Operasi | HTTP Endpoint | SDK Method |
|---------|--------------|------------|
| List networks | `GET /networks` | `NetworkList()` |
| Create network | `POST /networks` | `NetworkCreate()` |
| Connect to network | `POST /networks/:id/connect` | `NetworkConnect()` |
| Disconnect from network | `POST /networks/:id/disconnect` | `NetworkDisconnect()` |
| Remove network | `DELETE /networks/:id` | `NetworkRemove()` |

#### Volume Management

| Operasi | HTTP Endpoint | SDK Method |
|---------|--------------|------------|
| List volumes | `GET /volumes` | `VolumeList()` |
| Create volume | `POST /volumes` | `VolumeCreate()` |
| Remove volume | `DELETE /volumes/:name` | `VolumeRemove()` |
| Inspect volume | `GET /volumes/:name` | `VolumeInspect()` |

#### Docker Compose Stack Management

| Operasi | HTTP Endpoint | Keterangan |
|---------|--------------|------------|
| List stacks | `GET /stacks` | Daftar compose files + status |
| Deploy stack | `POST /stacks` | Parse compose → create containers |
| Update stack | `PUT /stacks/:name` | Re-deploy (rolling update) |
| Destroy stack | `DELETE /stacks/:name` | Stop + remove all containers |
| Stack status | `GET /stacks/:name` | Status semua services |
| Stack logs | `GET /stacks/:name/logs` | Logs dari semua services |

**Alur Deploy Stack:**
```
1. Receive compose file content + stack name
2. Parse compose file → identifikasi services, networks, volumes
3. Buat networks jika belum ada
4. Buat volumes jika belum ada
5. Untuk setiap service:
   a. Pull image jika belum ada
   b. Stop + remove container lama (jika update)
   c. Create container baru
   d. Start container
   e. Connect ke networks
6. Update state stack di Hub
7. Audit log: stack deployed/updated
```

**Alur Update Stack (Rolling):**
```
1. Receive compose file content baru
2. Parse → bandingkan dengan state saat ini
3. Untuk setiap service yang berubah:
   a. Create container baru dengan config baru
   b. Start container baru
   c. Stop container lama
   d. Remove container lama
4. Untuk service yang tidak berubah: skip
5. Update state + audit log
```

#### PostgreSQL Tables (via Hub API)

| Tabel | Operasi | Keterangan |
|-------|---------|------------|
| `ict_docker_host` | UPSERT | Host info (hostname, IP, OS, Docker version) |
| `ict_docker_container` | UPSERT | Container state (id, name, image, status, stats) |
| `ict_docker_container_port` | UPSERT | Port mapping per container |
| `ict_docker_container_env` | UPSERT | Environment variables per container |
| `ict_docker_container_mount` | UPSERT | Volume mounts per container |
| `ict_docker_container_network` | UPSERT | Network attachments per container |
| `ict_docker_image` | UPSERT | Image info (id, tags, size, created) |
| `ict_docker_network` | UPSERT | Network info (id, name, driver, containers) |
| `ict_docker_volume` | UPSERT | Volume info (name, driver, mountpoint) |
| `ict_docker_compose` | UPSERT | Stack state (name, compose content, status) |
| `ict_docker_compose_service` | UPSERT | Service state (name, image, replicas, status) |
| `ict_docker_stats` | INSERT | Container stats historical |
| `ict_docker_host_stats` | INSERT | Host system stats historical |
| `ict_audit_trail` | INSERT | Audit trail (entity_type=docker_host, action, user, note) |

#### In-Memory Cache

| Data | Sumber | Fungsi |
|------|--------|--------|
| Container State | Docker API | Quick reference untuk operasi |
| Compose File Hash | Filesystem | Deteksi perubahan compose file |
| Stack Mapping | `ict_docker_stack` | Map stack name → containers |

#### Konfigurasi (Environment Variable)

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `DOCKER_HOST` | `unix:///var/run/docker.sock` | Docker socket path |
| `AGENT_PORT` | `9100` | HTTP API port (localhost only) |
| `HUB_URL` | `http://localhost:8080` | URL Hub API (`ict_rest`) |
| `HUB_TOKEN` | — | Token autentikasi ke Hub |
| `COMPOSE_DIR` | `/opt/omnisight/stacks` | Direktori compose files |
| `COLLECT_STATS_INTERVAL` | `10` | Detik antara container stats |
| `COLLECT_HOST_INTERVAL` | `30` | Detik antara host stats |
| `SYNC_INTERVAL` | `60` | Detik antara sync container list |
| `COMPOSE_SCAN_INTERVAL` | `300` | Detik antara scan compose files |

#### Audit Logging

Setiap operasi management di-log ke `ict_audit_trail` (entity_type = `docker_host`):

```json
{
  "company_id": "comp-001",
  "entity_type": "docker_host",
  "entity_id": "docker-host-01",
  "user_id": "user-admin",
  "action": "stack.deploy",
  "old_value": null,
  "new_value": "omnisight-stack (5 services)",
  "note": "/opt/omnisight/stacks/docker-compose.yml",
  "attachment_url": null
}
```

**Action values:**
- `stack.deploy` — Deploy stack baru
- `stack.update` — Update stack (rolling)
- `stack.destroy` — Hapus stack
- `container.start` — Start container
- `container.stop` — Stop container
- `container.restart` — Restart container
- `container.remove` — Hapus container
- `image.pull` — Pull image baru
- `image.remove` — Hapus image
- `network.create` — Buat network baru
- `network.remove` — Hapus network
- `volume.create` — Buat volume baru
- `volume.remove` — Hapus volume
- `exec.command` — Eksekusi command di container

#### Struktur Folder

```
ict_auto/ict_docker_agent/
├── main.go              ← Entry point + scheduler (3 worker loops)
├── config.go            ← Konfigurasi dari env
├── api.go               ← HTTP API server (port 9100)
├── collector/
│   ├── stats.go         ← Container stats collection
│   ├── host.go          ← Host stats collection
│   └── list.go          ← Container/image/network/volume list
├── manager/
│   ├── container.go     ← Container lifecycle (create/start/stop/restart/remove)
│   ├── image.go         ← Image management (pull/remove)
│   ├── network.go       ← Network management (create/connect/disconnect)
│   └── volume.go        ← Volume management (create/remove)
├── compose/
│   ├── scanner.go       ← Scan COMPOSE_DIR untuk compose files
│   ├── parser.go        ← Parse compose file → service definitions
│   └── deployer.go      ← Deploy/update/destroy stacks via SDK
├── hub/
│   ├── client.go        ← HTTP client ke Hub API
│   └── sync.go          ← Sync state ke Hub
├── go.mod / go.sum
├── Dockerfile
├── docker-compose.yml   ← Deploy agent di setiap host
└── .env / .env.dev / .env.pub

```

#### Deployment per Host

```bash
# 1. Clone repo ke host
git clone https://github.com/omnisight/ict_auto.git /opt/ict_auto

# 2. Setup user & permissions
useradd -r -g docker omnisight
mkdir -p /opt/omnisight/stacks
chown -R omnisight:docker /opt/omnisight/stacks

# 3. Copy compose files ke stacks
cp /path/to/your/stacks/*.yml /opt/omnisight/stacks/

# 4. Buat .env
cat > /opt/ict_auto/ict_docker_agent/.env << EOF
DOCKER_HOST=unix:///var/run/docker.sock
AGENT_PORT=9100
HUB_URL=http://<omnisight-server>:8080
HUB_TOKEN=<token>
COMPOSE_DIR=/opt/omnisight/stacks
EOF

# 5. Jalankan via docker compose
cd /opt/ict_auto/ict_docker_agent
docker compose up -d
```

#### Security Considerations

| Aspek | Implementasi |
|-------|-------------|
| Docker Socket | Mount read-write, group `docker` only (660) |
| Agent User | Non-root (`omnisight:docker`), no login shell |
| HTTP API | Bind localhost only (`127.0.0.1:9100`), tidak exposed ke network |
| Compose Deploy | Hanya folder `COMPOSE_DIR` yang boleh di-deploy |
| Audit Trail | Semua operasi management di-log ke database |
| Hub Auth | Token-based auth ke Hub API, tidak plain |
| Container Exec | Audit log semua exec commands, restrict via Hub API |
| Image Pull | Hanya dari registry yang di-whitelist (optional) |

---

### 4.2 SIEM Agents

| Agent | Fungsi | Agent Type | Status |
|-------|--------|-----------|--------|
| `ict_siem_evaluator` | Menjalankan rule evaluation terhadap log masuk | siem | ⬜ Direncanakan |
| `ict_siem_notifier` | Mengirim notifikasi alert (email, webhook, dll) | siem | ⬜ Direncanakan |

#### 4.2.1 SIEM Evaluator (`ict_siem_evaluator`)

Agent Go yang berjalan sebagai daemon untuk mengevaluasi rule SIEM terhadap log yang masuk ke PostgreSQL. Mendukung threshold, pattern, anomaly, dan compound rules.

**Tech Stack:** Go 1.26.4, PostgreSQL (direct), godotenv

**Alur Kerja (setiap 30 detik):**
1. Ambil rule aktif dari `ict_siem_rule`
2. Query log terbaru dari `ict_monitor_log_entry` / `ict_nginx_log` (semenit terakhir)
3. Untuk setiap rule, eksekusi query rule terhadap log
4. Jika match → insert alert ke `ict_siem_alert`
5. Kirim alert ke `ict_siem_notifier` via channel/queue

**PostgreSQL Tables:**

| Tabel | Operasi |
|-------|---------|
| `ict_siem_rule` | SELECT (read-only) |
| `ict_monitor_log_entry` | SELECT (read-only) |
| `ict_nginx_log` | SELECT (read-only) |
| `ict_siem_alert` | INSERT |

#### 4.2.2 SIEM Notifier (`ict_siem_notifier`)

Agent Go yang berjalan sebagai daemon untuk mengirim notifikasi alert dari SIEM ke berbagai channel (email, webhook, Telegram, dll).

**Tech Stack:** Go 1.26.4, PostgreSQL (direct), godotenv, SMTP client

**Alur Kerja (setiap 10 detik):**
1. Ambil alert yang belum di-notifikasi dari `ict_siem_alert`
2. Untuk setiap alert, cek `ict_siem_notification_config` untuk channel yang aktif
3. Kirim notifikasi via channel yang sesuai
4. Update `ict_siem_alert` dengan status notified

**PostgreSQL Tables:**

| Tabel | Operasi |
|-------|---------|
| `ict_siem_alert` | SELECT, UPDATE |
| `ict_siem_notification_config` | SELECT (read-only) |

### SIEM Rule Types

#### Threshold Rule
Mencapai threshold dalam waktu tertentu.
```json
{
  "field": "status",
  "operator": "equals",
  "value": "401",
  "threshold": 10,
  "timeframe_minutes": 5,
  "group_by": "client_ip"
}
```

#### Pattern Rule
Mendeteksi pola string tertentu.
```json
{
  "field": "url",
  "operator": "contains",
  "value": ["../", "etc/passwd", "cmd.exe", "eval("],
  "match_count": 1,
  "timeframe_minutes": 1
}
```

#### Anomaly Rule
Mendeteksi deviation dari baseline.
```json
{
  "field": "total_requests",
  "operator": "greater_than",
  "value": "baseline * 3",
  "baseline_window_hours": 24,
  "timeframe_minutes": 15
}
```

#### Compound Rule
Kombinasi beberapa rule.
```json
{
  "operator": "AND",
  "rules": ["RULE_001", "RULE_002"],
  "timeframe_minutes": 10
}
```

### 4.3 Certificate Renewal Agent

Agent Go yang berjalan sebagai daemon untuk otomasi perpanjangan sertifikat SSL/TLS menggunakan certbot/ACME.

**Tech Stack:** Go 1.26.4, PostgreSQL (direct), godotenv, os/exec (certbot)

**Alur Kerja (setiap 24 jam):**
1. Query sertifikat yang akan expired dari `ict_certificate` (dalam 30 hari)
2. Untuk setiap sertifikat, jalankan `certbot renew --cert-name {name}`
3. Log hasil renewal ke `ict_audit_trail` (entity_type=certificate)
4. Jika gagal → insert alert ke `ict_notification`

**PostgreSQL Tables:**

| Tabel | Operasi |
|-------|---------|
| `ict_certificate` | SELECT, UPDATE (expiry_date) |
| `ict_audit_trail` | INSERT |
| `ict_notification` | INSERT (pada error) |

### 4.4 Service Request Notification Agent

Agent Go yang berjalan sebagai daemon untuk notifikasi otomatis pada Service Request.

**Tech Stack:** Go 1.26.4, PostgreSQL (direct), godotenv, SMTP client

**Alur Kerja (setiap 60 detik):**
1. Query SR dengan status berubah dalam 1 menit terakhir dari `ict_sr_request`
2. Tentukan notifikasi berdasarkan status change
3. Kirim email/notifikasi ke user yang terkait
4. Log ke `ict_audit_trail` (entity_type=sr_request)

**Trigger Notifikasi:**

| Status Change | Penerima | Channel |
|--------------|----------|---------|
| submitted → pending_approval | Approver | Email |
| pending_approval → approved | Requester + Assigned | Email |
| approved → in_progress | Requester | Email |
| SLA approaching (80%) | Requester + Manager | Email |
| SLA breached | Requester + Manager + Director | Email + Telegram |
| in_progress → fulfilled | Requester | Email |

**PostgreSQL Tables:**

| Tabel | Operasi |
|-------|---------|
| `ict_sr_request` | SELECT (status changes) |
| `ict_audit_trail` | INSERT |
| `ict_notification` | INSERT |

### 4.5 Preventive Maintenance Reminder Agent

Agent Go yang berjalan sebagai daemon untuk pengingat jadwal pengecekan preventive maintenance.

**Tech Stack:** Go 1.26.4, PostgreSQL (direct), godotenv, SMTP client

**Alur Kerja (setiap 1 jam):**
1. Query jadwal PM yang upcoming dalam 1 hari dari `ict_pm_schedule`
2. Query checklist yang overdue dari `ict_pm_checklist`
3. Kirim reminder/alert ke assigned_to
4. Generate weekly summary setiap hari Senin

**Trigger Notifikasi:**

| Trigger | Penerima | Channel |
|---------|----------|---------|
| Schedule upcoming (1 hari sebelum) | Assigned_to | Email |
| Checklist overdue (hari ini) | Assigned_to + Manager | Email |
| Checklist overdue (>3 hari) | Manager + Director | Email + Telegram |
| Weekly summary (Senin 08:00) | Manager | Email |

**PostgreSQL Tables:**

| Tabel | Operasi |
|-------|---------|
| `ict_pm_schedule` | SELECT (upcoming) |
| `ict_pm_checklist` | SELECT (overdue) |
| `ict_notification` | INSERT |

### 4.6 Host Monitor Agent (`ict_host_monitor`)

Agent Go yang berjalan sebagai daemon di setiap server untuk mengumpulkan metrik host (CPU, memory, disk, network, load) dan mengirimkannya ke PostgreSQL via Hub API. Menggunakan **Beszel agent** (lightweight) atau **Wazuh agent** untuk data collection.

#### Arsitektur Deployment

```
Server A                          Server B
┌─────────────────────┐          ┌─────────────────────┐
│ Beszel/Wazuh Agent  │          │ Beszel/Wazuh Agent  │
│ (system metrics)    │          │ (system metrics)    │
└─────────┬───────────┘          └─────────┬───────────┘
          │                                │
┌─────────┴───────────┐          ┌─────────┴───────────┐
│ ict_host_monitor    │          │ ict_host_monitor    │
│ (per-host, non-root)│          │ (per-host, non-root)│
└─────────┬───────────┘          └─────────┬───────────┘
          │                                │
          └──────────────┬─────────────────┘
                         ▼
           ┌──────────────────────┐
           │  PostgreSQL (Hub)    │
           └──────────────────────┘
```

**Prinsip desain:**
- **Per-host** — 1 agent instance per server
- **Non-root** — berjalan sebagai service user
- **Lightweight** — minimal resource usage
- **Heartbeat** — kirim heartbeat setiap 30 detik

#### Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Agent | Go 1.26.4 (single binary, `package main`) |
| Metrics | `github.com/shirou/gopsutil/v4` (CPU, memory, disk, network) |
| Database | PostgreSQL via Hub API (REST) |
| Config | `godotenv` (`.env` file) |

#### Dependencies

| Library | Versi | Fungsi |
|---------|-------|--------|
| `github.com/shirou/gopsutil/v4` | v4.x | System metrics collection |
| `github.com/joho/godotenv` | v1.5.1 | Loading file `.env` |

#### Alur Kerja

**Worker 1: Metrics Collector** (setiap 30 detik)
1. Collect CPU usage (per-core + overall)
2. Collect memory usage (total, used, available, swap)
3. Collect disk usage (per-mount point)
4. Collect network I/O (per-interface)
5. Collect load average (1m, 5m, 15m)
6. Collect system uptime
7. Kirim metrics ke Hub API

**Worker 2: Heartbeat** (setiap 30 detik)
1. Update `ict_monitor_agent.last_heartbeat_at`
2. Update `ict_monitor_host_snapshot` dengan data terbaru

#### PostgreSQL Tables

| Tabel | Operasi | Keterangan |
|-------|---------|------------|
| `ict_monitor_agent` | UPDATE | Heartbeat + status |
| `ict_monitor_host_snapshot` | INSERT | Host metrics snapshot |
| `ict_monitor_metric_point` | INSERT | Time-series metrics |

#### Konfigurasi

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `HUB_URL` | `http://localhost:8080` | URL Hub API |
| `HUB_TOKEN` | — | Token autentikasi |
| `AGENT_PORT` | `9101` | HTTP API port (localhost) |
| `COLLECT_INTERVAL` | `30` | Detik antara collection |
| `HEARTBEAT_INTERVAL` | `30` | Detik antara heartbeat |

#### Struktur Folder

```
ict_auto/ict_host_monitor/
├── main.go
├── config.go
├── collector/
│   ├── cpu.go
│   ├── memory.go
│   ├── disk.go
│   ├── network.go
│   └── load.go
├── hub/
│   ├── client.go
│   └── sync.go
├── go.mod / go.sum
├── Dockerfile
└── .env / .env.dev / .env.pub
```

---

### 4.7 Network Poller Agent (`ict_network_poller`)

Agent Go yang berjalan sebagai daemon terpusat (*centralized*) untuk melakukan SNMP polling terhadap perangkat jaringan non-Mikrotik (switch, router, firewall generic). Berbeda dengan `ict_mikrotik_agent` yang khusus untuk Mikrotik, agent ini menangani device SNMP standar.

#### Arsitektur Deployment

```
┌─────────────────────────────────────────┐
│        OmniSight Server (1 VM)          │
│                                         │
│  ┌──────────────────────────────────┐   │
│  │ ict_network_poller               │   │
│  │ (centralized, SNMPv3)            │   │
│  └──────────────┬───────────────────┘   │
│                 │                       │
└─────────────────┼───────────────────────┘
                  │ SNMPv3 (UDP 161)
     ┌────────────┼────────────┐
     │            │            │
┌────┴────┐   ┌───┴────┐   ┌───┴────┐
│ Cisco   │   │ HP     │   │ Juniper│
│ Switch  │   │ Switch │   │ Router │
└─────────┘   └────────┘   └────────┘
```

#### Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Agent | Go 1.26.4 (single binary, `package main`) |
| SNMP | `github.com/gosnmp/gosnmp` (SNMPv3 AuthPriv) |
| Database | PostgreSQL via Hub API (REST) |
| Config | `godotenv` (`.env` file) |

#### Dependencies

| Library | Versi | Fungsi |
|---------|-------|--------|
| `github.com/gosnmp/gosnmp` | v1.4.0 | Klien SNMPv3 |
| `github.com/joho/godotenv` | v1.5.1 | Loading file `.env` |

#### Alur Kerja

**Worker 1: SNMP Poller** (setiap 60 detik)
1. Query `ict_monitor_network_device` untuk daftar device aktif
2. Untuk setiap device, SNMP walk:
   - `sysDescr`, `sysName`, `sysUpTime` → identifikasi device
   - `ifTable` → interface status + traffic counters
   - `hrProcessorLoad` → CPU usage (jika tersedia)
   - `hrStorageTable` → Memory usage (jika tersedia)
3. Insert data ke `ict_monitor_netif_snapshot`
4. Update device status di `ict_monitor_network_device`

**Worker 2: SNMP Trap Receiver** (port 162, UDP)
1. Listen untuk SNMP traps
2. Parse trap data
3. Insert ke `ict_monitor_snmp_trap`
4. Jika trap kritis → insert alert ke `ict_siem_alert`

#### PostgreSQL Tables

| Tabel | Operasi | Keterangan |
|-------|---------|------------|
| `ict_monitor_network_device` | SELECT, UPDATE | Device list + status |
| `ict_monitor_netif_snapshot` | INSERT | Interface metrics |
| `ict_monitor_snmp_trap` | INSERT | SNMP traps |

#### Konfigurasi

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `HUB_URL` | `http://localhost:8080` | URL Hub API |
| `HUB_TOKEN` | — | Token autentikasi |
| `POLL_INTERVAL` | `60` | Detik antara SNMP poll |
| `TRAP_PORT` | `162` | UDP port untuk SNMP trap |
| `SNMP_COMMUNITY` | — | Default community (v2c fallback) |

#### Struktur Folder

```
ict_auto/ict_network_poller/
├── main.go
├── config.go
├── poller/
│   ├── snmp.go
│   ├── interface.go
│   └── system.go
├── trap/
│   └── receiver.go
├── hub/
│   ├── client.go
│   └── sync.go
├── go.mod / go.sum
├── Dockerfile
└── .env / .env.dev / .env.pub
```

---

### 4.8 FIM Agent (`ict_fim_agent`)

Agent Go yang berjalan sebagai daemon di setiap server untuk File Integrity Monitoring — mendeteksi perubahan file kritis (permission, ownership, hash) pada file-system yang dipantau.

#### Arsitektur Deployment

```
Server A                          Server B
┌─────────────────────┐          ┌─────────────────────┐
│ ict_fim_agent       │          │ ict_fim_agent       │
│ (per-host, root)    │          │ (per-host, root)    │
│                     │          │                     │
│ /etc/               │          │ /etc/               │
│ /usr/bin/           │          │ /usr/bin/           │
│ /var/www/           │          │ /opt/               │
└─────────┬───────────┘          └─────────┬───────────┘
          │                                │
          └──────────────┬─────────────────┘
                         ▼
           ┌──────────────────────┐
           │  PostgreSQL (Hub)    │
           └──────────────────────┘
```

**Catatan:** Agent ini berjalan sebagai **root** karena butuh akses ke semua file untuk scanning. Ini adalah pengecualian dari aturan non-root.

#### Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Agent | Go 1.26.4 (single binary, `package main`) |
| Hashing | `crypto/sha256` (built-in) |
| File Watch | `github.com/fsnotify/fsnotify` (real-time) |
| Database | PostgreSQL via Hub API (REST) |
| Config | `godotenv` (`.env` file) |

#### Dependencies

| Library | Versi | Fungsi |
|---------|-------|--------|
| `github.com/fsnotify/fsnotify` | v1.8.x | Real-time file system events |
| `github.com/joho/godotenv` | v1.5.1 | Loading file `.env` |

#### Alur Kerja

**Worker 1: Baseline Scanner** (setiap 24 jam / on-demand)
1. Scan direktori yang dipantau (dari config)
2. Hitung SHA-256 untuk setiap file
3. Simpan baseline ke `ict_monitor_fim_baseline`
4. Compare dengan baseline sebelumnya
5. Jika ada perubahan → insert ke `ict_monitor_fim_change`

**Worker 2: Real-time Watcher** (continuous)
1. Watch file system events (create, modify, delete, chmod, chown)
2. Untuk setiap event, hash file yang berubah
3. Compare dengan baseline
4. Jika berbeda → insert ke `ict_monitor_fim_change`
5. Kirim alert ke Hub API

#### PostgreSQL Tables

| Tabel | Operasi | Keterangan |
|-------|---------|------------|
| `ict_monitor_fim_baseline` | UPSERT | File hash baseline |
| `ict_monitor_fim_change` | INSERT | Detected changes |
| `ict_monitor_fim_alert` | INSERT | FIM alerts |

#### Konfigurasi

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `HUB_URL` | `http://localhost:8080` | URL Hub API |
| `HUB_TOKEN` | — | Token autentikasi |
| `WATCH_PATHS` | `/etc,/usr/bin` | Direktori yang dipantau (comma-separated) |
| `SCAN_INTERVAL` | `86400` | Detik antara full scan |
| `EXCLUDE_PATTERNS` | `*.log,*.tmp` | Pattern yang dikecualikan |

#### Struktur Folder

```
ict_auto/ict_fim_agent/
├── main.go
├── config.go
├── scanner/
│   ├── baseline.go
│   ├── hasher.go
│   └── comparator.go
├── watcher/
│   └── fsnotify.go
├── hub/
│   ├── client.go
│   └── sync.go
├── go.mod / go.sum
├── Dockerfile
└── .env / .env.dev / .env.pub
```

---

### 4.9 UptimeRobot Sync Agent (`ict_uptimerobot_agent`)

Agent Go yang berjalan sebagai daemon terpusat (*centralized*) untuk menarik data monitoring dari UptimeRobot API dan menyimpannya ke PostgreSQL.

#### Arsitektur Deployment

```
┌─────────────────────────────────────────┐
│        OmniSight Server (1 VM)          │
│                                         │
│  ┌──────────────────────────────────┐   │
│  │ ict_uptimerobot_agent            │   │
│  │ (centralized, HTTPS API)         │   │
│  └──────────────┬───────────────────┘   │
│                 │                       │
└─────────────────┼───────────────────────┘
                  │ HTTPS (api.uptimerobot.com)
                  ▼
         ┌─────────────────┐
         │  UptimeRobot    │
         │  API v2         │
         └─────────────────┘
```

#### Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Agent | Go 1.26.4 (single binary, `package main`) |
| HTTP Client | `net/http` (bawaan Go) |
| Database | PostgreSQL via Hub API (REST) |
| Config | `godotenv` (`.env` file) |

#### Dependencies

| Library | Versi | Fungsi |
|---------|-------|--------|
| `github.com/joho/godotenv` | v1.5.1 | Loading file `.env` |

#### Alur Kerja

**Worker 1: Monitor Sync** (setiap 5 menit)
1. GET `api.uptimerobot.com/v2/getMonitors` → daftar semua monitor
2. Upsert ke `ict_uptimerobot_log` (data monitor state)

**Worker 2: Incident Sync** (setiap 1 menit)
1. GET `api.uptimerobot.com/v2/getAlerts` → alert terbaru
2. Filter alert baru (since last sync)
3. Upsert ke `ict_uptimerobot_log` (incident data)
4. Update `ict_uptimerobot_sum` (daily summary)

**Worker 3: SLA Calculator** (setiap 1 jam)
1. Hitung SLA harian dari `ict_uptimerobot_sum`
2. Upsert ke `ict_uptimerobot_sla`

#### PostgreSQL Tables

| Tabel | Operasi | Keterangan |
|-------|---------|------------|
| `ict_uptimerobot_log` | UPSERT | Log incident/alert |
| `ict_uptimerobot_sum` | UPSERT | Ringkasan harian per monitor |
| `ict_uptimerobot_sla` | UPSERT | SLA harian global |

#### Konfigurasi

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `HUB_URL` | `http://localhost:8080` | URL Hub API |
| `HUB_TOKEN` | — | Token autentikasi |
| `UPTIMEROBOT_API_KEY` | — | API key UptimeRobot |
| `MONITOR_SYNC_INTERVAL` | `300` | Detik antara monitor sync |
| `INCIDENT_SYNC_INTERVAL` | `60` | Detik antara incident sync |
| `SLA_CALC_INTERVAL` | `3600` | Detik antara SLA calculation |

#### Struktur Folder

```
ict_auto/ict_uptimerobot_agent/
├── main.go
├── config.go
├── syncer/
│   ├── monitor.go
│   ├── incident.go
│   └── sla.go
├── hub/
│   ├── client.go
│   └── sync.go
├── go.mod / go.sum
├── Dockerfile
└── .env / .env.dev / .env.pub
```

---

### 4.10 SCA Agent (`ict_sca_agent`)

Agent Go yang berjalan sebagai daemon di setiap server untuk Security Configuration Assessment — audit konfigurasi keamanan server terhadap baseline/standar (CIS Benchmark, hardening guide).

#### Arsitektur Deployment

```
Server A                          Server B
┌─────────────────────┐          ┌─────────────────────┐
│ ict_sca_agent       │          │ ict_sca_agent       │
│ (per-host, root)    │          │ (per-host, root)    │
│                     │          │                     │
│ Scan:               │          │ Scan:               │
│ - SSH config        │          │ - SSH config        │
│ - Firewall rules    │          │ - Firewall rules    │
│ - User/group        │          │ - User/group        │
│ - File permission   │          │ - File permission   │
│ - Service status    │          │ - Service status    │
└─────────┬───────────┘          └─────────┬───────────┘
          │                                │
          └──────────────┬─────────────────┘
                         ▼
           ┌──────────────────────┐
           │  PostgreSQL (Hub)    │
           └──────────────────────┘
```

**Catatan:** Agent ini berjalan sebagai **root** karena butuh akses untuk membaca konfigurasi sistem.

#### Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Agent | Go 1.26.4 (single binary, `package main`) |
| Config Parser | Custom parsers (SSH, sysctl, passwd, shadow) |
| Database | PostgreSQL via Hub API (REST) |
| Config | `godotenv` (`.env` file) |

#### Dependencies

| Library | Versi | Fungsi |
|---------|-------|--------|
| `github.com/joho/godotenv` | v1.5.1 | Loading file `.env` |

#### Alur Kerja

**Worker 1: Policy Scanner** (setiap 24 jam)
1. Load SCA policies dari `ict_monitor_sca_policy`
2. Untuk setiap policy, eksekusi checks:
   - Baca file konfigurasi
   - Parse isi file
   - Evaluate rules terhadap isi
   - Tentukan pass/fail
3. Insert hasil ke `ict_monitor_sca_result`
4. Hitung compliance score
5. Jika score turun >10% → alert

**Worker 2: Check Executor** (on-demand via API)
1. Terima request check tertentu dari Hub
2. Eksekusi check spesifik
3. Return result

#### PostgreSQL Tables

| Tabel | Operasi | Keterangan |
|-------|---------|------------|
| `ict_monitor_sca_policy` | SELECT | SCA policies |
| `ict_monitor_sca_result` | INSERT | Scan results |
| `ict_monitor_sca_check` | SELECT | Individual checks |

#### Konfigurasi

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `HUB_URL` | `http://localhost:8080` | URL Hub API |
| `HUB_TOKEN` | — | Token autentikasi |
| `SCAN_INTERVAL` | `86400` | Detik antara full scan |
| `POLICY_DIR` | `/etc/ict_sca/policies` | Direktori policy files |

#### Struktur Folder

```
ict_auto/ict_sca_agent/
├── main.go
├── config.go
├── scanner/
│   ├── policy.go
│   ├── ssh.go
│   ├── firewall.go
│   ├── users.go
│   └── permissions.go
├── hub/
│   ├── client.go
│   └── sync.go
├── go.mod / go.sum
├── Dockerfile
└── .env / .env.dev / .env.pub
```

---

### 4.11 Mikrotik Device Sync Agent (`ict_mikrotik_agent`)

Agent Go yang berjalan sebagai daemon terpusat (*centralized*) untuk menyinkronkan data dari perangkat Mikrotik (RouterOS) ke PostgreSQL. Menggunakan **SNMPv3** untuk data monitoring (read-only) dan **RouterOS API v3** (TCP 8729 TLS) untuk data konfigurasi (CRUD).

#### Arsitektur Deployment

```
┌─────────────────────────────────────────────────┐
│            OmniSight Server (1 VM)              │
│                                                 │
│  ┌──────────┐ ┌──────────┐  ┌────────────────┐  │
│  │ ict_rest │ │ ict_site │  │ ict_mikrotik_  │  │
│  │          │ │          │  │ agent          │  │
│  └──────────┘ └──────────┘  └───────┬────────┘  │
│              ┌─────────────────┐    │           │
│              │  PostgreSQL     │    │           │
│              └─────────────────┘    │           │
└─────────────────────────────────────┼───────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    │                 │                 │
              ┌─────┴──────┐   ┌──────┴──────┐   ┌──────┴──────┐
              │ VPN/Tunnel │   │ VPN/Tunnel  │   │ VPN/Tunnel  │
              │ ke Site A  │   │ ke Site B   │   │ ke Site C   │
              └─────┬──────┘   └──────┬──────┘   └──────┬──────┘
                    │                 │                 │
              ┌─────┴──────┐   ┌──────┴──────┐   ┌──────┴──────┐
              │ Mikrotik   │   │ Mikrotik    │   │ Mikrotik    │
              │ Site A     │   │ Site B      │   │ Site C      │
              └────────────┘   └─────────────┘   └─────────────┘
```

**Prinsip desain:**
- **Centralized** — 1 agent instance di server utama, menjangkau semua Mikrotik via VPN/Tunnel
- **Hybrid protocol** — SNMPv3 untuk monitoring, API untuk konfigurasi
- **Database-first** — device list diambil dari `ict_mikrotik_device`, bukan dari env var
- **TLS required** — koneksi API wajib TLS (port 8729), tidak plain text

#### Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Agent | Go 1.26.4 (single binary, `package main`) |
| Database | PostgreSQL (via `database/sql` + `lib/pq`) |
| Mikrotik API | `github.com/go-routeros/routeros/v3` (TCP 8729 TLS) |
| SNMP | `github.com/gosnmp/gosnmp` (UDP 161, SNMPv3) |
| Config | `godotenv` (`.env` file) |

#### Dependencies

| Library | Versi | Fungsi |
|---------|-------|--------|
| `github.com/go-routeros/routeros/v3` | v3.0.2 | Klien RouterOS API v3 (TLS) |
| `github.com/gosnmp/gosnmp` | v1.4.0 | Klien SNMPv3 |
| `github.com/joho/godotenv` | v1.5.1 | Loading file `.env` |
| `github.com/lib/pq` | v1.12.3 | Driver PostgreSQL |

#### Alur Kerja

Agent menjalankan **dua worker loop** secara paralel:

**Worker 1: SNMP Collector** (setiap 60 detik)
1. Baca daftar device aktif dari `ict_mikrotik_device` (status = `connected`)
2. Untuk setiap device, kirim SNMP walk:
   - `sysDescr`, `sysName`, `sysUpTime` → update `ict_mikrotik_device`
   - `hrProcessorLoad` → CPU usage
   - `hrStorageUsed` / `hrStorageSize` → Memory usage
   - `ifInOctets` / `ifOutOctets` → Traffic counters
3. Update `status`, `last_seen_at`, `uptime` di database
4. Jika SNMP gagal 3x berturut-turut → set `status = disconnected`

**Worker 2: API Syncer** (setiap 5 menit)
1. Baca daftar device aktif dari `ict_mikrotik_device`
2. Untuk setiap device, hubungi via RouterOS API (TLS):
   - `/ip/hotspot/user/print` → sinkron ke `ict_mikrotik_hotspot_user`
   - `/ip/hotspot/user/profile/print` → sinkron ke `ict_mikrotik_hotspot_profile`
   - `/ip/dhcp-server/lease/print` → sinkron ke `ict_mikrotik_dhcp_lease`
   - `/ip/firewall/filter/print` + `/ip/firewall/nat/print` → sinkron ke `ict_mikrotik_firewall_rule`
   - `/queue/simple/print` → sinkron ke `ict_mikrotik_queue`
   - `/interface/print` → identifikasi port
   - `/vlan/print` → sinkron ke `ict_mikrotik_vlan_group`
3. Gunakan **upsert** (INSERT ON CONFLICT UPDATE) untuk semua data
4. Mark device `status = connected` + update `last_seen_at`
5. Jika API gagal → mark `status = error`, log error, lanjut device berikutnya

#### PostgreSQL Tables

| Tabel | Operasi | Oleh Worker |
|-------|---------|-------------|
| `ict_mikrotik_device` | UPDATE (status, uptime, version) | SNMP + API |
| `ict_mikrotik_hotspot_profile` | UPSERT | API |
| `ict_mikrotik_hotspot_user` | UPSERT | API |
| `ict_mikrotik_dhcp_lease` | UPSERT | API |
| `ict_mikrotik_firewall_rule` | UPSERT | API |
| `ict_mikrotik_queue` | UPSERT | API |
| `ict_mikrotik_vlan_group` | UPSERT | API |
| `ict_mikrotik_port_group` | UPSERT | API |
| `ict_mikrotik_client` | UPSERT | API |
| `ict_mikrotik_access_normalization` | UPSERT | API |
| `ict_mikrotik_access_list` | UPSERT | API |

#### In-Memory Cache

Cache di-refresh setiap 5 menit dari database:

| Data | Sumber | Fungsi |
|------|--------|--------|
| Device List | `ict_mikrotik_device` | Daftar device aktif untuk polling |
| SNMP Credentials | `ict_mikrotik_device` | SNMPv3 auth/priv keys per device |

#### Konfigurasi (Environment Variable)

| Variabel | Default | Keterangan |
|----------|---------|------------|
| `PG_HOST` | `localhost` | Host PostgreSQL |
| `PG_PORT` | `5432` | Port PostgreSQL |
| `PG_USER` | `dbe` | User PostgreSQL |
| `PG_PASS` | `rahasia` | Password PostgreSQL |
| `PG_DATA` | `erp` | Nama database |
| `IS_POOL` | `false` | Gunakan connection pool |
| `SNMP_INTERVAL` | `60` | Interval SNMP polling (detik) |
| `API_INTERVAL` | `300` | Interval API sync (detik) |
| `API_TLS` | `true` | Wajib TLS untuk RouterOS API |

#### Struktur Folder

```
ict_auto/ict_mikrotik_agent/
├── main.go              ← Entry point + scheduler (2 worker loop)
├── config.go            ← Konfigurasi dari env + device list loader
├── collector/
│   ├── snmp.go          ← SNMPv3 data collection
│   └── api.go           ← RouterOS API data collection
├── syncer/
│   ├── device.go        ← Sync device info → ict_mikrotik_device
│   ├── hotspot.go       ← Sync hotspot profiles + users
│   ├── dhcp.go          ← Sync DHCP leases
│   ├── firewall.go      ← Sync firewall rules (filter + NAT)
│   ├── queue.go         ← Sync bandwidth queues
│   ├── port.go          ← Sync port groups
│   ├── vlan.go          ← Sync VLAN groups
│   ├── client.go        ← Sync clients
│   └── access.go        ← Sync access normalization + access list
├── go.mod / go.sum
├── Dockerfile
└── .env / .env.dev / .env.pub
```

#### Dockerfile

Multi-stage build:
- **Stage 1 (builder):** `golang:1.26.4-alpine` — kompilasi statik dengan `CGO_ENABLED=0 GOOS=linux`
- **Stage 2 (runtime):** `alpine:3.20` — minimal image dengan `ca-certificates` dan `tzdata`

#### Deployment

```bash
# Docker
docker run -d \
  --name ict-mikrotik-agent \
  -e PG_HOST=erp_pool \
  -e PG_PORT=5432 \
  -e PG_USER=dbe \
  -e PG_PASS=rahasia \
  -e PG_DATA=erp \
  -e IS_POOL=true \
  -e SNMP_INTERVAL=60 \
  -e API_INTERVAL=300 \
  -e API_TLS=true \
  --network=omnisight \
  ict_auto/ict_mikrotik_agent
```

#### Alur Data: Mikrotik → Agent → Database → Frontend

```
Mikrotik Device (RouterOS)
  │
  ├── SNMPv3 (UDP 161) ──→ ict_mikrotik_agent
  │   sysName, uptime,       │
  │   CPU, memory, traffic   │
  │                          │
  └── API TLS (TCP 8729) ───→│
      hotspot, DHCP,         │
      firewall, queue        │
                             ▼
                    PostgreSQL (ict_mikrotik_*)
                             │
                             ▼
                    ict_rest API (MK01-MK10)
                             │
                             ▼
                    ict_site Frontend
                    (DataTable + Detail views)
```

#### Security Considerations

| Aspek | Implementasi |
|-------|-------------|
| API Connection | TLS required (port 8729), plain text rejected |
| SNMP | SNMPv3 dengan AuthPriv (SHA + AES), tidak SNMPv1/v2c |
| Credentials | Disimpan di `ict_mikrotik_device` (encrypted dengan AES-GCM) |
| Network | Agent hanya bisa diakses dari internal network |
| Logging | Semua koneksi成功/gagal di-log dengan Zerolog |

---

## Ringkasan Agent

| Agent | Fungsi Utama | Agent Type | Status | Prioritas |
|-------|-------------|-----------|--------|-----------|
| `ict_log_nginx` | Log sync + WAF + threat classification + auto-ban + SLA | siem | ✅ Sudah ada | — |
| `ict_log_rotate` | Log archive + retention + JSONL export | siem | ✅ Sudah ada | — |
| `ict_mikrotik_agent` | Mikrotik device sync (SNMP + API) → PostgreSQL | network_poller | ⬜ Direncanakan | Tinggi |
| `ict_docker_agent` | Docker monitoring + container management + compose deploy | docker_collector | ⬜ Direncanakan | Tinggi |
| `ict_host_monitor` | Host metrics collection (CPU, memory, disk, network) | host_monitor | ⬜ Direncanakan | Tinggi |
| `ict_network_poller` | SNMP polling untuk device non-Mikrotik | network_poller | ⬜ Direncanakan | Tinggi |
| `ict_siem_evaluator` | SIEM rule evaluation | siem | ⬜ Direncanakan | Sedang |
| `ict_siem_notifier` | SIEM alert notification | siem | ⬜ Direncanakan | Sedang |
| `ict_fim_agent` | File Integrity Monitoring | fim | ⬜ Direncanakan | Sedang |
| `ict_uptimerobot_agent` | UptimeRobot data sync → PostgreSQL | combined | ⬜ Direncanakan | Sedang |
| `ict_sca_agent` | Security Configuration Assessment | combined | ⬜ Direncanakan | Sedang |
| Certificate Renewal | Auto-renewal SSL/TLS | combined | ⬜ Direncanakan | Rendah |
| SR Notification | Service Request notifikasi | combined | ⬜ Direncanakan | Rendah |
| PM Reminder | Preventive maintenance reminder | combined | ⬜ Direncanakan | Rendah |

### Catatan Arsitektur

**Agent yang sudah ada** (`ict_log_nginx`, `ict_log_rotate`):
- Tidak menggunakan framework CLI — binary langsung dengan `main()` yang menjalankan worker loop
- Tidak menjalankan HTTP server — murni background workers
- Tidak menggunakan framework `ict_rest` — berdiri sendiri menggunakan `database/sql` mentah

**Agent per-host** (`ict_docker_agent`, `ict_host_monitor`, `ict_fim_agent`, `ict_sca_agent`):
- Berjalan di setiap server/host yang perlu dimonitor
- Docker agent berjalan sebagai **non-root** (user `omnisight:docker`) dengan Docker socket di-mount
- FIM dan SCA agent berjalan sebagai **root** (butuh akses file system lengkap)
- Menggunakan **Go SDK** untuk komunikasi dengan sistem lokal

**Agent centralized** (`ict_mikrotik_agent`, `ict_network_poller`, `ict_uptimerobot_agent`):
- Berjalan di server utama OmniSight
- Menjangkau devices via VPN/Tunnel atau API eksternal
- Menggunakan **SNMPv3** untuk polling device jaringan
- Menggunakan **HTTPS API** untuk layanan eksternal (UptimeRobot)

**Agent notification** (`ict_siem_evaluator`, `ict_siem_notifier`, Certificate Renewal, SR Notification, PM Reminder):
- Berjalan sebagai daemon di server utama
- Mengirim notifikasi via email, webhook, atau Telegram
- Beberapa berjalan sebagai **root** (Certificate Renewal butuh akses certbot)

**Semua agent:**
- Menggunakan multi-stage Docker build (golang:alpine → alpine)
- Menggunakan pola file `.env` / `.env.dev` / `.env.pub` untuk konfigurasi
- Berdiri sendiri (*standalone*) tanpa shared library
