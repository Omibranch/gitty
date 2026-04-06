//go:build !windows

package main

// KernelHandle is the platform handle type for Unix console operations.
// On Unix this is a stub — no actual kernel handle is needed.
type KernelHandle = *consoleDLL
