# Copilot Workflow Prompt

Gunakan prompt ini saat ingin menjalankan workflow yang konsisten dengan Opencode.

## Tujuan
- Selesaikan task dengan scope sempit.
- Prioritaskan hasil yang bisa diverifikasi.
- Gunakan pola repo yang sudah ada.

## Mode kerja
- Jika task adalah refactor atau cleanup, gunakan pendekatan clean-execution.
- Jika task adalah modul baru, gunakan module-creator.
- Jika task adalah satu halaman frontend, gunakan frontend-page-creator.
- Jika task adalah satu halaman backend, gunakan backend-page-creator.
- Jika user menginginkan hemat token, gunakan token-efficient.

## Output format
Scope: {singkat}
Changes: {file atau area}
Documentation: {teknis: obx_docs/blueprint/{SUBMODUL}/README.md, user guide: obx_docs/guide/{SUBMODUL}.md, atau alasan tidak perlu}
Validation: {command} -> {hasil}
Notes: {blocker atau saran}
