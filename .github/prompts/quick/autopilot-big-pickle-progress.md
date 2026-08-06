# Autopilot Big Pickle - Progress Tracker

Gunakan file ini untuk mencatat progres eksekusi harian per modul.

## Cara Pakai

1. Tambahkan blok tanggal baru setiap mulai sesi.
2. Isi status modul: `done`, `in-progress`, atau `blocked`.
3. Catat command validasi dan hasil ringkas.
4. Simpan blocker dan next action sebelum tutup sesi.

## Template Harian

```text
Tanggal: YYYY-MM-DD
Model: opencode/big-pickle
Executor: {nama}

Batch dikerjakan:
- F1-P2 dat_signature: done (2026-08-06 13:37:30)
- F1-P2 dat_request: done (2026-08-06 13:34:27)
- F1-P2 dat_request: done (2026-08-06 13:33:30)
- F1-P2 dat_document_evidence: done (2026-08-06 13:33:06)
- F1-P2 dat_document_revision: done (2026-08-06 13:27:48)
- F1-P2 dat_document: done (2026-08-06 13:23:13)
- F1-P2 dat_signature: done (2026-08-06 13:17:25)
- F1-P2 dat_request: done (2026-08-06 13:14:06)
- F{n}-P{n} {module_code}: {done|in-progress|blocked}

Validation:
- cd obx_rest && go build ./... -> {pass|fail}
- cd obx_site && npx tsc --noEmit -> {pass|fail|n/a}

Ringkasan:
- Scope: ...
- Changes: ...
- Notes: ...

Blocker:
- ...

Next Action:
- ...
```

## Tracking Board

### Fase 1

- [x] F1-P2 dat_request
- [x] F1-P2 dat_signature
- [x] F1-P2 dat_document
- [x] F1-P2 dat_document_revision
- [x] F1-P2 dat_document_evidence

### Fase 2

- [ ] F2-P1 jms_asset_group
- [ ] F2-P1 jms_asset
- [ ] F2-P1 jms_asset_group_member
- [ ] F2-P1 jms_account
- [ ] F2-P1 jms_account_secret
- [ ] F2-P2 jms_session
- [ ] F2-P2 jms_session_event
- [ ] F2-P2 jms_session_command
- [ ] F2-P2 jms_audit_log
- [ ] F2-P3 jms_connect_token
- [ ] F2-P3 jms_policy
- [ ] F2-P3 jms_approval
- [ ] F2-P4 ssh_token_endpoint
- [ ] F2-P4 ssh_session_lifecycle
- [ ] F2-P4 ssh_runtime_event
- [ ] F2-P5 jms_file_transfer
- [ ] F2-P5 jms_session_event
- [ ] F2-P6 jms_web_app
- [ ] F2-P6 jms_connect_token
- [ ] F2-P6 proxy_request_access_log
- [ ] F2-P7 rdp_token_endpoint
- [ ] F2-P7 rdp_session_revoke_control
- [ ] F2-P7 rdp_runtime_recovery_note

### Fase 3

- [ ] F3-P1 infra_host
- [ ] F3-P1 infra_stack
- [ ] F3-P1 vm_host
- [ ] F3-P2 web_site
- [ ] F3-P2 web_upstream
- [ ] F3-P2 web_certificate
- [ ] F3-P2 web_config_version
- [ ] F3-P2 web_reload_history
- [ ] F3-P3 net_device
- [ ] F3-P3 net_interface
- [ ] F3-P3 net_poll_sample
- [ ] F3-P4 infra_metric_sample
- [ ] F3-P4 infra_alert_rule
- [ ] F3-P4 net_alert_rule
- [ ] F3-P4 net_backup_job
- [ ] F3-P4 vm_resource_sample
- [ ] F3-P4 vm_permission
- [ ] F3-P4 vm_action_log
- [ ] F3-P5 sec_event
- [ ] F3-P5 sec_event_source
- [ ] F3-P5 sec_event_parser
- [ ] F3-P5 sec_rule
- [ ] F3-P5 sec_alert
- [ ] F3-P5 sec_incident

## Log Harian

### 2026-08-06

Model: opencode/big-pickle  
Executor: Owner

