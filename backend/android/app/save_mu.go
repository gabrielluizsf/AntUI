//go:build android

package app

import (
	"sync"
)

var (
	saveMu sync.Mutex
	saveFn func() []byte
)
