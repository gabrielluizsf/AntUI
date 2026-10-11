//go:build windows

package windows

import (
	"sync"
)

// windowsFromHandle maps a window handle back to the Go window it belongs to.
// Keeping it here rather than in GWLP_USERDATA sidesteps the 32/64-bit split
// in the window-long calls, and keeps a Go pointer out of a C-visible slot.
var (
	windowsMu sync.RWMutex
	windows   = map[uintptr]*Driver{}
)
