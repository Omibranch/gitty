<div align="center">

  <img src="https://files.catbox.moe/sjgj3f.png" alt="gitty" width="1280" height="720" />

  <h3>gitty — Git CLI that speaks human</h3>

  [![Windows](https://img.shields.io/badge/winget-Omibranch.Gitty-blue?logo=windows)](https://github.com/microsoft/winget-pkgs)
  [![AUR](https://img.shields.io/aur/version/gitty-cli?label=AUR&logo=archlinux)](https://aur.archlinux.org/packages/gitty-cli)
  [![Release](https://img.shields.io/github/v/release/Omibranch/gitty)](https://github.com/Omibranch/gitty/releases)
  [![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)](https://golang.org)

  **[English](#english) · [Русский](#русский)**

</div>

---

<a name="english"></a>

# gitty — English

gitty is a single-binary CLI that wraps Git and GitHub into short, human-readable commands. You never type `git add`, `git commit`, or remember branch flags — gitty does it for you.

**Core idea:** `gitty up` replaces the entire add → commit → push cycle in one word.

---

## Installation

**Windows**
```sh
winget install Omibranch.Gitty
```

**Arch Linux (via yay or paru)**
```sh
yay -S gitty-cli
```
> Binary installs as `gitty`. The package name is `gitty-cli` because `gitty` was already taken in AUR by a different project.

**Linux / macOS — install script**
```sh
curl -fsSL https://raw.githubusercontent.com/Omibranch/gitty/master/install.sh | sh
```

**Manual** — download the binary from [Releases](https://github.com/Omibranch/gitty/releases) and put it somewhere on your `PATH`.

---

## Quick start

```sh
# 1. Auth with GitHub (once)
gitty auth

# 2. Create a repo and link current folder to it
gitty add repo "my-project"

# 3. Stage, commit and push — all at once
gitty up

# 4. Push with a custom commit message
gitty up --commit "fix: typo in README"
```

---

## All commands

### Everyday workflow

| Command | What it does |
|---|---|
| `gitty up` | Stage everything → auto-commit (timestamp) → push |
| `gitty up --commit "msg"` | Same, but with your own commit message |
| `gitty add .` | Stage everything + auto-commit (no push) |
| `gitty add . --commit "msg"` | Stage everything + commit with custom message |
| `gitty push <branch>` | Push committed changes to branch |
| `gitty push <branch> --share` | Push + copy GitHub URL to clipboard |
| `gitty push <branch> --force` | Force push |
| `gitty pull <branch>` | Safe pull — only adds new files, never overwrites |
| `gitty pull <branch> --hard` | Pull + overwrite files that exist locally |
| `gitty pull <branch> --hard-reset` | Mirror remote exactly (deletes local-only files) |
| `gitty undo` | Soft-reset last commit, keep changes staged |

### Shorthand syntax

gitty uses `=` for "to" and `~` for "from":

```sh
gitty push=main           # → gitty push main
gitty push to main        # → gitty push main
gitty pull~main           # → gitty pull main
gitty pull from main      # → gitty pull main
```

### Repositories & branches

| Command | What it does |
|---|---|
| `gitty init "https://github.com/u/r.git"` | Init git + set remote origin |
| `gitty clone "https://github.com/u/r.git"` | Clone repo into current folder |
| `gitty clone "https://github.com/u/r.git" "folder"` | Clone into specific folder |
| `gitty add repo "name"` | Create private GitHub repo + link folder |
| `gitty add repo "name" --public` | Create public GitHub repo |
| `gitty add branch "name"` | Create new local branch (stay on current) |
| `gitty rename branch "old"=="new"` | Rename branch locally + on remote |
| `gitty rename repo "new-name"` | Rename GitHub repo |

### Selective commit

| Command | What it does |
|---|---|
| `gitty pick main.go 10-20` | Send only lines 10–20 of a file |
| `gitty pick main.go 10-*` | Send from line 10 to end of file |
| `gitty pick i.py start1-end1` | Send by code markers `#gitty:start1` / `#gitty:end1` |

> The file must be staged first (`gitty add .`). The full file stays intact locally — only the selected lines go to the repo.

### Conflict resolution

```sh
gitty fix <file>
```
Finds `<<<<<<<` / `=======` / `>>>>>>>` markers in the file and asks interactively:
- **Keep mine** — discard the incoming changes
- **Take theirs** — discard your local changes
- **Merge both** — stitch one block under the other

Cleans up the conflict markers and stages the file automatically.

### History tools

| Command | What it does |
|---|---|
| `gitty erase <path>` | Remove a file/folder from **all** past commits forever |
| `gitty back <file> <N>` | Restore a file to its state N commits ago |
| `gitty undo` | Soft-reset last commit (changes stay staged) |

> `gitty erase` is useful when you accidentally pushed passwords, keys, or a huge `node_modules` folder.

### Checkpoints (tags)

```sh
gitty checkpoint "v1" in main   # tag the tip of main as "v1" and push
gitty checkpoint "v1"*main      # shorthand
gitty checkpoint "v1"           # use current branch
gitty restore "v1"              # reset working tree to that tag (detached HEAD)
```

### Log & stats

| Command | What it does |
|---|---|
| `gitty log` | Pretty graph log — last 1 week by default |
| `gitty log --5h` | Last 5 hours |
| `gitty log --3day` | Last 3 days |
| `gitty log --2week` | Last 2 weeks |
| `gitty log --1month` | Last month |
| `gitty state` | Commits, branches, file count for current repo |
| `gitty state <URL>` | Same stats for any public GitHub repo |
| `gitty state --branches` | Branch list only |
| `gitty state --commits` | Commit count only |
| `gitty state --files` | File count only |
| `gitty status` | Show linked GitHub account + remote + active branch |

### Aliases

```sh
gitty alias save "up --commit"          # define alias
gitty alias deploy "up and push=prod"   # chain with 'and'
gitty save "hotfix deployed"            # use it
```

Aliases are stored in `.gittyconf` in the repo root. The file is created automatically on `gitty install` with commented examples.

### Utilities

| Command | What it does |
|---|---|
| `gitty gitignore` | Interactive picker for GitHub's official .gitignore templates |
| `gitty install` | Install git + gh CLI (if missing) and add gitty to PATH |
| `gitty auth` | Run `gh auth login` to connect GitHub account |
| `gitty clear` | Clear terminal screen |
| `gitty help` | Full built-in manual (language picker: EN / RU) |
| `gitty --version` | Print version |

### Chaining commands

Use `and` to chain multiple commands:

```sh
gitty add . and push main
gitty up and push=staging
```

### Proxy support

```sh
gitty up --proxy "http://1.2.3.4:8080"
gitty up --proxy "http://user:pass@1.2.3.4:8080"
```

---

## Build from source

```sh
git clone https://github.com/Omibranch/gitty
cd gitty/source

# Windows
go build -ldflags="-s -w" -o ../gitty.exe .

# Linux / macOS
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../gitty .
```

Requires Go 1.21+. No external dependencies.

---

---

<a name="русский"></a>

# gitty — Русский

gitty — однофайловая утилита командной строки, которая оборачивает Git и GitHub в короткие, человекочитаемые команды. Вы больше не пишете `git add`, `git commit` и не запоминаете флаги веток — gitty делает это за вас.

**Главная идея:** `gitty up` заменяет весь цикл add → commit → push одним словом.

---

## Установка

**Windows**
```sh
winget install Omibranch.Gitty
```

**Arch Linux (через yay или paru)**
```sh
yay -S gitty-cli
```
> Бинарник устанавливается как `gitty`. Пакет называется `gitty-cli`, потому что имя `gitty` в AUR уже занято другим проектом.

**Linux / macOS — скрипт установки**
```sh
curl -fsSL https://raw.githubusercontent.com/Omibranch/gitty/master/install.sh | sh
```

**Вручную** — скачайте бинарник из [Releases](https://github.com/Omibranch/gitty/releases) и положите его в папку из `PATH`.

---

## Быстрый старт

```sh
# 1. Авторизация в GitHub (один раз)
gitty auth

# 2. Создать репозиторий и привязать текущую папку
gitty add repo "my-project"

# 3. Сохранить и отправить всё — одной командой
gitty up

# 4. С произвольным сообщением коммита
gitty up --commit "fix: опечатка в README"
```

---

## Все команды

### Ежедневная работа

| Команда | Что делает |
|---|---|
| `gitty up` | Стейдж → автокоммит (временна́я метка) → push |
| `gitty up --commit "текст"` | То же, но со своим сообщением коммита |
| `gitty add .` | Стейдж + автокоммит (без push) |
| `gitty add . --commit "текст"` | Стейдж + коммит с вашим сообщением |
| `gitty push <ветка>` | Отправить закоммиченные изменения в ветку |
| `gitty push <ветка> --share` | Push + скопировать GitHub-ссылку в буфер |
| `gitty push <ветка> --force` | Принудительный push |
| `gitty pull <ветка>` | Безопасный pull — добавляет только новые файлы |
| `gitty pull <ветка> --hard` | Pull + перезаписать файлы, которые есть локально |
| `gitty pull <ветка> --hard-reset` | Зеркало remote (удаляет локальные файлы, которых нет на remote) |
| `gitty undo` | Мягкий откат последнего коммита, изменения остаются в staged |

### Сокращённый синтаксис

gitty использует `=` вместо «куда» и `~` вместо «откуда»:

```sh
gitty push=main           # → gitty push main
gitty push to main        # → gitty push main
gitty pull~main           # → gitty pull main
gitty pull from main      # → gitty pull main
```

### Репозитории и ветки

| Команда | Что делает |
|---|---|
| `gitty init "https://github.com/u/r.git"` | Инициализировать git + задать remote origin |
| `gitty clone "https://github.com/u/r.git"` | Клонировать репо в текущую папку |
| `gitty clone "https://github.com/u/r.git" "папка"` | Клонировать в конкретную папку |
| `gitty add repo "название"` | Создать приватный репо на GitHub + привязать папку |
| `gitty add repo "название" --public` | Создать публичный репо |
| `gitty add branch "название"` | Создать новую локальную ветку (остаться на текущей) |
| `gitty rename branch "старое"=="новое"` | Переименовать ветку локально + на remote |
| `gitty rename repo "новое-имя"` | Переименовать репо на GitHub |

### Выборочный коммит

| Команда | Что делает |
|---|---|
| `gitty pick main.go 10-20` | Отправить только строки 10–20 файла |
| `gitty pick main.go 10-*` | Отправить от строки 10 до конца файла |
| `gitty pick i.py start1-end1` | Отправить по меткам в коде `#gitty:start1` / `#gitty:end1` |

> Файл должен быть добавлен через `gitty add .`. Локально файл остаётся полным — в репо улетает только выбранный фрагмент.

### Разрешение конфликтов

```sh
gitty fix <файл>
```
Находит в файле маркеры `<<<<<<<` / `=======` / `>>>>>>>` и спрашивает интерактивно:
- **Оставить моё** — отбросить входящие изменения
- **Взять чужое** — отбросить ваши локальные изменения
- **Объединить** — склеить оба блока один под другим

Автоматически очищает маркеры конфликта и делает `git add` файла.

### Работа с историей

| Команда | Что делает |
|---|---|
| `gitty erase <путь>` | Удалить файл/папку из **всех** прошлых коммитов навсегда |
| `gitty back <файл> <N>` | Вернуть файл к версии N коммитов назад |
| `gitty undo` | Мягкий откат последнего коммита (изменения остаются в staged) |

> `gitty erase` нужен, если случайно залили пароли, ключи или огромную папку `node_modules`.

### Чекпоинты (теги)

```sh
gitty checkpoint "v1" in main   # пометить кончик ветки main тегом "v1" и запушить
gitty checkpoint "v1"*main      # сокращённо
gitty checkpoint "v1"           # использовать текущую ветку
gitty restore "v1"              # откатить рабочую директорию к тегу (detached HEAD)
```

### Лог и статистика

| Команда | Что делает |
|---|---|
| `gitty log` | Красивый граф-лог — последняя неделя по умолчанию |
| `gitty log --5h` | Последние 5 часов |
| `gitty log --3day` | Последние 3 дня |
| `gitty log --2week` | Последние 2 недели |
| `gitty log --1month` | Последний месяц |
| `gitty state` | Коммиты, ветки, файлы текущего репо |
| `gitty state <URL>` | То же для любого публичного репо на GitHub |
| `gitty state --branches` | Только список веток |
| `gitty state --commits` | Только количество коммитов |
| `gitty state --files` | Только количество файлов |
| `gitty status` | Показать привязанный аккаунт GitHub + remote + активную ветку |

### Алиасы

```sh
gitty alias save "up --commit"          # задать алиас
gitty alias deploy "up and push=prod"   # цепочка через 'and'
gitty save "hotfix deployed"            # использовать
```

Алиасы хранятся в `.gittyconf` в корне репо. Файл создаётся автоматически при `gitty install` с примерами в комментариях.

### Утилиты

| Команда | Что делает |
|---|---|
| `gitty gitignore` | Интерактивный выбор шаблона .gitignore от GitHub |
| `gitty install` | Установить git + gh CLI (если нет) и добавить gitty в PATH |
| `gitty auth` | Запустить `gh auth login` для подключения аккаунта GitHub |
| `gitty clear` | Очистить экран терминала |
| `gitty help` | Полное встроенное руководство (выбор языка: EN / RU) |
| `gitty --version` | Версия |

### Цепочки команд

Используйте `and` для объединения нескольких команд:

```sh
gitty add . and push main
gitty up and push=staging
```

### Прокси

```sh
gitty up --proxy "http://1.2.3.4:8080"
gitty up --proxy "http://user:pass@1.2.3.4:8080"
```

---

## Сборка из исходников

```sh
git clone https://github.com/Omibranch/gitty
cd gitty/source

# Windows
go build -ldflags="-s -w" -o ../gitty.exe .

# Linux / macOS
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../gitty .
```

Требуется Go 1.21+. Внешних зависимостей нет.

---

<div align="center">Made with ❤️ in Go</div>
