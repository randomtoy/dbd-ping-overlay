$ErrorActionPreference = "Stop"

try {
    $distDir = "dist"
    $output = Join-Path $distDir "dbd-ping-overlay.exe"

    if (-not (Test-Path $distDir)) {
        New-Item -ItemType Directory -Path $distDir | Out-Null
    }

    Write-Host "Building $output..."

    go build -trimpath -ldflags="-s -w -H windowsgui" -o $output ./cmd/dbd-ping-overlay
    if ($LASTEXITCODE -ne 0) {
        throw "go build failed with exit code $LASTEXITCODE"
    }

    Write-Host "Build succeeded: $output"
}
catch {
    Write-Host "Build failed: $_"
    exit 1
}
