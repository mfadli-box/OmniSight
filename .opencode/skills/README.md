# Panduan Skill `.opencode/skills`

Panduan ini membantu memilih skill yang paling tepat agar workflow agent tetap rapi, sempit, dan konsisten.

## Ringkasan cepat

| Skill | Gunakan saat | Hindari saat |
|---|---|---|
| `clean-execution` | Implementasi umum, refactor, penyederhanaan, atau perapihan eksekusi | Task sudah jelas termasuk workflow modul, backend page, atau frontend page yang lebih spesifik |
| `module-creator` | Membuat modul baru, menambah halaman modul, atau menambah sub-entity dengan alur design-first | Hanya mengubah satu halaman frontend atau satu halaman backend yang struktur modulnya sudah jelas |
| `frontend-page-creator` | Membuat atau merapikan satu halaman frontend di `obx_site` | Task perlu design doc modul atau perubahan backend sebagai fokus utama |
| `backend-page-creator` | Membuat atau merapikan satu halaman backend di `obx_rest/skeleton/{KODE}/` | Task perlu design doc modul atau perubahan frontend sebagai fokus utama |
| `token-efficient` | Menghemat token, menjaga output ringkas, dan menjaga scope tetap sempit | Task yang membutuhkan desain panjang, penjelasan ekstensif, atau investigasi luas |

## Matriks keputusan cepat

| Jika request utama adalah... | Skill utama |
|---|---|
| Modul baru atau penambahan halaman resmi modul | `module-creator` |
| Satu halaman frontend dengan endpoint sudah ada | `frontend-page-creator` |
| Satu halaman backend atau skeleton route per kode halaman | `backend-page-creator` |
| Refactor, simplifikasi, atau perapihan umum | `clean-execution` |
| Hemat token, respons singkat, atau scope sempit | `token-efficient` (sebagai skill pendamping atau overlay) |

## Cara memilih

### 1. Gunakan `module-creator` jika:

- request menyebut modul baru
- request menambah halaman resmi baru ke modul
- request menambah sub-entity atau child flow
- request butuh rancangan sebelum implementasi

### 2. Gunakan `frontend-page-creator` jika:

- request fokus pada satu halaman frontend
- backend endpoint sudah ada atau bukan fokus utama
- task berkaitan dengan `DataTable`, `DataDialog`, form, tabs, select flow, atau mobile card view

### 3. Gunakan `backend-page-creator` jika:

- request fokus pada satu halaman backend
- task berkaitan dengan `template.go`, `repository.go`, `usecase.go`, `handler.go`, atau wiring route untuk satu kode halaman
- desain modul sudah jelas dan tidak perlu design doc formal

### 4. Gunakan `clean-execution` jika:

- task adalah refactor, simplifikasi, atau implementasi umum
- tidak ada workflow domain khusus yang lebih cocok
- Anda ingin perubahan tetap kecil, terkontrol, dan mudah diverifikasi

### 5. Gunakan `token-efficient` jika:

- user meminta hemat token atau respons singkat
- task sudah jelas dan dapat diselesaikan dengan scope sempit
- Anda ingin mengurangi pemborosan konteks dan menghindari output verbose

Jika task juga butuh workflow domain spesifik, gunakan skill spesifik tersebut lalu terapkan prinsip `token-efficient` di atasnya.

## Hubungan antar-skill

- `clean-execution` adalah aturan gaya eksekusi dasar.
- `module-creator` adalah workflow modul lintas desain dan implementasi.
- `frontend-page-creator` adalah workflow page-level frontend.
- `backend-page-creator` adalah workflow page-level backend.

Jika ada skill yang lebih spesifik, prioritaskan skill spesifik tersebut, lalu jalankan dengan prinsip `clean-execution` dan `token-efficient`.

## Standar output agent

Setiap eksekusi skill wajib menghasilkan output yang ringkas dan terukur, dengan prinsip `token-efficient` sebagai overlay wajib:

## Aturan dokumentasi

Saat skill dipakai, wajib melengkapi dokumen yang dirujuk jika perubahan memengaruhi:
- pola kerja atau workflow yang dijelaskan di obx_docs/blueprint/{SUBMODUL}/README.md atau skill docs
- struktur modul, route, atau skema yang dicantumkan di dokumentasi
- behavior atau alur implementasi yang sebelumnya dijelaskan

Prinsip ini berlaku untuk skill apapun, termasuk `module-creator`, `frontend-page-creator`, `backend-page-creator`, dan `clean-execution`.

1. **Scope singkat**: apa yang dikerjakan dan mengapa.
2. **Daftar file yang berubah**: file utama yang terpengaruh.
3. **Verifikasi**: command yang dijalankan dan hasilnya.
4. **Blocker jika ada**: kendala desain, schema, atau privilege yang menghambat.

Format respons yang disarankan:

```text
Scope: {singkat}
Changes: {file atau area yang berubah}
Validation: {command} -> {hasil}
Notes: {blocker atau catatan penting}
```

## Contoh pemetaan request

- "Buat modul baru NM dengan halaman NM01" → `module-creator`
- "Buat page frontend NM01 pakai DataTable dan DataDialog" → `frontend-page-creator`
- "Buat backend NM01 lengkap sampai route" → `backend-page-creator`
- "Rapikan flow form dan validasi biar konsisten" → `clean-execution`

## Workflow end-to-end

Contoh alur yang disarankan saat request berkembang dari level modul ke level implementasi:

### Contoh 1: Modul baru sampai implementasi penuh

1. Mulai dengan `module-creator` untuk sinkronisasi konteks, rancangan modul, dan approval user.
2. Setelah rancangan disetujui, lanjut implementasi backend page per kode halaman dengan pola `backend-page-creator`.
3. Setelah backend untuk halaman target siap, lanjut implementasi UI halaman dengan `frontend-page-creator`.
4. Selama semua fase, jalankan prinsip `clean-execution`: perubahan kecil, validasi cepat, dan tidak melebar.

### Contoh 2: Backend sudah ada, tinggal halaman frontend

1. Lewati `module-creator` jika desain modul dan endpoint sudah final.
2. Gunakan `frontend-page-creator` untuk membangun atau merapikan halaman.
3. Gunakan prinsip `clean-execution` untuk menjaga scope tetap sempit dan tervalidasi.

### Contoh 3: Hanya satu halaman backend yang belum ada

1. Lewati `module-creator` jika desain modul, privilege, dan schema sudah jelas.
2. Gunakan `backend-page-creator` untuk file skeleton dan wiring route halaman itu.
3. Gunakan prinsip `clean-execution` untuk validasi cepat dengan `go build ./...`.
