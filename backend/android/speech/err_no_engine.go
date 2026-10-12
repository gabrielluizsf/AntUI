//go:build android

package speech

import (
	"errors"
)

// ErrNoEngine is a device with no speech engine installed. It is not rare:
// the engine is an app, and a device without Google's services may have none.
var ErrNoEngine = errors.New("speech: this device has no text-to-speech engine")
