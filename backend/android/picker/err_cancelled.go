//go:build android

package picker

import "errors"

// ErrCancelled is the user backing out without choosing.
var ErrCancelled = errors.New("picker: the user chose nothing")
