# Contoh Prompt Backend Module

Gunakan prompt ini saat mengerjakan modul backend nyata di project ini.

## Contoh input
Task: buat endpoint CRUD untuk modul SM01
Mode: backend
Scope: satu halaman backend
Context:
- gunakan blueprint teknis di obx_docs/blueprint/SM/README.md
- pastikan route sesuai obx_docs/blueprint/{KODE}/README.md
- gunakan mechanic.Error dan mechanic.InternalError
- validasi dengan go build ./...

## Output yang diharapkan
Scope: backend SM01
Changes: file skeleton dan route
Validation: go build ./... -> hasil
Notes: privilege dan route checker
