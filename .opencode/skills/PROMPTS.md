# Koleksi Prompt Cepat

File ini berisi contoh prompt siap pakai untuk memanggil skill yang tepat dengan input yang ringkas dan konsisten.

## 1. Clean Execution

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

## 2. Module Creator

```text
Scope: module
Goal: buat modul baru
Module code: NM
Page code: NM01
Entity name: Notification
Route group: admin
Database table: dat_notification
Prisma schema: belum ada
Fields:
- id: text, primary key, default gen_random_uuid()
- company_id: text, not null, relation dat_company
- code: text, not null, unique per company
- message: text, not null
- severity: enum(critical, warning, info), default info
- is_active: boolean, default true
- created_at: timestamptz, default now()
- updated_at: timestamptz, default now()
Sub-entity: tidak ada
Kebutuhan UI: CRUD standar, select endpoint, mobile card view
```

## 3. Frontend Page Creator

```text
Scope: frontend-page
Goal: buat halaman baru
Page code: NM01
Page purpose: CRUD Notification
Page type: CRUD standar
Primary data: notification
Field form:
- code: text, required
- message: textarea, required
- severity: select, required
- is_active: checkbox, optional
Kebutuhan UI:
- DataTable: ya
- DataDialog: ya
- Tabs: tidak
- SearchSelect/Combobox: tidak
- Mobile card view: ya
Backend endpoint: sudah ada
```

## 4. Backend Page Creator

```text
Scope: backend-page
Goal: buat backend baru
Page code: NM01
Entity name: Notification
Route group: admin
Database table: dat_notification
Schema/table: sudah ada
Fields response:
- id: string
- code: string
- message: string
- severity: string
- is_active: boolean
- created_at: string
Request fields:
- code: required
- message: required
- severity: required
- is_active: optional
Sub-entity: tidak ada
Select endpoint: ya
Company scope: ya
```

## 5. Token Efficient

```text
Scope: efficiency
Goal: hemat token
Target: {file | folder | task}
Context: lakukan task dengan scope sempit, hasil ringkas, dan validasi minimal
Constraints:
- keep changes minimal
- use existing pattern
- avoid unnecessary explanation
- validate after changes
```

## 6. Alur cepat memilih prompt

- Jika task berupa refactor atau cleanup umum, mulai dari `Clean Execution`.
- Jika task dimulai dari desain modul atau penambahan halaman resmi modul, mulai dari `Module Creator`.
- Jika task hanya frontend satu halaman dan backend sudah siap, mulai dari `Frontend Page Creator`.
- Jika task hanya backend satu halaman dan desain modul sudah jelas, mulai dari `Backend Page Creator`.
- Jika user menginginkan hasil yang lebih hemat token atau singkat, pakai `Token Efficient` sebagai skill pendamping.

## 7. Format output yang disarankan

```text
Scope: {ringkas}
Changes: {file atau area yang berubah}
Validation: {command} -> {hasil}
Notes: {blocker, keputusan penting, atau saran lanjutan}
```
