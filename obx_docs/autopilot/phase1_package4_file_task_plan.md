# Rencana Task File-by-File - Paket 4 Fase 1

Dokumen ini menurunkan Paket 4 Fase 1 menjadi daftar task file-by-file yang siap dipakai untuk audit readiness dan dokumentasi operasional.

## Ruang Lingkup

- Fokus: template SOP/ISO, evidence checklist, owner/reviewer/approver, dan indeks dokumentasi
- Tujuan: memastikan Fase 1 siap untuk audit internal dasar dan handoff ke Fase 2

## Urutan Task

### 1. Template and Sample Review

#### Target File
- `obx_docs/autopilot/template_iso_policy.md`
- `obx_docs/autopilot/template_iso_sop.md`
- `obx_docs/autopilot/template_iso_work_instruction.md`
- `obx_docs/autopilot/template_iso_form_record.md`
- `obx_docs/autopilot/template_iso_evidence_checklist.md`
- `obx_docs/autopilot/template_approval_matrix.md`
- `obx_docs/autopilot/sample_iso_sop_request_approval.md`
- `obx_docs/autopilot/sample_iso_policy_bastion_access.md`
- `obx_docs/autopilot/sample_approval_matrix_sop_iso.md`

#### Task
1. Perbarui template agar konsisten dengan alur approval terkini.
2. Pastikan sample dokumen dan approval matrix mengikuti format template.
3. Verifikasi bahasa dan struktur dokumentasi tetap konsisten.

#### Output
- Template dan sample dokumen terkendali sinkron.

### 2. Index and Guide Review

#### Target File
- `obx_docs/autopilot/README.md`
- `obx_docs/blueprint/README.md`
- `obx_docs/guide/README.md`
- `obx_docs/blueprint/PLAN/README.md`

#### Task
1. Pastikan indeks dokumen memuat template, sample, dan backlog Fase 1.
2. Pastikan referensi fase dan paket tidak bertabrakan.
3. Pastikan status draft/active tetap jelas.

#### Output
- Indeks dokumentasi rapi dan konsisten.

### 3. Audit Readiness Review

#### Target File
- `obx_docs/autopilot/checklist_sm04_signature_rollout.md`
- `obx_docs/autopilot/checklist_platform_replacement_execution.md`
- `obx_docs/autopilot/phase1_execution_backlog.md`
- `obx_docs/autopilot/phase1_backend_execution_backlog.md`
- `obx_docs/autopilot/phase1_backend_coding_plan.md`

#### Task
1. Review checklist agar selaras dengan paket coding dan backend backlog.
2. Pastikan gate Fase 1 jelas dan tidak tumpang tindih dengan Fase 2.
3. Verifikasi bahwa evidence dan approval flow disebut dengan konsisten.

#### Output
- Audit readiness untuk Fase 1 stabil.

### 4. Documentation Sync for Fase 1

#### Target File
- `obx_docs/blueprint/BASE/README.md`
- `obx_docs/guide/BASE.md`
- `obx_docs/blueprint/REST/README.md`
- `obx_docs/guide/REST.md`
- `obx_docs/blueprint/SITE/README.md`
- `obx_docs/guide/SITE.md`

#### Task
1. Sinkronkan catatan update terakhir dan referensi Fase 1.
2. Pastikan dokumen inti menyebut template, sample, dan checklist yang relevan.
3. Pastikan tidak ada referensi placeholder yang tersisa.

#### Output
- Dokumentasi operasional Fase 1 lengkap.

### 5. Verification Batch

#### Command
```bash
cd obx_base
npx prisma generate
cd ../obx_rest
go build ./...
cd ../obx_site
npm run lint
```

#### Smoke Check
- Review sample SOP approval
- Review policy bastion sample
- Review approval matrix sample
- Review evidence checklist sample

#### Output
- Bukti verifikasi Paket 4 tersedia.

## Risiko

1. Template dan sample tidak konsisten satu sama lain.
2. Indeks dokumentasi tidak sinkron dengan file baru.
3. Checklist audit readiness terlalu longgar atau tumpang tindih dengan fase lain.

## Mitigasi

1. Update template dan sample dalam batch yang sama.
2. Cek indeks dokumentasi setelah setiap perubahan.
3. Pertahankan gate Fase 1 agar tetap tegas sebelum masuk Fase 2.

## Referensi

- phase1_backend_coding_plan.md
- phase1_backend_execution_backlog.md
- phase1_execution_backlog.md
- ai_runbook.md
