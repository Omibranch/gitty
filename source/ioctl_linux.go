//go:build linux

package main

// TCGETS / TCSETS ioctl numbers for Linux (architecture-independent values).
const (
	ioctlTcgeta = 0x5401 // TCGETS
	ioctlTcseta = 0x5402 // TCSETS
)
