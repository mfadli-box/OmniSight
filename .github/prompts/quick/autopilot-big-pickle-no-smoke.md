# Autopilot Big Pickle (No Smoke Test)

Prompt ini dipakai untuk eksekusi autopilot dengan model `opencode/big-pickle` tanpa menjalankan smoke test runtime.

Command Opencode CLI: lihat `opencode-cli-big-pickle-commands.md` pada folder yang sama.

## Kapan Dipakai

- Saat fokus hanya pada perubahan kode dan validasi statik.
- Saat environment runtime belum siap untuk pengujian endpoint interaktif.
- Saat batch kecil butuh verifikasi cepat sebelum review lanjutan.

## Batasan

- Tidak menjalankan smoke test endpoint/UI.
- Tetap wajib build/typecheck.
- Status akhir harus mencantumkan bahwa smoke test dilewati.

## Template Utama

```text
Task: Eksekusi autopilot {fase} paket {paket} modul {module_code}
Mode: backend
Scope: satu modul, small-batch
Model: opencode/big-pickle
Constraints:
- keep changes minimal
- follow existing repository patterns
- validate after changes
- skip smoke test for this run
Output:
- Scope
- Changes
- Validation
- Notes (include: smoke test skipped)
```

## Validasi Wajib (No Smoke)

```text
cd obx_rest
go build ./...
```

Jika terdampak frontend:

```text
cd obx_site
npx tsc --noEmit
```

## One-Liner Cepat

```text
Task: Eksekusi autopilot {fase}-{paket} modul {module_code}; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate, skip smoke test; Output: Scope/Changes/Validation/Notes (smoke test skipped).
```

## Contoh Siap Pakai

```text
Task: Eksekusi autopilot F1-P2 modul dat_request; Mode: backend; Scope: single-module small-batch; Model: opencode/big-pickle; Constraints: minimal changes, follow repo pattern, validate, skip smoke test; Output: Scope/Changes/Validation/Notes (smoke test skipped).
```