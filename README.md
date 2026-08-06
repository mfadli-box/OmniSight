# OmniSight

Monorepo OmniSight untuk operasional DevSecOps dengan tiga komponen aplikasi utama:
- obx_base: skema database PostgreSQL dengan Prisma
- obx_rest: backend API Go (Gin)
- obx_site: frontend Next.js

## Struktur Root

| Path | Keterangan |
|---|---|
| obx_base/ | Database schema, migrasi, seed |
| obx_rest/ | API backend dan middleware |
| obx_site/ | Aplikasi web frontend |
| obx_auto/ | Agent otomasi |
| obx_docs/ | Pusat dokumentasi blueprint, guide, dan autopilot |
| docker-compose.yml | Compose produksi |
| docker-compose.dev.yml | Compose pengembangan |

## Status Modul Saat Ini (Workspace)

### Backend (obx_rest/skeleton)
- SP00
- SP01
- SP02
- SP03
- SM01
- SM02
- SM03
- SM04
- SM05
- XX99

### Frontend (obx_site/src/app/board/pages)
- SP01
- SP02
- SP03
- SM01
- SM02
- SM03
- SM04
- SM05
- XX99

## Dokumentasi Acuan

- Indeks dokumentasi teknis: obx_docs/blueprint/README.md
- Rancangan aktif (draft): obx_docs/blueprint/PLAN/README.md
- Indeks dokumen autopilot: obx_docs/autopilot/README.md
- Roadmap rancangan JMS MVP: obx_docs/autopilot/roadmap_JMS_MVP.md

## Aturan Penting

- Bahasa kode aplikasi: Inggris.
- Bahasa dokumentasi: Indonesia.
- Folder terlindungi yang dilarang diakses AI agent: obx_docs/config/

## Menjalankan Layanan Lokal

### Database (obx_base)

```bash
cd obx_base
npx prisma generate
npx prisma migrate dev --name init
npx prisma db seed
```

### Backend (obx_rest)

```bash
cd obx_rest
go run main.go
```

### Frontend (obx_site)

```bash
cd obx_site
npm run dev
```

## Catatan

Untuk detail arsitektur dan pola implementasi per submodul, selalu rujuk dokumen di obx_docs terlebih dahulu sebelum mengubah kode.

## Update Terakhir

- Tanggal: 6 Agustus 2026
- Oleh: opencode/big-pickle
- Keterangan: Penambahan paket prompt cepat autopilot, varian no-smoke, tracker progres harian, dan runner PowerShell untuk eksekusi batch Opencode CLI.