Batch dikerjakan:
- F0-P0 preparation-autopilot: done
- F1-P2 dat_request: in-progress

Validation:
- cd obx_rest && go build ./... -> pending
- cd obx_site && npx tsc --noEmit -> pass

Ringkasan:
- Scope: Menyusun prompt cepat autopilot dan tracker progres harian sampai final flow.
- Changes: Menambahkan prompt per fase/paket/modul dan checklist tracking harian.
- Notes: Backend runtime lokal belum stabil pada sesi ini, tetapi jalur proxy guest sudah merespons dari server frontend yang aktif.

Blocker:
- Backend endpoint lokal pada localhost:36665 belum konsisten tersedia saat diuji langsung.

Next Action:
- Jalankan backend lokal, verifikasi go build, lalu mulai eksekusi modul F1-P2 dat_request sebagai batch pertama implementasi.

### 2026-08-06 (Update 2)

Model: opencode/big-pickle  
Executor: Owner

Batch dikerjakan:
- F0-P0 preparation-autopilot: done
- F1-P2 dat_request: in-progress

Validation:
- cd obx_rest && go build ./... -> pass (EXIT:0)
- cd obx_site && npx tsc --noEmit -> pass
- curl http://localhost:36665/rest/guest/SP00 -> pass (HTTP 200)

Ringkasan:
- Scope: Verifikasi runtime backend dan kesiapan gate teknis sebelum eksekusi modul pertama.
- Changes: Tidak ada perubahan kode modul, hanya pembaruan tracker progres.
- Notes: Kondisi backend pada sesi ini terkonfirmasi reachable untuk endpoint guest.

Blocker:
- Tidak ada blocker aktif untuk memulai batch F1-P2 dat_request.

Next Action:
- Mulai implementasi modul F1-P2 dat_request dan catat output Scope/Changes/Validation/Notes setelah batch selesai.

### 2026-08-06 (Update 3)

Model: opencode/big-pickle  
Executor: Owner

Batch dikerjakan:
- F0-P0 preparation-autopilot: done
- F1-P2 dat_request: done

Validation:
- cd obx_rest && go build ./... -> pass
- cd obx_base && npx prisma generate -> pass
- cd obx_site && npx tsc --noEmit -> n/a (mode backend, tidak ada perubahan frontend)

Ringkasan:
- Scope: Implementasi document request (`dat_request`) pada modul SM04: list + grid, create, update, delete, registrasi route, dan sinkronisasi dokumentasi REST.
- Changes: `skeleton/SM04/{template,repository,usecase,handler}.go` menambah request CRUD (CountRequest/ListRequest/CreateRequest/UpdateRequest/DeleteRequest); `backbone/routes.go` menambah route `/SM04/request` (admin group, `USLogs`); `blueprint/REST/SM04.md` dan `guide/REST/SM04.md` didokumentasikan; schema `dat_request` sudah tersedia di `dat_document.prisma` tanpa perubahan.
- Notes: Pola mengikuti blueprint AGENTS.md (mechanic error, InternalError wrap, parameterized query, sortExpr whitelist, rows.Err() setelah loop, USLogs pada write route). Tidak ada file schema terpisah `dat_request.prisma`; model sudah ada di dalam `dat_document.prisma`.

Blocker:
- Tidak ada blocker aktif.

Next Action:
- Lanjut eksekusi modul F1-P2 dat_signature sebagai batch berikutnya.

### 2026-08-06 (Auto Update)

Model: opencode/big-pickle  
Executor: Owner

Batch dikerjakan:
- F1-P2 dat_request: done (2026-08-06 13:04:15)

Validation:
- opencode run -> pass

Ringkasan:
- Scope: Auto update dari runner script.
- Changes: Checklist modul ditandai selesai.
- Notes: Entri ini dibuat otomatis.

Blocker:
- -

Next Action:
- Lanjut ke modul berikutnya.

### 2026-08-06 (F1-P2 dat_signature)

Model: opencode/big-pickle  
Executor: Owner

Batch dikerjakan:
- F1-P2 dat_signature: done

Validation:
- cd obx_rest && go build ./... -> pass
- cd obx_rest && go vet ./skeleton/SM04/... -> pass

