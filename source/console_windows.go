//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	enableVirtualTerminalProcessing = 0x0004
	stdOutputHandle                 = ^uintptr(10) // (DWORD)-11
	stdInputHandle                  = ^uintptr(9)  // (DWORD)-10
)

// Virtual-key codes returned by ReadConsoleInput
const (
	vkLeft  = 0x25
	vkRight = 0x27
	vkUp    = 0x26
	vkDown  = 0x28
	vkEnter = 0x0D
	vkEsc   = 0x1B
)

// keyXxx constants used by pickLanguage / readKey (cross-platform names)
const (
	keyLeft  = 1
	keyRight = 2
	keyUp    = 3
	keyDown  = 4
	keyEnter = 5
	keyEsc   = 6
	keyQ     = 7
	keyBackspace = 8
	keyChar  = 9 // printable character – see readKeyOrChar
)

// INPUT_RECORD and KEY_EVENT_RECORD structures for ReadConsoleInput
type keyEventRecord struct {
	bKeyDown          int32
	wRepeatCount      uint16
	wVirtualKeyCode   uint16
	wVirtualScanCode  uint16
	unicodeChar       uint16
	dwControlKeyState uint32
}

type inputRecord struct {
	eventType uint16
	_         [2]byte
	// KEY_EVENT is the largest event type (20 bytes)
	event [20]byte
}

var (
	kernel32dll        *syscall.DLL
	procGetStdHandle   *syscall.Proc
	procGetConsoleMode *syscall.Proc
	procSetConsoleMode *syscall.Proc
	procReadConsoleInput *syscall.Proc
)

// loadKernel32 loads kernel32.dll and resolves the required procs.
func loadKernel32() (*syscall.DLL, error) {
	dll, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return nil, fmt.Errorf("could not load kernel32.dll: %w", err)
	}
	kernel32dll = dll

	procGetStdHandle, _ = dll.FindProc("GetStdHandle")
	procGetConsoleMode, _ = dll.FindProc("GetConsoleMode")
	procSetConsoleMode, _ = dll.FindProc("SetConsoleMode")
	procReadConsoleInput, _ = dll.FindProc("ReadConsoleInputW")
	return dll, nil
}

// setConsoleMode enables ANSI virtual terminal processing on stdout.
func setConsoleMode(dll *syscall.DLL) error {
	if procGetStdHandle == nil || procGetConsoleMode == nil || procSetConsoleMode == nil {
		return fmt.Errorf("required procs not loaded")
	}

	handle, _, _ := procGetStdHandle.Call(stdOutputHandle)
	if handle == 0 {
		return fmt.Errorf("invalid stdout handle")
	}

	var mode uint32
	r, _, err2 := procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if r == 0 {
		return fmt.Errorf("GetConsoleMode: %w", err2)
	}

	mode |= enableVirtualTerminalProcessing
	r, _, err2 = procSetConsoleMode.Call(handle, uintptr(mode))
	if r == 0 {
		return fmt.Errorf("SetConsoleMode: %w", err2)
	}
	return nil
}

// readKey reads a single keypress from the console and returns a keyXxx constant.
// It uses ReadConsoleInputW so it can distinguish arrow keys without needing
// any external library.
func readKey() (int, error) {
	k, _, err := readKeyOrChar()
	return k, err
}

// getConsoleMode returns the current stdin console mode.
func getConsoleMode(dll *syscall.DLL) (uint32, error) {
	handle, _, _ := procGetStdHandle.Call(stdInputHandle)
	if handle == 0 {
		return 0, fmt.Errorf("invalid stdin handle")
	}
	var mode uint32
	r, _, err := procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if r == 0 {
		return 0, err
	}
	return mode, nil
}

// setConsoleModeRaw disables line-input and echo on stdin (raw mode).
func setConsoleModeRaw(dll *syscall.DLL) error {
	handle, _, _ := procGetStdHandle.Call(stdInputHandle)
	if handle == 0 {
		return fmt.Errorf("invalid stdin handle")
	}
	var mode uint32
	procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	rawMode := mode &^ uint32(0x0006) // clear ENABLE_LINE_INPUT | ENABLE_ECHO_INPUT
	r, _, err := procSetConsoleMode.Call(handle, uintptr(rawMode))
	if r == 0 {
		return err
	}
	return nil
}

// restoreConsoleMode restores stdin to a previously saved mode.
func restoreConsoleMode(dll *syscall.DLL, mode uint32) error {
	handle, _, _ := procGetStdHandle.Call(stdInputHandle)
	if handle == 0 {
		return fmt.Errorf("invalid stdin handle")
	}
	procSetConsoleMode.Call(handle, uintptr(mode))
	return nil
}

// readKeyOrChar reads a single event and returns (keyCode, char, error).
// If the event is a printable character, keyCode == keyChar and char holds the rune.
// Otherwise char == 0 and keyCode is one of the keyXxx constants.
func readKeyOrChar() (int, rune, error) {
	if procReadConsoleInput == nil || procGetStdHandle == nil {
		return keyEnter, 0, nil
	}

	handle, _, _ := procGetStdHandle.Call(stdInputHandle)
	if handle == 0 {
		return keyEnter, 0, nil
	}

	for {
		var rec inputRecord
		var numRead uint32
		r, _, _ := procReadConsoleInput.Call(
			handle,
			uintptr(unsafe.Pointer(&rec)),
			1,
			uintptr(unsafe.Pointer(&numRead)),
		)
		if r == 0 || numRead == 0 {
			return 0, 0, fmt.Errorf("ReadConsoleInput failed")
		}

		if rec.eventType != 1 {
			continue
		}

		kr := (*keyEventRecord)(unsafe.Pointer(&rec.event[0]))
		if kr.bKeyDown == 0 {
			continue
		}

		switch kr.wVirtualKeyCode {
		case vkLeft:
			return keyLeft, 0, nil
		case vkRight:
			return keyRight, 0, nil
		case vkUp:
			return keyUp, 0, nil
		case vkDown:
			return keyDown, 0, nil
		case vkEnter:
			return keyEnter, 0, nil
		case vkEsc:
			return keyEsc, 0, nil
		case 0x08: // VK_BACK
			return keyBackspace, 0, nil
		default:
			ch := rune(kr.unicodeChar)
			if ch == 'q' || ch == 'Q' {
				return keyQ, 0, nil
			}
			if ch >= 32 && ch < 127 {
				return keyChar, ch, nil
			}
		}
	}
}
