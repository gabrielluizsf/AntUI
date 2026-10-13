//go:build android

package ndk

import "errors"

// ErrNoBatching is what [SensorQueue.Register] reports on a platform too old
// to have it.
var ErrNoBatching = errors.New("ndk: this platform has no sensor batching; " +
	"it arrived in API 26")
