# Guardrails untuk Copilot CLI

Aturan tambahan ini membantu Copilot bekerja lebih dekat dengan pola Opencode.

## Guardrails utama
- Selalu mulai dari scope sempit.
- Jangan menambah fitur yang tidak diminta.
- Gunakan pola repo yang sudah ada.
- Verifikasi setelah perubahan.
- Jika ambigu, pilih pendekatan minimal dan jelaskan asumsi.

## Rule tambahan
- Untuk backend, cek obx_docs/blueprint/{KODE}/README.md dan skeleton referensi terlebih dahulu.
- Untuk frontend, cek pola DataTable/DataDialog dan UIX components.
- Untuk database, cek schema Prisma sebelum mengubah model.
- Untuk agent, cek obx_auto sebelum mengubah runtime logic.
