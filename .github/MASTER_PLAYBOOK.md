# Master Playbook Copilot CLI untuk OmniSight

Dokumen ini adalah ringkasan utama untuk memakai Copilot CLI di workspace ini agar perilakunya lebih konsisten dengan Opencode.

## 1. Tujuan
Gunakan workflow ini untuk kerja yang:
- terarah
- hemat token
- konsisten dengan pola repo
- mudah diverifikasi

## 2. Instruksi utama
- Prioritaskan perubahan kecil dan terfokus.
- Gunakan pola yang sudah ada di repo.
- Untuk backend, ikuti pola skeleton di obx_rest/skeleton/{KODE}/.
- Untuk frontend, ikuti pola DataTable + DataDialog di obx_site.
- Untuk database, ikuti pola schema di obx_base/prisma/schema.
- Untuk agent, ikuti pola di obx_auto.
- Selalu verifikasi dengan command yang relevan.

## 3. Mode kerja
Pilih salah satu mode berikut sesuai kebutuhan:
- backend
- frontend
- database
- agent
- debug
- token-efficient

## 4. Format input yang disarankan
```text
Task: {jelaskan task singkat}
Mode: {backend | frontend | database | agent | debug | token-efficient}
Scope: {sempit / satu file / satu modul}
Constraints:
- keep changes minimal
- follow existing pattern
- validate after changes
```

## 5. Format output yang disarankan
```text
Scope: {singkat}
Changes: {file atau area}
Documentation: {teknis: obx_docs/blueprint/{SUBMODUL}/README.md, user guide: obx_docs/guide/{SUBMODUL}.md, atau alasan tidak perlu}
Validation: {command} -> {hasil}
Notes: {blocker atau saran}
```

## 6. Guardrails
- Jangan menambah fitur yang tidak diminta.
- Jangan melebar dari scope yang sudah ditentukan.
- Periksa referensi repo sebelum edit.
- Validasi hasil setelah perubahan.
- Jika task memengaruhi dokumen yang dirujuk, lengkapi dokumentasi yang relevan sebagai bagian dari pekerjaan (teknis di obx_docs/blueprint/{SUBMODUL}/README.md dan user guide di obx_docs/guide/{SUBMODUL}.md).
- Jika ambigu, pilih solusi minimal dan jelaskan asumsi.

## 7. Referensi file
- .copilot-instructions.md
- .github/instructions/copilot-workflow.instructions.md
- .github/instructions/database-agent.instructions.md
- .github/prompts/modes/*.md
- .github/prompts/examples/*.md
- .github/HELPER.md
- .github/TECHNICAL.md
