# gitty installer
# Usage: irm https://raw.githubusercontent.com/Omibranch/gitty/main/install.ps1 | iex

$ErrorActionPreference = 'Stop'

$repo    = 'Omibranch/gitty'
$installDir = "$env:LOCALAPPDATA\gitty"

Write-Host ""
Write-Host "  gitty installer" -ForegroundColor Cyan
Write-Host "  ───────────────────────────────────" -ForegroundColor DarkGray
Write-Host ""

# ── 1. Fetch latest release ───────────────────────────────────────────────────
Write-Host "[INFO] Fetching latest release from GitHub..." -ForegroundColor Cyan
$release = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest" -UseBasicParsing
$tag     = $release.tag_name
$asset   = $release.assets | Where-Object { $_.name -eq 'gitty.exe' } | Select-Object -First 1

if (-not $asset) {
    Write-Host "[ERROR] Could not find gitty.exe in release $tag" -ForegroundColor Red
    exit 1
}

$downloadUrl = $asset.browser_download_url
Write-Host "[INFO] Found $tag — downloading gitty.exe..." -ForegroundColor Cyan

# ── 2. Download gitty.exe ─────────────────────────────────────────────────────
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

$dest = Join-Path $installDir 'gitty.exe'
Invoke-WebRequest -Uri $downloadUrl -OutFile $dest -UseBasicParsing
Write-Host "[SUCCESS] Downloaded to $dest" -ForegroundColor Green

# ── 3. Add to User PATH ───────────────────────────────────────────────────────
$currentPath = [Environment]::GetEnvironmentVariable('PATH', 'User')
if ($currentPath -split ';' | Where-Object { $_ -ieq $installDir }) {
    Write-Host "[INFO] $installDir is already in PATH" -ForegroundColor Cyan
} else {
    $newPath = if ($currentPath) { "$currentPath;$installDir" } else { $installDir }
    [Environment]::SetEnvironmentVariable('PATH', $newPath, 'User')
    Write-Host "[SUCCESS] Added $installDir to User PATH" -ForegroundColor Green
}

# Make available in the current session immediately
$env:PATH = "$env:PATH;$installDir"

# ── 4. Run gitty install (sets up git + gh) ───────────────────────────────────
Write-Host ""
Write-Host "[INFO] Running 'gitty install' to set up git and gh CLI..." -ForegroundColor Cyan
Write-Host ""
& $dest install

Write-Host ""
Write-Host "[SUCCESS] gitty $tag is installed!" -ForegroundColor Green
Write-Host "[HINT] Restart your terminal, then type: gitty help" -ForegroundColor Yellow
Write-Host ""
