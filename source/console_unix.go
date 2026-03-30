//go:build !windows

package main

import "syscall"

// keyXxx constants (mirrored from console_windows.go for non-Windows builds).
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

// loadKernel32 is a no-op stub on non-Windows platforms.
func loadKernel32() (*syscall.DLL, error) { return nil, nil }

// setConsoleMode is a no-op stub on non-Windows platforms.
func setConsoleMode(_ *syscall.DLL) error { return nil }

// getConsoleMode is a no-op stub on non-Windows platforms.
func getConsoleMode(_ *syscall.DLL) (uint32, error) { return 0, nil }

// setConsoleModeRaw is a no-op stub on non-Windows platforms.
func setConsoleModeRaw(_ *syscall.DLL) error { return nil }

// restoreConsoleMode is a no-op stub on non-Windows platforms.
func restoreConsoleMode(_ *syscall.DLL, _ uint32) error { return nil }

// readKey is a stub for non-Windows platforms.
func readKey() (int, error) { return keyEnter, nil }

// readKeyOrChar is a stub for non-Windows platforms.
func readKeyOrChar() (int, rune, error) { return keyEnter, 0, nil }
