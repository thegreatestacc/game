# Лаунчер и приватные обновления GitHub

`dist/SpellGame.exe` теперь открывает **лаунчер**, а не игру напрямую.
В нём две кнопки:

- **Играть** — заново проверить последний стабильный GitHub Release, скачать
  игру при необходимости, проверить размер и SHA-256, затем запустить её.
- **Обновить игру** — выполнить те же проверки и установку без запуска.

При открытии лаунчер также проверяет доступный релиз. Ошибка сети, доступа,
отсутствующий релиз или повреждённая загрузка блокируют запуск. Офлайн-обхода
проверки в лаунчере нет. Повторная попытка — любой из двух кнопок.
Во время проверки, загрузки и работы игры повторные действия блокируются.

## Что уже подготовлено

- Репозиторий: `thegreatestacc/game` (приватный).
- `cmd/launcher` — графический лаунчер.
- `internal/updater` — клиент приватных Releases, загрузка и проверка файлов.
- `.github/workflows/release.yml` — тесты, сборка и публикация по тегу `v*`.
- `scripts/build-windows.ps1` — локальная сборка обоих EXE и манифеста.

Публикация и проверка настоящего приватного релиза из среды разработки пока
не выполнены: подключение GitHub возвращает 404 для этого репозитория.
Наличие `origin` в Git не предоставляет приложению доступ к приватным Releases.
Проверки загрузки и отказов выполнены на локальном HTTP-сервере в тестах.

## 1. Отправить исходники и выпустить первый релиз

В обычном PowerShell, из папки проекта:

```powershell
cd E:\projects\game
git add .gitignore README.md LAUNCHER.md go.mod go.sum cmd internal images scripts .github spellgame_mvp.md
git commit -m "Add launcher and private GitHub release updates"
git push -u origin HEAD
git tag v0.2.0
git push origin v0.2.0
```

Если удалённая ветка уже содержит коммиты, сначала синхронизируйте её обычным
Git-процессом; не используйте принудительную отправку для обхода конфликта.
Если файлы уже закоммичены, повторный коммит не нужен.

Откройте вкладку **Actions** репозитория и дождитесь успешного workflow
**Windows release**. Он создаст черновик, загрузит все файлы, затем опубликует
релиз как последний. До завершения загрузки черновик не виден лаунчеру.
В репозитории должны быть разрешены GitHub Actions и запись `contents` для workflow.
Токен для workflow GitHub предоставляет автоматически — PAT игрока ему не нужен.

В релизе будут:

- `SpellGame.exe` — лаунчер для пользователя.
- `SpellGame-game.exe` — исполняемый файл игры.
- `release-manifest.json` — версия, имя файла, размер и SHA-256 игры.
- `FONT-LICENSE.txt` — лицензия шрифта.

Теги стабильных релизов: `vMAJOR.MINOR.PATCH`, например `v0.2.0`, `v0.2.1`.
Изменение исходников само по себе не выпускает обновление: нужен новый тег.
Повторно не меняйте уже опубликованный релиз — создавайте следующую версию.
Если публикация упала и оставила черновик, проверьте его в Releases перед повтором.

## 2. Настроить доступ на компьютере игрока

Для приватного репозитория у пользователя должен быть доступ к нему в GitHub.
В настройках GitHub создайте **fine-grained personal access token**:

1. **Resource owner** — владелец `thegreatestacc/game`.
2. **Repository access** → **Only select repositories** → `game`.
3. **Repository permissions** → **Contents: Read-only**.
4. Задайте срок действия. Если организация требует одобрения, дождитесь его.

Токен нельзя встраивать в код/EXE или отправлять в репозиторий. На каждом
компьютере используйте личный токен пользователя с доступом к репозиторию.
Для установки переменной без вывода токена на экран и без токена в истории команд:

```powershell
$secret = Read-Host 'GitHub token' -AsSecureString
$launcherToken = [Net.NetworkCredential]::new('', $secret).Password
[Environment]::SetEnvironmentVariable('SPELLGAME_GITHUB_TOKEN', $launcherToken, 'User')
$env:SPELLGAME_GITHUB_TOKEN = $launcherToken
Remove-Variable launcherToken, secret
Start-Process -FilePath 'E:\projects\game\dist\SpellGame.exe'
```

Переменная сохраняется в пользовательских настройках Windows. Это не зашифрованное
хранилище секретов; не используйте общий Windows-аккаунт для чужих пользователей.
Уже открытые терминалы/IDE могут иметь старое окружение — перезапустите их.
После изменения или истечения токена обновите переменную и перезапустите лаунчер.

Для удаления токена:

```powershell
[Environment]::SetEnvironmentVariable('SPELLGAME_GITHUB_TOKEN', $null, 'User')
Remove-Item Env:SPELLGAME_GITHUB_TOKEN -ErrorAction SilentlyContinue
```

## 3. Выпуск следующих версий

Сохраните изменения в Git, отправьте ветку, создайте и отправьте новый тег:

```powershell
git push
git tag v0.2.1
git push origin v0.2.1
```

После успешного workflow лаунчер увидит новую версию при следующем запуске,
нажатии «Играть» или «Обновить игру». Автообновление самого лаунчера пока
не реализовано: при изменениях его протокола скачайте новый `SpellGame.exe` из Releases.

## Локальная сборка и проверка

```powershell
powershell -ExecutionPolicy Bypass -File scripts/build-windows.ps1 -Version v0.2.0
```

Версия по умолчанию — `dev`. Скрипт создаёт `dist/SpellGame.exe`,
`dist/SpellGame-game.exe` и `dist/release-manifest.json`.
Для разработки без публикации релиза запускайте `go run ./cmd/game` или
`dist/SpellGame-game.exe` напрямую. Путь через лаунчер всегда требует успешной
проверки релиза, включая локальную сборку лаунчера.

Игра устанавливается в `%LOCALAPPDATA%\SpellGame\versions\<SHA-256>\SpellGame-game.exe`.
Проверенная старая версия сохраняется. Новый файл сначала скачивается во временный
файл в своей папке, проверяется и только потом получает имя EXE.
Повреждённый файл и временная загрузка не запускаются. Оставшиеся старые версии
можно удалить вручную, когда игра закрыта; автоматической очистки пока нет.
Токен используется только лаунчером: он не передаётся процессу игры и снимается
с запросов при перенаправлении на файловое хранилище GitHub.

```powershell
go test ./internal/updater
go test ./...
```

Тесты проверяют авторизацию приватных assets, установку, повторное использование
проверенного файла, восстановление повреждённого файла, неправильный SHA-256,
прерванную загрузку, отмену запроса, HTTP-ошибки и отсутствие утечки токена при redirect.

Документация GitHub: [Private release assets](https://docs.github.com/en/rest/releases/assets#get-a-release-asset),
[Personal access tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens).
