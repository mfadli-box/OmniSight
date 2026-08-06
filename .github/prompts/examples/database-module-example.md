# Contoh Prompt Database Module

Gunakan prompt ini saat mengerjakan perubahan schema atau migrasi.

## Contoh input
Task: tambah kolom baru ke schema Prisma untuk modul tertentu
Mode: database
Scope: satu model atau satu schema file
Context:
- gunakan pola di obx_base/prisma/schema
- pastikan perubahan minimal
- validasi dengan prisma generate

## Output yang diharapkan
Scope: database schema update
Changes: schema file dan migrasi terkait
Validation: npx prisma generate -> hasil
Notes: dampak terhadap model dan migrasi
