//go:build !windows

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// copyToClipboard copies text to the system clipboard on Linux/macOS.
// Tries multiple clipboard backends in order of preference:
//   - wl-copy   (Wayland)
//   - xclip     (X11)
//   - xsel      (X11 fallback)
//   - pbcopy    (macOS)
func copyToClipboard(text string) error {
	backends := [][]string{
		{"wl-copy"},                           // Wayland
		{"xclip", "-selection", "clipboard"},  // X11
		{"xsel", "--clipboard", "--input"},    // X11 fallback
		{"pbcopy"},                            // macOS
	}

	var lastErr error
	for _, args := range backends {
		if _, err := exec.LookPath(args[0]); err != nil {
			continue // not installed
		}
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}

	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("no clipboard utility found (install xclip, xsel, or wl-clipboard)")
}
