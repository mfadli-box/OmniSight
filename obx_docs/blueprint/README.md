# Blueprint Dokumentasi Teknis

Folder ini berisi dokumentasi teknis per submodul dengan pola wajib:

- obx_docs/blueprint/{SUBMODUL}/README.md

## Tujuan

- Menjadi sumber acuan teknis implementasi backend dan frontend.
- Menyamakan route, privilege, struktur data, dan checklist implementasi.
- Menjadi referensi utama untuk skill Copilot/OpenCode.

## Struktur Standar

Setiap submodul minimal memiliki bagian berikut:

1. Overview
2. Database Schema
3. API Endpoints
4. Frontend Pages
5. Implementation Checklist
6. Relationships
7. Backend Pattern Reference

## Daftar Submodul

| Submodul | File | Status |
|---|---|---|
| AUTO | ../blueprint/AUTO/README.md | active |
| BASE | ../blueprint/BASE/README.md | active |
| PLAN | ../blueprint/PLAN/README.md | draft |
| PLAN gap analysis | ../blueprint/PLAN/gap_analysis_platform_replacement.md | draft |
| REST | ../blueprint/REST/README.md | active |
| SITE | ../blueprint/SITE/README.md | active |

## Catatan

- Gunakan template pada obx_docs/blueprint/_template/README.md saat menambah submodul baru.
- Dokumen harus ringkas, operasional, dan konsisten dengan pola repository.
- Untuk perubahan schema Prisma, perbarui blueprint BASE terlebih dahulu sebelum blueprint submodul lain.
