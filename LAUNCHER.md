# Лаунчер SpellGame

Репозиторий [`thegreatestacc/game`](https://github.com/thegreatestacc/game) публичный. Лаунчер получает опубликованные релизы анонимно: пользователям не нужны токены, GitHub-аккаунты или настройка Windows. «Играть» всегда проверяет последний стабильный релиз; «Обновить игру» устанавливает его без запуска. Размер и SHA-256 сверяются с манифестом. Если GitHub недоступен или проверка не проходит, игра не запускается.

## Опубликовать первый релиз

Из PowerShell в корне проекта отправьте исходники:

```powershell
git add .gitignore README.md LAUNCHER.md go.mod go.sum cmd internal images scripts .github spellgame_mvp.md
git commit -m "Add launcher and public GitHub updates"
git push -u origin HEAD
```

Создайте тег и отправьте его:

```powershell
git tag v0.2.0
git push origin v0.2.0
```

GitHub Actions прогонит тесты, соберёт Windows x64 лаунчер и игру, загрузит файлы и опубликует Release. Следите за результатом во вкладке **Actions**. Первый Release должен успешно опубликоваться, прежде чем лаунчер сможет найти игру.

В Release будут `SpellGame-Launcher.exe` (лаунчер), `SpellGame-game.exe`, `release-manifest.json` и `FONT-LICENSE.txt`. В репозитории должны быть разрешены GitHub Actions и запись `contents` для workflow.

## Выпуск обновлений

После отправки изменений создайте новый тег SemVer:

```powershell
git push
git tag v0.2.1
git push origin v0.2.1
```

Опубликованные Releases не переиспользуйте: выпускайте следующий номер. Если workflow упал и оставил черновик Release, удалите незавершённый черновик на GitHub и отправьте новый тег версии.

## Локальная сборка

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-windows.ps1 -Version v0.2.0
```

Скрипт создаёт оба EXE в `dist`. Версия локальной сборки по умолчанию — `dev`.

Игра устанавливается в `%LOCALAPPDATA%\SpellGame\versions\<SHA-256>`. Она сначала скачивается во временный файл и проверяется. Существующие версии не перезаписываются. Лаунчер автоматически не обновляется: при изменении его протокола загрузите новый `SpellGame-Launcher.exe` из Releases.

Проверка:

```powershell
go test ./...
go vet ./...
```