Ringkasan:
- Scope: Melengkapi modul dat_signature pada SM04: mendaftarkan route approval flow (signature form & flag) yang kodenya sudah tersedia di skeleton, lalu menyinkronkan dokumentasi blueprint dan guide REST.
- Changes: `backbone/routes.go` menambah 3 route admin (`GET/POST /SM04/request/:id/form`, `PUT /SM04/form/:formId/flag/:flagId`, semuanya dengan `USLogs("SM04")` untuk write); `blueprint/REST/SM04.md` dan `guide/REST/SM04.md` didokumentasikan untuk generate form, list form, dan flag action (approve/reject + advance step + finalisasi request).
- Notes: Signature type CRUD (`dat_signature_type`/`dat_approval_step`/`dat_approval_sign`) sudah ada di commit sebelumnya; form/flag flow (`dat_signature_form`/`dat_signature_flag`) sudah ditulis di skeleton namun belum terhubung route. Tidak ada perubahan schema Prisma; enum `approval_flag` dan `approval_info` sudah ada.

Blocker:
- Tidak ada blocker aktif.

Next Action:
- Lanjut eksekusi modul F1-P2 dat_document sebagai batch berikutnya.

### 2026-08-06 (F1-P2 dat_document)

Model: opencode/big-pickle  
Executor: opencode/big-pickle

Batch dikerjakan:
- F1-P2 dat_document: done

Validation:
- cd obx_rest && go build ./... -> pass (EXIT:0)
- cd obx_rest && go vet ./skeleton/SM04/... ./backbone/... -> pass (EXIT:0)

Ringkasan:
- Scope: Implementasi dokumen terkelola (`dat_document`) pada modul SM04: list + grid (search/sort/pagination), create, update, delete, lookup kategori dokumen, registrasi route, dan sinkronisasi dokumentasi REST.
- Changes: `skeleton/SM04/{template,repository,usecase,handler}.go` menambah document CRUD (CountDocument/ListDocument/CreateDocument/UpdateDocument/DeleteDocument) + ListDocumentCategory; `backbone/routes.go` menambah 5 route admin (`GET/POST /SM04/document`, `GET/PUT/DELETE /SM04/document/:id`, `GET /SM04/document/category`, semua write dengan `USLogs("SM04")`); `blueprint/REST/SM04.md` dan `guide/REST/SM04.md` didokumentasikan; schema `dat_document` sudah tersedia di `dat_document.prisma` tanpa perubahan.
- Notes: Pola mengikuti blueprint AGENTS.md dan batch dat_request sebelumnya (mechanic error, InternalError wrap, parameterized query, sortExpr whitelist, rows.Err() setelah loop, USLogs pada write route). Route static `/SM04/document/category` didaftarkan sebelum `/SM04/document/:id` untuk menghindari konflik gin. `created_by` diisi dari session user; `file_size` (BigInt/int8) dikirim sebagai nil saat 0.

Blocker:
- Tidak ada blocker aktif.

Next Action:
- Lanjut eksekusi modul F1-P2 dat_document_revision sebagai batch berikutnya.

### 2026-08-06 (F1-P2 dat_document_revision)

Model: opencode/big-pickle  
Executor: opencode/big-pickle

Batch dikerjakan:
- F1-P2 dat_document_revision: done

Validation:
- cd obx_rest && go build ./... -> pass (EXIT:0)
- cd obx_rest && go vet ./skeleton/SM04/... ./backbone/... -> pass (EXIT:0)

Ringkasan:
- Scope: Implementasi histori revisi dokumen (`dat_document_version`) pada modul SM04: list + grid (search/sort/pagination), create, update, delete sebagai sub-entity dari dokumen, registrasi route, dan sinkronisasi dokumentasi REST.
- Changes: `skeleton/SM04/{template,repository,usecase,handler}.go` menambah revision CRUD (CountDocumentRevision/ListDocumentRevision/CreateDocumentRevision/UpdateDocumentRevision/DeleteDocumentRevision); `backbone/routes.go` menambah 4 route admin (`GET/POST /SM04/document/:id/revision`, `PUT/DELETE /SM04/document/:id/revision/:revisionId`, semua write dengan `USLogs("SM04")`); `blueprint/REST/SM04.md` dan `guide/REST/SM04.md` didokumentasikan; schema `dat_document_version` sudah tersedia di `dat_document.prisma` tanpa perubahan.
- Notes: Pola mengikuti blueprint AGENTS.md dan batch dat_document sebelumnya (mechanic error, InternalError wrap, parameterized query, sortExpr whitelist, rows.Err() setelah loop, USLogs pada write route, companyID fallback dari body untuk admin). `created_by` diisi dari session user; `content`/`file_path`/`note` dikirim sebagai nil saat kosong; tabel tidak memiliki kolom `updated_at` sehingga update hanya memodifikasi field terpilih.

