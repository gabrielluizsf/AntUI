//go:build android

package torch

import "sync"

var (
	idOnce  sync.Once
	flashID string
	idErr   error
)
