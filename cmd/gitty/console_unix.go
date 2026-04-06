//go:build !windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

// keyXxx constants (mirrored from console_windows.go for cross-platform use).
const (
	keyLeft      = 1
	keyRight     = 2
	keyUp        = 3
	keyDown      = 4
	keyEnter     = 5
	keyEsc       = 6
	keyQ         = 7
	keyBackspace = 8
	keyChar      = 9
)

// ──────────────────────────────────────────────────────────────────
// Stub type replacing *syscall.DLL (doesn't exist on Linux/macOS)
// ──────────────────────────────────────────────────────────────────

// consoleDLL is a stub type used in place of *syscall.DLL on Unix.
// All functions that accept it are no-ops or use termios directly.
type consoleDLL struct{}

// ──────────────────────────────────────────────────────────────────
// termios raw-mode helpers
// ──────────────────────────────────────────────────────────────────

// termiosState mirrors the C termios struct for direct syscall access.
type termiosState struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]uint8
	Ispeed uint32
	Ospeed uint32
}

// saved terminal state for restore
var (
	savedState  termiosState
	stdinFd     uintptr
	stateValid  bool
)

func tcgetattr(fd uintptr, t *termiosState) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, ioctlTcgeta, uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return errno
	}
	return nil
}

func tcsetattr(fd uintptr, t *termiosState) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, ioctlTcseta, uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return errno
	}
	return nil
}

// ──────────────────────────────────────────────────────────────────
// Console API (mirrors console_windows.go signatures)
// ──────────────────────────────────────────────────────────────────

// loadKernel32 returns a stub — no kernel32 on Unix.
func loadKernel32() (*consoleDLL, error) { return &consoleDLL{}, nil }

// setConsoleMode is a no-op on Unix: ANSI colours work natively.
func setConsoleMode(_ *consoleDLL) error { return nil }

// getConsoleMode returns 0 on Unix (not used).
func getConsoleMode(_ *consoleDLL) (uint32, error) { return 0, nil }

// setConsoleModeRaw puts stdin into raw mode (no echo, no line buffering).
func setConsoleModeRaw(_ *consoleDLL) error {
	fd := os.Stdin.Fd()
	var t termiosState
	if err := tcgetattr(fd, &t); err != nil {
		return err
	}
	savedState = t
	stdinFd = fd
	stateValid = true

	// Raw mode: disable canonical mode and echo; read 1 byte at a time
	t.Lflag &^= syscall.ICANON | syscall.ECHO
	t.Cc[syscall.VMIN] = 1
	t.Cc[syscall.VTIME] = 0
	return tcsetattr(fd, &t)
}

// restoreConsoleMode restores the terminal to the state before setConsoleModeRaw.
func restoreConsoleMode(_ *consoleDLL, _ uint32) error {
	if stateValid {
		err := tcsetattr(stdinFd, &savedState)
		stateValid = false
		return err
	}
	return nil
}

// ──────────────────────────────────────────────────────────────────
// Key reading (ANSI escape sequences)
// ──────────────────────────────────────────────────────────────────

// readKey reads one keypress and returns a keyXxx constant.
func readKey() (int, error) {
	k, _, err := readKeyOrChar()
	return k, err
}

// readKeyOrChar reads a single keypress and returns (keyCode, char, error).
// Printable characters: keyCode == keyChar, char holds the rune.
// Special keys: char == 0, keyCode is one of the keyXxx constants.
func readKeyOrChar() (int, rune, error) {
	buf := make([]byte, 8)
	n, err := os.Stdin.Read(buf)
	if err != nil || n == 0 {
		return keyEnter, 0, err
	}

	// ANSI escape sequence: ESC [ A/B/C/D  (arrow keys)
	if n >= 3 && buf[0] == 0x1b && buf[1] == '[' {
		switch buf[2] {
		case 'A':
			return keyUp, 0, nil
		case 'B':
			return keyDown, 0, nil
		case 'C':
			return keyRight, 0, nil
		case 'D':
			return keyLeft, 0, nil
		}
	}

	// Single-byte keys
	switch buf[0] {
	case 0x1b: // ESC alone (or unrecognised escape sequence)
		return keyEsc, 0, nil
	case '\r', '\n': // Enter
		return keyEnter, 0, nil
	case 0x7f, 0x08: // DEL / Backspace
		return keyBackspace, 0, nil
	case 'q', 'Q':
		return keyQ, 0, nil
	}

	// Printable ASCII
	ch := rune(buf[0])
	if ch >= 32 && ch < 127 {
		return keyChar, ch, nil
	}

	return keyEnter, 0, nil
}
