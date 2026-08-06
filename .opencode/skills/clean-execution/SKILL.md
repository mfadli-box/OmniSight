---
name: clean-execution
description: Default execution skill for disciplined, minimal, and consistent implementation or refactor work. Triggers: "rapih", "bersih", "konsisten", "refactor", "simplify", "structured", "clean workflow". If a task is specifically about module, frontend-page, or backend-page creation, prefer the dedicated skill.
---

# Clean Execution Skill

Tujuan skill ini adalah membuat workflow agent lebih rapi, konsisten, dan tidak bertele-tele.

## Referensi wajib

Gunakan konteks repo yang relevan dan pilih referensi minimum yang benar-benar dibutuhkan untuk task yang sedang dikerjakan.

## Integrasi skill lain

Skill ini berlaku sebagai aturan eksekusi dasar saat dipasangkan dengan skill lain, terutama `module-creator`.

- Jika task memakai `module-creator`, pendekatan design-first tetap WAJIB dijalankan dengan perubahan minimal, scope sempit, dan validasi bertahap.
- Jangan biarkan workflow modul melebar menjadi eksplorasi panjang yang tidak perlu.
- Saat membuat modul atau halaman baru, ikuti blueprint dan checklist proyek, tetapi tetap kerjakan per irisan kecil yang bisa divalidasi segera.
- Approval user untuk rancangan tetap wajib sebelum implementasi, dan setelah implementasi dimulai validasi harus dilakukan segera setelah edit substantif pertama.

## Kapan dipakai

Gunakan skill ini jika:

- task utamanya adalah implementasi, refactor, penyederhanaan, atau perapihan eksekusi
- user ingin perubahan tetap kecil, konsisten, dan mudah diverifikasi
- tidak ada workflow domain khusus yang lebih tepat

Jangan jadikan skill ini sebagai skill utama jika:

- user secara spesifik meminta pembuatan modul baru, halaman modul, atau sub-entity
- task membutuhkan workflow design-first dengan approval rancangan
- task secara spesifik berfokus pada pembuatan atau restrukturisasi satu halaman frontend dengan pola UI proyek yang jelas

Dalam kasus tersebut, pakai skill yang lebih spesifik dan biarkan `clean-execution` mengatur cara eksekusinya.

## Format prompt yang disarankan

Gunakan format singkat seperti ini jika ingin hasil cepat dan terarah:

```text
Scope: execution
Goal: refactor | simplify | cleanup | implement
Target: {file | folder | page | component}
Context: {masalah atau tujuan utama}
Constraints:
- keep changes minimal
- follow existing pattern
- validate after changes
```

Contoh prompt siap pakai:

```text
Scope: execution
Goal: cleanup
Target: obx_site/src/uix
Context: rapikan komponen agar konsisten dengan pola proyek
Constraints:
- keep changes minimal
- follow existing pattern
- validate after changes
```

## Prinsip utama

1. Pahami request secara utuh sebelum mulai bekerja.
2. Fokus pada tujuan utama; jangan menambah fitur yang tidak diminta.
3. Gunakan pola yang sudah ada di repository, jangan menciptakan solusi baru jika tidak perlu.
4. Prioritaskan perubahan kecil, tepat, dan mudah dipelihara.
5. Selalu verifikasi sebelum menyatakan selesai.

## Output contract

Saat eksekusi selesai, agent wajib melaporkan dengan prinsip `token-efficient`:

- scope singkat dari pekerjaan
- file atau area yang berubah
- command verifikasi yang dijalankan
- hasil verifikasi yang segar
- blocker atau keputusan penting jika ada

Jika task tidak bisa selesai karena konteks belum cukup, sebutkan batasan itu sebelum memperluas scope.

## Workflow ringkas

1. Klarifikasi seperlunya: tanya hanya poin yang benar-benar memblokir.
2. Rencanakan singkat: 3 sampai 5 langkah, fokus satu slice dulu.
3. Implementasi sempit: edit file yang relevan, ikuti pola repo, hindari perubahan sampingan.
4. Verifikasi segera: setelah edit substantif pertama, jalankan validasi paling murah yang relevan.
5. Perbaiki root cause: jika gagal, perbaiki slice yang sama sebelum memperluas scope.

## Standar kualitas

- Kode harus rapi, terstruktur, dan mudah dibaca.
- Dokumentasi harus ringkas dan konsisten.
- Hindari logika yang berbelit-belit atau solusi over-engineered.
- Jika ada opsi sederhana dan opsi kompleks, pilih yang sederhana.

## Verifikasi

- Jalankan validasi paling murah yang relevan setelah edit substantif pertama.
- Jangan mengklaim selesai tanpa bukti verifikasi yang segar.
- Jika ada skill yang lebih spesifik, ikuti validasi layer yang diwajibkan skill tersebut.

## Gaya komunikasi

- Singkat, jelas, dan langsung ke inti.
- Jangan bertele-tele.
- Beri update progres seperlunya.
- Saat melaporkan hasil, sertakan status dan bukti verifikasi.

## Rule khusus

- Jangan mengubah banyak file jika satu file cukup.
- Jangan membuat file baru jika sudah ada yang bisa dipakai.
- Jangan mengulang pekerjaan yang sudah selesai.
- Saat ada ketidaksesuaian, perbaiki sesuai standar proyek yang sudah ada.
- Selalu utamakan kualitas, konsistensi, dan efisiensi.
- Jika perubahan memengaruhi dokumen yang dirujuk, lengkapi dokumentasi yang relevan sebagai bagian dari eksekusi (utama: `obx_docs/blueprint/{SUBMODUL}/README.md` untuk teknis dan `obx_docs/guide/{SUBMODUL}.md` untuk user guide).
- Jika dipakai bersama `module-creator`, jangan lompat ke implementasi tanpa rancangan yang sudah disetujui.
