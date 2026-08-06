# Rencana Task File-by-File Paket 1 Fase 2

Dokumen ini memecah Paket 1 Fase 2 menjadi task file-by-file untuk klaster Asset dan Account Inventory bastion.

## Scope Paket 1

- Modul domain: jms_asset_group, jms_asset, jms_asset_group_member, jms_account, jms_account_secret.
- Tujuan: inventory asset dan account bastion siap dipakai oleh session, approval, dan connect token.

## Batch Task A - Finalisasi Schema JMS

### File Target
- `obx_base/prisma/schema/jsm_stack.prisma`
- `obx_base/prisma/schema/schema.prisma`

### Task
1. Review field wajib untuk asset, account, dan secret.
2. Review index yang diperlukan untuk query list dan filter.
3. Pastikan enum/properti status konsisten untuk pilot bastion.
4. Verifikasi relasi asset-account-secret.

### Verifikasi
- `npx prisma generate` di `obx_base` lulus.

## Batch Task B - API Inventory Bastion

### File Target
- `obx_rest/backbone/routes.go`
- `obx_rest/skeleton/{KODE_BASTION_ASSET}/template.go`
- `obx_rest/skeleton/{KODE_BASTION_ASSET}/repository.go`
- `obx_rest/skeleton/{KODE_BASTION_ASSET}/usecase.go`
- `obx_rest/skeleton/{KODE_BASTION_ASSET}/handler.go`

### Task
1. Definisikan DTO request/response untuk asset dan account.
2. Implement list/count/create/update/delete asset dengan scope company.
3. Implement list/count/create/update/delete account.
4. Batasi expose secret hanya pada jalur yang aman.
5. Pastikan write routes memakai logging yang konsisten.

### Verifikasi
- `go build ./...` di `obx_rest` lulus.
- Smoke test endpoint inventory bastion lulus.

## Batch Task C - Frontend Board Bastion Inventory

### File Target
- `obx_site/src/app/board/pages/{KODE_BASTION_ASSET}/page.tsx`
- `obx_site/src/app/board/model/module.ts`

### Task
1. Buat halaman board inventory bastion berbasis DataTable + DataDialog.
2. Pastikan form validasi memakai zod dan menampilkan FieldError.
3. Sediakan status active/inactive dan metadata minimum.
4. Sinkronkan menu board untuk modul inventory bastion.

### Verifikasi
- `npm run lint` di `obx_site` lulus.
- Halaman board inventory dapat list/create/update data pilot.

## Batch Task D - Sinkronisasi Dokumentasi

### File Target
- `obx_docs/blueprint/BASE/jsm_stack.md`
- `obx_docs/guide/BASE/jsm_stack.md`
- `obx_docs/blueprint/PLAN/README.md`

### Task
1. Sinkronkan penjelasan domain inventory bastion dengan schema aktual.
2. Tambahkan alur penggunaan minimal untuk asset dan account.
3. Pastikan catatan batasan secret reference tertulis jelas.
4. Tautkan ke roadmap dan backlog Fase 2.

### Verifikasi
- Dokumen BASE dan PLAN konsisten dengan route serta payload API.

## Urutan Eksekusi yang Disarankan

1. Batch A - Schema
2. Batch B - API
3. Batch C - Frontend
4. Batch D - Dokumentasi

## Risiko

1. Schema inventory terlalu longgar untuk policy access.
2. Secret terekspos pada endpoint list/detail.
3. Frontend inventory tidak sinkron dengan route backend.

## Mitigasi

1. Tetapkan field minimum dan status wajib di awal.
2. Gunakan jalur secret terpisah dengan kontrol akses ketat.
3. Lakukan smoke test API dan board dalam batch yang sama.

## Referensi

- phase2_backend_coding_plan.md
- phase2_backend_execution_backlog.md
- phase2_execution_backlog.md
- ai_runbook.md
