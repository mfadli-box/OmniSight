# dat_company - Company Domain

## Ringkasan

- File schema: `obx_base/prisma/schema/dat_company.prisma`
- Isi utama: company, module mapping, area mapping, request mapping

## Model Inti

- dat_company
- dat_company_module
- dat_company_area
- dat_request (relasi dari company)

## Detail Field Kunci

### dat_company

- id (uuid, primary key)
- code (unique global)
- name
- vat_id (optional)
- reg_no (optional)
- address (optional)
- valuta (default: IDR)
- hris_link (optional)
- is_active
- created_at
- updated_at

### Relasi dat_company

- users -> dat_user_company[]
- modules -> dat_company_module[]
- areas -> dat_company_area[]
- doc_requests -> dat_request[]

## Constraint dan Index Penting

- dat_company.code menggunakan unique global
- dat_company_module menggunakan unique gabungan [company_id, module_id]
- dat_company_area menggunakan unique gabungan [company_id, code]
- dat_company_area memiliki index pada company_id

## Catatan Teknis

- Pola tenancy banyak memakai company_id
- Relasi company menjadi basis akses user dan module
- Relasi request kini merujuk ke model `dat_request`

## Checklist

- [ ] Relasi company sinkron
- [ ] company_id tetap konsisten
- [ ] relasi ke dat_request tervalidasi
