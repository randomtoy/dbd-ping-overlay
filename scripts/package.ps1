$ErrorActionPreference = "Stop"

try {
    $distDir = "dist"
    $exePath = Join-Path $distDir "dbd-ping-overlay.exe"
    $packageDir = Join-Path $distDir "package"
    $zipPath = Join-Path $distDir "dbd-ping-overlay-windows-amd64.zip"

    if (-not (Test-Path $exePath)) {
        throw "Build output not found: $exePath. Run scripts/build.ps1 first."
    }

    if (Test-Path $packageDir) {
        Remove-Item -Recurse -Force $packageDir
    }
    New-Item -ItemType Directory -Path $packageDir | Out-Null

    Copy-Item $exePath $packageDir
    Copy-Item "README.md" $packageDir

    if (Test-Path "LICENSE") {
        Copy-Item "LICENSE" $packageDir
    }

    if (Test-Path $zipPath) {
        Remove-Item -Force $zipPath
    }

    Compress-Archive -Path (Join-Path $packageDir "*") -DestinationPath $zipPath

    Write-Host "Package created: $zipPath"
}
catch {
    Write-Host "Packaging failed: $_"
    exit 1
}
