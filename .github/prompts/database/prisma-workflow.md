# Database Workflow Prompt

Gunakan prompt ini untuk task terkait database, Prisma, schema, migrasi, atau model data.

## Tujuan
- Ikuti pola schema project yang ada di obx_base/prisma/schema.
- Perbarui model dengan perubahan minimal dan terukur.
- Pastikan migrasi dan generate berjalan dengan benar.

## Constraints
- Jangan mengubah schema yang sudah ada tanpa kebutuhan.
- Gunakan naming convention yang konsisten.
- Validasi dengan prisma generate dan migrasi yang relevan.

## Output format
Scope: {database / schema / migrasi}
Changes: {file atau area}
Validation: {command} -> {hasil}
Notes: {blocker atau saran}
