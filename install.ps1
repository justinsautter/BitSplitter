$ErrorActionPreference = "Stop"

$Repo = "justinsautter/bitsplitter"
$InstallDir = "$env:LOCALAPPDATA\bitsplitter"

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Error "32-bit Windows is not supported."
    exit 1
}

# Get latest release tag
$Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
$Tag = $Release.tag_name

$Archive = "bitsplitter_windows_$Arch.zip"
$Url = "https://github.com/$Repo/releases/download/$Tag/$Archive"

Write-Host "Downloading bitsplitter $Tag for windows/$Arch..."

$TmpDir = Join-Path ([System.IO.Path]::GetTempPath()) "bitsplitter-install"
if (Test-Path $TmpDir) { Remove-Item -Recurse -Force $TmpDir }
New-Item -ItemType Directory -Path $TmpDir | Out-Null

$ZipPath = Join-Path $TmpDir $Archive
Invoke-WebRequest -Uri $Url -OutFile $ZipPath
Expand-Archive -Path $ZipPath -DestinationPath $TmpDir

# Install
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}
Copy-Item (Join-Path $TmpDir "bitsplitter.exe") -Destination $InstallDir -Force

# Add to PATH if not already present
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    Write-Host "Added $InstallDir to your PATH. Restart your terminal for it to take effect."
}

# Cleanup
Remove-Item -Recurse -Force $TmpDir

Write-Host "bitsplitter installed successfully to $InstallDir\bitsplitter.exe"
