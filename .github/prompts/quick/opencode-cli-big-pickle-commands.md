# Opencode CLI Commands - Big Pickle

Panduan ini berisi command siap pakai untuk menjalankan prompt autopilot lewat Opencode CLI.

## Prasyarat

- Jalankan dari root workspace: `D:/devops/project`
- Model tersedia: `opencode/big-pickle`
- File prompt sudah ada di folder `.github/prompts/quick/`

## Cek Model

```powershell
opencode models
```

Pastikan ada `opencode/big-pickle` di daftar.

## Mode Interaktif (TUI)

```powershell
opencode . -m opencode/big-pickle
```

Lanjutkan sesi terakhir:

```powershell
opencode . -m opencode/big-pickle --continue
```

## Mode Non-Interaktif (Single Run)

### 1) Prompt inline

```powershell
opencode run -m opencode/big-pickle "Task: Eksekusi autopilot F1-P2 modul dat_request; Mode: backend; Scope: single-module small-batch; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes."
```

### 2) Prompt dari file

```powershell
$prompt = Get-Content -Raw ".github/prompts/quick/autopilot-big-pickle.md"
opencode run -m opencode/big-pickle $prompt
```

### 3) Prompt no-smoke dari file

```powershell
$prompt = Get-Content -Raw ".github/prompts/quick/autopilot-big-pickle-no-smoke.md"
opencode run -m opencode/big-pickle $prompt
```

## Jalankan Per Modul (Disarankan)

Gunakan one-liner per modul dari file prompt utama:

```powershell
opencode run -m opencode/big-pickle "Task: Eksekusi autopilot F2-P1 modul jms_asset; Mode: backend; Scope: single-module small-batch; Constraints: minimal changes, follow repo pattern, validate; Output: Scope/Changes/Validation/Notes."
```

## One-Click PowerShell Runner

Gunakan script ini untuk eksekusi per batch atau per modul:

```powershell
.\.github\prompts\quick\run-opencode-autopilot.ps1 -Batch F1-P2
```

No-smoke batch run:

```powershell
.\.github\prompts\quick\run-opencode-autopilot.ps1 -Batch F2-P1 -NoSmoke
```

Single module run:

```powershell
.\.github\prompts\quick\run-opencode-autopilot.ps1 -Phase 1 -Package 2 -Module dat_request
```

Preview prompt tanpa eksekusi:

```powershell
.\.github\prompts\quick\run-opencode-autopilot.ps1 -Batch F3-P2 -DryRun
```

Lanjutkan walau ada batch yang gagal:

```powershell
.\.github\prompts\quick\run-opencode-autopilot.ps1 -Batch F2-P2 -ContinueOnError
```

Update tracker otomatis saat modul sukses:

```powershell
.\.github\prompts\quick\run-opencode-autopilot.ps1 -Batch F1-P2 -UpdateTracker
```

Update tracker dengan executor custom:

```powershell
.\.github\prompts\quick\run-opencode-autopilot.ps1 -Batch F1-P2 -UpdateTracker -Executor "Owner"
```

Gunakan tracker path custom:

```powershell
.\.github\prompts\quick\run-opencode-autopilot.ps1 -Batch F1-P2 -UpdateTracker -TrackerPath ".github/prompts/quick/autopilot-big-pickle-progress.md"
```

## Pola No-Smoke Test

Gunakan constraint ini pada message:

```text
skip smoke test for this run
```

Contoh:

```powershell
opencode run -m opencode/big-pickle "Task: Eksekusi autopilot F1-P2 modul dat_signature; Mode: backend; Scope: single-module small-batch; Constraints: minimal changes, follow repo pattern, validate, skip smoke test for this run; Output: Scope/Changes/Validation/Notes (smoke test skipped)."
```

## Validasi Minimum Setelah Run

```powershell
cd obx_rest; go build ./...
cd ../obx_site; npx tsc --noEmit
```

## Referensi Cepat

- Prompt utama: `.github/prompts/quick/autopilot-big-pickle.md`
- Prompt no-smoke: `.github/prompts/quick/autopilot-big-pickle-no-smoke.md`
- Tracker progres: `.github/prompts/quick/autopilot-big-pickle-progress.md`
- Runner script: `.github/prompts/quick/run-opencode-autopilot.ps1`