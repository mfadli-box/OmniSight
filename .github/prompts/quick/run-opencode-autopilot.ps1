param(
  [string]$Model = "opencode/big-pickle",
  [string]$Phase,
  [string]$Package,
  [string]$Module,
  [string]$Batch,
  [switch]$NoSmoke,
  [switch]$DryRun,
  [switch]$ContinueOnError,
  [switch]$UpdateTracker,
  [string]$TrackerPath = ".github/prompts/quick/autopilot-big-pickle-progress.md",
  [string]$Executor = "Owner"
)

$ErrorActionPreference = "Stop"

$batches = @{
  "F1-P2" = @("dat_request", "dat_signature", "dat_document", "dat_document_revision", "dat_document_evidence")
  "F2-P1" = @("jms_asset_group", "jms_asset", "jms_asset_group_member", "jms_account", "jms_account_secret")
  "F2-P2" = @("jms_session", "jms_session_event", "jms_session_command", "jms_audit_log")
  "F2-P3" = @("jms_connect_token", "jms_policy", "jms_approval")
  "F2-P4" = @("ssh_token_endpoint", "ssh_session_lifecycle", "ssh_runtime_event")
  "F2-P5" = @("jms_file_transfer", "jms_session_event")
  "F2-P6" = @("jms_web_app", "jms_connect_token", "proxy_request_access_log")
  "F2-P7" = @("rdp_token_endpoint", "rdp_session_revoke_control", "rdp_runtime_recovery_note")
  "F3-P1" = @("infra_host", "infra_stack", "vm_host")
  "F3-P2" = @("web_site", "web_upstream", "web_certificate", "web_config_version", "web_reload_history")
  "F3-P3" = @("net_device", "net_interface", "net_poll_sample")
  "F3-P4" = @("infra_metric_sample", "infra_alert_rule", "net_alert_rule", "net_backup_job", "vm_resource_sample", "vm_permission", "vm_action_log")
  "F3-P5" = @("sec_event", "sec_event_source", "sec_event_parser", "sec_rule", "sec_alert", "sec_incident")
}

function New-TaskPrompt {
  param(
    [string]$PhaseValue,
    [string]$PackageValue,
    [string]$ModuleValue,
    [bool]$SkipSmoke
  )

  $constraints = "minimal changes, follow repo pattern, validate"
  $notes = "Output: Scope/Changes/Validation/Notes."

  if ($SkipSmoke) {
    $constraints = "$constraints, skip smoke test for this run"
    $notes = "Output: Scope/Changes/Validation/Notes (smoke test skipped)."
  }

  return "Task: Eksekusi autopilot F$PhaseValue-P$PackageValue modul $ModuleValue; Mode: backend; Scope: single-module small-batch; Constraints: $constraints; $notes"
}

function Print-Usage {
  Write-Host "Usage examples:" -ForegroundColor Yellow
  Write-Host "  .\\run-opencode-autopilot.ps1 -Batch F1-P2"
  Write-Host "  .\\run-opencode-autopilot.ps1 -Batch F2-P1 -NoSmoke"
  Write-Host "  .\\run-opencode-autopilot.ps1 -Phase 1 -Package 2 -Module dat_request"
  Write-Host "  .\\run-opencode-autopilot.ps1 -Batch F3-P2 -DryRun"
  Write-Host "  .\\run-opencode-autopilot.ps1 -Batch F1-P2 -UpdateTracker"
}

function Update-ProgressTracker {
  param(
    [string]$FilePath,
    [string]$PhaseValue,
    [string]$PackageValue,
    [string]$ModuleValue,
    [string]$ExecutorValue
  )

  if (-not (Test-Path $FilePath)) {
    throw "Tracker file not found: $FilePath"
  }

  $key = "F$PhaseValue-P$PackageValue $ModuleValue"
  $content = Get-Content -Raw $FilePath

  $unchecked = "- [ ] $key"
  $checked = "- [x] $key"

  if ($content.Contains($unchecked)) {
    $content = $content.Replace($unchecked, $checked)
  }

  $dateText = (Get-Date).ToString("yyyy-MM-dd")
  $timestamp = (Get-Date).ToString("yyyy-MM-dd HH:mm:ss")
  $logHeader = "### $dateText (Auto Update)"

  if ($content.Contains($logHeader)) {
    $insertToken = "Batch dikerjakan:"
    $pos = $content.IndexOf($insertToken, [StringComparison]::Ordinal)
    if ($pos -ge 0) {
      $line = "`r`n- F$PhaseValue-P$PackageValue $($ModuleValue): done ($timestamp)"
      $content = $content.Insert($pos + $insertToken.Length, $line)
    }
  } else {
    $entry = @"

### $dateText (Auto Update)

Model: opencode/big-pickle  
Executor: $ExecutorValue

Batch dikerjakan:
- F$PhaseValue-P$PackageValue $($ModuleValue): done ($timestamp)

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
"@
    $content += $entry
  }

  Set-Content -Path $FilePath -Value $content -NoNewline
}

if (-not (Get-Command opencode -ErrorAction SilentlyContinue)) {
  throw "opencode CLI not found in PATH."
}

if ($UpdateTracker -and $DryRun) {
  Write-Host "UpdateTracker ignored in DryRun mode." -ForegroundColor Yellow
}

$targets = @()

if ($Batch) {
  if (-not $batches.ContainsKey($Batch)) {
    throw "Unknown batch '$Batch'. Available: $($batches.Keys -join ', ')"
  }

  if ($Batch -notmatch '^F(?<phase>\d+)-P(?<package>\d+)$') {
    throw "Invalid batch format '$Batch'. Use format F{n}-P{n}, for example F2-P1."
  }

  $batchPhase = $Matches.phase
  $batchPackage = $Matches.package

  foreach ($item in $batches[$Batch]) {
    $targets += [PSCustomObject]@{
      Phase = $batchPhase
      Package = $batchPackage
      Module = $item
    }
  }
} elseif ($Phase -and $Package -and $Module) {
  $targets += [PSCustomObject]@{
    Phase = $Phase
    Package = $Package
    Module = $Module
  }
} else {
  Print-Usage
  throw "Provide either -Batch or all of -Phase, -Package, and -Module."
}

$runCount = 0

foreach ($target in $targets) {
  $runCount += 1
  $message = New-TaskPrompt -PhaseValue $target.Phase -PackageValue $target.Package -ModuleValue $target.Module -SkipSmoke $NoSmoke.IsPresent

  Write-Host "[$runCount/$($targets.Count)] Running module $($target.Module) with model $Model" -ForegroundColor Cyan

  if ($DryRun) {
    Write-Host "DRY RUN PROMPT:" -ForegroundColor DarkYellow
    Write-Host $message
    continue
  }

  & opencode run -m $Model $message
  $exitCode = $LASTEXITCODE

  if ($exitCode -ne 0) {
    Write-Host "Command failed for module $($target.Module) with exit code $exitCode" -ForegroundColor Red
    if (-not $ContinueOnError) {
      exit $exitCode
    }
  } else {
    if ($UpdateTracker -and -not $DryRun) {
      Update-ProgressTracker -FilePath $TrackerPath -PhaseValue $target.Phase -PackageValue $target.Package -ModuleValue $target.Module -ExecutorValue $Executor
      Write-Host "Tracker updated for F$($target.Phase)-P$($target.Package) $($target.Module)" -ForegroundColor DarkGreen
    }
  }
}

if ($DryRun) {
  Write-Host "Dry run completed for $($targets.Count) target(s)." -ForegroundColor Green
} else {
  Write-Host "Execution completed for $($targets.Count) target(s)." -ForegroundColor Green
}