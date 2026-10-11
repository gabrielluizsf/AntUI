//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"sync"
)

var (
	themeMu sync.Mutex
	theme   themeState
)
