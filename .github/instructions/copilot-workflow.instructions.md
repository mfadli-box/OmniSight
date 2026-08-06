# Copilot CLI Workspace Instructions

Gunakan instruksi ini untuk menjaga workflow tetap konsisten dengan project OmniSight.

## Prinsip utama
- Prioritaskan perubahan kecil, terfokus, dan mudah diverifikasi.
- Gunakan pola yang sudah ada di repo, jangan membuat solusi baru tanpa kebutuhan.
- Untuk task backend, ikuti pola skeleton di obx_rest/skeleton/{KODE}/.
- Untuk task frontend, ikuti pola DataTable + DataDialog di obx_site.
- Selalu verifikasi dengan command yang relevan sebelum menyatakan selesai.
- Terapkan prinsip `token-efficient` di setiap workflow: output ringkas, fokus, dan hindari verbose yang tidak perlu.
- Saat skill dipakai, lengkapi dokumentasi yang relevan yang dirujuk (utama: obx_docs/blueprint/{SUBMODUL}/README.md untuk teknis dan obx_docs/guide/{SUBMODUL}.md untuk user guide) jika task memengaruhi pola, workflow, atau struktur yang dijelaskan di sana.

## Skill / workflow yang disarankan
- `clean-execution`: untuk refactor, cleanup, atau implementasi umum.
- `module-creator`: untuk modul baru atau halaman modul baru.
- `frontend-page-creator`: untuk satu halaman frontend.
- `backend-page-creator`: untuk satu halaman backend.
- `token-efficient`: untuk hemat token, output singkat, dan scope sempit.

## Output format
Gunakan format singkat berikut saat melaporkan hasil:

Scope: {singkat}
Changes: {file atau area}
Validation: {command} -> {hasil}
Notes: {blocker atau saran}
