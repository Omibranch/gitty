# gitty

A Windows CLI tool that wraps Git and GitHub CLI into simple, human-friendly commands. Zero manual commits — gitty handles staging, committing, and pushing for you.

## Quick Install

### Option 1 — PowerShell one-liner (recommended)
```powershell
irm https://raw.githubusercontent.com/Omibranch/gitty/main/install.ps1 | iex
```
Downloads `gitty.exe`, adds it to your User PATH, then runs `gitty install` to set up git and gh automatically.

### Option 2 — winget
```
winget install Omibranch.Gitty
```

### Option 3 — Manual
Download `gitty.exe` from [Releases](https://github.com/Omibranch/gitty/releases), place it anywhere in your PATH, then run:
```
gitty install
```

---

## Commands

| Command | What it does |
|---|---|
| `gitty install` | Install git + gh CLI, add gitty to PATH |
| `gitty auth` | Sign into GitHub (gh auth login) |
| `gitty init "url"` | Init git repo and set remote origin |
| `gitty add repo "name"` | Create GitHub repo and link current folder |
| `gitty add repo "name" --public` | Same but public |
| `gitty add branch "name"` | Create local branch without switching |
| `gitty add .` | Stage all changes and commit |
| `gitty push->branch` | Push committed changes to remote branch |
| `gitty pull~branch` | Safe pull — only adds missing files |
| `gitty pull~branch --hard` | Overwrite existing files from remote |
| `gitty pull~branch --hard-reset` | Mirror remote exactly (destructive) |
| `gitty reset~branch` | Wipe all content and history from a branch |
| `gitty status` | Show linked repo, branch, and GitHub account |
| `gitty gitignore` | Interactive picker for .gitignore templates |
| `gitty clear` | Clear terminal screen |
| `gitty help` | Full manual (EN / RU) |
| `gitty --v` | Show version |

## Chaining with `and`

```
gitty add . and push->main
gitty install and auth and add repo "my-project" and add . and push->main
```

## Proxy support

```
gitty <command> --proxy "http://ip:port"
gitty <command> --proxy "http://user:pass@ip:port"
```

Works for all network operations including git, gh, and direct downloads.

## Build from source

Requires Go 1.21+, no external dependencies.

```
git clone https://github.com/Omibranch/gitty.git
cd gitty
go build -ldflags="-s -w" -o gitty.exe .
```

## License

MIT (c) Omibranch
