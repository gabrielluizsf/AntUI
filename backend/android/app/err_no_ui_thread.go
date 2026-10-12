//go:build android

package app

import (
	"errors"
)

// ErrNoUIThread is what the runner reports before the activity exists.
var ErrNoUIThread = errors.New("app: there is no UI thread yet")
