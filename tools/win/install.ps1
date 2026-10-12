# Dro1dRa1n Windows installer — no admin needed.
# Installs dro1drain.exe to %LOCALAPPDATA%\Dro1dRa1n and adds it to your user PATH.
# Run from the release zip folder:  powershell -ExecutionPolicy Bypass -File install.ps1
$ErrorActionPreference = "Stop"

$dest = Join-Path $env:LOCALAPPDATA "Dro1dRa1n"
$exe  = Join-Path $PSScriptRoot "dro1drain.exe"

if (-not (Test-Path $exe)) {
    Write-Host "dro1drain.exe not found next to this script. Re-download the release zip." -ForegroundColor Red
    exit 1
}

New-Item -ItemType Directory -Force -Path $dest | Out-Null
Copy-Item $exe (Join-Path $dest "dro1drain.exe") -Force
Write-Host "Installed to $dest" -ForegroundColor Green

# user PATH (no admin)
$path = [Environment]::GetEnvironmentVariable("Path", "User")
if ($path -notlike "*$dest*") {
    [Environment]::SetEnvironmentVariable("Path", "$path;$dest", "User")
    Write-Host "Added to user PATH (restart your terminal to pick it up)." -ForegroundColor Green
} else {
    Write-Host "Already on user PATH."
}

# platform-tools check
if (-not (Get-Command adb -ErrorAction SilentlyContinue)) {
    Write-Host ""
    Write-Host "adb not found. Install Android platform-tools:" -ForegroundColor Yellow
    Write-Host "  winget install Google.PlatformTools"
    Write-Host "Samsung phones also need the Samsung USB driver:"
    Write-Host "  (bundled in the repo: android-hacking/Drivers/SAMSUNG-USB-Driver-for-Mobile-Phones.exe)"
} else {
    Write-Host "adb found: $((adb --version | Select-Object -First 1))"
}

Write-Host ""
Write-Host "Next: open a NEW terminal and run  dro1drain doctor"
