# Panduan Penggunaan Copilot CLI di Workspace

Dokumen ini menjelaskan cara memakai Copilot CLI dengan workflow custom yang sudah disiapkan untuk project OmniSight.

## 1. Tujuan
Workflow ini dibuat agar Copilot CLI bekerja lebih konsisten seperti Opencode, terutama untuk:
- backend di obx_rest
- frontend di obx_site
- database di obx_base/prisma
- agent di obx_auto
- debugging cepat
- pekerjaan hemat token

## 2. Struktur file yang tersedia
- .github/instructions/: instruksi umum workspace
- .github/prompts/: prompt dan mode khusus
- .github/HELPER.md: ringkasan helper
- .copilot-instructions.md: instruksi utama Copilot

## 3. Cara memakai
### Mode backend
Gunakan prompt seperti:
```text
Task: buat atau perbaiki endpoint backend
Mode: backend
Scope: satu halaman atau satu file
```

### Mode frontend
```text
Task: buat halaman frontend CRUD
Mode: frontend
Scope: satu halaman
```

### Mode database
```text
Task: tambah field di schema prisma
Mode: database
Scope: satu model
```

### Mode agent
```text
Task: perbaiki logic agent automation
Mode: agent
Scope: satu service
```

### Mode debug
```text
Task: cari akar masalah error login
Mode: debug
Scope: satu file atau satu area
```

### Mode token-efficient
```text
Task: lakukan task ini hemat token
Mode: token-efficient
Scope: sempit
```

## 4. Contoh pemanggilan CLI
Contoh dasar:
```powershell
& "$env:APPDATA\Code\User\globalStorage\github.copilot-chat\copilotCli\copilot.ps1" -p "Task: periksa struktur repo. Mode: token-efficient. Scope: ringkas."
```

## 5. Format respons yang diharapkan
Copilot diharapkan memberikan respons dengan format:
```text
Scope: {singkat}
Changes: {file atau area}
Validation: {command} -> {hasil}
Notes: {blocker atau saran}
```

## 6. Catatan penting
- Selalu gunakan pola repo yang sudah ada.
- Fokus pada perubahan kecil dan terverifikasi.
- Jangan menambah fitur yang tidak diminta.
