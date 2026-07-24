# 🚀 OmniSight

Integrated documentation for the architecture management, development, and deployment of all ICT services.

---

## 📌 Publish Stack

Follow these steps sequentially to provision and deploy the stack into production.

### Step 1: Create Custom Isolated Network
```bash
docker network create --attachable --driver=bridge --subnet=172.99.66.0/24 --gateway=172.99.66.254 blackbox
```

### Step 2: Clone Repository & Provision Directories
```bash
git clone https://github.com/mfadli-box/OmniSight.git
cd OmniSight
git pull origin main

# Setup user mapping for pgbouncer
vi ict_base/pgbouncer.txt

"dbe" "rahasia"
```

### Step 3: Populate Global Environment Variables
```bash
vi .env
```
```env
TZ=Asia/Jakarta
PG_HOST=ict_base
PG_POOL=ict_pool
PG_PORT=5432
PG_DATA=ict
PG_BACK=ict_backup
PG_USER=dbe
PG_PASS=rahasia
IS_POOL=true
AD_MAIL=admin@localhost
AD_PASS=rahasia
BE_POOL=http://ict_rest:36665
ES_LINK=http://elasticsearch:9200
WS_NAME=OmniSight
WS_META=OmniSight
WS_DESC=OmniSight
FT_HTTP=150
BG_TIME=1
RE_PATH=archive
RE_NORMAL=7
RE_ATTACK=90
RE_SECRET=rahasia
```
```bash
vi ict_base/.env
```
```env
DATABASE_URL="postgresql://dbe:rahasia@ict_base:5432/ict?sslmode=disable&schema=public"
AD_MAIL="admin@localhost"
AD_PASS="rahasia"
```

### Step 4: Pull Latest Core Artifacts & Trigger Production Build
```bash
git pull origin main
docker-compose up -d --build

# Initialize Database Schema & Client Generation
cd ict_base
npx prisma generate
npx prisma migrate dev --name init
npx prisma db seed

# Restart core stack to pick up initial migrations
cd ..
docker-compose down
docker-compose up -d


api http://localhost:36665
web http://localhost:36666
```

### For PgAdmin
```
volumes:
  pgadmin:
    driver: local

services:
  ict_tool:
    image: dpage/pgadmin4:latest
    container_name: ict_tool
    logging:
      options:
        max-size: "10m"
        max-file: "3"
    environment:
      - PGADMIN_DEFAULT_EMAIL=${AD_MAIL:-admin@localhost}
      - PGADMIN_DEFAULT_PASSWORD=${AD_PASS:-rahasia}
      - TZ=${TZ:-Asia/Jakarta}
    restart: always
    volumes:
      - pgadmin:/var/lib/pgadmin
    ports:
      - 36663:80
    networks:
      blackbox:
        ipv4_address: 172.99.66.3
    depends_on:
      - ict_pool
```

### For Optimize postgresql Change docker-compose
```
services:
  ict_base:
    image: postgres:18-alpine
    container_name: ict_base
    shm_size: 10gb
    deploy:
      resources:
        limits:
          cpus: '14.0'
          memory: 10g
    command: >
      postgres 
      -c shared_buffers=8GB 
      -c effective_cache_size=8GB 
      -c work_mem=16MB 
      -c maintenance_work_mem=2GB 
      -c max_connections=150 
      -c checkpoint_completion_target=0.9 
      -c wal_buffers=16MB 
      -c random_page_cost=1.1 
      -c effective_io_concurrency=200
```
