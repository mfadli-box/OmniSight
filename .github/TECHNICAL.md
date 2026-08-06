# Dokumentasi Teknis Copilot Workspace Helper

Dokumen ini menjelaskan arsitektur dan mekanisme konfigurasi Copilot CLI yang disiapkan untuk workspace OmniSight.

## 1. Tujuan arsitektur
Konfigurasi ini bertujuan untuk menyelaraskan perilaku Copilot CLI dengan workflow Opencode melalui:
- instruksi workspace yang konsisten
- prompt mode terstruktur
- format output standar
- panduan khusus untuk area backend, frontend, database, dan agent

## 2. Komponen utama
### 2.1 Instruksi workspace
File utama:
- .copilot-instructions.md
- .github/instructions/copilot-workflow.instructions.md
- .github/instructions/database-agent.instructions.md

Fungsi:
- memberikan aturan kerja umum
- menjaga konsistensi pola pengembangan
- mengarahkan agent ke workflow yang tepat

### 2.2 Prompt mode
Folder:
- .github/prompts/

Isi:
- prompt umum seperti clean-execution, module-creator, token-efficient
- prompt mode spesifik seperti backend, frontend, database, agent, debug
- prompt helper seperti master-template dan parity

### 2.3 Dokumentasi pendukung
- .github/HELPER.md: indeks helper
- .github/USAGE.md: panduan penggunaan
- .github/TECHNICAL.md: dokumentasi teknis

## 3. Alur kerja yang direkomendasikan
1. Tentukan area kerja: backend, frontend, database, agent, atau debug.
2. Pilih mode yang sesuai.
3. Berikan task singkat, scope, dan constraint.
4. Copilot merespons dengan output ringkas dan terverifikasi.

## 4. Format output standar
Semua workflow disarankan memakai format:
```text
Scope: {singkat}
Changes: {file atau area}
Validation: {command} -> {hasil}
Notes: {blocker atau saran}
```

## 5. Kesesuaian dengan Opencode
Perbedaan utama antara Opencode dan Copilot CLI adalah mekanisme eksekusi, tetapi struktur workflow yang disediakan sudah dibuat setara melalui:
- instruksi yang konsisten
- prompt mode yang terstruktur
- output format yang seragam

## 6. Keunggulan konfigurasi ini
- memudahkan penggunaan harian
- mengurangi hasil yang terlalu verbose
- menjaga konsistensi dengan project structure
- memudahkan onboarding pengguna baru
