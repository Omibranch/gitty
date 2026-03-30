package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// ─────────────────────────────────────────────
// Version — change this when releasing a new build
// ─────────────────────────────────────────────

const gittyVersion = "1.0"

// ─────────────────────────────────────────────
// ANSI colour helpers
// ─────────────────────────────────────────────
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
)

func success(msg string) { fmt.Printf("%s[SUCCESS]%s %s\n", colorGreen, colorReset, msg) }
func fail(msg string)    { fmt.Printf("%s[ERROR]%s %s\n", colorRed, colorReset, msg) }
func hint(msg string)    { fmt.Printf("%s[HINT]%s %s\n", colorYellow, colorReset, msg) }
func info(msg string)    { fmt.Printf("%s[INFO]%s %s\n", colorCyan, colorReset, msg) }

// ─────────────────────────────────────────────
// Global proxy setting (set once in main via --proxy flag)
// ─────────────────────────────────────────────

// proxyURL holds the user-supplied proxy, e.g. "http://ip:port" or
// "http://login:pass@ip:port". Empty means no proxy override.
var proxyURL string

// proxyHint prints a single consistent hint when a network error occurs.
func proxyHint() {
	hint("If you are behind a proxy, use the --proxy flag:")
	hint("  gitty <command> --proxy \"http://ip:port\"")
	hint("  gitty <command> --proxy \"http://login:pass@ip:port\"")
}

// ─────────────────────────────────────────────
// Shell helpers
// ─────────────────────────────────────────────

// proxyEnv returns the current process environment with HTTPS_PROXY / HTTP_PROXY
// injected when proxyURL is set. git and gh both honour these variables.
func proxyEnv() []string {
	env := os.Environ()
	if proxyURL == "" {
		return env
	}
	// Overwrite or append proxy variables
	out := make([]string, 0, len(env)+2)
	for _, e := range env {
		k := strings.SplitN(e, "=", 2)[0]
		switch strings.ToUpper(k) {
		case "HTTP_PROXY", "HTTPS_PROXY", "ALL_PROXY":
			// drop existing values; we inject our own below
		default:
			out = append(out, e)
		}
	}
	out = append(out,
		"HTTP_PROXY="+proxyURL,
		"HTTPS_PROXY="+proxyURL,
		"ALL_PROXY="+proxyURL,
	)
	return out
}

