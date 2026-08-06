# Database and Agent Instructions

Gunakan instruksi ini untuk task yang berkaitan dengan database, Prisma, dan automation agent.

## Database
- Ikuti pola schema di obx_base/prisma/schema.
- Perubahan model harus konsisten dan minimal.
- Pastikan migrasi dan generate terverifikasi.
- Terapkan prinsip `token-efficient`: jangan membuat penjelasan panjang untuk task yang sederhana.
- Saat task database/agent dipakai, lengkapi dokumentasi yang relevan yang dirujuk, terutama jika perubahan memengaruhi schema, workflow, atau operasi runtime.

## Agent
- Ikuti pola implementasi di obx_auto.
- Fokus pada runtime behavior, error handling, dan integrasi service.
- Hindari perubahan luas jika task bisa diselesaikan dengan patch kecil.
- Terapkan prinsip `token-efficient`: solusi singkat, verifikasi minimal, hasil fokus.
- Saat task agent dipakai, lengkapi dokumentasi yang relevan yang dirujuk apabila perilaku atau alur kerja berubah.

## Output format
Scope: {database / agent / runtime}
Changes: {file atau area}
Validation: {command} -> {hasil}
Notes: {blocker atau saran}
