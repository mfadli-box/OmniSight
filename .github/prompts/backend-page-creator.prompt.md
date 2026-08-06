---
mode: agent
---

Buat atau rapikan satu halaman backend di obx_rest/skeleton/{KODE}/ dengan pola template → repository → usecase → handler.

Constraints:
- follow backend rules in obx_docs/blueprint/{KODE}/README.md
- use mechanic.Error and mechanic.InternalError properly
- keep route and privilege aligned with obx_docs/blueprint/{KODE}/README.md
- update related docs in obx_docs/blueprint/{KODE}/README.md and obx_docs/guide/{KODE}.md when behavior changes
- validate with go build ./...

Format output:
Scope: {halaman backend}
Changes: {file yang berubah}
Validation: {command} -> {hasil}
Notes: {blocker atau saran}
