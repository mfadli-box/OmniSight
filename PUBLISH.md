# Publish dan Deployment OmniSight

Dokumen ini adalah panduan ringkas untuk publish dan deploy OmniSight sesuai kondisi repository saat ini.

## Stack Runtime

- Database: PostgreSQL 18 (obx_base)
- Backend API: Go + Gin (obx_rest)
- Frontend: Next.js (obx_site)
- Agent: obx_auto

## Berkas Compose

| Berkas | Fungsi |
|---|---|
| docker-compose.yml | Lingkungan produksi (utama) |
| docker-compose.dev.yml | Override pengembangan (opsional) |

## Prasyarat

1. Docker dan Docker Compose tersedia.
2. File .env root terisi.
3. File obx_base/.env terisi untuk Prisma.

## Deploy Produksi

```bash
docker-compose up -d --build
```

Inisialisasi database setelah container aktif:

```bash
cd obx_base
npx prisma generate
npx prisma migrate dev --name init
npx prisma db seed
cd ..
```

Restart stack setelah inisialisasi:

```bash
docker-compose down
docker-compose up -d --build
```

## Deploy Pengembangan

```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build
```

Inisialisasi database:

```bash
cd obx_base
npx prisma generate
npx prisma migrate dev --name init
npx prisma db seed
cd ..
```

Restart stack dev:

```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml down
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build
```

## Verifikasi Minimum

```bash
# API health
curl http://localhost:36665/rest

# Frontend
# buka http://localhost:36666 di browser
```

## Referensi Dokumentasi

- Dokumen teknis utama: obx_docs/blueprint/README.md
- Rancangan aktif (draft): obx_docs/blueprint/PLAN/README.md
- Runbook autopilot: obx_docs/autopilot/README.md

## Catatan Keamanan

- Jangan menyimpan secret produksi di dokumentasi.
- Pastikan akses ke obx_docs/config/ dibatasi.
- Untuk route write backend, gunakan logging sesuai standar project.

## Update Terakhir

- Tanggal: 5 Agustus 2026
- Oleh: Owner
- Keterangan: Sinkronisasi panduan publish dengan dokumentasi aktif saat ini.
