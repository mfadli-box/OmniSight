# Daftar Perubahan Implementasi Pertama Paket 1 Fase 2

Dokumen ini adalah jembatan dari perencanaan ke eksekusi coding untuk modul inventory bastion pada Paket 1 Fase 2.

## Tujuan

1. Menentukan urutan perubahan kode pertama yang paling aman.
2. Menjaga perubahan kecil, terukur, dan mudah diverifikasi.
3. Menyediakan acceptance check per langkah.

## Prasyarat

1. Kode modul bastion final sudah ditetapkan untuk placeholder KODE_BASTION_ASSET.
2. Schema JMS terbaru sudah sinkron.
3. Dokumen peta file nyata Paket 1 Fase 2 telah ditinjau.

## Langkah Implementasi Pertama

### Langkah 1 - Finalisasi Kontrak Data di Schema

#### File
- obx_base/prisma/schema/jsm_stack.prisma
- obx_base/prisma/schema/schema.prisma

#### Perubahan Minimum
1. Konfirmasi field wajib untuk jms_asset_group, jms_asset, jms_account, jms_account_secret.
2. Konfirmasi index untuk company_id, asset_id, username, status.
3. Konfirmasi relasi asset ke account dan account ke secret.

#### Acceptance Check
- npx prisma generate di obx_base lulus.

### Langkah 2 - Definisi DTO dan Interface Backend

#### File
- obx_rest/skeleton/{KODE_BASTION_ASSET}/template.go

#### Perubahan Minimum
1. Tambahkan DTO list/create/update untuk asset group, asset, account.
2. Tambahkan interface repository dan usecase sesuai kebutuhan CRUD awal.
3. Pastikan nama field payload konsisten dengan schema.

#### Acceptance Check
- Go compiler mengenali tipe DTO dan interface tanpa error.

### Langkah 3 - Implement Repository Layer

#### File
- obx_rest/skeleton/{KODE_BASTION_ASSET}/repository.go

#### Perubahan Minimum
1. Implement Count dan List berbasis company scope.
2. Implement Create, Update, Delete untuk asset dan account.
3. Pastikan query parameterized dan whitelist sort field.
4. Pastikan secret tidak ikut di payload list atau detail publik.

#### Acceptance Check
- go build ./... di obx_rest lulus.

### Langkah 4 - Implement Usecase Layer

#### File
- obx_rest/skeleton/{KODE_BASTION_ASSET}/usecase.go

#### Perubahan Minimum
1. Validasi input wajib sebelum operasi repository.
2. Wrap error repository dengan InternalError yang konsisten.
3. Standarkan pesan validasi untuk field wajib.

#### Acceptance Check
- Unit flow list/create/update/delete berjalan tanpa panic.

### Langkah 5 - Implement Handler Layer dan Route

#### File
- obx_rest/skeleton/{KODE_BASTION_ASSET}/handler.go
- obx_rest/backbone/routes.go

#### Perubahan Minimum
1. Tambahkan endpoint list/create/update/delete.
2. Pastikan bind query dan bind JSON tervalidasi.
3. Pastikan write route memakai logging middleware yang konsisten.

#### Acceptance Check
- go build ./... di obx_rest lulus.
- Smoke test API CRUD lulus.

### Langkah 6 - Implement Halaman Board

#### File
- obx_site/src/app/board/pages/{KODE_BASTION_ASSET}/page.tsx
- obx_site/src/app/board/model/module.ts

#### Perubahan Minimum
1. Tambahkan DataTable untuk list asset dan account.
2. Tambahkan DataDialog untuk create dan update.
3. Tambahkan validasi zod dan error display.
4. Sinkronkan menu modul di board.

#### Acceptance Check
- npm run lint di obx_site lulus.
- UI dapat list, create, update data pilot.

### Langkah 7 - Sinkronisasi Dokumentasi

#### File
- obx_docs/blueprint/BASE/jsm_stack.md
- obx_docs/guide/BASE/jsm_stack.md
- obx_docs/blueprint/PLAN/README.md

#### Perubahan Minimum
1. Catat endpoint awal yang sudah aktif.
2. Catat batasan secret exposure.
3. Catat dependency ke Paket 2 Fase 2.

#### Acceptance Check
- Dokumen BASE dan PLAN sinkron dengan implementasi aktual.

## Catatan Risiko

1. Scope awal melebar ke fitur session atau proxy terlalu cepat.
2. Secret handling belum dipisah jelas dari payload publik.
3. Route dan menu tidak sinkron saat modul pertama aktif.

## Mitigasi

1. Kunci scope pada inventory bastion dahulu.
2. Audit payload response setiap endpoint.
3. Lakukan smoke test API dan UI pada batch yang sama.

## Referensi

- phase2_package1_file_task_plan.md
- phase2_package1_actual_file_map.md
- phase2_package1_module_actual_file_map.md
- phase2_backend_coding_plan.md
- ai_runbook.md
