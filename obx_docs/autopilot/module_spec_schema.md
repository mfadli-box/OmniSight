# Module Spec Schema

Dokumen ini mendefinisikan format spesifikasi modul yang harus diisi sebelum AI membuat modul.

## Tujuan

- Menyamakan input antara product, developer, dan AI.
- Mengurangi ambiguity pada route, tabel, dan privilege.
- Menjadi sumber data untuk generate backend, frontend, test, dan docs.

## Format Wajib

Gunakan YAML dengan struktur berikut:

```yaml
module_code: "MD01"
module_name: "Module Display Name"
module_type: "list" # list | link | form | workflow
owner_submodule: "REST" # AUTO | BASE | REST | SITE

business_goal: "Satu kalimat tujuan modul"

backend:
  route_group: "/rest/pages"
  route_base: "/rest/pages/MD01"
  handler_code: "MD01"
  endpoints:
    - method: "GET"
      path: "/rest/pages/MD01"
      purpose: "List data"
    - method: "POST"
      path: "/rest/pages/MD01"
      purpose: "Create data"

database:
  primary_table: "dat_module_sample"
  prisma_schema_file: "obx_base/prisma/schema/dat_module.prisma"
  fields:
    - name: "id"
      type: "string"
      required: true
      key: "pk"
    - name: "name"
      type: "string"
      required: true
      key: "normal"
  indexes:
    - "name"
  relations:
    - table: "dat_company"
      local_key: "company_id"
      target_key: "id"

frontend:
  page_code: "MD01"
  page_path: "/board/MD01"
  ui_pattern: "DataTable+DataDialog"
  fields:
    - name: "name"
      component: "Input"
      required: true
  mobile_behavior: "DataTableCard"

authorization:
  roles_allowed: ["admin", "manager"]
  privilege_code: "MD01"

logging:
  write_log_code: "MD01"

test_minimum:
  backend_cases: ["list_ok", "create_ok", "validation_error", "unauthorized"]
  frontend_cases: ["render_list", "open_dialog", "submit_success", "submit_error"]

documentation:
  blueprint_file: "obx_docs/blueprint/REST/MD01.md"
  guide_file: "obx_docs/guide/REST/MD01.md"
```

## Aturan Validasi Spec

- module_code harus unik dan uppercase.
- route_base harus konsisten dengan page_code.
- primary_table harus ada di schema Prisma.
- privilege_code harus sama dengan write_log_code kecuali ada alasan khusus.
- blueprint_file dan guide_file wajib diisi sebelum implementasi final.

## Output Turunan dari Spec

Spec ini dipakai untuk menghasilkan:

1. Skeleton endpoint backend.
2. Struktur halaman frontend.
3. Daftar test minimal.
4. Draft dokumen blueprint dan guide.