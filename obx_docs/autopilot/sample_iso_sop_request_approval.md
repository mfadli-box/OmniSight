# SOP - Proses Persetujuan Request Operasional

## Metadata Dokumen

- document_code: SOP-OPS-REQ-001
- document_title: Proses Persetujuan Request Operasional
- document_type: SOP
- version: 1.0
- status: DRAFT
- process_owner: IT Operation Lead
- reviewer_name: Security Lead
- approver_name: IT Manager
- effective_date: 06 Agustus 2026
- next_review_date: 06 Februari 2027

## Tujuan

Menetapkan alur baku pembuatan, review, persetujuan, penolakan, dan penyelesaian request operasional agar jejak audit lengkap dan terukur.

## Ruang Lingkup

- Request operasional internal yang melalui workflow approval.
- Modul dan data terkait pada domain request, signature, dan user action.

## Prasyarat

- Akses/privilege: user memiliki akses create request; approver memiliki akses approve/reject.
- Data yang harus tersedia: company aktif, signature type aktif, user aktif.
- Tool yang digunakan: backend REST API, frontend board, database log/audit.

## Input dan Output

| Input | Sumber | Output | Tujuan |
|---|---|---|---|
| Form request | User requester | Record request | Memulai alur approval |
| Decision approval | Approver | Status approved/rejected | Kontrol otorisasi |
| Completion note | Executor/owner | Status completed | Penutupan request |

## Langkah Prosedur

1. Requester membuat request baru dan melengkapi field wajib.
2. Sistem mengaitkan request ke signature type yang relevan.
3. Approver melakukan review lalu approve/reject dengan catatan.
4. Jika approved, owner proses mengeksekusi lalu menandai completion.
5. Sistem menyimpan histori status, actor, dan timestamp.

## Kontrol dan Checkpoint

- Checkpoint 1: request tidak boleh diproses jika data wajib tidak lengkap.
- Checkpoint 2: approval hanya oleh user yang berhak.
- Checkpoint 3: completion wajib punya completion note.

## KPI / SLA

- SLA review approval: maksimal 1 hari kerja.
- SLA completion setelah approved: maksimal 3 hari kerja.

## Eskalasi dan Penanganan Error

- Kondisi eskalasi: request kritis melewati SLA atau approver tidak merespons.
- PIC eskalasi: IT Manager.
- Jalur komunikasi: tiket internal + notifikasi grup operasional.

## Referensi Form / Evidence

- Form terkait: FORM-REQ-OPS-001.
- Evidence yang wajib disimpan: request payload, approval decision, completion note, log actor.

## Riwayat Revisi

| Versi | Tanggal | Perubahan | Oleh |
|---|---|---|---|
| 1.0 | 06-08-2026 | Draft awal SOP approval request | opencode/GPT-5.3-Codex |
