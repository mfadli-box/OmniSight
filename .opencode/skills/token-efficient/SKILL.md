---
name: token-efficient
description: Token-saving workflow for concise, focused execution. Use when the user wants lower token usage, shorter responses, or when the repo context is already clear. Best paired with clean-execution, frontend-page-creator, backend-page-creator, or module-creator.
---

# Token-Efficient Skill

Tujuan skill ini adalah menekan penggunaan token tanpa mengorbankan kualitas hasil.

## Kapan dipakai

Gunakan skill ini saat:

- user meminta hemat token atau respons singkat
- task rutin, sudah jelas, dan tidak memerlukan eksplorasi luas
- repo sudah punya pola yang jelas sehingga konteks dapat dipangkas
- Anda perlu bekerja cepat dengan scope sempit dan verifikasi minimal

Jangan gunakan sebagai skill utama saat:

- task butuh desain mendalam atau review multi-file
- user meminta penjelasan ekstensif atau dokumentasi detail
- task memerlukan investigasi akar masalah yang belum jelas

## Prinsip utama

1. Gunakan konteks yang paling relevan saja.
2. Jangan membaca file lebih banyak dari yang dibutuhkan.
3. Prioritaskan satu hipotesis, satu perubahan, satu verifikasi.
4. Hindari penjelasan berulang; langsung ke tindakan.
5. Kalau task sudah cocok dengan skill lain, pakai skill spesifik itu dan terapkan aturan ini sebagai overlay.

## Workflow ringkas

1. Identifikasi scope sebenarnya.
2. Baca file anchor yang paling dekat dengan kebutuhan.
3. Implementasi kecil dan terfokus.
4. Verifikasi cepat dengan command paling murah.
5. Laporkan hasil secara ringkas, tidak verbose.

## Format respons yang disarankan

```text
Scope: {singkat}
Changes: {file atau area}
Validation: {command} -> {hasil}
Notes: {blocker atau saran}
```

## Integrasi dengan skill lain

- Dipakai bersama `clean-execution` untuk pekerjaan umum.
- Dipakai bersama `frontend-page-creator` atau `backend-page-creator` saat halaman sudah jelas.
- Dipakai bersama `module-creator` saat desain modul sudah stabil dan Anda ingin menghindari pemborosan konteks.

## Rule khusus

- Jangan mengulang konteks yang sudah ada.
- Jangan bertanya jika informasi yang dibutuhkan sudah bisa diperkirakan dari pola repo.
- Jika ada ambiguity yang tidak mengganggu, lanjutkan dengan asumsi minimal dan sebutkan asumsi itu.
- Fokus pada hasil yang bisa diverifikasi, bukan penjelasan yang panjang.