// run executes a command streaming output to stdout/stderr.
// Proxy is injected via environment variables (git, gh).
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = proxyEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// runSilent executes a command and returns combined output as string.
func runSilent(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = proxyEnv()
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// runInteractive runs a command with full stdin/stdout/stderr passthrough.
func runInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = proxyEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// installGhFallback downloads and installs the latest gh CLI release directly
// via PowerShell + Invoke-WebRequest. Used when winget fails (e.g. proxy issues).
func installGhFallback() error {
	info("Trying fallback: downloading gh directly via PowerShell...")
	ps := buildPsFallback(
		"https://api.github.com/repos/cli/cli/releases/latest",
		// assemble filename and URL from release tag
		"$ver=$rel.tag_name.TrimStart('v');"+
			"$msi=\"gh_${ver}_windows_amd64.msi\";"+
			"$dlUrl=\"https://github.com/cli/cli/releases/download/$($rel.tag_name)/${msi}\";",
		// install command
		"Start-Process msiexec.exe -ArgumentList \"/i $dest /qn /norestart\" -Wait -NoNewWindow;",
	)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", ps)
	cmd.Env = proxyEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// installGitFallback downloads and installs the latest Git for Windows directly
// via PowerShell + Invoke-WebRequest. Used when winget fails (e.g. proxy issues).
func installGitFallback() error {
	info("Trying fallback: downloading Git directly via PowerShell...")
	ps := buildPsFallback(
		"https://api.github.com/repos/git-for-windows/git/releases/latest",
		"$ver=($rel.tag_name -replace '^v','') -replace '\\.windows\\.\\d+$','';"+
			"$exe=\"Git-${ver}-64-bit.exe\";"+
			"$dlUrl=\"https://github.com/git-for-windows/git/releases/download/$($rel.tag_name)/${exe}\";",
		"Start-Process $dest -ArgumentList '/VERYSILENT /NORESTART /NOCANCEL /SP-' -Wait -NoNewWindow;",
	)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", ps)
	cmd.Env = proxyEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// buildPsFallback returns a one-liner PowerShell script that:
//  1. GETs the GitHub releases API (with proxy if set, or -NoProxy to bypass system proxy)
//  2. Runs urlBuilder to compute $dlUrl and $dest filename
//  3. Downloads the asset
//  4. Runs installCmd to install it
func buildPsFallback(apiURL, urlBuilder, installCmd string) string {
	proxyFrag := ""
	credSetup := ""
	if proxyURL != "" {
		// Parse credentials out of the URL: http://user:pass@host:port
		credPart := ""
		bare := proxyURL
		if idx := strings.Index(bare, "://"); idx >= 0 {
			scheme := bare[:idx+3]
			rest := bare[idx+3:]
			if at := strings.LastIndex(rest, "@"); at >= 0 {
				credPart = rest[:at]
				bare = scheme + rest[at+1:]
			}
		}
		proxyFrag = fmt.Sprintf(" -Proxy '%s'", bare)
		if credPart != "" {
			parts := strings.SplitN(credPart, ":", 2)
			user, pass := parts[0], ""
			if len(parts) == 2 {
				pass = parts[1]
			}
			credSetup = fmt.Sprintf(
				"$proxyPass=ConvertTo-SecureString '%s' -AsPlainText -Force;"+
					"$proxyCred=New-Object System.Management.Automation.PSCredential('%s',$proxyPass);",
				pass, user)
			proxyFrag += " -ProxyCredential $proxyCred"
		}
	} else {
		// No --proxy flag: explicitly bypass any system/WinHTTP proxy so a
		// misconfigured system proxy (e.g. 407) cannot block the download.
		proxyFrag = " -NoProxy"
	}

	return "$ErrorActionPreference='Stop';" +
		credSetup +
		fmt.Sprintf("$rel=Invoke-RestMethod -Uri '%s' -UseBasicParsing%s;", apiURL, proxyFrag) +
		urlBuilder +
		"$dest=Join-Path $env:TEMP (Split-Path $dlUrl -Leaf);" +
		"Write-Host \"[INFO] Downloading $dlUrl\";" +
		fmt.Sprintf("Invoke-WebRequest -Uri $dlUrl -OutFile $dest -UseBasicParsing%s;", proxyFrag) +
		"Write-Host \"[INFO] Installing $(Split-Path $dest -Leaf)\";" +
		installCmd +
		"Remove-Item $dest -Force"
}

func toolExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// prompt asks a yes/no question, returns true for 'y'.
func prompt(question string) bool {
	fmt.Printf("%s (y/n): ", question)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

func timestamp() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05 UTC")
}

// ─────────────────────────────────────────────
// Arrow-key language selector
// ─────────────────────────────────────────────

// pickLanguage renders a two-option selector (EN / RU) navigable with
// ← → arrow keys and confirmed with Enter. Returns "en" or "ru".
func pickLanguage() string {
	langs := []string{"EN", "RU"}
	sel := 0

	// Hide cursor while selecting
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	renderLangPicker := func() {
		fmt.Print("\r\033[2K") // go to col 0, clear line
		fmt.Print("  Language / Язык:  ")
		for i, l := range langs {
			if i == sel {
				fmt.Printf("%s%s[ %s ]%s", colorBold, colorCyan, l, colorReset)
			} else {
				fmt.Printf("%s  %s  %s", colorDim, l, colorReset)
			}
			if i < len(langs)-1 {
				fmt.Print("  ")
			}
		}
		fmt.Printf("   %s←→%s move  %sEnter%s confirm",
			colorYellow, colorReset, colorGreen, colorReset)
	}

	renderLangPicker()

	for {
		b, err := readKey()
		if err != nil {
			break
		}
		switch b {
		case keyLeft:
			if sel > 0 {
				sel--
			}
		case keyRight:
			if sel < len(langs)-1 {
				sel++
			}
		case keyEnter:
			fmt.Println()
			return strings.ToLower(langs[sel])
		case keyEsc, keyQ:
			fmt.Println()
			return "en"
		}
		renderLangPicker()
	}
	fmt.Println()
	return "en"
}

// ─────────────────────────────────────────────
// PATH helpers
// ─────────────────────────────────────────────

func addToUserPath(dir string) error {
	currentPath, err := runSilent("powershell", "-NoProfile", "-Command",
		`[Environment]::GetEnvironmentVariable('PATH','User')`)
	if err != nil {
		currentPath = ""
	}
	for _, p := range strings.Split(currentPath, ";") {
		if strings.EqualFold(strings.TrimSpace(p), dir) {
			info("Directory already in PATH.")
			return nil
		}
	}
	newPath := currentPath
	if newPath != "" {
		newPath += ";"
	}
	newPath += dir
	_, err = runSilent("powershell", "-NoProfile", "-Command",
		fmt.Sprintf(`[Environment]::SetEnvironmentVariable('PATH','%s','User')`, newPath))
	if err != nil {
		return fmt.Errorf("failed to update PATH: %w", err)
	}
	return nil
}

// ─────────────────────────────────────────────
// Commands
// ─────────────────────────────────────────────

func cmdInstall() {
	info("Checking dependencies...")
	if !toolExists("git") {
		info("git not found. Installing...")
		if err := installGitFallback(); err != nil {
			fail("Failed to install git: " + err.Error())
			proxyHint()
			hint("Or run this terminal as Administrator.")
			os.Exit(1)
		}
		success("git installed.")
	} else {
		success("git is already installed.")
	}
	if !toolExists("gh") {
		info("GitHub CLI (gh) not found. Installing...")
		if err := installGhFallback(); err != nil {
			fail("Failed to install GitHub CLI: " + err.Error())
			proxyHint()
			hint("Or run this terminal as Administrator.")
			os.Exit(1)
		}
		success("GitHub CLI (gh) installed.")
	} else {
		success("GitHub CLI (gh) is already installed.")
	}
	exe, err := os.Executable()
	if err != nil {
		fail("Could not determine gitty executable path: " + err.Error())
		os.Exit(1)
	}
	dir := filepath.Dir(exe)
	info(fmt.Sprintf("Adding %s to User PATH...", dir))
	if err := addToUserPath(dir); err != nil {
		fail(err.Error())
		hint("You can manually add the directory to your PATH.")
		os.Exit(1)
	}
	success(fmt.Sprintf("'%s' added to User PATH.", dir))
	hint("Restart your terminal session for the PATH change to take effect.")
	success("gitty install complete.")
}

func cmdAuth() {
	if !toolExists("gh") {
		fail("GitHub CLI (gh) is not installed.")
		hint("Run 'gitty install' to set it up.")
		os.Exit(1)
	}
	info("Starting GitHub authentication...")
	if err := runInteractive("gh", "auth", "login"); err != nil {
		fail("Authentication failed: " + err.Error())
		proxyHint()
		os.Exit(1)
	}
	username, err := runSilent("gh", "api", "user", "--jq", ".login")
	if err != nil || username == "" {
		fail("Could not retrieve GitHub username after login.")
		os.Exit(1)
	}
	_ = run("git", "config", "--local", "user.gitty-gh-user", username)
	success(fmt.Sprintf("Authenticated as '%s' and linked to this repository.", username))
}

func cmdInit(url string) {
	if url == "" {
		fail("No URL provided.")
		hint("Usage: gitty init \"https://github.com/user/repo.git\"")
		os.Exit(1)
	}
	info("Initialising git repository...")
	if err := run("git", "init"); err != nil {
		fail("git init failed: " + err.Error())
		os.Exit(1)
	}
	info(fmt.Sprintf("Setting remote origin to: %s", url))
	if err := run("git", "remote", "add", "origin", url); err != nil {
		info("Remote 'origin' already exists. Updating URL...")
		if err2 := run("git", "remote", "set-url", "origin", url); err2 != nil {
			fail("Failed to set remote: " + err2.Error())
			os.Exit(1)
		}
	}
	success(fmt.Sprintf("Repository initialised with remote origin: %s", url))
}

// cmdAddRepo creates a GitHub repo and links the CWD to it.
// Pass public=true for a public repo, false for private.
// If the folder is already linked to a remote, the user is prompted to
// choose: [1] Replace (reinitialise this folder), [2] Keep (create on
// GitHub only, leave the local config untouched), or [3] Cancel.
func cmdAddRepo(name string, public bool) {
	if name == "" {
		fail("No repository name provided.")
		hint("Usage: gitty add repo \"my-project\"")
		hint("       gitty add repo \"my-project\" --public")
		os.Exit(1)
	}
	if !toolExists("gh") {
		fail("GitHub CLI (gh) is not installed.")
		hint("Run 'gitty install' to set it up.")
		os.Exit(1)
	}
	if _, err := runSilent("gh", "auth", "status"); err != nil {
		fail("Not authenticated with GitHub.")
		hint("Run 'gitty auth' first.")
		os.Exit(1)
	}

	// ── Detect existing remote ─────────────────────────────────────────────
	existingRemote, remoteErr := runSilent("git", "remote", "get-url", "origin")
	hasExistingRepo := remoteErr == nil && existingRemote != ""

	replaceLocal := false

	if hasExistingRepo {
		fmt.Printf("\n%s[WARNING]%s This folder is already linked to:\n",
			colorYellow, colorReset)
		fmt.Printf("  %s%s%s\n\n", colorCyan, existingRemote, colorReset)
		fmt.Println("Choose an action:")
		fmt.Printf("  %s[1]%s Replace  – create \"%s\" and rewire THIS folder to it\n",
			colorRed, colorReset, name)
		fmt.Printf("  %s[2]%s Keep     – create \"%s\" on GitHub only, this folder stays as-is\n",
			colorGreen, colorReset, name)
		fmt.Printf("  %s[3]%s Cancel\n\n", colorDim, colorReset)
		fmt.Print("Your choice (1/2/3): ")

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			replaceLocal = true
		case "2":
			replaceLocal = false
		default:
			info("Operation cancelled.")
			return
		}
	}

	// ── Ensure git is initialised ──────────────────────────────────────────
	if _, err := runSilent("git", "rev-parse", "--git-dir"); err != nil {
		info("No git repository found. Running git init first...")
		if err := run("git", "init"); err != nil {
			fail("git init failed: " + err.Error())
			os.Exit(1)
		}
	}

	// ── Ensure at least one commit exists (gh repo create --push requires it) ──
	if !hasExistingRepo || replaceLocal {
		// hasCommits: try to resolve HEAD — fails on a brand-new repo with no commits
		_, headErr := runSilent("git", "rev-parse", "HEAD")
		if headErr != nil {
			info("No commits yet. Creating initial commit...")
			// stage everything that exists in the folder
			_ = run("git", "add", ".")
			// if nothing is staged, create a placeholder so the commit is non-empty
			staged, _ := runSilent("git", "diff", "--cached", "--name-only")
			if strings.TrimSpace(staged) == "" {
				if f, err := os.Create(".gitkeep"); err == nil {
					f.Close()
				}
				_ = run("git", "add", ".gitkeep")
			}
			if err := run("git", "commit", "-m", "gitty: initial commit"); err != nil {
				fail("Failed to create initial commit: " + err.Error())
				os.Exit(1)
			}
		}
	}

	// ── Create repo on GitHub ──────────────────────────────────────────────
	visibility := "--private"
	visLabel := "private"
	if public {
		visibility = "--public"
		visLabel = "public"
	}
	info(fmt.Sprintf("Creating %s GitHub repository \"%s\"...", visLabel, name))

	if !hasExistingRepo || replaceLocal {
		// Link and push current folder to new repo
		if err := run("gh", "repo", "create", name,
			visibility, "--source=.", "--remote=origin", "--push"); err != nil {
			fail("Failed to create repository: " + err.Error())
			hint("Ensure you are authenticated. Run 'gitty auth'.")
			proxyHint()
			os.Exit(1)
		}
		if replaceLocal {
			success(fmt.Sprintf("Repository \"%s\" created. This folder is now linked to it.", name))
			info(fmt.Sprintf("Previous remote was: %s", existingRemote))
		} else {
			success(fmt.Sprintf("Repository \"%s\" created and linked to this folder.", name))
		}
	} else {
		// Create on GitHub only – no local changes
		if err := run("gh", "repo", "create", name, visibility); err != nil {
			fail("Failed to create repository: " + err.Error())
			hint("Ensure you are authenticated. Run 'gitty auth'.")
			proxyHint()
			os.Exit(1)
		}
		success(fmt.Sprintf("Repository \"%s\" created on GitHub.", name))
		info(fmt.Sprintf("This folder remains linked to: %s", existingRemote))
		hint(fmt.Sprintf("To link this folder later: gitty init \"https://github.com/<user>/%s.git\"", name))
	}
}

// cmdAddBranch creates a new local branch without switching to it.
// Branch name must be provided in quotes.
func cmdAddBranch(name string) {
	if name == "" {
		fail("No branch name provided.")
		hint("Usage: gitty add branch \"<branch-name>\"")
		os.Exit(1)
	}
	// Strip any residual quotes the shell may have passed through
	name = strings.Trim(name, "\"'")

	// git branch requires at least one commit to exist.
	// If the repo is empty, create an initial commit first.
	if _, err := runSilent("git", "rev-parse", "HEAD"); err != nil {
		info("Repository has no commits yet. Creating initial commit before branching...")
		ensureInitialCommit()
	}

	info(fmt.Sprintf("Creating branch \"%s\" (without switching)...", name))
	if err := run("git", "branch", name); err != nil {
		fail(fmt.Sprintf("Failed to create branch \"%s\": %s", name, err.Error()))
		os.Exit(1)
	}
	success(fmt.Sprintf("Branch \"%s\" created. You remain on your current branch.", name))
}

// ensureInitialCommit makes sure there is at least one commit in the current repo.
// If nothing is staged, it creates a .gitkeep so the commit is not empty.
func ensureInitialCommit() {
	// Stage everything that exists
	_, _ = runSilent("git", "add", ".")
	// If still nothing staged, create .gitkeep
	staged, _ := runSilent("git", "diff", "--cached", "--name-only")
	if strings.TrimSpace(staged) == "" {
		if err := os.WriteFile(".gitkeep", []byte(""), 0644); err != nil {
			fail("Could not create .gitkeep: " + err.Error())
			os.Exit(1)
		}
		_, _ = runSilent("git", "add", ".gitkeep")
	}
	if out, err := runSilent("git", "commit", "-m", "gitty: initial commit"); err != nil {
		fail("Initial commit failed: " + out)
		os.Exit(1)
	}
	success("Initial commit created.")
}

// cmdAddDot stages all changes and creates a commit.
// It does NOT push — use gitty push <branch> for that.
func cmdAddDot() {
	// If repo is empty (no commits yet), create initial commit
	_, headErr := runSilent("git", "rev-parse", "HEAD")
	if headErr != nil {
		info("Repository has no commits yet. Creating initial commit...")
		ensureInitialCommit()
		return
	}

	info("Staging all changes...")
	addOut, addErr := runSilent("git", "add", ".")
	if addErr != nil {
		badPaths := extractUnbornSubmodulePaths(addOut)
		if len(badPaths) > 0 {
			hint("Detected nested repository path(s) without commits. Retrying while skipping them:")
			for _, p := range badPaths {
				hint("  - " + p)
			}

			retryArgs := []string{"add", "."}
			for _, p := range badPaths {
				clean := strings.TrimSuffix(strings.TrimSpace(p), "/")
				if clean != "" {
					retryArgs = append(retryArgs, ":(exclude)"+clean)
				}
			}

			retryOut, retryErr := runSilent("git", retryArgs...)
			if retryErr != nil {
				fail("git add failed: " + retryOut)
				os.Exit(1)
			}
			info("Staging completed with exclusions for nested repo paths.")
		} else {
			fail("git add failed: " + addOut)
			os.Exit(1)
		}
	}
	status, _ := runSilent("git", "status", "--porcelain")
	if strings.TrimSpace(status) == "" {
		info("Nothing to commit – working tree is clean.")
		return
	}
	msg := fmt.Sprintf("gitty auto-sync [%s]", timestamp())
	info(fmt.Sprintf("Committing: %s", msg))
	commitOut, commitErr := runSilent("git", "commit", "-m", msg)
	if commitErr != nil {
		if strings.Contains(commitOut, "nothing to commit") ||
			strings.Contains(commitOut, "nothing added to commit") {
			info("Nothing new to commit.")
		} else {
			fail("git commit failed: " + commitOut)
			os.Exit(1)
		}
	} else {
		success("Changes committed.")
	}
}

func extractUnbornSubmodulePaths(addOutput string) []string {
	re := regexp.MustCompile(`'([^']+)' does not have a commit checked out`)
	seen := map[string]bool{}
	paths := []string{}
	for _, m := range re.FindAllStringSubmatch(addOutput, -1) {
		if len(m) < 2 {
			continue
		}
		p := strings.TrimSpace(m[1])
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		paths = append(paths, p)
	}
	return paths
}

// cmdPush pushes committed changes to the given remote branch.
// It does NOT stage or commit — use gitty add . for that first.
func cmdPush(branch string) {
	if branch == "" {
		fail("No target branch specified.")
		hint("Usage: gitty push <branch>  (or gitty push%<branch>)")
		os.Exit(1)
	}
	info(fmt.Sprintf("Pushing to origin/%s...", branch))
	pushOut, pushErr := runSilent("git", "push", "origin", branch)
	if pushErr != nil {
		if strings.Contains(pushOut, "does not exist") ||
			strings.Contains(pushOut, "has no upstream") ||
			strings.Contains(pushErr.Error(), "error") {
			if prompt(fmt.Sprintf("Branch '%s' does not exist on remote. Create and push?", branch)) {
				if err := run("git", "push", "--set-upstream", "origin", branch); err != nil {
					fail("Push failed: " + err.Error())
					os.Exit(1)
				}
				success(fmt.Sprintf("Branch '%s' created on remote and changes pushed.", branch))
				return
			}
			info("Push cancelled.")
			return
		}
		if strings.Contains(pushOut, "conflict") || strings.Contains(pushOut, "rejected") {
			fail("Conflict detected. Manual resolution required or use 'gitty pull~" + branch + " --hard' to overwrite.")
			os.Exit(1)
		}
		fail("git push failed: " + pushOut)
		proxyHint()
		os.Exit(1)
	}
	success(fmt.Sprintf("Changes pushed to origin/%s.", branch))
}

// cmdResetBranch deletes all commits on a branch by replacing it with an empty
// orphan commit, effectively wiping its entire history and content.
// Requires arrow-key Yes/No confirmation.
func cmdResetBranch(branch string) {
	if branch == "" {
		fail("No branch specified.")
		hint("Usage: gitty reset~<branch>")
		os.Exit(1)
	}

	// ── Arrow-key Yes / No confirmation ─────────────────────────────────
	opts := []string{"No", "Yes"}
	sel := 0 // default: No (safe)

	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	render := func() {
		fmt.Print("\r\033[2K")
		fmt.Printf("  %s[WARN]%s Delete ALL content and history of branch %s\"%s\"%s?   ",
			colorRed, colorReset, colorBold, branch, colorReset)
		for i, o := range opts {
			if i == sel {
				col := colorGreen
				if o == "Yes" {
					col = colorRed
				}
				fmt.Printf("%s%s[ %s ]%s", colorBold, col, o, colorReset)
			} else {
				fmt.Printf("%s  %s  %s", colorDim, o, colorReset)
			}
			if i < len(opts)-1 {
				fmt.Print("  ")
			}
		}
		fmt.Printf("   %s←→%s  %sEnter%s", colorYellow, colorReset, colorGreen, colorReset)
	}

	render()
	confirmed := false
	for {
		k, err := readKey()
		if err != nil {
			break
		}
		switch k {
		case keyLeft:
			if sel > 0 {
				sel--
			}
		case keyRight:
			if sel < len(opts)-1 {
				sel++
			}
		case keyEnter:
			fmt.Println()
			confirmed = opts[sel] == "Yes"
			goto done
		case keyEsc, keyQ:
			fmt.Println()
			goto done
		}
		render()
	}
done:
	if !confirmed {
		info("Reset cancelled.")
		return
	}

	// ── Find out which branch we are currently on ────────────────────────
	currentBranch, err := runSilent("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		fail("Could not determine current branch: " + err.Error())
		os.Exit(1)
	}

	onTarget := strings.EqualFold(currentBranch, branch)

	if onTarget {
		// We are on the branch to be reset — create orphan in place
		info(fmt.Sprintf("Resetting branch \"%s\" (currently checked out)...", branch))

		// Switch to a detached orphan state
		if err := run("git", "checkout", "--orphan", "__gitty_tmp_orphan__"); err != nil {
			fail("Could not create orphan branch: " + err.Error())
			os.Exit(1)
		}
		// Remove everything from the index and working tree
		_ , _ = runSilent("git", "rm", "-rf", ".")
		// Empty commit so the branch has a root
		if _, err := runSilent("git", "commit", "--allow-empty", "-m", "gitty: branch reset"); err != nil {
			fail("Empty commit failed.")
			os.Exit(1)
		}
		// Force-move the target branch name to this new orphan HEAD
		if err := run("git", "branch", "-M", branch); err != nil {
			fail("Could not rename orphan branch: " + err.Error())
			os.Exit(1)
		}
	} else {
		// We are on a different branch — reset target branch without switching
		info(fmt.Sprintf("Resetting branch \"%s\"...", branch))

		// Save current HEAD so we can return
		// Create a temporary orphan on a throwaway name
		if err := run("git", "checkout", "--orphan", "__gitty_tmp_orphan__"); err != nil {
			fail("Could not create orphan: " + err.Error())
			os.Exit(1)
		}
		_, _ = runSilent("git", "rm", "-rf", ".")
		if _, err := runSilent("git", "commit", "--allow-empty", "-m", "gitty: branch reset"); err != nil {
			fail("Empty commit failed.")
			os.Exit(1)
		}
		// Force-point the target branch here
		if err := run("git", "branch", "-f", branch, "HEAD"); err != nil {
			fail("Could not reset target branch: " + err.Error())
			os.Exit(1)
		}
		// Return to original branch
		if err := run("git", "checkout", currentBranch); err != nil {
			fail("Could not return to original branch: " + err.Error())
			os.Exit(1)
		}
		// Clean up the temp orphan
		_, _ = runSilent("git", "branch", "-D", "__gitty_tmp_orphan__")
	}

	// Push force to remote if remote branch exists
	remoteCheck, _ := runSilent("git", "ls-remote", "--heads", "origin", branch)
	if strings.TrimSpace(remoteCheck) != "" {
		info(fmt.Sprintf("Force-pushing empty branch to origin/%s...", branch))
		if err := run("git", "push", "origin", branch, "--force"); err != nil {
			fail("Force push failed: " + err.Error())
			proxyHint()
			os.Exit(1)
		}
	}

	success(fmt.Sprintf("Branch \"%s\" has been completely reset (all content and history removed).", branch))
}

func cmdPull(branch string, flag string) {
	if branch == "" {
		fail("No source branch specified.")
		hint("Usage: gitty pull <branch> [--hard | --hard-reset]  (or gitty pull~<branch>)")
		os.Exit(1)
	}
	switch flag {
	case "":
		info(fmt.Sprintf("Fetching origin/%s (safe mode – will not overwrite local files)...", branch))
		if err := run("git", "fetch", "origin", branch); err != nil {
			fail("git fetch failed: " + err.Error())
			proxyHint()
			os.Exit(1)
		}
		remoteFiles, err := runSilent("git", "ls-tree", "-r", "--name-only",
			fmt.Sprintf("origin/%s", branch))
		if err != nil {
			fail("Could not list remote files: " + err.Error())
			os.Exit(1)
		}
		copied, skipped := 0, 0
		for _, f := range strings.Split(remoteFiles, "\n") {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			if _, statErr := os.Stat(f); os.IsNotExist(statErr) {
				if checkoutErr := run("git", "checkout",
					fmt.Sprintf("origin/%s", branch), "--", f); checkoutErr != nil {
					fail(fmt.Sprintf("Failed to checkout '%s': %s", f, checkoutErr.Error()))
				} else {
					info(fmt.Sprintf("  + Added: %s", f))
					copied++
				}
			} else {
				skipped++
			}
		}
		success(fmt.Sprintf("Safe pull complete. %d file(s) added, %d existing file(s) untouched.", copied, skipped))

	case "--hard":
		info(fmt.Sprintf("Fetching origin/%s (--hard – will overwrite existing files)...", branch))
		if err := run("git", "fetch", "origin", branch); err != nil {
			fail("git fetch failed: " + err.Error())
			proxyHint()
			os.Exit(1)
		}
		if err := run("git", "checkout", fmt.Sprintf("origin/%s", branch), "--", "."); err != nil {
			fail("Hard checkout failed: " + err.Error())
			os.Exit(1)
		}
		success(fmt.Sprintf("Hard pull complete from origin/%s. Local-only files preserved.", branch))

	case "--hard-reset":
		fmt.Printf("%s[WARNING]%s --hard-reset will DELETE all local files not on origin/%s.\n",
			colorRed, colorReset, branch)
		if !prompt("Are you absolutely sure?") {
			info("Operation cancelled.")
			return
		}
		info(fmt.Sprintf("Fetching origin/%s...", branch))
		if err := run("git", "fetch", "origin", branch); err != nil {
			fail("git fetch failed: " + err.Error())
			proxyHint()
			os.Exit(1)
		}
		if err := run("git", "reset", "--hard", fmt.Sprintf("origin/%s", branch)); err != nil {
			fail("git reset failed: " + err.Error())
			os.Exit(1)
		}
		if err := run("git", "clean", "-fd"); err != nil {
			fail("git clean failed: " + err.Error())
			os.Exit(1)
		}
		success(fmt.Sprintf("Hard-reset complete. Working tree now mirrors origin/%s.", branch))

	default:
		fail(fmt.Sprintf("Unknown flag '%s'.", flag))
		hint("Valid flags: (none), --hard, --hard-reset")
		os.Exit(1)
	}
}

// ─────────────────────────────────────────────
// Clear
// ─────────────────────────────────────────────

// cmdClear clears the terminal screen. The standard 'cls' only works in
// cmd.exe; this uses the ANSI escape sequence which works everywhere.
func cmdClear() {
	fmt.Print("\033[H\033[2J\033[3J")
}

// ─────────────────────────────────────────────
// Gitignore
// ─────────────────────────────────────────────

// gitignoreTemplates is the full list from https://api.github.com/gitignore/templates
// embedded at compile time so no network call is needed for browsing.
var gitignoreTemplates = []string{
	"AL", "Actionscript", "Ada", "AdventureGameStudio", "Agda", "Android",
	"Angular", "AppEngine", "AppceleratorTitanium", "ArchLinuxPackages",
	"Autotools", "Ballerina", "C", "C++", "CFWheels", "CMake", "CUDA",
	"CakePHP", "ChefCookbook", "Clojure", "CodeIgniter", "CommonLisp",
	"Composer", "Concrete5", "Coq", "CraftCMS", "D", "DM", "Dart", "Delphi",
	"Dotnet", "Drupal", "EPiServer", "Eagle", "Elisp", "Elixir", "Elm",
	"Erlang", "ExpressionEngine", "ExtJs", "Fancy", "Finale", "Firebase",
	"FlaxEngine", "Flutter", "ForceDotCom", "Fortran", "FuelPHP", "GWT",
	"Gcov", "GitBook", "GitHubPages", "Gleam", "Go", "Godot", "Gradle",
	"Grails", "HIP", "Haskell", "Haxe", "IAR", "IGORPro", "Idris", "JBoss",
	"JENKINS_HOME", "Java", "Jekyll", "Joomla", "Julia", "Katalon", "KiCad",
	"Kohana", "Kotlin", "LabVIEW", "LangChain", "Laravel", "Leiningen",
	"LemonStand", "Lilypond", "Lithium", "Lua", "Luau", "Magento", "Maven",
	"Mercury", "MetaProgrammingSystem", "Modelica", "Nanoc", "Nestjs",
	"Nextjs", "Nim", "Nix", "Node", "OCaml", "Objective-C", "Opa",
	"OpenCart", "OracleForms", "Packer", "Perl", "Phalcon", "PlayFramework",
	"Plone", "Prestashop", "Processing", "PureScript", "Python", "Qooxdoo",
	"Qt", "R", "ROS", "Racket", "Rails", "Raku", "ReScript",
	"RhodesRhomobile", "Ruby", "Rust", "SCons", "SSDT-sqlproj", "Sass",
	"Scala", "Scheme", "Scrivener", "Sdcc", "SeamGen", "SketchUp",
	"Smalltalk", "Solidity-Remix", "Stella", "SugarCRM", "Swift", "Symfony",
	"SymphonyCMS", "TeX", "Terraform", "TestComplete", "Textpattern",
	"TurboGears2", "TwinCAT3", "Typo3", "Unity", "UnrealEngine", "VBA",
	"VVVV", "VisualStudio", "Waf", "WordPress", "Xojo", "Yeoman", "Yii",
	"ZendFramework", "Zephir", "Zig", "ecu.test",
}

// fetchGitignore downloads the gitignore template content from GitHub API.
func fetchGitignore(name string) (string, error) {
	apiURL := "https://api.github.com/gitignore/templates/" + url.PathEscape(name)
	client := &http.Client{}
	if proxyURL != "" {
		pu, err := url.Parse(proxyURL)
		if err == nil {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(pu)}
		}
	}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		Source string `json:"source"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Source, nil
}

// cmdGitignore shows an interactive fuzzy picker for GitHub gitignore templates.
// Type to filter, ↑↓ to navigate, Enter to select, Esc/q to cancel.
func cmdGitignore() {
	kernel32, _ := loadKernel32()

	query := ""
	cursor := 0
	const maxVisible = 12

	filterList := func(q string) []string {
		q = strings.ToLower(q)
		var out []string
		for _, t := range gitignoreTemplates {
			if q == "" || strings.Contains(strings.ToLower(t), q) {
				out = append(out, t)
			}
		}
		return out
	}

	clearLines := func(n int) {
		for i := 0; i < n; i++ {
			fmt.Print("\033[A\033[2K")
		}
	}

	// Enter raw mode
	var oldMode uint32
	if m, err := getConsoleMode(kernel32); err == nil {
		oldMode = m
		_ = setConsoleModeRaw(kernel32)
	}
	restore := func() { _ = restoreConsoleMode(kernel32, oldMode) }
	defer restore()

	fmt.Print("\033[?25l") // hide cursor
	defer fmt.Print("\033[?25h")

	fmt.Println()
	fmt.Printf("  %sSelect .gitignore template%s\n", colorBold, colorReset)
	fmt.Printf("  %sType to filter  ↑↓ navigate  Enter select  Esc cancel%s\n\n", colorDim, colorReset)

	prevLines := 0

	for {
		list := filterList(query)
		if len(list) > 0 && cursor >= len(list) {
			cursor = len(list) - 1
		}
		if cursor < 0 {
			cursor = 0
		}

		clearLines(prevLines)

		// Search box line
		fmt.Printf("  %s›%s %s%s%s\n", colorCyan, colorReset, colorBold, query, colorReset)
		lines := 1

		if len(list) == 0 {
			fmt.Printf("  %s(no matches)%s\n", colorDim, colorReset)
			lines++
		} else {
			start := cursor - maxVisible/2
			if start < 0 {
				start = 0
			}
			end := start + maxVisible
			if end > len(list) {
				end = len(list)
				start = end - maxVisible
				if start < 0 {
					start = 0
				}
			}
			if start > 0 {
				fmt.Printf("  %s  ↑ %d more%s\n", colorDim, start, colorReset)
				lines++
			}
			for i := start; i < end; i++ {
				if i == cursor {
					fmt.Printf("  %s%s▶  %s%s\n", colorGreen, colorBold, list[i], colorReset)
				} else {
					fmt.Printf("  %s   %s%s\n", colorDim, list[i], colorReset)
				}
				lines++
			}
			if end < len(list) {
				fmt.Printf("  %s  ↓ %d more%s\n", colorDim, len(list)-end, colorReset)
				lines++
			}
		}
		prevLines = lines

		k, ch, err := readKeyOrChar()
		if err != nil {
			break
		}

		switch k {
		case keyUp:
			if cursor > 0 {
				cursor--
			}
		case keyDown:
			list2 := filterList(query)
			if cursor < len(list2)-1 {
				cursor++
			}
		case keyBackspace:
			if len(query) > 0 {
				query = query[:len([]rune(query))-1]
				cursor = 0
			}
		case keyChar:
			query += string(ch)
			cursor = 0
		case keyEnter:
			list2 := filterList(query)
			if len(list2) == 0 {
				continue
			}
			chosen := list2[cursor]
			clearLines(prevLines)
			restore()
			fmt.Print("\033[?25h")

			info(fmt.Sprintf("Downloading template: %s", chosen))
			content, dlErr := fetchGitignore(chosen)
			if dlErr != nil {
				fail("Failed to download template: " + dlErr.Error())
				proxyHint()
				os.Exit(1)
			}

			dest := ".gitignore"
			if _, statErr := os.Stat(dest); statErr == nil {
				fmt.Printf("\n  %s[!] .gitignore already exists.%s Overwrite? (y/n): ", colorYellow, colorReset)
				// restore mode so the user can type normally
				_ = restoreConsoleMode(kernel32, oldMode)
				reader := bufio.NewReader(os.Stdin)
				ans, _ := reader.ReadString('\n')
				if strings.ToLower(strings.TrimSpace(ans)) != "y" {
					info("Cancelled.")
					return
				}
				// back to raw for safety (we're about to return anyway)
			}

			if err := os.WriteFile(dest, []byte(content), 0644); err != nil {
				fail("Failed to write .gitignore: " + err.Error())
				os.Exit(1)
			}
			success(fmt.Sprintf(".gitignore created from \"%s\" template.", chosen))
			return

		case keyEsc, keyQ:
			clearLines(prevLines)
			info("Cancelled.")
			return
		}
	}
}

// ─────────────────────────────────────────────
// Status
// ─────────────────────────────────────────────

// cmdStatus shows which GitHub account and repository are linked to the CWD.
func cmdStatus() {
	fmt.Println()

	// ── git remote ────────────────────────────────────────────────────────
	remote, err := runSilent("git", "remote", "get-url", "origin")
	remote = strings.TrimSpace(remote)
	if err != nil || remote == "" {
		info("No git remote configured in this folder.")
		hint("Run 'gitty init \"<url>\"' or 'gitty add repo \"name\"' to link one.")
	} else {
		success("Remote origin:  " + remote)
	}

	// ── current branch ────────────────────────────────────────────────────
	branch, err2 := runSilent("git", "rev-parse", "--abbrev-ref", "HEAD")
	branch = strings.TrimSpace(branch)
	if err2 == nil && branch != "" && branch != "HEAD" {
		info("Current branch: " + branch)
	}

	// ── gh auth status ────────────────────────────────────────────────────
	if !toolExists("gh") {
		hint("gh (GitHub CLI) is not installed. Run 'gitty install'.")
		fmt.Println()
		return
	}
	authOut, authErr := runSilent("gh", "auth", "status")
	authOut = strings.TrimSpace(authOut)
	if authErr != nil || authOut == "" {
		info("GitHub account: not authenticated.")
		hint("Run 'gitty auth' to sign in.")
	} else {
		// gh auth status output contains the account name on a line like:
		// ✓ Logged in to github.com account USERNAME (...)
		account := ""
		for _, line := range strings.Split(authOut, "\n") {
			line = strings.TrimSpace(line)
			lower := strings.ToLower(line)
			if strings.Contains(lower, "logged in") || strings.Contains(lower, "account") {
				account = line
				break
			}
		}
		if account != "" {
			success("GitHub account: " + account)
		} else {
			success("GitHub account: authenticated")
			info(authOut)
		}
	}
	fmt.Println()
}

// ─────────────────────────────────────────────
// Help – bilingual with arrow-key language picker
// ─────────────────────────────────────────────

func cmdHelp() {
	bold := func(s string) string { return colorBold + s + colorReset }
	cyan := func(s string) string { return colorCyan + s + colorReset }
	green := func(s string) string { return colorGreen + s + colorReset }
	yellow := func(s string) string { return colorYellow + s + colorReset }

	fmt.Println()
	lang := pickLanguage()
	fmt.Println()

	if lang == "ru" {
		printHelpRU(bold, cyan, green, yellow)
	} else {
		printHelpEN(bold, cyan, green, yellow)
	}
}

func printHelpEN(bold, cyan, green, yellow func(string) string) {
	dim := func(s string) string { return colorDim + s + colorReset }
	red := func(s string) string { return colorRed + s + colorReset }
	fmt.Printf(`
%s

  %s  means  %s  (push destination)
  %s  means  %s  (pull source)

%s
  gitty handles all the boring Git work: staging, committing, pushing.
  Zero-Commit principle — you never type 'git commit' manually.
  You always stay on your own branch. The target branch is specified
  inline using %s and %s syntax.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  %s
    gitty install

    Installs git and gh (GitHub CLI) if missing.
    Downloads them directly — no winget, no admin rights needed.
    Adds gitty.exe location to your User PATH.

  %s
    gitty auth

    Runs gh auth login to sign into GitHub.
    Required before creating repositories.

  %s
    gitty init "https://github.com/user/repo.git"

    Initialises git in the current folder and sets the remote
    origin to the given URL. URL must be in quotes.

  %s
    gitty add repo "name"
    gitty add repo "name" %s

    Creates a GitHub repository under your account.
    Default: %s. Add %s to make it %s.

    If this folder is already linked to a repo, you are asked:
      %s  Replace  — rewire this folder to the new repo
      %s  Keep     — create on GitHub only, folder stays as-is
      %s  Cancel

    Requires: gitty auth

  %s
    gitty add branch "name"

    Creates a new local branch (quotes required).
    You stay on your current branch — no checkout.
    If the repo has no commits yet, an initial commit is created automatically.

  %s
    gitty add .

    Stages all changes and creates a commit:
      1. git add .         — stage everything (.gitignore respected)
      2. git commit -m "gitty auto-sync [UTC timestamp]"

    If the repo has no commits yet, an initial commit is created automatically.

	%s
	gitty push %s<branch>%s
		Compact: gitty push%%<branch>

    Pushes committed changes to the given remote branch.
    Run %s first to stage and commit.

    Examples:
	gitty push main
	gitty push dev
	gitty push feature/login

    If the remote branch does not exist, gitty will ask to create it.
	Chain with add: gitty add . and push main

	%s
	gitty pull %s<branch>%s [--hard | --hard-reset]
		Compact: gitty pull~<branch> [--hard | --hard-reset]

    %s (default)   Copy only files missing locally. Never overwrites.
					gitty pull main

    %s             Overwrite files that exist on remote too.
                    Local-only files are kept.
					gitty pull staging --hard

    %s  %s Mirror remote exactly. Deletes local files not on remote.
                    Confirmation required.
					gitty pull main --hard-reset

  %s
    gitty reset~%s<branch>%s

    Wipes ALL content and commit history from a branch.
    An empty orphan commit replaces the entire branch.
    If the branch exists on the remote, it is force-pushed.
    Requires arrow-key %s / %s confirmation before executing.

    Example:
      gitty reset~second
      gitty reset~old-feature

  %s
    gitty help

    Shows this manual. Language picker appears first.

  %s
    gitty clear

    Clears the terminal screen (works in PowerShell, cmd, Windows Terminal).

  %s
    gitty status

    Shows which GitHub account is authenticated and which remote repository
    is linked to the current folder, plus the active branch.

  %s
    gitty gitignore

    Interactive fuzzy picker for GitHub's official .gitignore templates.
    Type to filter, ↑↓ to navigate, Enter to download and save .gitignore.
    Uses the same templates available on github.com/new.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

%s

  %s  Success       %s  Error       %s  Hint

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

%s

	%s  push      TO a branch       gitty push main
	%s  pull      FROM a branch     gitty pull main
  %s  and         chain commands            gitty auth and add repo "name"
											gitty install and auth and add repo "x" and add . and push main
  %s  --public    create a public repo    gitty add repo "name" --public
  %s  --proxy     set proxy               gitty <cmd> --proxy "http://ip:port"
                                          gitty <cmd> --proxy "http://user:pass@ip:port"
  %s  --hard      overwrite on pull       gitty pull~main --hard
  %s  --hard-reset  mirror remote %s   gitty pull~main --hard-reset

`,
		bold(cyan("╔══════════════════════════════════════════╗\n║          GITTY  MANUAL  (EN)             ║\n╚══════════════════════════════════════════╝")),
		green("%"), bold("TO"),
		yellow("~"), bold("FROM"),
		bold("OVERVIEW"),
		green("%"), yellow("~"),
		bold("gitty install"),
		bold("gitty auth"),
		bold("gitty init"),
		bold("gitty add repo"),
		bold("--public"),
		bold("private"), bold("--public"), bold("public"),
		dim("[1]"), dim("[2]"), dim("[3]"),
		bold("gitty add branch"),
		bold("gitty add ."),
		bold("gitty push"),
		green(""), green(""),
		bold("gitty add ."),
		bold("gitty pull"),
		yellow(""), yellow(""),
		dim("(no flag)"),
		bold("--hard"),
		bold("--hard-reset"), red("DESTRUCTIVE"),
		bold("gitty reset~"),
		red(""), red(""),
		red("Yes"), red("No"),
		bold("gitty help"),
		bold("gitty clear"),
		bold("gitty status"),
		bold("gitty gitignore"),
		bold("OUTPUT PREFIXES"),
		green("[SUCCESS]"), red("[ERROR]"), yellow("[HINT]"),
		bold("FLAGS & SYNTAX"),
		green(""), yellow(""), dim(""), dim(""), dim(""), dim(""), dim(""), red(""),
	)
}

func printHelpRU(bold, cyan, green, yellow func(string) string) {
	dim := func(s string) string { return colorDim + s + colorReset }
	red := func(s string) string { return colorRed + s + colorReset }
	fmt.Printf(`
%s

  %s  означает  %s  (куда отправить)
  %s  означает  %s  (откуда получить)

%s
  gitty берёт рутину Git на себя: стейджинг, коммит, пуш.
  Принцип нулевого коммита — вы никогда не вводите 'git commit'.
  Вы всегда остаётесь на своей ветке. Целевая ветка задаётся
  прямо в команде через %s и %s.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  %s
    gitty install

    Устанавливает git и gh (GitHub CLI), если их нет.
    Скачивает напрямую — без winget и прав администратора.
    Добавляет папку с gitty.exe в PATH пользователя.

  %s
    gitty auth

    Запускает gh auth login для входа в GitHub.
    Обязателен перед созданием репозиториев.

  %s
    gitty init "https://github.com/user/repo.git"

    Инициализирует git в текущей папке и устанавливает
    remote origin на указанный URL. URL нужно взять в кавычки.

  %s
    gitty add repo "название"
    gitty add repo "название" %s

    Создаёт репозиторий на GitHub под вашим аккаунтом.
    По умолчанию: %s. Добавьте %s чтобы сделать %s.

    Если папка уже привязана к другому репо, будет предложено:
      %s  Заменить — перепривязать папку к новому репо
      %s  Оставить — создать на GitHub, папку не трогать
      %s  Отмена

    Требует: gitty auth

  %s
    gitty add branch "название"

    Создаёт локальную ветку (кавычки обязательны).
    Вы остаётесь на текущей ветке — переключения не происходит.
    Если коммитов ещё нет, начальный коммит создаётся автоматически.

  %s
    gitty add .

    Стейджит все изменения и создаёт коммит:
      1. git add .         — добавить всё (.gitignore учитывается)
      2. git commit -m "gitty auto-sync [UTC время]"

    Если коммитов ещё нет, начальный коммит создаётся автоматически.

	%s
		gitty push %s<ветка>%s
		Коротко: gitty push%%<ветка>

    Отправляет закоммиченные изменения в указанную ветку на remote.
    Сначала выполните %s для стейджинга и коммита.

    Примеры:
	gitty push main
	gitty push dev
	gitty push feature/login

    Если ветки нет на remote, gitty предложит её создать.
	Цепочка: gitty add . and push main

	%s
		gitty pull %s<ветка>%s [--hard | --hard-reset]
		Коротко: gitty pull~<ветка> [--hard | --hard-reset]

    %s (без флага)  Копирует только файлы, которых нет локально.
                     Существующие не трогает.
					 gitty pull main

    %s              Перезаписывает файлы с remote. Локальные
                     уникальные файлы сохраняются.
					 gitty pull staging --hard

    %s  %s Приводит папку в точное соответствие с remote.
                     Локальные файлы, которых нет на remote, удаляются.
                     Требует подтверждения.
					 gitty pull main --hard-reset

  %s
    gitty reset~%s<ветка>%s

    Удаляет ВСЁ содержимое и историю коммитов указанной ветки.
    Вместо них создаётся пустой сиротский коммит.
    Если ветка есть на remote — выполняется принудительный push.
    Перед выполнением требуется подтверждение стрелками: %s / %s.

    Примеры:
      gitty reset~second
      gitty reset~old-feature

  %s
    gitty help

    Показывает это руководство. Сначала выбор языка.

  %s
    gitty clear

    Очищает экран терминала (работает в PowerShell, cmd, Windows Terminal).

  %s
    gitty status

    Показывает привязанный GitHub аккаунт, remote-репозиторий
    и текущую ветку в этой папке.

  %s
    gitty gitignore

    Интерактивный поиск по официальным шаблонам .gitignore с GitHub.
    Вводите название — список фильтруется на лету.
    ↑↓ — навигация, Enter — скачать и сохранить .gitignore в папку.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

%s

  %s  Успешно       %s  Ошибка       %s  Подсказка

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

%s

	%s  push      отправить В ветку   gitty push main
	%s  pull      получить ИЗ ветки   gitty pull main
  %s  and         цепочка команд            gitty auth and add repo "название"
											gitty install and auth and add repo "x" and add . and push main
  %s  --public    публичный репо           gitty add repo "название" --public
  %s  --proxy     прокси                   gitty <команда> --proxy "http://ip:port"
                                           gitty <команда> --proxy "http://user:pass@ip:port"
  %s  --hard      перезапись при pull      gitty pull~main --hard
  %s  --hard-reset  зеркало remote %s  gitty pull~main --hard-reset

`,
		bold(cyan("╔══════════════════════════════════════════╗\n║       РУКОВОДСТВО  GITTY  (RU)           ║\n╚══════════════════════════════════════════╝")),
		green("%"), bold("ТУДА"),
		yellow("~"), bold("ОТТУДА"),
		bold("ОБЗОР"),
		green("%"), yellow("~"),
		bold("gitty install"),
		bold("gitty auth"),
		bold("gitty init"),
		bold("gitty add repo"),
		bold("--public"),
		bold("приватный"), bold("--public"), bold("публичным"),
		dim("[1]"), dim("[2]"), dim("[3]"),
		bold("gitty add branch"),
		bold("gitty add ."),
		bold("gitty push"),
		green(""), green(""),
		bold("gitty add ."),
		bold("gitty pull"),
		yellow(""), yellow(""),
		dim("(нет флага)"),
		bold("--hard"),
		bold("--hard-reset"), red("ДЕСТРУКТИВНО"),
		bold("gitty reset~"),
		red(""), red(""),
		red("Yes"), red("No"),
		bold("gitty help"),
		bold("gitty clear"),
		bold("gitty status"),
		bold("gitty gitignore"),
		bold("ПРЕФИКСЫ ВЫВОДА"),
		green("[SUCCESS]"), red("[ERROR]"), yellow("[HINT]"),
		bold("ФЛАГИ И СИНТАКСИС"),
		green(""), yellow(""), dim(""), dim(""), dim(""), dim(""), dim(""), red(""),
	)
}

// ─────────────────────────────────────────────
// Argument parser / main
// ─────────────────────────────────────────────

func main() {
	if runtime.GOOS == "windows" {
		_ = enableWindowsANSI()
	}

	args := os.Args[1:]

	// ── --v / --version / -v — print version and exit ────────────────────
	for _, a := range args {
		if a == "--v" || a == "--version" || a == "-v" {
			fmt.Printf("gitty version %s\n", gittyVersion)
			return
		}
	}

	// ── Pre-parse global --proxy flag ─────────────────────────────────────
	// Supported anywhere in the argument list:
	//   gitty install --proxy "http://ip:port"
	//   gitty install --proxy "http://login:pass@ip:port"
	//   gitty --proxy "http://ip:port" install   (also works)
	filtered := args[:0:0] // empty slice sharing no backing array
	for i := 0; i < len(args); i++ {
		if args[i] == "--proxy" {
			if i+1 >= len(args) {
				fail("--proxy requires a value.")
				hint("Example: gitty install --proxy \"http://ip:port\"")
				os.Exit(1)
			}
			proxyURL = strings.Trim(args[i+1], "\"'")
			i++ // skip the value token
			info(fmt.Sprintf("Proxy set: %s", proxyURL))
		} else {
			filtered = append(filtered, args[i])
		}
	}
	args = filtered
	// ─────────────────────────────────────────────────────────────────────

	if len(args) == 0 {
		cmdHelp()
		return
	}

	// ── Split on "and" to support chaining: gitty auth and add repo "x" ──
	var segments [][]string
	cur := []string{}
	for _, a := range args {
		if strings.ToLower(a) == "and" {
			if len(cur) > 0 {
				segments = append(segments, cur)
				cur = []string{}
			}
		} else {
			cur = append(cur, a)
		}
	}
	if len(cur) > 0 {
		segments = append(segments, cur)
	}

	for _, seg := range segments {
		dispatch(seg)
	}
}

func dispatch(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {

	case "install":
		cmdInstall()

	case "auth":
		cmdAuth()

	case "clear":
		cmdClear()

	case "status":
		cmdStatus()

	case "push":
		if len(args) < 2 {
			fail("No target branch specified.")
			hint("Usage: gitty push <branch>  (or gitty push%<branch>)")
			os.Exit(1)
		}
		cmdPush(strings.Trim(args[1], "\"'"))

	case "pull":
		if len(args) < 2 {
			fail("No source branch specified.")
			hint("Usage: gitty pull <branch> [--hard | --hard-reset]  (or gitty pull~<branch>)")
			os.Exit(1)
		}
		flag := ""
		if len(args) > 2 {
			flag = strings.ToLower(strings.TrimSpace(args[2]))
		}
		cmdPull(strings.Trim(args[1], "\"'"), flag)

	case "gitignore":
		cmdGitignore()

	case "help", "--help", "-h":
		cmdHelp()

	case "init":
		url := ""
		if len(args) > 1 {
			url = strings.Trim(args[1], "\"'")
		}
		cmdInit(url)

	case "add":
		if len(args) < 2 {
			fail("Incomplete 'add' command.")
			hint("Usage: gitty add repo \"name\" | gitty add branch \"name\" | gitty add .")
			os.Exit(1)
		}
		switch args[1] {

		case "repo":
			name := ""
			public := false
			// collect name (first non-flag token) and --public flag
			for _, a := range args[2:] {
				if a == "--public" {
					public = true
				} else if name == "" {
					name = strings.Trim(a, "\"'")
				}
			}
			cmdAddRepo(name, public)

		case "branch":
			name := ""
			if len(args) > 2 {
				name = strings.Trim(args[2], "\"'")
			}
			cmdAddBranch(name)

		case ".":
			cmdAddDot()

		default:
			fail(fmt.Sprintf("Unknown 'add' sub-command: '%s'", args[1]))
			hint("Valid sub-commands: repo, branch, .")
			os.Exit(1)
		}

	default:
		rePush  := regexp.MustCompile(`^push%(.+)$`)
		rePull  := regexp.MustCompile(`^pull~(.+)$`)
		reReset := regexp.MustCompile(`^reset~(.+)$`)
		if mp := rePush.FindStringSubmatch(args[0]); mp != nil {
			cmdPush(mp[1])
		} else if mp := rePull.FindStringSubmatch(args[0]); mp != nil {
			flag := ""
			if len(args) > 1 {
				flag = strings.ToLower(strings.TrimSpace(args[1]))
			}
			cmdPull(mp[1], flag)
		} else if mp := reReset.FindStringSubmatch(args[0]); mp != nil {
			cmdResetBranch(mp[1])
		} else if args[0] == "push-" {
			currentBranch, err := runSilent("git", "rev-parse", "--abbrev-ref", "HEAD")
			if err != nil || strings.TrimSpace(currentBranch) == "" {
				fail("Could not infer target branch for 'push-'.")
				hint("Use: gitty push <branch>  (recommended)")
				os.Exit(1)
			}
			hint("It looks like '>' was interpreted by shell redirection. Using current branch instead.")
			hint("Recommended syntax without special characters: gitty push <branch>")
			cmdPush(strings.TrimSpace(currentBranch))
		} else {
			fail(fmt.Sprintf("Unknown command: '%s'", args[0]))
			hint("Run 'gitty help' for a full list of commands.")
			os.Exit(1)
		}
	}
}

// ─────────────────────────────────────────────
// Windows ANSI
// ─────────────────────────────────────────────

func enableWindowsANSI() error {
	if runtime.GOOS != "windows" {
		return nil
	}
	kernel32, err := loadKernel32()
	if err != nil {
		return err
	}
	return setConsoleMode(kernel32)
}
