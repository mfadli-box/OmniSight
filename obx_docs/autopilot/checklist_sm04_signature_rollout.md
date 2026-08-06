# Checklist Implementasi SM04 - Rollout Signature Type SOP/ISO

## Tujuan

Dokumen ini menjadi panduan eksekusi untuk merealisasikan signature type usulan pada modul SM04 agar dapat dipakai pada workflow persetujuan dokumen fase 1.

## Scope Rollout

- Modul utama: SM04 (Signature Type)
- Integrasi pendukung: SP03 (audit action)
- Domain dokumen: policy, SOP, work instruction, form/record

## Daftar Signature Type Target

| Code | Nama | Kategori Dokumen | Versi |
|---|---|---|---|
| SIG_POLICY_V1 | Signature Policy V1 | POLICY | V1 |
| SIG_SOP_V1 | Signature SOP V1 | SOP | V1 |
| SIG_WI_V1 | Signature Work Instruction V1 | WORK_INSTRUCTION | V1 |
| SIG_FORM_V1 | Signature Form Record V1 | FORM_RECORD | V1 |

## Prasyarat

- Akses admin aktif ke halaman Signature Type.
- User signer sudah ada dan aktif.
- Dokumen baseline tersedia:
  - sample_iso_policy_bastion_access.md
  - sample_iso_sop_request_approval.md
  - sample_approval_matrix_sop_iso.md

## Tahap A - Persiapan Data Signer

1. Ambil daftar user signer dari endpoint list signer.
2. Verifikasi role minimal untuk signer:
   - Reviewer role tersedia.
   - Approver role tersedia.
3. Buat daftar kandidat signer per kategori dokumen.
4. Tetapkan fallback signer untuk kondisi cuti/tidak tersedia.

## Tahap B - Pembuatan Signature Type di SM04

Lakukan untuk setiap code pada tabel target.

1. Buat master signature type baru.
2. Isi metadata wajib:
   - code
   - name
   - is_active=true
3. Tambahkan approval step minimal:
   - Step 1: Reviewer
   - Step 2: Approver
4. Tambahkan signer per step (minimal 1 user per step).
5. Simpan konfigurasi.

## Tahap C - Validasi Fungsional

1. Buka kembali detail signature type.
2. Pastikan step dan signer tersimpan.
3. Lakukan update kecil (misal nama display) untuk memastikan endpoint PUT berjalan.
4. Pastikan status aktif tampil di list.
5. Ulangi untuk semua signature type target.

## Tahap D - Uji Integrasi Workflow Dokumen

1. Pilih 1 dokumen sample per kategori.
2. Kaitkan dokumen ke signature type yang sesuai.
3. Simulasikan alur:
   - submit
   - review
   - approve/reject
4. Pastikan perubahan status terdokumentasi di histori.
5. Verifikasi jejak audit tersedia untuk setiap aksi.

## Tahap E - Kontrol Audit dan Evidence

Kumpulkan bukti berikut untuk setiap signature type:

- Screenshot list signature type aktif.
- Screenshot detail step + signer.
- Catatan hasil simulasi review/approve/reject.
- Log audit action (SP03) untuk aksi create/update.
- Ringkasan gap/temuan jika ada kegagalan.

## Definition of Done

- 4 signature type target berhasil dibuat dan aktif.
- Setiap signature type punya minimal 2 step (reviewer, approver).
- Setiap step punya signer aktif.
- Uji integrasi dokumen lulus untuk 4 kategori.
- Evidence checklist lengkap dan tersimpan.

## Risiko dan Mitigasi

| Risiko | Dampak | Mitigasi |
|---|---|---|
| Signer tidak tersedia | Approval macet | Tetapkan fallback signer per step |
| Mapping kategori salah | Workflow tidak sesuai | Review silang dengan approval matrix |
| Step tersimpan tidak lengkap | Audit mismatch | Wajib cek detail pasca-save |
| Audit tidak tercatat | Kehilangan traceability | Verifikasi SP03 setelah setiap uji |

## Runbook Eksekusi Singkat

1. Siapkan kandidat signer.
2. Buat 4 signature type target.
3. Validasi tiap record.
4. Jalankan simulasi workflow dokumen.
5. Simpan evidence dan hasil akhir rollout.

## Riwayat Revisi

| Versi | Tanggal | Perubahan | Oleh |
|---|---|---|---|
| 1.0 | 06-08-2026 | Draft awal checklist rollout SM04 | opencode/GPT-5.3-Codex |
