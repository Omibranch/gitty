//go:build darwin

package main

// TIOCGETA / TIOCSETA ioctl numbers for macOS.
const (
	ioctlTcgeta = 0x40487413 // TIOCGETA
	ioctlTcseta = 0x80487414 // TIOCSETA
)
