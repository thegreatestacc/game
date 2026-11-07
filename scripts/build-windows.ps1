param(
    [string]$Version = 'dev',
    [string]$Repository = 'thegreatestacc/game'
)
$ErrorActionPreference = 'Stop'
if ($Version -notmatch '^(dev|v[0-9]+\.[0-9]+\.[0-9]+)$') { throw 'Version must be dev or vMAJOR.MINOR.PATCH.' }
if ($Repository -notmatch '^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$') { throw 'Repository must be owner/name.' }
Push-Location (Split-Path $PSScriptRoot -Parent)
$oldGOOS = $env:GOOS
$oldGOARCH = $env:GOARCH
$oldCGO = $env:CGO_ENABLED
try {
    $env:GOOS = 'windows'
    $env:GOARCH = 'amd64'
    $env:CGO_ENABLED = '0'
    go mod tidy
    if ($LASTEXITCODE -ne 0) { throw 'Dependency download failed.' }
    if ([Environment]::OSVersion.Platform -eq 'Win32NT') {
        go test ./...
        if ($LASTEXITCODE -ne 0) { throw 'Tests failed.' }
    }
    New-Item -ItemType Directory -Force dist | Out-Null
    $linkFlags = "-H=windowsgui -X game/internal/buildinfo.Version=$Version -X game/internal/buildinfo.Repository=$Repository"
    go build -trimpath "-ldflags=$linkFlags" -o dist/SpellGame-game.exe ./cmd/game
    if ($LASTEXITCODE -ne 0) { throw 'Game build failed.' }
    go build -trimpath "-ldflags=$linkFlags" -o dist/SpellGame.exe ./cmd/launcher
    if ($LASTEXITCODE -ne 0) { throw 'Launcher build failed.' }
    $manifest = [ordered]@{
        version = $Version
        asset = 'SpellGame-game.exe'
        sha256 = (Get-FileHash -LiteralPath dist/SpellGame-game.exe -Algorithm SHA256).Hash.ToLowerInvariant()
        size = (Get-Item -LiteralPath dist/SpellGame-game.exe).Length
    }
    $manifestJson = $manifest | ConvertTo-Json
    [IO.File]::WriteAllText((Join-Path (Get-Location) 'dist/release-manifest.json'), $manifestJson, (New-Object Text.UTF8Encoding($false)))
    Copy-Item internal/assets/fonts/LICENSE.txt dist/FONT-LICENSE.txt -Force
    Write-Host 'Built launcher: dist/SpellGame.exe; game: dist/SpellGame-game.exe; release-manifest.json'
} finally {
    $env:GOOS = $oldGOOS
    $env:GOARCH = $oldGOARCH
    $env:CGO_ENABLED = $oldCGO
    Pop-Location
}
