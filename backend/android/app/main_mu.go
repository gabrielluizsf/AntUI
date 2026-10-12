//go:build android

package app

import (
	"sync"
)

var (
	mainMu   sync.Mutex
	mainFn   func(*App)
	fallback func(*App)
)
