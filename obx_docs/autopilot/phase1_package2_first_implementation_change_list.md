# Daftar Perubahan Implementasi Pertama Paket 2 Fase 1

Dokumen ini menjadi jembatan dari perencanaan ke eksekusi coding untuk modul workflow dokumen terkendali pada Paket 2 Fase 1.

## Prasyarat

1. Kode modul final untuk Paket 2 Fase 1 telah ditetapkan: dat_request, dat_signature, dat_document, dat_document_revision, dat_document_evidence.
2. Peta file nyata dan task plan Paket 2 Fase 1 telah ditinjau.
3. Kontrak enum approval (`approval_info`, `approval_flag`) sudah tersedia pada schema aktif.

## Langkah Implementasi Pertama

### Langkah 1 - Verifikasi Schema Dokumen dan Approval

#### File
- obx_base/prisma/schema/all_enum_dat.prisma
- obx_base/prisma/schema/dat_document.prisma
- obx_base/prisma/schema/dat_signature.prisma
- obx_base/prisma/schema/schema.prisma

#### Acceptance Check
- Relasi `dat_request` -> `dat_signature_type` valid.
- Relasi `dat_document` -> `dat_document_version` valid.
- Enum approval dapat dipakai pada flow request dan flag.

### Langkah 2 - Definisi DTO dan Interface di SM04

#### File
- obx_rest/skeleton/SM04/template.go

#### Acceptance Check
- DTO untuk request, document, revision, dan signature form tersedia.
- Interface repository/usecase selaras dengan kebutuhan route.

### Langkah 3 - Implement Repository dan Usecase SM04

#### File
- obx_rest/skeleton/SM04/repository.go
- obx_rest/skeleton/SM04/usecase.go

#### Acceptance Check
- Pattern query parameterized dipakai di semua operasi.
- Sort field di-whitelist melalui map ekspresi.
- Error repository dibungkus `mechanic.InternalError(...)` di usecase.

### Langkah 4 - Implement Handler dan Registrasi Route

#### File
- obx_rest/skeleton/SM04/handler.go
- obx_rest/backbone/routes.go

#### Acceptance Check
- Endpoint list/create/update/delete aktif untuk request, document, revision.
- Endpoint approval flow aktif untuk form generate/list dan flag action.
- Semua write route memakai `USLogs("SM04")`.

### Langkah 5 - Sinkronisasi Blueprint dan Guide REST

#### File
- obx_docs/blueprint/REST/SM04.md
- obx_docs/guide/REST/SM04.md

#### Acceptance Check
- Endpoint dan validasi sesuai implementasi backend terbaru.
- Catatan alur approval (advance step/finalisasi) terdokumentasi.

## Referensi

- phase1_package2_file_task_plan.md
- phase1_package2_actual_file_map.md
- phase1_package2_module_actual_file_map.md
- phase1_backend_coding_plan.md
- phase1_backend_execution_backlog.md
- ai_runbook.md

## Status Eksekusi Saat Ini

1. Implementasi dat_request selesai.
2. Implementasi dat_signature (form dan flag flow) selesai.
3. Implementasi dat_document selesai.
4. Implementasi dat_document_revision selesai.
5. Sinkronisasi route dan dokumentasi REST SM04 selesai.
6. Validasi backend (`go build` dan `go vet` scope SM04) dinyatakan lulus pada batch eksekusi sebelumnya.
7. Modul dat_document_evidence menjadi lanjutan berikutnya.
