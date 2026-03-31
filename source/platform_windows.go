//go:build windows

package main

import "syscall"

// KernelHandle is the platform handle type for Windows console operations.
// On Windows this wraps *syscall.DLL.
type KernelHandle = *syscall.DLL
