#!/usr/bin/env sh
set -eu
cd "$(dirname "$0")/.."
go mod tidy
mkdir -p dist
version="${1:-dev}"
repository="${2:-thegreatestacc/game}"
case "$version" in *[!a-zA-Z0-9.-]*|'') echo 'Invalid version'; exit 1;; esac
case "$repository" in *[!a-zA-Z0-9_./-]*|'') echo 'Invalid repository'; exit 1;; esac
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
  go build -trimpath -ldflags="-H=windowsgui -X game/internal/buildinfo.Version=$version -X game/internal/buildinfo.Repository=$repository" -o dist/SpellGame-game.exe ./cmd/game
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
  go build -trimpath -ldflags="-H=windowsgui -X game/internal/buildinfo.Version=$version -X game/internal/buildinfo.Repository=$repository" -o dist/SpellGame-Launcher.exe ./cmd/launcher
go run ./cmd/release-manifest "$version" dist/SpellGame-game.exe dist/release-manifest.json
cp internal/assets/fonts/LICENSE.txt dist/FONT-LICENSE.txt
echo 'Built dist/SpellGame-Launcher.exe (launcher) and dist/SpellGame-game.exe; verify on Windows.'

