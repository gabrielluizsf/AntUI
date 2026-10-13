//go:build android

package ndk

import "errors"

// ErrNoImageReader is a device below API 24, where video frames cannot be
// got at from native code at all.
var ErrNoImageReader = errors.New("ndk: this device cannot hand decoded video " +
	"to native code; the image reader arrived in API 24")
