<div align="center"><div align="center">



<img src="https://files.catbox.moe/sjgj3f.png" alt="gitty" width="1280" height="720" />



**gitty**



[EN](#-guide-en) | [RU](#руководство-ru)



------



*Zero-commit Git workflow for humans.*  *Zero-commit Git workflow for humans.*  

Stage. Commit. Push. All in one command.Stage. Commit. Push. — All in one command.



[![Release](https://img.shields.io/github/v/release/Omibranch/gitty?color=black&style=flat-square)](https://github.com/Omibranch/gitty/releases)[![Release](https://img.shields.io/github/v/release/Omibranch/gitty?color=black&style=flat-square)](https://github.com/Omibranch/gitty/releases)

[![License: MIT](https://img.shields.io/badge/license-MIT-black?style=flat-square)](LICENSE)[![License: MIT](https://img.shields.io/badge/license-MIT-black?style=flat-square)](LICENSE)

[![Platform](https://img.shields.io/badge/platform-Windows-black?style=flat-square)](https://github.com/Omibranch/gitty/releases)[![Platform](https://img.shields.io/badge/platform-Windows-black?style=flat-square)](https://github.com/Omibranch/gitty/releases)



</div></div>



------



## Guide (EN)## Guide (EN)



What is gitty?



**gitty** is a Windows CLI tool that wraps Git and GitHub CLI into simple, human-friendly commands.  **gitty** is a Windows CLI tool that wraps Git and GitHub CLI into simple, human-friendly commands.  

The core idea: **you never type `git commit`**. gitty handles staging, committing, and pushing automatically.The core idea: **you never type `git commit`**. gitty handles staging, committing, and pushing automatically.



| Symbol | Meaning | Alias |`%` means **TO** (push to a branch) &nbsp;|&nbsp; `~` means **FROM** (pull from a branch)

|---|---|---|

| `%` | TO (push destination) | `to` |---

| `~` | FROM (pull source) | `from` |

| `#` | IN (branch specifier) | `in` |### Quick Install



Semantic aliases let you write natural English:  **Option 1 — PowerShell one-liner (recommended)**

`gitty push to main` = `gitty push main` = `gitty push%main`  ```powershell

`gitty pull from dev` = `gitty pull dev` = `gitty pull~dev`  irm https://raw.githubusercontent.com/Omibranch/gitty/main/install.ps1 | iex

`gitty checkpoint "v1" in main` = `gitty checkpoint "v1"#main````

Downloads `gitty.exe`, adds it to your User PATH, then runs `gitty install` to set up Git and GitHub CLI automatically. No admin rights required.

---

**Option 2 — winget**

### Quick Install```

winget install Omibranch.Gitty

**Option 1 — PowerShell one-liner (recommended)**```

```powershell

irm https://raw.githubusercontent.com/Omibranch/gitty/main/install.ps1 | iex**Option 3 — Manual**  

```Download `gitty.exe` from [Releases](https://github.com/Omibranch/gitty/releases), place it anywhere in your PATH, then run:

Downloads `gitty.exe`, adds it to your User PATH, then runs `gitty install` to set up Git and GitHub CLI automatically. No admin rights required.```

gitty install

**Option 2 — winget**```

```

winget install Omibranch.Gitty---

```

### Commands

**Option 3 — Manual**  

Download `gitty.exe` from [Releases](https://github.com/Omibranch/gitty/releases), place it anywhere in your PATH, then run:| Command | Description |

```|---|---|

gitty install| `gitty install` | Install Git + GitHub CLI, add gitty to PATH (no admin needed) |

```| `gitty auth` | Sign in to GitHub via `gh auth login` |

| `gitty init "url"` | Init a git repo and set remote origin |

---| `gitty add repo "name"` | Create a private GitHub repo and link current folder |

| `gitty add repo "name" --public` | Same, but public |

### Commands| `gitty add branch "name"` | Create a local branch without switching to it |

| `gitty add .` | Stage all changes and auto-commit |

| Command | Description || `gitty push branch` | Push committed changes to a remote branch (`push%branch` also works) |

|---|---|| `gitty pull branch` | Safe pull — adds only files missing locally (`pull~branch` still works) |

| `gitty install` | Install Git + GitHub CLI, add gitty to PATH (no admin needed) || `gitty pull branch --hard` | Overwrite local files from remote (keeps unique local files) |

| `gitty auth` | Sign in to GitHub via `gh auth login` || `gitty pull branch --hard-reset` | Mirror remote exactly — destructive, requires confirmation |

| `gitty init "url"` | Init a git repo and set remote origin || `gitty reset~branch` | Wipe all content and history from a branch |

| `gitty add repo "name"` | Create a private GitHub repo and link current folder || `gitty migration branch1%branch2` | Replace all files in `branch1` with files from `branch2` (with Yes/No confirmation) |

| `gitty add repo "name" --public` | Same, but public || `gitty status` | Show linked repo, current branch, and GitHub account |

| `gitty add branch "name"` | Create a local branch without switching to it || `gitty gitignore` | Interactive picker for `.gitignore` templates |

| `gitty add .` | Stage all changes and auto-commit || `gitty clear` | Clear the terminal screen |

| `gitty push branch` | Push committed changes to a remote branch || `gitty help` | Full manual (EN / RU) |

| `gitty push branch --share` | Push + copy GitHub branch link to clipboard || `gitty --v` | Show version |

| `gitty pull branch` | Safe pull — adds only files missing locally |

| `gitty pull branch --hard` | Overwrite local files from remote (keeps unique local files) |---

| `gitty pull branch --hard-reset` | Mirror remote exactly — destructive, requires confirmation |

| `gitty reset~branch` | Wipe all content and history from a branch |### Command Chaining with `and`

| `gitty migration branch1%branch2` | Replace all files in `branch1` with files from `branch2` |

| `gitty undo` | Undo last commit, keep changes staged |Run multiple commands in sequence:

| `gitty log` | Show git log (last week by default) |

| `gitty log --3day` | Show commits from last 3 days (`--Nh`, `--Nday`, `--Nweek`, `--Nmonth`) |```

| `gitty checkpoint "name" in branch` | Create a named checkpoint (git tag) on a branch |gitty add . and push main

| `gitty restore "name"` | Revert working tree to a checkpoint |gitty install and auth and add repo "my-project" and add . and push main

| `gitty rename branch "old"%"new"` | Rename a branch locally and on remote |```

| `gitty rename repo "new"` | Rename the current linked repository on GitHub |

| `gitty rename repo "old"%"new"` | Rename a specific repository on GitHub |---

| `gitty status` | Show linked repo, current branch, and GitHub account |

| `gitty gitignore` | Interactive picker for `.gitignore` templates |### Command Details

| `gitty clear` | Clear the terminal screen |

| `gitty help` | Full manual (EN / RU) |<details>

| `gitty --v` | Show version |<summary><code>gitty install</code></summary>



---Installs Git and GitHub CLI if they are not present. Downloads directly — no winget, no admin rights. Adds the gitty folder to User PATH.



### Semantic Aliases</details>



You can use natural words instead of symbols in any command:<details>

<summary><code>gitty auth</code></summary>

| Alias | Symbol | Example |

|---|---|---|Runs `gh auth login` to authenticate with GitHub. Required before creating repositories.

| `to` | `%` | `gitty push to main` |

| `from` | `~` | `gitty pull from dev` |</details>

| `in` | `#` | `gitty checkpoint "v1" in main` |

| `to` (migration) | `%` | `gitty migration main to develop` |<details>

<summary><code>gitty init "url"</code></summary>

---

Initializes git in the current folder and sets the remote origin to the given URL. The URL must be in quotes.

### Command Chaining with `and`

```

Run multiple commands in sequence:gitty init "https://github.com/user/repo.git"

```

```

gitty add . and push main</details>

gitty install and auth and add repo "my-project" and add . and push main

```<details>

<summary><code>gitty add repo "name"</code></summary>

---

Creates a GitHub repository under your account. Private by default. Use `--public` to make it public.

### Command Details

If the folder is already linked to another repo, gitty will ask:

<details>- **[1] Replace** — relink the folder to the new repo

<summary><code>gitty install</code></summary>- **[2] Keep** — create on GitHub, do not touch the local folder

- **[3] Cancel**

Installs Git and GitHub CLI if they are not present. Downloads directly — no winget, no admin rights. Adds the gitty folder to User PATH.

Requires: `gitty auth`

</details>

</details>

<details>

<summary><code>gitty auth</code></summary><details>

<summary><code>gitty add branch "name"</code></summary>

Runs `gh auth login` to authenticate with GitHub. Required before creating repositories.

Creates a local branch (quotes required). You stay on your current branch — no switching occurs. If there are no commits yet, an initial commit is created automatically.

</details>

</details>

<details>

<summary><code>gitty init "url"</code></summary><details>

<summary><code>gitty add .</code></summary>

Initializes git in the current folder and sets the remote origin to the given URL. The URL must be in quotes.

Stages all changes and creates a commit:

```

gitty init "https://github.com/user/repo.git"1. `git add .` — stage everything (`.gitignore` is respected)

```2. `git commit -m "gitty auto-sync [UTC time]"`



</details>If there are no commits yet, an initial commit is created automatically.



<details>If git encounters nested repository folders without commits (for example `proxy/`),

<summary><code>gitty add repo "name"</code></summary>gitty retries automatically and stages everything else, skipping only those paths.



Creates a GitHub repository under your account. Private by default. Use `--public` to make it public.</details>



If the folder is already linked to another repo, gitty will ask:<details>

- **[1] Replace** — relink the folder to the new repo<summary><code>gitty push branch</code></summary>

- **[2] Keep** — create on GitHub, do not touch the local folder

- **[3] Cancel**Pushes committed changes to the specified remote branch. Run `gitty add .` first to stage and commit.



Requires: `gitty auth`Preferred syntax (works without quotes in cmd/PowerShell):



</details>```

gitty push main

<details>gitty push dev

<summary><code>gitty add branch "name"</code></summary>gitty push feature/login

```

Creates a local branch (quotes required). You stay on your current branch — no switching occurs. If there are no commits yet, an initial commit is created automatically.

Compact syntax is also supported:

</details>

```

<details>gitty push%main

<summary><code>gitty add .</code></summary>```



Stages all changes and creates a commit:If the branch doesn't exist on the remote, gitty will offer to create it.



1. `git add .` — stage everything (`.gitignore` is respected)</details>

2. `git commit -m "gitty auto-sync [UTC time]"`

<details>

If there are no commits yet, an initial commit is created automatically.<summary><code>gitty pull branch</code></summary>



If git encounters nested repository folders without commits (for example `proxy/`), gitty retries automatically and stages everything else, skipping only those paths.Three modes:



</details>```

gitty pull main              # Safe: copies only files missing locally

<details>gitty pull main --hard    # Overwrites files from remote; unique local files are kept

<summary><code>gitty push branch</code></summary>gitty pull main --hard-reset # DESTRUCTIVE: mirrors remote exactly, deletes local-only files

```

Pushes committed changes to the specified remote branch. Run `gitty add .` first to stage and commit.

Legacy syntax (`pull~branch`) is still supported.

Preferred syntax:

`--hard-reset` requires confirmation.

```

gitty push main</details>

gitty push to main

gitty push dev<details>

gitty push feature/login<summary><code>gitty reset~branch</code></summary>

gitty push%main

```Wipes all content and commit history from the specified branch. An empty orphan commit is created in its place. If the branch exists on the remote, a force push is performed. Requires arrow-key confirmation: **Yes / No**.



If the branch doesn't exist on the remote, gitty will offer to create it.```

gitty reset~second

</details>gitty reset~old-feature

```

<details>

<summary><code>gitty push branch --share</code></summary></details>



Pushes the branch and automatically copies the GitHub branch URL to the clipboard.<details>

<summary><code>gitty migration branch1%branch2</code></summary>

```

gitty push main --shareDeletes all files in `branch1` and replaces them with files from `branch2`.

```Before running, gitty shows arrow-key confirmation: **Yes / No**.



</details>```

gitty migration main%develop

<details>```

<summary><code>gitty pull branch</code></summary>

</details>

Three modes:

<details>

```<summary><code>gitty gitignore</code></summary>

gitty pull main              # Safe: copies only files missing locally

gitty pull staging --hard    # Overwrites files from remote; unique local files are keptInteractive search across official `.gitignore` templates from GitHub. Type a name — the list filters live. Use `↑↓` to navigate, `Enter` to download and save `.gitignore` to the current folder.

gitty pull main --hard-reset # DESTRUCTIVE: mirrors remote exactly, deletes local-only files

```</details>



Aliases:---

```

gitty pull from main### Proxy Support

gitty pull~main

```Works for all network operations including git, gh, and direct downloads:



`--hard-reset` requires confirmation.```

gitty <command> --proxy "http://ip:port"

</details>gitty <command> --proxy "http://user:pass@ip:port"

```

<details>

<summary><code>gitty reset~branch</code></summary>---



Wipes all content and commit history from the specified branch. An empty orphan commit is created in its place. If the branch exists on the remote, a force push is performed. Requires arrow-key confirmation: **Yes / No**.### Build from Source



```Requires Go 1.21+, no external dependencies.

gitty reset~second

gitty reset~old-feature```

```git clone https://github.com/Omibranch/gitty.git

cd gitty

</details>go build -ldflags="-s -w" -o gitty.exe .

```

<details>

<summary><code>gitty migration branch1%branch2</code></summary>---



Deletes all files in `branch1` and replaces them with files from `branch2`. Before running, gitty shows arrow-key confirmation: **Yes / No**.### Output Prefixes



```| Prefix | Meaning |

gitty migration main%develop|---|---|

gitty migration main to develop| `[SUCCESS]` | Operation completed successfully |

```| `[ERROR]` | Something went wrong |

| `[HINT]` | Helpful tip or suggestion |

</details>

---

<details>

<summary><code>gitty undo</code></summary>### Flags Reference



Reverts the last commit, keeping all changes staged and ready to re-commit.  | Flag | Usage |

Equivalent to `git reset HEAD~1 --soft`.|---|---|

| `push` | Push TO a branch — `gitty push main` |

```| `pull` | Pull FROM a branch — `gitty pull main` |

gitty undo| `and` | Chain commands — `gitty add . and push main` |

```| `--public` | Create a public repo — `gitty add repo "name" --public` |

| `--proxy` | Use a proxy — `gitty <cmd> --proxy "http://ip:port"` |

</details>| `--hard` | Overwrite on pull — `gitty pull main --hard` |

| `--hard-reset` | Mirror remote on pull — `gitty pull main --hard-reset` |

<details>| `%` | Compact push separator and migration separator — `gitty push%main`, `gitty migration main%develop` |

<summary><code>gitty log</code></summary>

---

Shows a formatted git log with graph and branch decorations. Default range: last 1 week.

### License

```

gitty logMIT © [Omibranch](https://github.com/Omibranch)

gitty log --1h        # last 1 hour

gitty log --3day      # last 3 days---

gitty log --2week     # last 2 weeks---

gitty log --1month    # last 1 month

```<div align="center">



</details>[EN](#-guide-en) | **RU**



<details></div>

<summary><code>gitty checkpoint "name" in branch</code></summary>

## Руководство (RU)

Creates a git tag `name` pointing at the tip of `branch` and pushes it to origin. Use checkpoints to mark stable states before experimenting.

### Что такое gitty?

```

gitty checkpoint "v1-stable" in main**gitty** — CLI-инструмент для Windows, который оборачивает Git и GitHub CLI в простые, человекопонятные команды.  

gitty checkpoint "before-refactor"Главный принцип: **вы никогда не вводите `git commit`**. gitty берёт стейджинг, коммит и пуш на себя.

gitty checkpoint "v1-stable"#main

````%` означает **ТУДА** (отправить в ветку) &nbsp;|&nbsp; `~` означает **ОТТУДА** (получить из ветки)



</details>---



<details>### Быстрая установка

<summary><code>gitty restore "name"</code></summary>

**Вариант 1 — PowerShell одной строкой (рекомендуется)**

Reverts the working tree to a previously created checkpoint (tag). Puts the repository into detached HEAD state at that tag. Requires arrow-key confirmation: **Yes / No**.```powershell

irm https://raw.githubusercontent.com/Omibranch/gitty/main/install.ps1 | iex

``````

gitty restore "v1-stable"Скачивает `gitty.exe`, добавляет в PATH пользователя, затем запускает `gitty install` для автоматической настройки Git и GitHub CLI. Права администратора не нужны.

```

**Вариант 2 — winget**

After restoring, to continue developing:```

```winget install Omibranch.Gitty

git checkout -b new-branch-name```

```

**Вариант 3 — Вручную**  

</details>Скачайте `gitty.exe` из [Releases](https://github.com/Omibranch/gitty/releases), положите в любую папку из PATH, затем выполните:

```

<details>gitty install

<summary><code>gitty rename branch "old"%"new"</code></summary>```



Renames a branch both locally and on the remote. Deletes the old remote branch name and pushes the new one.---



```### Команды

gitty rename branch "feature-x"%"feature-login"

```| Команда | Описание |

|---|---|

</details>| `gitty install` | Установить Git + GitHub CLI, добавить gitty в PATH |

| `gitty auth` | Войти в GitHub через `gh auth login` |

<details>| `gitty init "url"` | Инициализировать репо и привязать remote origin |

<summary><code>gitty rename repo "new"</code></summary>| `gitty add repo "название"` | Создать приватный репо на GitHub и привязать папку |

| `gitty add repo "название" --public` | То же, но публичный |

Renames the GitHub repository linked to the current folder. If the current folder's remote is affected, the local remote URL is updated automatically.| `gitty add branch "название"` | Создать локальную ветку без переключения на неё |

| `gitty add .` | Стейджинг всех изменений и авто-коммит |

```| `gitty push ветка` | Отправить закоммиченные изменения в remote-ветку (`push%ветка` тоже поддерживается) |

gitty rename repo "my-new-name"| `gitty pull ветка` | Безопасный pull — добавляет только отсутствующие файлы (`pull~ветка` тоже поддерживается) |

gitty rename repo "old-project"%"new-project"| `gitty pull ветка --hard` | Перезаписать файлы с remote (уникальные локальные сохраняются) |

```| `gitty pull ветка --hard-reset` | Зеркало remote — деструктивно, требует подтверждения |

| `gitty reset~ветка` | Удалить всё содержимое и историю ветки |

</details>| `gitty migration ветка1%ветка2` | Полностью заменить файлы в `ветка1` файлами из `ветка2` (с подтверждением Yes/No) |

| `gitty status` | Показать привязанный репо, ветку и аккаунт GitHub |

<details>| `gitty gitignore` | Интерактивный выбор шаблонов `.gitignore` |

<summary><code>gitty gitignore</code></summary>| `gitty clear` | Очистить экран терминала |

| `gitty help` | Полное руководство (EN / RU) |

Interactive search across official `.gitignore` templates from GitHub. Type a name — the list filters live. Use `↑↓` to navigate, `Enter` to download and save `.gitignore` to the current folder.| `gitty --v` | Показать версию |



</details>---



---### Цепочки команд через `and`



### Proxy Support```

gitty add . and push main

Works for all network operations including git, gh, and direct downloads:gitty install and auth and add repo "мой-проект" and add . and push main

```

```

gitty <command> --proxy "http://ip:port"---

gitty <command> --proxy "http://user:pass@ip:port"

```### Детали команд



---<details>

<summary><code>gitty install</code></summary>

### Build from Source

Устанавливает Git и GitHub CLI, если их нет. Скачивает напрямую — без winget и прав администратора. Добавляет папку с gitty.exe в PATH пользователя.

Requires Go 1.21+, no external dependencies.

</details>

```

git clone https://github.com/Omibranch/gitty.git<details>

cd gitty/source<summary><code>gitty auth</code></summary>

go build -ldflags="-s -w" -o ../gitty.exe .

```Запускает `gh auth login` для входа в GitHub. Обязателен перед созданием репозиториев.



---</details>



### Output Prefixes<details>

<summary><code>gitty init "url"</code></summary>

| Prefix | Meaning |

|---|---|Инициализирует git в текущей папке и устанавливает remote origin на указанный URL. URL нужно взять в кавычки.

| `[SUCCESS]` | Operation completed successfully |

| `[ERROR]` | Something went wrong |```

| `[HINT]` | Helpful tip or suggestion |gitty init "https://github.com/user/repo.git"

```

---

</details>

### Flags Reference

<details>

| Flag / Symbol | Usage |<summary><code>gitty add repo "название"</code></summary>

|---|---|

| `%` / `to` | Push TO a branch — `gitty push main` / `gitty push to main` |Создаёт репозиторий на GitHub под вашим аккаунтом. По умолчанию приватный. Добавьте `--public` для публичного.

| `~` / `from` | Pull FROM a branch — `gitty pull main` / `gitty pull from main` |

| `#` / `in` | Branch IN checkpoint — `gitty checkpoint "v1" in main` |Если папка уже привязана к другому репо, gitty предложит:

| `and` | Chain commands — `gitty add . and push main` |- **[1] Заменить** — перепривязать папку к новому репо

| `--public` | Create a public repo — `gitty add repo "name" --public` |- **[2] Оставить** — создать на GitHub, папку не трогать

| `--share` | Push + copy link — `gitty push main --share` |- **[3] Отмена**

| `--proxy` | Use a proxy — `gitty <cmd> --proxy "http://ip:port"` |

| `--hard` | Overwrite on pull — `gitty pull main --hard` |Требует: `gitty auth`

| `--hard-reset` | Mirror remote on pull — `gitty pull main --hard-reset` |

</details>

---

<details>

### License<summary><code>gitty add branch "название"</code></summary>



MIT © [Omibranch](https://github.com/Omibranch)Создаёт локальную ветку (кавычки обязательны). Вы остаётесь на текущей ветке — переключения не происходит. Если коммитов ещё нет, начальный коммит создаётся автоматически.



---</details>

---

<details>

<div align="center"><summary><code>gitty add .</code></summary>



[EN](#guide-en) | **RU**Стейджит все изменения и создаёт коммит:



</div>1. `git add .` — добавить всё (`.gitignore` учитывается)

2. `git commit -m "gitty auto-sync [UTC время]"`

## Руководство (RU)

Если коммитов ещё нет, начальный коммит создаётся автоматически.

### Что такое gitty?

Если git встречает вложенные папки-репозитории без коммитов (например `proxy/`),

**gitty** — CLI-инструмент для Windows, который оборачивает Git и GitHub CLI в простые, человекопонятные команды.  gitty автоматически повторяет стейджинг и пропускает только эти пути.

Главный принцип: **вы никогда не вводите `git commit`**. gitty берёт стейджинг, коммит и пуш на себя.

</details>

| Символ | Значение | Псевдоним |

|---|---|---|<details>

| `%` | ТУДА (push destination) | `to` |<summary><code>gitty push ветка</code></summary>

| `~` | ОТТУДА (pull source) | `from` |

| `#` | В (branch specifier) | `in` |Отправляет закоммиченные изменения в указанную ветку на remote. Сначала выполните `gitty add .` для стейджинга и коммита.



Псевдонимы позволяют писать естественным языком:  Предпочтительный синтаксис (работает без кавычек в cmd/PowerShell):

`gitty push to main` = `gitty push main` = `gitty push%main`  

`gitty pull from dev` = `gitty pull dev` = `gitty pull~dev`  ```

`gitty checkpoint "v1" in main` = `gitty checkpoint "v1"#main`gitty push main

gitty push dev

---gitty push feature/login

```

### Быстрая установка

Короткий синтаксис тоже поддерживается:

**Вариант 1 — PowerShell одной строкой (рекомендуется)**

```powershell```

irm https://raw.githubusercontent.com/Omibranch/gitty/main/install.ps1 | iexgitty push%main

``````

Скачивает `gitty.exe`, добавляет в PATH пользователя, затем запускает `gitty install` для автоматической настройки Git и GitHub CLI. Права администратора не нужны.

Если ветки нет на remote, gitty предложит её создать.

**Вариант 2 — winget**

```</details>

winget install Omibranch.Gitty

```<details>

<summary><code>gitty pull ветка</code></summary>

**Вариант 3 — Вручную**  

Скачайте `gitty.exe` из [Releases](https://github.com/Omibranch/gitty/releases), положите в любую папку из PATH, затем выполните:Три режима:

```

gitty install```

```gitty pull main              # Безопасный: копирует только отсутствующие файлы

gitty pull main --hard    # Перезаписывает файлы с remote; уникальные локальные сохраняются

---gitty pull main --hard-reset # ДЕСТРУКТИВНО: зеркало remote, локальные уникальные файлы удаляются

```

### Команды

Старый синтаксис (`pull~ветка`) тоже поддерживается.

| Команда | Описание |

|---|---|`--hard-reset` требует подтверждения.

| `gitty install` | Установить Git + GitHub CLI, добавить gitty в PATH |

| `gitty auth` | Войти в GitHub через `gh auth login` |</details>

| `gitty init "url"` | Инициализировать репо и привязать remote origin |

| `gitty add repo "название"` | Создать приватный репо на GitHub и привязать папку |<details>

| `gitty add repo "название" --public` | То же, но публичный |<summary><code>gitty reset~ветка</code></summary>

| `gitty add branch "название"` | Создать локальную ветку без переключения на неё |

| `gitty add .` | Стейджинг всех изменений и авто-коммит |Удаляет всё содержимое и историю коммитов указанной ветки. Создаётся пустой сиротский коммит. Если ветка есть на remote — выполняется принудительный push. Требует подтверждения стрелками: **Yes / No**.

| `gitty push ветка` | Отправить закоммиченные изменения в remote-ветку |

| `gitty push ветка --share` | Пуш + скопировать ссылку на ветку в буфер обмена |```

| `gitty pull ветка` | Безопасный pull — добавляет только отсутствующие файлы |gitty reset~second

| `gitty pull ветка --hard` | Перезаписать файлы с remote (уникальные локальные сохраняются) |gitty reset~old-feature

| `gitty pull ветка --hard-reset` | Зеркало remote — деструктивно, требует подтверждения |```

| `gitty reset~ветка` | Удалить всё содержимое и историю ветки |

| `gitty migration ветка1%ветка2` | Заменить файлы в `ветка1` файлами из `ветка2` |</details>

| `gitty undo` | Отменить последний коммит, сохранив изменения в staged |

| `gitty log` | Показать git log (по умолчанию — последняя неделя) |<details>

| `gitty log --3day` | Показать коммиты за последние 3 дня (`--Nh`, `--Nday`, `--Nweek`, `--Nmonth`) |<summary><code>gitty migration ветка1%ветка2</code></summary>

| `gitty checkpoint "название" in ветка` | Создать именованный чекпоинт (git tag) на ветке |

| `gitty restore "название"` | Откатиться к чекпоинту |Удаляет все файлы из `ветка1` и заменяет их файлами из `ветка2`.

| `gitty rename branch "старое"%"новое"` | Переименовать ветку локально и на remote |Перед выполнением gitty показывает подтверждение стрелками: **Yes / No**.

| `gitty rename repo "новое"` | Переименовать текущий репозиторий на GitHub |

| `gitty rename repo "старое"%"новое"` | Переименовать конкретный репозиторий на GitHub |```

| `gitty status` | Показать привязанный репо, ветку и аккаунт GitHub |gitty migration main%develop

| `gitty gitignore` | Интерактивный выбор шаблонов `.gitignore` |```

| `gitty clear` | Очистить экран терминала |

| `gitty help` | Полное руководство (EN / RU) |</details>

| `gitty --v` | Показать версию |

<details>

---<summary><code>gitty gitignore</code></summary>



### Семантические псевдонимыИнтерактивный поиск по официальным шаблонам `.gitignore` с GitHub. Вводите название — список фильтруется на лету. `↑↓` — навигация, `Enter` — скачать и сохранить `.gitignore` в папку.



Вместо символов можно использовать обычные слова:</details>



| Псевдоним | Символ | Пример |---

|---|---|---|

| `to` | `%` | `gitty push to main` |### Поддержка прокси

| `from` | `~` | `gitty pull from dev` |

| `in` | `#` | `gitty checkpoint "v1" in main` |Работает для всех сетевых операций — git, gh и прямые загрузки:

| `to` (migration) | `%` | `gitty migration main to develop` |

```

---gitty <команда> --proxy "http://ip:port"

gitty <команда> --proxy "http://user:pass@ip:port"

### Цепочки команд через `and````



```---

gitty add . and push main

gitty install and auth and add repo "мой-проект" and add . and push main### Сборка из исходников

```

Требуется Go 1.21+, внешние зависимости отсутствуют.

---

```

### Детали командgit clone https://github.com/Omibranch/gitty.git

cd gitty

<details>go build -ldflags="-s -w" -o gitty.exe .

<summary><code>gitty install</code></summary>```



Устанавливает Git и GitHub CLI, если их нет. Скачивает напрямую — без winget и прав администратора. Добавляет папку с gitty.exe в PATH пользователя.---



</details>### Префиксы вывода



<details>| Префикс | Значение |

<summary><code>gitty auth</code></summary>|---|---|

| `[SUCCESS]` | Операция выполнена успешно |

Запускает `gh auth login` для входа в GitHub. Обязателен перед созданием репозиториев.| `[ERROR]` | Произошла ошибка |

| `[HINT]` | Подсказка или совет |

</details>

---

<details>

<summary><code>gitty init "url"</code></summary>### Флаги и синтаксис



Инициализирует git в текущей папке и устанавливает remote origin на указанный URL. URL нужно взять в кавычки.| Флаг | Использование |

|---|---|

```| `push` | Отправить В ветку — `gitty push main` |

gitty init "https://github.com/user/repo.git"| `pull` | Получить ИЗ ветки — `gitty pull main` |

```| `and` | Цепочка команд — `gitty add . and push main` |

| `--public` | Публичный репо — `gitty add repo "название" --public` |

</details>| `--proxy` | Прокси — `gitty <команда> --proxy "http://ip:port"` |

| `--hard` | Перезапись при pull — `gitty pull main --hard` |

<details>| `--hard-reset` | Зеркало remote — `gitty pull main --hard-reset` |

<summary><code>gitty add repo "название"</code></summary>| `%` | Разделитель в compact push и migration — `gitty push%main`, `gitty migration main%develop` |



Создаёт репозиторий на GitHub под вашим аккаунтом. По умолчанию приватный. Добавьте `--public` для публичного.---



Если папка уже привязана к другому репо, gitty предложит:### Лицензия

- **[1] Заменить** — перепривязать папку к новому репо

- **[2] Оставить** — создать на GitHub, папку не трогатьMIT © [Omibranch](https://github.com/Omibranch)

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

Если git встречает вложенные папки-репозитории без коммитов (например `proxy/`), gitty автоматически повторяет стейджинг и пропускает только эти пути.

</details>

<details>
<summary><code>gitty push ветка</code></summary>

Отправляет закоммиченные изменения в указанную ветку на remote. Сначала выполните `gitty add .` для стейджинга и коммита.

Предпочтительный синтаксис:

```
gitty push main
gitty push to main
gitty push dev
gitty push feature/login
gitty push%main
```

Если ветки нет на remote, gitty предложит её создать.

</details>

<details>
<summary><code>gitty push ветка --share</code></summary>

Пушит ветку и автоматически копирует ссылку на GitHub-ветку в буфер обмена.

```
gitty push main --share
```

</details>

<details>
<summary><code>gitty pull ветка</code></summary>

Три режима:

```
gitty pull main              # Безопасный: копирует только отсутствующие файлы
gitty pull staging --hard    # Перезаписывает файлы с remote; уникальные локальные сохраняются
gitty pull main --hard-reset # ДЕСТРУКТИВНО: зеркало remote, локальные уникальные файлы удаляются
```

Псевдонимы:
```
gitty pull from main
gitty pull~main
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
<summary><code>gitty migration ветка1%ветка2</code></summary>

Удаляет все файлы из `ветка1` и заменяет их файлами из `ветка2`. Перед выполнением gitty показывает подтверждение стрелками: **Yes / No**.

```
gitty migration main%develop
gitty migration main to develop
```

</details>

<details>
<summary><code>gitty undo</code></summary>

Отменяет последний коммит, оставляя все изменения в состоянии staged для повторного коммита.  
Эквивалентно `git reset HEAD~1 --soft`.

```
gitty undo
```

</details>

<details>
<summary><code>gitty log</code></summary>

Показывает форматированный git-лог с графом и метками веток. По умолчанию: последняя неделя.

```
gitty log
gitty log --1h        # последний 1 час
gitty log --3day      # последние 3 дня
gitty log --2week     # последние 2 недели
gitty log --1month    # последний 1 месяц
```

</details>

<details>
<summary><code>gitty checkpoint "название" in ветка</code></summary>

Создаёт git-тег `название` на кончике `ветки` и пушит его на origin. Используйте чекпоинты для фиксации стабильного состояния перед экспериментами.

```
gitty checkpoint "v1-stable" in main
gitty checkpoint "before-refactor"
gitty checkpoint "v1-stable"#main
```

</details>

<details>
<summary><code>gitty restore "название"</code></summary>

Откатывает рабочую директорию к ранее созданному чекпоинту (тегу). Репозиторий переходит в состояние detached HEAD. Требует подтверждения стрелками: **Yes / No**.

```
gitty restore "v1-stable"
```

После восстановления, чтобы продолжить разработку:
```
git checkout -b новое-название-ветки
```

</details>

<details>
<summary><code>gitty rename branch "старое"%"новое"</code></summary>

Переименовывает ветку локально и на remote. Старая ветка на remote удаляется, новая создаётся.

```
gitty rename branch "feature-x"%"feature-login"
```

</details>

<details>
<summary><code>gitty rename repo "новое"</code></summary>

Переименовывает репозиторий GitHub, привязанный к текущей папке. Если текущий remote затрагивается — локальный URL обновляется автоматически.

```
gitty rename repo "my-new-name"
gitty rename repo "old-project"%"new-project"
```

</details>

<details>
<summary><code>gitty gitignore</code></summary>

Интерактивный поиск по официальным шаблонам `.gitignore` с GitHub. Вводите название — список фильтруется на лету. `↑↓` — навигация, `Enter` — скачать и сохранить `.gitignore` в папку.

</details>

---

### Поддержка прокси

Работает для всех сетевых операций — git, gh и прямые загрузки:

```
gitty <команда> --proxy "http://ip:port"
gitty <команда> --proxy "http://user:pass@ip:port"
```

---

### Сборка из исходников

Требуется Go 1.21+, внешние зависимости отсутствуют.

```
git clone https://github.com/Omibranch/gitty.git
cd gitty/source
go build -ldflags="-s -w" -o ../gitty.exe .
```

---

### Префиксы вывода

| Префикс | Значение |
|---|---|
| `[SUCCESS]` | Операция выполнена успешно |
| `[ERROR]` | Произошла ошибка |
| `[HINT]` | Подсказка или совет |

---

### Флаги и синтаксис

| Флаг / Символ | Использование |
|---|---|
| `%` / `to` | Отправить В ветку — `gitty push main` / `gitty push to main` |
| `~` / `from` | Получить ИЗ ветки — `gitty pull main` / `gitty pull from main` |
| `#` / `in` | Ветка В checkpoint — `gitty checkpoint "v1" in main` |
| `and` | Цепочка команд — `gitty add . and push main` |
| `--public` | Публичный репо — `gitty add repo "название" --public` |
| `--share` | Пуш + скопировать ссылку — `gitty push main --share` |
| `--proxy` | Прокси — `gitty <команда> --proxy "http://ip:port"` |
| `--hard` | Перезапись при pull — `gitty pull main --hard` |
| `--hard-reset` | Зеркало remote — `gitty pull main --hard-reset` |

---

### Лицензия

MIT © [Omibranch](https://github.com/Omibranch)
