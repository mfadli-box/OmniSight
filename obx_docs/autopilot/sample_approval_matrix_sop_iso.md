# Approval Matrix - SOP/ISO (Draft Awal)

## Tujuan

Mendefinisikan role approval dokumen terkendali dan pemetaan awal ke signature type pada modul SM04.

## Matrix

| Document Type | Creator | Reviewer | Approver | Publisher | Archive Owner |
|---|---|---|---|---|---|
| POLICY | Security Lead | Infrastructure Lead | IT Manager | Document Controller | Quality Assurance |
| SOP | Process Owner | Team Lead terkait | IT Manager | Document Controller | Quality Assurance |
| WORK_INSTRUCTION | Engineer/Operator | Process Owner | Team Lead | Document Controller | Quality Assurance |
| FORM_RECORD | Operator | Process Owner | Team Lead | Document Controller | Quality Assurance |

## Aturan Approval

- Minimal jumlah reviewer: 1
- Minimal approver: 1
- SLA review: 1 hari kerja
- SLA approval: 1 hari kerja setelah review selesai
- Kondisi reject: data tidak lengkap, risiko tidak terkontrol, atau scope tidak jelas
- Kondisi urgent change: boleh fast-track dengan approval IT Manager + Security Lead

## Mapping ke Sistem

- Signature type code:
  - SIG_POLICY_V1 (usulan)
  - SIG_SOP_V1 (usulan)
  - SIG_WI_V1 (usulan)
  - SIG_FORM_V1 (usulan)
- Module code terkait:
  - SM04 untuk konfigurasi signature type dan signer
  - SP03 untuk audit action pengguna
  - Domain request/document untuk workflow persetujuan
- Kategori dokumen:
  - POLICY
  - SOP
  - WORK_INSTRUCTION
  - FORM_RECORD

## Catatan Implementasi

- Kode signature di atas adalah usulan awal dan harus dibuat di modul SM04 sebelum dipakai di workflow produksi.
- Matrix ini harus ditautkan ke template dokumen dan evidence checklist pada fase 1.
- Perubahan matrix wajib melalui approval yang sama agar governance konsisten.