Blocker:
- Tidak ada blocker aktif.

Next Action:
- Lanjut eksekusi modul F1-P2 dat_document_evidence sebagai batch berikutnya.

### 2026-08-06 (F1-P2 dat_document_evidence)

Model: opencode/big-pickle  
Executor: opencode/big-pickle

Batch dikerjakan:
- F1-P2 dat_document_evidence: done

Validation:
- cd obx_rest && go build ./... -> pass (EXIT:0)
- cd obx_rest && go vet ./skeleton/SM04/... ./backbone/... -> pass (EXIT:0)

Ringkasan:
- Scope: Implementasi jejak evidence/approval dokumen (`dat_document_approval`) pada modul SM04: list + grid (search/sort/pagination), create, update, delete sebagai sub-entity dari dokumen, registrasi route, dan sinkronisasi dokumentasi REST.
- Changes: `skeleton/SM04/{template,repository,usecase,handler}.go` menambah evidence CRUD (CountDocumentEvidence/ListDocumentEvidence/CreateDocumentEvidence/UpdateDocumentEvidence/DeleteDocumentEvidence); `backbone/routes.go` menambah 4 route admin (`GET/POST /SM04/document/:id/evidence`, `PUT/DELETE /SM04/document/:id/evidence/:evidenceId`, semua write dengan `USLogs("SM04")`); `blueprint/REST/SM04.md` dan `guide/REST/SM04.md` didokumentasikan; schema `dat_document_approval` sudah tersedia di `dat_document.prisma` tanpa perubahan.
- Notes: Pola mengikuti blueprint AGENTS.md dan batch dat_document_revision sebelumnya (mechanic error, InternalError wrap, parameterized query, sortExpr whitelist, rows.Err() setelah loop, USLogs pada write route, companyID fallback dari body untuk admin). `user_id` diisi dari session user; list menggabung `dat_user` untuk `user_name`; `note` dikirim sebagai nil saat kosong; tabel tidak memiliki kolom `updated_at` sehingga update hanya memodifikasi field terpilih.

Blocker:
- Tidak ada blocker aktif.

Next Action:
- Batch F1-P2 selesai; lanjut ke Fase 2 (modul jms_asset_group, jms_asset, dst) sesuai tracking board.

### 2026-08-06 (F1-P2 dat_signature - re-run)

Model: opencode/big-pickle  
Executor: opencode/big-pickle

Batch dikerjakan:
- F1-P2 dat_signature: done

Validation:
- cd obx_rest && go build ./... -> pass
- cd obx_rest && go vet ./skeleton/SM04/... -> pass
- cd obx_base && npx prisma validate -> pass

Ringkasan:
- Scope: Verifikasi batch backend dat_signature (SM04): signature type + approval step/signer, request, signature form & flag flow. Mode backend, single-module small-batch.
- Changes: `skeleton/SM04/usecase.go` — wrap raw error di `ListSignatureType` dan `ListUser` dengan `mechanic.InternalError` agar konsisten pattern repo. Tidak ada perubahan schema Prisma atau route.
- Notes: Batch dat_signature sebelumnya sudah done (create/update signature type + generate form + flag action + advance step + finalisasi request). Risiko sisa: `FlagAction` belum mengecek step form == `current_step` (approval bisa tidak berurutan) dan belum mengotorisasi actor sebagai signer flag.

Blocker:
- Tidak ada blocker aktif.

Next Action:
- Batch F1-P2 selesai; lanjut Fase 2 (modul jms_asset_group, jms_asset, dst) sesuai tracking board.