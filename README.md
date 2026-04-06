<div align="center">

  <img src="https://files.catbox.moe/sjgj3f.png" alt="gitty" width="1280" height="720" />

  <h3>gitty - Git CLI that speaks human</h3>

  [![Windows](https://img.shields.io/badge/winget-Omibranch.Gitty-blue?logo=windows)](https://github.com/microsoft/winget-pkgs)
  [![AUR](https://img.shields.io/aur/version/gitty-cli?label=AUR&logo=archlinux)](https://aur.archlinux.org/packages/gitty-cli)
  [![Release](https://img.shields.io/github/v/release/Omibranch/gitty)](https://github.com/Omibranch/gitty/releases)
  [![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)](https://golang.org)
  [![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

  **[English](#english) · [Русский](#русский)**

</div>

---

<a name="english"></a>

# gitty - English

`gitty` is a single-binary CLI for Git and GitHub with short, human-readable commands.

Core idea:
- `gitty up` replaces `git add . && git commit -m "..." && git push`

If gitty saves you time, a ⭐ on GitHub means a lot — it helps others find it too.

## Why gitty

- Single binary, no external runtime dependencies
- Faster daily flow for commit/push operations
- Semantic command shortcuts (`to`, `from`, `in`)
- Works with GitHub CLI (`gh`) and supports proxy mode

## Installation

**Windows**
```sh
winget install Omibranch.Gitty
```

**Arch Linux (yay / paru)**
```sh
yay -S gitty-cli
```

> Binary name is `gitty`. AUR package is `gitty-cli` because `gitty` is already used by another project.

**macOS / Linux - Homebrew**
```sh
brew tap Omibranch/gitty
brew install gitty
```

**Linux / macOS - install script**
```sh
curl -fsSL https://raw.githubusercontent.com/Omibranch/gitty/master/install.sh | sh
```

**Manual**
- Download from [Releases](https://github.com/Omibranch/gitty/releases)
- Put binary into a folder from your `PATH`

## Quick start

```sh
# 1) Install git + gh if needed
# (safe to run multiple times)
gitty install

# 2) Login to GitHub once
gitty auth

# 3) Create and link a GitHub repo from current folder
gitty add repo "my-project"

# 4) Stage + commit + push
gitty up

# 5) Same, with explicit commit message
gitty up --commit "feat: initial CLI"
```

## Command syntax shortcuts

`gitty` supports semantic shortcuts:

- `=` means **to**
- `~` means **from**
- `*` means **in**
- words `to`, `from`, `in` are auto-converted to these symbols

Examples:

```sh
gitty push=main
gitty push to main

gitty pull~main
gitty pull from main

gitty checkpoint "v1"*main
gitty checkpoint "v1" in main
```

## Full command list

### Setup and info

| Command | Description |
|---|---|
| `gitty help` | Open full interactive manual (EN/RU picker) |
| `gitty install` | Install missing `git` and `gh`, add gitty location to PATH |
| `gitty auth` | Run `gh auth login` |
| `gitty status` | Show auth state, remote URL, current branch |
| `gitty state` | Show commits + branches + file count for current repo |
| `gitty state <owner/repo>` | Same stats for any public GitHub repo |
| `gitty state <url> --commits` | Show only commits count |
| `gitty state --branches` | Show only branches |
| `gitty state --files` | Show only files count |
| `gitty gitignore` | Interactive picker for official GitHub .gitignore templates |
| `gitty --version` / `gitty --v` / `gitty -v` | Print version |
| `gitty clear` | Clear terminal screen |

### Repository and branch operations

| Command | Description |
|---|---|
| `gitty init "https://github.com/u/r.git"` | Initialize git and set `origin` |
| `gitty clone "https://github.com/u/r.git"` | Clone repository |
| `gitty clone "https://github.com/u/r.git" "folder"` | Clone to custom folder |
| `gitty add repo "name"` | Create private GitHub repo and link folder |
| `gitty add repo "name" --public` | Create public GitHub repo |
| `gitty add branch "name"` | Create local branch (stays on current branch) |
| `gitty switch to branch <name>` | Switch to an existing branch |
| `gitty switch to repo <owner/repo>` | Rewire this folder to a different GitHub repo (affects `gitty state` and `gitty up`) |
| `gitty rename branch "old"="new"` | Rename branch locally and remotely |
| `gitty rename repo "new-name"` | Rename current linked GitHub repo |
| `gitty rename repo "old"="new"` | Rename specific GitHub repo |

### Daily workflow

| Command | Description |
|---|---|
| `gitty add .` | Stage all + auto commit |
| `gitty add . --commit "msg"` | Stage all + commit with custom message |
| `gitty up` | `add .` + push current branch |
| `gitty up --commit "msg"` | Same with custom commit message |
| `gitty push <branch>` | Push branch |
| `gitty push <branch> --share` | Push and copy GitHub link |
| `gitty push <branch> --force` | Force push |
| `gitty pull <branch>` | Safe pull (adds only missing files) |
| `gitty pull <branch> --hard` | Pull and overwrite existing files |
| `gitty pull <branch> --hard-reset` | Mirror remote exactly |
| `gitty undo` | Soft reset last commit, keep staged changes |

### History, checkpoints, conflict tools

| Command | Description |
|---|---|
| `gitty log` | Pretty graph log (last 1 week by default) |
| `gitty log --6h` | Last 6 hours |
| `gitty log --3day` | Last 3 days |
| `gitty log --2week` | Last 2 weeks |
| `gitty log --1month` | Last month |
| `gitty checkpoint "name"` | Create tag on current branch and push |
| `gitty checkpoint "name" in <branch>` | Create tag on a specific branch |
| `gitty restore "name"` | Restore working tree to tag (detached HEAD) |
| `gitty fix <file>` | Resolve merge markers interactively and stage file |
| `gitty erase <path>` | Remove file/folder from full git history |
| `gitty back <file> <N>` | Restore file from N commits ago |

### Advanced branch rewrite commands

| Command | Description |
|---|---|
| `gitty reset~<branch>` | Completely wipe branch history/content and recreate empty branch |
| `gitty migration <target>=<source>` | Replace all files in target branch with source branch files |
| `gitty migration <target> to <source>` | Same as above via semantic alias |

### Partial commit tool

| Command | Description |
|---|---|
| `gitty pick <file> 10-20` | Commit only lines 10-20 from file |
| `gitty pick <file> 10-*` | Commit from line 10 to EOF |
| `gitty pick <file> start1-end1` | Commit between markers `#gitty:start1` and `#gitty:end1` |

> Note: for `pick`, file should be staged first (`gitty add .`).

### Alias system (`.gittyconf`)

| Command | Description |
|---|---|
| `gitty alias` | List configured aliases |
| `gitty alias save "up --commit"` | Create/update alias |
| `gitty alias deploy "up and push=main"` | Create chain alias |
| `gitty alias <name> ""` | Remove alias |

Examples:

```sh
gitty alias quick "add . and push=main"
gitty quick
```

### Command chaining

Use `and` to execute multiple commands in one call:

```sh
gitty add . and push main
gitty up --commit "hotfix" and push=staging --share
```

### Proxy mode (global flag)

`--proxy` works with any command:

```sh
gitty up --proxy "http://1.2.3.4:8080"
gitty install --proxy "http://user:pass@1.2.3.4:8080"
```

## Safety notes

Potentially destructive commands:
- `gitty pull <branch> --hard-reset`
- `gitty reset~<branch>`
- `gitty migration <target>=<source>`
- `gitty erase <path>`

Use them only when you understand the impact.

## Build from source

```sh
git clone https://github.com/Omibranch/gitty
cd gitty/source

# Windows
go build -ldflags="-s -w" -o ../gitty.exe .

# Linux / macOS
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../gitty .
```

Requires Go 1.21+.

---

<a name="русский"></a>

# gitty - Русский

`gitty` - это однофайловый CLI для Git и GitHub с короткими, понятными командами.

Главная идея:
- `gitty up` заменяет `git add . && git commit -m "..." && git push`

Если gitty экономит тебе время — ⭐ на GitHub будет приятен и поможет другим найти проект.

## Зачем использовать gitty

- Один бинарник, без лишних зависимостей
- Быстрый повседневный workflow для Git
- Семантические сокращения (`to`, `from`, `in`)
- Работает с GitHub CLI (`gh`) и поддерживает прокси

## Установка

**Windows**
```sh
winget install Omibranch.Gitty
```

**Arch Linux (yay / paru)**
```sh
yay -S gitty-cli
```

> Бинарник называется `gitty`. Пакет в AUR называется `gitty-cli`, потому что имя `gitty` уже занято другим проектом.

**macOS / Linux - Homebrew**
```sh
brew tap Omibranch/gitty
brew install gitty
```

**Linux / macOS - скрипт установки**
```sh
curl -fsSL https://raw.githubusercontent.com/Omibranch/gitty/master/install.sh | sh
```

**Вручную**
- Скачайте из [Releases](https://github.com/Omibranch/gitty/releases)
- Положите бинарник в каталог из `PATH`

## Быстрый старт

```sh
# 1) Установить git и gh, если их нет
gitty install

# 2) Авторизоваться в GitHub (один раз)
gitty auth

# 3) Создать и привязать репозиторий
gitty add repo "my-project"

# 4) Сохранить и отправить изменения одной командой
gitty up

# 5) То же с явным сообщением коммита
gitty up --commit "feat: initial CLI"
```

## Сокращения синтаксиса

В `gitty` есть семантические сокращения:

- `=` означает **to** (куда)
- `~` означает **from** (откуда)
- `*` означает **in** (в какой ветке)
- слова `to`, `from`, `in` автоматически конвертируются в эти символы

Примеры:

```sh
gitty push=main
gitty push to main

gitty pull~main
gitty pull from main

gitty checkpoint "v1"*main
gitty checkpoint "v1" in main
```

## Полный список команд

### Установка и информация

| Команда | Что делает |
|---|---|
| `gitty help` | Открыть полное интерактивное руководство (выбор EN/RU) |
| `gitty install` | Поставить недостающие `git` и `gh`, добавить путь в PATH |
| `gitty auth` | Запустить `gh auth login` |
| `gitty status` | Показать авторизацию, remote URL, текущую ветку |
| `gitty state` | Показать коммиты + ветки + число файлов для текущего репо |
| `gitty state <owner/repo>` | То же для любого публичного репозитория |
| `gitty state <url> --commits` | Показать только число коммитов |
| `gitty state --branches` | Показать только ветки |
| `gitty state --files` | Показать только число файлов |
| `gitty gitignore` | Интерактивно выбрать и скачать шаблон .gitignore от GitHub |
| `gitty --version` / `gitty --v` / `gitty -v` | Показать версию |
| `gitty clear` | Очистить экран терминала |

### Репозитории и ветки

| Команда | Что делает |
|---|---|
| `gitty init "https://github.com/u/r.git"` | Инициализировать git и установить `origin` |
| `gitty clone "https://github.com/u/r.git"` | Клонировать репозиторий |
| `gitty clone "https://github.com/u/r.git" "folder"` | Клонировать в указанную папку |
| `gitty add repo "name"` | Создать приватный репозиторий на GitHub и привязать папку |
| `gitty add repo "name" --public` | Создать публичный репозиторий |
| `gitty add branch "name"` | Создать локальную ветку (без переключения) |
| `gitty switch to branch <имя>` | Переключиться на существующую ветку |
| `gitty switch to repo <owner/repo>` | Переключить папку на другой репозиторий GitHub (влияет на `gitty state` и `gitty up`) |
| `gitty rename branch "old"="new"` | Переименовать ветку локально и на remote |
| `gitty rename repo "new-name"` | Переименовать текущий привязанный репозиторий |
| `gitty rename repo "old"="new"` | Переименовать конкретный репозиторий |

### Ежедневный workflow

| Команда | Что делает |
|---|---|
| `gitty add .` | Добавить все изменения в stage + автокоммит |
| `gitty add . --commit "msg"` | Stage + коммит с вашим сообщением |
| `gitty up` | `add .` + push текущей ветки |
| `gitty up --commit "msg"` | То же с вашим сообщением коммита |
| `gitty push <ветка>` | Сделать push ветки |
| `gitty push <ветка> --share` | Push и скопировать ссылку GitHub |
| `gitty push <ветка> --force` | Принудительный push |
| `gitty pull <ветка>` | Безопасный pull (добавляет только недостающие файлы) |
| `gitty pull <ветка> --hard` | Pull с перезаписью существующих файлов |
| `gitty pull <ветка> --hard-reset` | Полное зеркалирование с удалением лишних локальных файлов |
| `gitty undo` | Мягкий откат последнего коммита, изменения остаются staged |

### История, чекпоинты, конфликты

| Команда | Что делает |
|---|---|
| `gitty log` | Красивый граф-лог (по умолчанию последняя неделя) |
| `gitty log --6h` | Лог за последние 6 часов |
| `gitty log --3day` | Лог за последние 3 дня |
| `gitty log --2week` | Лог за последние 2 недели |
| `gitty log --1month` | Лог за последний месяц |
| `gitty checkpoint "name"` | Создать тег в текущей ветке и запушить |
| `gitty checkpoint "name" in <ветка>` | Создать тег в указанной ветке |
| `gitty restore "name"` | Вернуть рабочее дерево к тегу (detached HEAD) |
| `gitty fix <файл>` | Найти и интерактивно разобрать merge-конфликт |
| `gitty erase <путь>` | Удалить файл/папку из всей истории git |
| `gitty back <файл> <N>` | Вернуть файл к состоянию N коммитов назад |

### Продвинутые команды перезаписи веток

| Команда | Что делает |
|---|---|
| `gitty reset~<ветка>` | Полностью очистить ветку (контент и история), создать пустую |
| `gitty migration <target>=<source>` | Заменить все файлы в target ветке файлами из source |
| `gitty migration <target> to <source>` | То же через семантический alias |

### Выборочный коммит по диапазону

| Команда | Что делает |
|---|---|
| `gitty pick <файл> 10-20` | Отправить только строки 10-20 |
| `gitty pick <файл> 10-*` | Отправить строки с 10-й до конца файла |
| `gitty pick <файл> start1-end1` | Отправить участок по меткам `#gitty:start1` и `#gitty:end1` |

> Для `pick` сначала добавьте изменения в stage (`gitty add .`).

### Система alias (`.gittyconf`)

| Команда | Что делает |
|---|---|
| `gitty alias` | Показать список алиасов |
| `gitty alias save "up --commit"` | Создать/обновить алиас |
| `gitty alias deploy "up and push=main"` | Создать цепочку команд как алиас |
| `gitty alias <name> ""` | Удалить алиас |

Пример:

```sh
gitty alias quick "add . and push=main"
gitty quick
```

### Цепочки команд

Используйте `and`, чтобы выполнять несколько команд подряд:

```sh
gitty add . and push main
gitty up --commit "hotfix" and push=staging --share
```

### Прокси-режим (глобальный флаг)

Флаг `--proxy` работает с любой командой:

```sh
gitty up --proxy "http://1.2.3.4:8080"
gitty install --proxy "http://user:pass@1.2.3.4:8080"
```

## Важные замечания по безопасности

Потенциально деструктивные команды:
- `gitty pull <ветка> --hard-reset`
- `gitty reset~<ветка>`
- `gitty migration <target>=<source>`
- `gitty erase <путь>`

Используйте их только если точно понимаете последствия.

## Сборка из исходников

```sh
git clone https://github.com/Omibranch/gitty
cd gitty/source

# Windows
go build -ldflags="-s -w" -o ../gitty.exe .

# Linux / macOS
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../gitty .
```

Требуется Go 1.21+.

---

<div align="center">Made with ❤ in Go</div>
