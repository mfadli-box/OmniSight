# Rancangan Submodul: {KODE} - {NAMA}

## 1. Overview

- Group:
- Icon:
- Akses:
- Tujuan:
- Referensi skema:

## 2. Database Schema

- Tabel utama:
- Relasi:
- Field kunci:
- Catatan migrasi:

## 3. API Endpoints

| Method | Endpoint | Handler | Keterangan |
|---|---|---|---|
| GET | /rest/pages/{KODE} | {KODE}Handler.Find | List data |
| GET | /rest/pages/{KODE}/:id | {KODE}Handler.FindByID | Detail data |
| POST | /rest/pages/{KODE} | {KODE}Handler.Create | Buat data |
| PUT | /rest/pages/{KODE}/:id | {KODE}Handler.Update | Ubah data |
| DELETE | /rest/pages/{KODE}/:id | {KODE}Handler.Delete | Hapus data |

## 4. Frontend Pages

- Halaman:
- Komponen utama:
- Field form:
- Mobile behavior:

## 5. Implementation Checklist

- [ ] Prisma schema siap
- [ ] Backend template/repository/usecase/handler selesai
- [ ] Route terdaftar
- [ ] Frontend page selesai
- [ ] Typecheck/build lulus
- [ ] User guide diperbarui

## 6. Relationships

- Relasi tabel:
- Dependensi halaman:

## 7. Backend Pattern Reference

- Error handling: mechanic.Error dan mechanic.InternalError
- Query: parameterized query dan whitelist sort
- Routing: static route sebelum dynamic route
- Logging: write route wajib USLogs("{KODE}")
