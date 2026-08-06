# dat_document - Document Domain

## Ringkasan

- File schema: `obx_base/prisma/schema/dat_document.prisma`
- Isi utama: category, document, version, approval, request

## Model Inti

- dat_document_category
- dat_document
- dat_document_version
- dat_document_approval
- dat_request

## Field Kunci dat_request

- id
- company_id
- type_id
- requester_id
- code
- title
- description
- priority
- status
- current_step
- completion_note
- completed_by
- completed_at
- created_at
- updated_at

## Relasi Utama dat_request

- company_id -> dat_company.id
- type_id -> dat_signature_type.id
- requester_id -> dat_user.id (RequestCreator)
- completed_by -> dat_user.id (RequestCompleter)

## Constraint dat_request

- unique gabungan [company_id, code]
- index [company_id, status]
- index [requester_id]
- index [type_id]

## Catatan Teknis

- Dipakai untuk document lifecycle, approval, dan request flow.
- Alur request (`dat_request`) dan alur dokumen berbagi scope company yang sama.

## Checklist

- [ ] Versioning dokumen sinkron
- [ ] Request flow terdokumentasi bersama approval flow
