//go:build windows

package windows

import (
	"sync"
	"syscall"
)

var (
	classOnce sync.Once
	classErr  error
	wndProcCB = syscall.NewCallback(wndProc)
)
