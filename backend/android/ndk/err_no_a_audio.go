//go:build android

package ndk

import (
	"errors"
)

// ErrNoAAudio is a device too old for AAudio, which is anything below
// Android 8.
var ErrNoAAudio = errors.New("ndk: this device has no AAudio; it arrived in API 26")
