//go:build android

package ndk

import (
	"errors"
)

// ErrNoCamera2 is a device too old for the NDK's camera API, which arrived
// in API 24.
var ErrNoCamera2 = errors.New("ndk: this device has no NDK camera; it arrived in API 24")
