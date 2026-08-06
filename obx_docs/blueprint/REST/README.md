# Rancangan Submodul: REST - obx_rest Backend API

## Ringkasan

- Scope: backend service API berbasis Go + Gin
- Tujuan: autentikasi, session, master data, dan robot endpoint
- Port default: 36665

## Struktur Dokumen

- Summary teknis: file ini
- Detail page teknis: `obx_docs/blueprint/REST/*.md`
- Summary user guide: `obx_docs/guide/REST.md`
- Detail page user guide: `obx_docs/guide/REST/*.md`

## Inti Backend

- Entry point: `obx_rest/main.go`
- Router: `obx_rest/backbone/routes.go`
- Auth/middleware: `obx_rest/backbone/memory.go`
- DB connector: `obx_rest/backbone/database.go`
- Shared error helper: `obx_rest/mechanic/helper.go`

## Route Domain

- `/rest/guest`: guest/auth flow
- `/rest/pages`: user pages dan admin pages
- `/rest/robot`: robot/agent flow

Catatan route aktual:
- Endpoint `SP01/privilege` aktif pada group `authu`.
- Skeleton `XX99` tersedia di source code, namun route `XX99` belum diregistrasi di `backbone/routes.go`.

## Pattern Summary

- Setiap modul mengikuti pola Repo -> UseCase -> Handler
- Handler memakai `mechanic.Error(c, err)`
- UseCase wrap error internal dengan `mechanic.InternalError`
- Write route memakai `USLogs("KODE")`
- Route static didaftarkan sebelum dynamic route

## Page Index

| Page | Detail Teknis | Detail Guide |
|---|---|---|
| SP00 | SP00.md | SP00.md |
| SP01 | SP01.md | SP01.md |
| SP02 | SP02.md | SP02.md |
| SP03 | SP03.md | SP03.md |
| SM01 | SM01.md | SM01.md |
| SM02 | SM02.md | SM02.md |
| SM03 | SM03.md | SM03.md |
| SM04 | SM04.md | SM04.md |
| SM05 | SM05.md | SM05.md |
| XX99 (template, route belum aktif) | XX99.md | XX99.md |

## Status Agent

- Middleware USRobot aktif
- Group `/rest/robot` aktif
- Endpoint turunan `/rest/robot/*` belum diregistrasi

## Checklist

- [ ] Endpoint terdaftar di route group yang benar
- [ ] Handler memakai mechanic.Error
- [ ] UseCase wrap error internal
- [ ] Write route memakai USLogs("KODE")
- [ ] Query parameterized
- [ ] Build lulus: go build ./...

## Referensi

- `obx_rest/main.go`
- `obx_rest/backbone/routes.go`
- `obx_rest/backbone/memory.go`
- `obx_rest/backbone/database.go`
- `obx_rest/mechanic/helper.go`
