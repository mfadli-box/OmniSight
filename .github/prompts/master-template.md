# Master Prompt Template

Gunakan prompt ini sebagai template utama saat bekerja dengan Copilot CLI di workspace ini.

## Mode kerja
Pilih salah satu mode sesuai kebutuhan:
- backend: buat atau perbaiki halaman backend di obx_rest/skeleton/{KODE}/
- frontend: buat atau perbaiki halaman frontend di obx_site/src/app/board/pages/{KODE}/
- database: perbaiki schema, prisma, atau migrasi di obx_base/prisma/schema/
- agent: perbaiki atau tambahkan logic automation di obx_auto/
- debug: cari akar masalah dan perbaiki bug kecil
- token-efficient: hemat token, singkat, dan fokus

## Format input
Task: {jelaskan task singkat}
Mode: {backend | frontend | database | agent | debug | token-efficient}
Scope: {sempit / satu file / satu modul}
Constraints:
- keep changes minimal
- follow existing pattern
- validate after changes
- if the task affects referenced documentation, update the relevant docs as part of the work

## Output format
Scope: {singkat}
Changes: {file atau area}
Validation: {command} -> {hasil}
Notes: {blocker atau saran}
