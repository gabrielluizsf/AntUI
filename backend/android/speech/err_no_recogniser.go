//go:build android

package speech

import "errors"

// ErrNoRecogniser is a device with nothing installed that can listen. It is
// common: recognition is a service, usually Google's, and a device without
// those services has none.
var ErrNoRecogniser = errors.New("speech: this device has nothing that can listen")
