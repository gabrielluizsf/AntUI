//go:build android

package torch

import "errors"

// ErrNoFlash is a device with no lamp on any of its cameras.
var ErrNoFlash = errors.New("torch: no camera on this device has a flash")
