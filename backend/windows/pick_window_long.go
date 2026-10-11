//go:build windows

package windows

import (
	"syscall"
)

func pickWindowLong(name string, fallback *syscall.LazyProc) *syscall.LazyProc {
	proc := user32.NewProc(name)
	if proc.Find() == nil {
		return proc
	}
	return fallback
}
