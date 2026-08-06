# User Guide REST - Backend API obx_rest

## Ringkasan

obx_rest adalah layanan backend utama untuk autentikasi, manajemen user/module/company, signature flow, dan session management.

## Prasyarat

- PostgreSQL sudah aktif dan dapat diakses
- File .env obx_rest berisi koneksi DB valid
- KEY_ROBOT tersedia jika endpoint robot dipakai
- Go modules sudah sinkron (go mod tidy)

## Konfigurasi Dasar

Env minimal:
- PG_HOST
- PG_PORT
- PG_USER
- PG_PASS
- PG_DATA

Env opsional:
- IS_POOL=true untuk koneksi pool lebih besar
- KEY_ROBOT untuk autentikasi route /rest/robot

Catatan status saat ini:
- Route group `/rest/robot` sudah ada.
- Endpoint turunan `/rest/robot/*` belum diregistrasi pada backend saat ini.
- Endpoint `GET /rest/pages/SP01/privilege` sudah aktif.
- Modul `XX99` masih template; route belum diregistrasi pada `backbone/routes.go`.

## Menjalankan Service

Dari folder obx_rest:

```sh
go mod tidy
go run main.go
```

Service akan listen di:
- http://localhost:36665

## Endpoint Cek Cepat

1. Health root:
- GET /
- expected: status 200

2. REST marker:
- GET /rest
- expected: {"message":"rest"}

3. Guest marker:
- GET /rest/guest/
- expected: {"message":"guest"}

4. Pages marker (butuh Authorization token valid):
- GET /rest/pages/
- expected: {"message":"pages"}

## Alur Login Dasar (SP00)

1. Ambil company HRIS:
- GET /rest/guest/SP00

2. Login:
- POST /rest/guest/SP00
- Simpan token Authorization dari response.

3. Pakai token untuk endpoint /rest/pages/*.

4. Logout:
- DELETE /rest/pages/SP00

## Scope Endpoint per Role

- guest:
  - login flow (SP00)
- auth user:
  - context company/module (SP01)
  - privilege current user (SP01/privilege)
  - ganti password (SP02)
  - lihat action log sendiri (SP03)
- admin:
  - CRUD user/module/company/signature/session (SM01-SM05)

## Detail Guide Pages

- SP00: obx_docs/guide/REST/SP00.md
- SP01: obx_docs/guide/REST/SP01.md
- SP02: obx_docs/guide/REST/SP02.md
- SP03: obx_docs/guide/REST/SP03.md
- SM01: obx_docs/guide/REST/SM01.md
- SM02: obx_docs/guide/REST/SM02.md
- SM03: obx_docs/guide/REST/SM03.md
- SM04: obx_docs/guide/REST/SM04.md
- SM05: obx_docs/guide/REST/SM05.md
- XX99: obx_docs/guide/REST/XX99.md (template autopilot, route belum aktif)

## Header yang Umum Dipakai

- Authorization: Bearer <token>
- X-Company-ID: company context aktif (opsional, divalidasi USLoad)
- Content-Type: application/json

## Troubleshooting

| Gejala | Penyebab Umum | Solusi |
|---|---|---|
| 401 Token required | Header Authorization kosong | Tambahkan Authorization token |
| 401 Session expired/invalid | Token tidak valid atau expired | Login ulang melalui SP00 |
| 403 User account deactivated | is_active=false di dat_user | Aktifkan user melalui SM01 |
| 403 Admin access only | Endpoint admin diakses non-admin | Gunakan user admin |
| 403 No company selected | Company context kosong | Set default company user atau kirim X-Company-ID valid |
| 500 Internal error | Error query/usecase | Cek log server berdasarkan request_id |
| Startup gagal connect DB | Env PG_* salah | Verifikasi PG_HOST/PORT/USER/PASS/DATA |

## Monitoring Operasional

- Cek log server untuk request_id dan error trace.
- Audit write actions tersimpan di dat_user_action melalui USLogs.
- Pantau sesi user aktif di SM05.

## Aktivasi Endpoint Agent (Saat Implementasi Ditambahkan)

1. Registrasikan endpoint di group `agent` pada `backbone/routes.go`.
2. Jalankan backend dengan:
```sh
go run main.go
```
3. Uji endpoint robot dengan header:
- Authorization: Bearer <KEY_ROBOT>
4. Jika response `401 Invalid robot token`, cek env KEY_ROBOT pada runtime backend.

## Praktik Aman

- Jangan expose KEY_ROBOT ke client frontend.
- Gunakan HTTPS di environment production.
- Batasi CORS origins sesuai domain produksi.
- Rotasi token sesi secara berkala.

## Validasi Setelah Perubahan Backend

1. Build check:
```sh
go build ./...
```

2. Smoke test endpoint kritikal:
- /rest
- /rest/guest/SP00
- /rest/pages/SP01/company
- /rest/pages/SP01/privilege
- /rest/pages/SM05

3. Verifikasi audit action untuk write route.
