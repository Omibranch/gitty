//go:build windows

package main

import (
	"fmt"
	"os/exec"
)

// copyToClipboard copies text to the Windows clipboard via PowerShell.
func copyToClipboard(text string) error {
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		fmt.Sprintf("Set-Clipboard -Value '%s'", text))
	return cmd.Run()
}
