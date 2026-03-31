<div align="center">

  <img src="https://files.catbox.moe/sjgj3f.png" alt="gitty" width="1280" height="720" />

  <h3>gitty — minimal Git CLI that speaks human</h3>

  [![Windows](https://img.shields.io/badge/winget-Omibranch.Gitty-blue?logo=windows)](https://github.com/microsoft/winget-pkgs)
  [![Release](https://img.shields.io/github/v/release/Omibranch/gitty)](https://github.com/Omibranch/gitty/releases)
  [![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)](https://golang.org)
</div>

---

## Install

**Windows**
```sh
winget install Omibranch.Gitty
```

**Linux / macOS**
```sh
curl -fsSL https://raw.githubusercontent.com/Omibranch/gitty/master/install.sh | sh
```

Or grab a binary from [Releases](https://github.com/Omibranch/gitty/releases) and place it on your `PATH`.

---

## Commands

### Basic workflow

| Command | What it does |
|---|---|
| `gitty add .` | Stage all + auto-commit |
| `gitty add . --commit "msg"` | Stage all + commit with your message |
| `gitty push <branch>` | Push to branch |
| `gitty push <branch> --share` | Push + copy GitHub link to clipboard |
| `gitty push <branch> --force` | Force push |
| `gitty pull <branch>` | Safe pull (only adds new files, skips existing) |
| `gitty pull <branch> --hard` | Pull + overwrite local files |
| `gitty pull <branch> --commit` | Pull + auto-commit any changes |
| `gitty undo` | Soft-reset last commit (keeps changes staged) |

### Shorthand syntax

```sh
gitty push=main           # same as: gitty push main
gitty push to main        # same as: gitty push main
gitty pull~main           # same as: gitty pull main
gitty pull from main      # same as: gitty pull main
gitty migration main=dev  # merge main into dev
gitty migration main to dev
```

### Repo & branches

| Command | What it does |
|---|---|
| `gitty add repo "name"` | Create GitHub repo + link current folder |
| `gitty add repo "name" --public` | Create public repo |
| `gitty add branch "name"` | Create new local branch |
| `gitty rename branch old=new` | Rename branch (local + remote) |
| `gitty rename repo new-name` | Rename GitHub repo |

### Log & state

| Command | What it does |
|---|---|
| `gitty log` | Pretty git log |
| `gitty log --5h` | Commits in last 5 hours |
| `gitty log --3day` | Commits in last 3 days |
| `gitty log --2week` | Commits in last 2 weeks |
| `gitty log --1month` | Commits in last month |
| `gitty state` | Local repo stats (commits, branches, files) |
| `gitty state <URL>` | Stats for any GitHub repo |
| `gitty state --branches` | Branch list only |
| `gitty state --commits` | Commit count only |
| `gitty state --files` | File count only |

### Checkpoints

```sh
gitty checkpoint "v1" in main    # freeze current remote state of main as tag
gitty checkpoint v1*main         # shorthand
gitty restore "v1"               # reset local branch to checkpoint
```

### Extras

| Command | What it does |
|---|---|
| `gitty gitignore` | Interactive .gitignore template picker (arrow keys) |
| `gitty --version` | Show version |
| `gitty help` | Full help in English |
| `gitty help ru` | Полная справка на русском |

---

## Build from source

```sh
git clone https://github.com/Omibranch/gitty
cd gitty/source

# Windows
go build -ldflags="-s -w" -o ../gitty.exe .

# Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../gitty .
```

---

## Linux package repositories

- **AUR**: package files are ready in `pkg/aur/` (`PKGBUILD`, `.SRCINFO`)
- **pacman official (Arch extra/community)**: requires Arch maintainer acceptance
- **apt official (Debian/Ubuntu)**: requires Debian/Ubuntu packaging review + sponsorship

Publication guide:
- `pkg/OFFICIAL_REPOS.md`
- `pkg/aur/README.md`
- `pkg/deb/build.sh`

---

<div align="center">Made with ❤️ in Go</div>
