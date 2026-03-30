<div align="center">

<img src="https://raw.githubusercontent.com/Omibranch/gitty/logo.png" alt="gitty" width="80" />

# gitty

**EN** | [RU](#-руководство-ru)

---

*Zero-commit Git workflow for humans.*  
Stage. Commit. Push. — All in one command.

[![Release](https://img.shields.io/github/v/release/Omibranch/gitty?color=black&style=flat-square)](https://github.com/Omibranch/gitty/releases)
[![License: MIT](https://img.shields.io/badge/license-MIT-black?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows-black?style=flat-square)](https://github.com/Omibranch/gitty/releases)

</div>

---

## 📖 Guide (EN)

### What is gitty?

**gitty** is a Windows CLI tool that wraps Git and GitHub CLI into simple, human-friendly commands.  
The core idea: **you never type `git commit`**. gitty handles staging, committing, and pushing automatically.

`->` means **TO** (push to a branch) &nbsp;|&nbsp; `~` means **FROM** (pull from a branch)

---

### ⚡ Quick Install

**Option 1 — PowerShell one-liner (recommended)**
```powershell
irm https://raw.githubusercontent.com/Omibranch/gitty/main/install.ps1 | iex
```
Downloads `gitty.exe`, adds it to your User PATH, then runs `gitty install` to set up Git and GitHub CLI automatically. No admin rights required.

**Option 2 — winget**
```
winget install Omibranch.Gitty
```

**Option 3 — Manual**  
Download `gitty.exe` from [Releases](https://github.com/Omibranch/gitty/releases), place it anywhere in your PATH, then run:
```
gitty install
```

---

### 📋 Commands

| Command | Description |
|---|---|
| `gitty install` | Install Git + GitHub CLI, add gitty to PATH (no admin needed) |
| `gitty auth` | Sign in to GitHub via `gh auth login` |
| `gitty init "url"` | Init a git repo and set remote origin |
| `gitty add repo "name"` | Create a private GitHub repo and link current folder |
| `gitty add repo "name" --public` | Same, but public |
| `gitty add branch "name"` | Create a local branch without switching to it |
| `gitty add .` | Stage all changes and auto-commit |
| `gitty push->branch` | Push committed changes to a remote branch |
| `gitty pull~branch` | Safe pull — adds only files missing locally |
| `gitty pull~branch --hard` | Overwrite local files from remote (keeps unique local files) |
| `gitty pull~branch --hard-reset` | Mirror remote exactly — destructive, requires confirmation |
| `gitty reset~branch` | Wipe all content and history from a branch |
| `gitty status` | Show linked repo, current branch, and GitHub account |
| `gitty gitignore` | Interactive picker for `.gitignore` templates |
| `gitty clear` | Clear the terminal screen |
| `gitty help` | Full manual (EN / RU) |
| `gitty --v` | Show version |

---

### 🔗 Command Chaining with `and`

Run multiple commands in sequence:

```
gitty add . and push->main
gitty install and auth and add repo "my-project" and add . and push->main
```

---

### 🔧 Command Details

<details>
<summary><code>gitty install</code></summary>

Installs Git and GitHub CLI if they are not present. Downloads directly — no winget, no admin rights. Adds the gitty folder to User PATH.

</details>

<details>
<summary><code>gitty auth</code></summary>

Runs `gh auth login` to authenticate with GitHub. Required before creating repositories.

</details>

<details>
<summary><code>gitty init "url"</code></summary>

Initializes git in the current folder and sets the remote origin to the given URL. The URL must be in quotes.

```
gitty init "https://github.com/user/repo.git"
```

</details>

<details>
<summary><code>gitty add repo "name"</code></summary>

Creates a GitHub repository under your account. Private by default. Use `--public` to make it public.

If the folder is already linked to another repo, gitty will ask:
- **[1] Replace** — relink the folder to the new repo
- **[2] Keep** — create on GitHub, do not touch the local folder
- **[3] Cancel**

Requires: `gitty auth`

</details>

<details>
<summary><code>gitty add branch "name"</code></summary>

Creates a local branch (quotes required). You stay on your current branch — no switching occurs. If there are no commits yet, an initial commit is created automatically.

</details>

<details>
<summary><code>gitty add .</code></summary>

Stages all changes and creates a commit:

1. `git add .` — stage everything (`.gitignore` is respected)
2. `git commit -m "gitty auto-sync [UTC time]"`

If there are no commits yet, an initial commit is created automatically.

</details>

<details>
<summary><code>gitty push->branch</code></summary>

Pushes committed changes to the specified remote branch. Run `gitty add .` first to stage and commit.

```
gitty push->main
gitty push->dev
gitty push->feature/login
```

If the branch doesn't exist on the remote, gitty will offer to create it.

</details>

<details>
<summary><code>gitty pull~branch</code></summary>

Three modes:

```
gitty pull~main              # Safe: copies only files missing locally
gitty pull~staging --hard    # Overwrites files from remote; unique local files are kept
gitty pull~main --hard-reset # DESTRUCTIVE: mirrors remote exactly, deletes local-only files
```

`--hard-reset` requires confirmation.

</details>

<details>
<summary><code>gitty reset~branch</code></summary>

Wipes all content and commit history from the specified branch. An empty orphan commit is created in its place. If the branch exists on the remote, a force push is performed. Requires arrow-key confirmation: **Yes / No**.

```
gitty reset~second
gitty reset~old-feature
```

</details>

<details>
<summary><code>gitty gitignore</code></summary>

Interactive search across official `.gitignore` templates from GitHub. Type a name — the list filters live. Use `↑↓` to navigate, `Enter` to download and save `.gitignore` to the current folder.

</details>

---

### 🌐 Proxy Support

Works for all network operations including git, gh, and direct downloads:

```
gitty <command> --proxy "http://ip:port"
gitty <command> --proxy "http://user:pass@ip:port"
```

---

### 🛠 Build from Source

Requires Go 1.21+, no external dependencies.

```
git clone https://github.com/Omibranch/gitty.git
cd gitty
go build -ldflags="-s -w" -o gitty.exe .
```

---

### 📤 Output Prefixes

| Prefix | Meaning |
|---|---|
| `[SUCCESS]` | Operation completed successfully |
| `[ERROR]` | Something went wrong |
| `[HINT]` | Helpful tip or suggestion |

---

### ⚙️ Flags Reference

| Flag | Usage |
|---|---|
| `->` | Push TO a branch — `gitty push->main` |
| `~` | Pull FROM a branch — `gitty pull~main` |
| `and` | Chain commands — `gitty add . and push->main` |
| `--public` | Create a public repo — `gitty add repo "name" --public` |
| `--proxy` | Use a proxy — `gitty <cmd> --proxy "http://ip:port"` |
| `--hard` | Overwrite on pull — `gitty pull~main --hard` |
| `--hard-reset` | Mirror remote on pull — `gitty pull~main --hard-reset` |

---

### 📄 License

MIT © [Omibranch](https://github.com/Omibranch)

---
---

<div align="center">

[EN](#-guide-en) | **RU**

</div>

## 📖 Руководство (RU)

### Что такое gitty?

**gitty** — CLI-инструмент для Windows, который оборачивает Git и GitHub CLI в простые, человекопонятные команды.  
Главный принцип: **вы никогда не вводите `git commit`**. gitty берёт стейджинг, коммит и пуш на себя.

`->` означает **ТУДА** (отправить в ветку) &nbsp;|&nbsp; `~` означает **ОТТУДА** (получить из ветки)

---

### ⚡ Быстрая установка

**Вариант 1 — PowerShell одной строкой (рекомендуется)**
```powershell
irm https://raw.githubusercontent.com/Omibranch/gitty/main/install.ps1 | iex
```
Скачивает `gitty.exe`, добавляет в PATH пользователя, затем запускает `gitty install` для автоматической настройки Git и GitHub CLI. Права администратора не нужны.

**Вариант 2 — winget**
```
winget install Omibranch.Gitty
```

**Вариант 3 — Вручную**  
Скачайте `gitty.exe` из [Releases](https://github.com/Omibranch/gitty/releases), положите в любую папку из PATH, затем выполните:
```
gitty install
```

---

### 📋 Команды

| Команда | Описание |
|---|---|
| `gitty install` | Установить Git + GitHub CLI, добавить gitty в PATH |
| `gitty auth` | Войти в GitHub через `gh auth login` |
| `gitty init "url"` | Инициализировать репо и привязать remote origin |
| `gitty add repo "название"` | Создать приватный репо на GitHub и привязать папку |
| `gitty add repo "название" --public` | То же, но публичный |
| `gitty add branch "название"` | Создать локальную ветку без переключения на неё |
| `gitty add .` | Стейджинг всех изменений и авто-коммит |
| `gitty push->ветка` | Отправить закоммиченные изменения в remote-ветку |
| `gitty pull~ветка` | Безопасный pull — добавляет только отсутствующие файлы |
| `gitty pull~ветка --hard` | Перезаписать файлы с remote (уникальные локальные сохраняются) |
| `gitty pull~ветка --hard-reset` | Зеркало remote — деструктивно, требует подтверждения |
| `gitty reset~ветка` | Удалить всё содержимое и историю ветки |
| `gitty status` | Показать привязанный репо, ветку и аккаунт GitHub |
| `gitty gitignore` | Интерактивный выбор шаблонов `.gitignore` |
| `gitty clear` | Очистить экран терминала |
| `gitty help` | Полное руководство (EN / RU) |
| `gitty --v` | Показать версию |

---

### 🔗 Цепочки команд через `and`

```
gitty add . and push->main
gitty install and auth and add repo "мой-проект" and add . and push->main
```

---

### 🔧 Детали команд

<details>
<summary><code>gitty install</code></summary>

Устанавливает Git и GitHub CLI, если их нет. Скачивает напрямую — без winget и прав администратора. Добавляет папку с gitty.exe в PATH пользователя.

</details>

<details>
<summary><code>gitty auth</code></summary>

Запускает `gh auth login` для входа в GitHub. Обязателен перед созданием репозиториев.

</details>

<details>
<summary><code>gitty init "url"</code></summary>

Инициализирует git в текущей папке и устанавливает remote origin на указанный URL. URL нужно взять в кавычки.

```
gitty init "https://github.com/user/repo.git"
```

</details>

<details>
<summary><code>gitty add repo "название"</code></summary>

Создаёт репозиторий на GitHub под вашим аккаунтом. По умолчанию приватный. Добавьте `--public` для публичного.

Если папка уже привязана к другому репо, gitty предложит:
- **[1] Заменить** — перепривязать папку к новому репо
- **[2] Оставить** — создать на GitHub, папку не трогать
- **[3] Отмена**

Требует: `gitty auth`

</details>

<details>
<summary><code>gitty add branch "название"</code></summary>

Создаёт локальную ветку (кавычки обязательны). Вы остаётесь на текущей ветке — переключения не происходит. Если коммитов ещё нет, начальный коммит создаётся автоматически.

</details>

<details>
<summary><code>gitty add .</code></summary>

Стейджит все изменения и создаёт коммит:

1. `git add .` — добавить всё (`.gitignore` учитывается)
2. `git commit -m "gitty auto-sync [UTC время]"`

Если коммитов ещё нет, начальный коммит создаётся автоматически.

</details>

<details>
<summary><code>gitty push->ветка</code></summary>

Отправляет закоммиченные изменения в указанную ветку на remote. Сначала выполните `gitty add .` для стейджинга и коммита.

```
gitty push->main
gitty push->dev
gitty push->feature/login
```

Если ветки нет на remote, gitty предложит её создать.

</details>

<details>
<summary><code>gitty pull~ветка</code></summary>

Три режима:

```
gitty pull~main              # Безопасный: копирует только отсутствующие файлы
gitty pull~staging --hard    # Перезаписывает файлы с remote; уникальные локальные сохраняются
gitty pull~main --hard-reset # ДЕСТРУКТИВНО: зеркало remote, локальные уникальные файлы удаляются
```

`--hard-reset` требует подтверждения.

</details>

<details>
<summary><code>gitty reset~ветка</code></summary>

Удаляет всё содержимое и историю коммитов указанной ветки. Создаётся пустой сиротский коммит. Если ветка есть на remote — выполняется принудительный push. Требует подтверждения стрелками: **Yes / No**.

```
gitty reset~second
gitty reset~old-feature
```

</details>

<details>
<summary><code>gitty gitignore</code></summary>

Интерактивный поиск по официальным шаблонам `.gitignore` с GitHub. Вводите название — список фильтруется на лету. `↑↓` — навигация, `Enter` — скачать и сохранить `.gitignore` в папку.

</details>

---

### 🌐 Поддержка прокси

Работает для всех сетевых операций — git, gh и прямые загрузки:

```
gitty <команда> --proxy "http://ip:port"
gitty <команда> --proxy "http://user:pass@ip:port"
```

---

### 🛠 Сборка из исходников

Требуется Go 1.21+, внешние зависимости отсутствуют.

```
git clone https://github.com/Omibranch/gitty.git
cd gitty
go build -ldflags="-s -w" -o gitty.exe .
```

---

### 📤 Префиксы вывода

| Префикс | Значение |
|---|---|
| `[SUCCESS]` | Операция выполнена успешно |
| `[ERROR]` | Произошла ошибка |
| `[HINT]` | Подсказка или совет |

---

### ⚙️ Флаги и синтаксис

| Флаг | Использование |
|---|---|
| `->` | Отправить В ветку — `gitty push->main` |
| `~` | Получить ИЗ ветки — `gitty pull~main` |
| `and` | Цепочка команд — `gitty add . and push->main` |
| `--public` | Публичный репо — `gitty add repo "название" --public` |
| `--proxy` | Прокси — `gitty <команда> --proxy "http://ip:port"` |
| `--hard` | Перезапись при pull — `gitty pull~main --hard` |
| `--hard-reset` | Зеркало remote — `gitty pull~main --hard-reset` |

---

### 📄 Лицензия

MIT © [Omibranch](https://github.com/Omibranch)
