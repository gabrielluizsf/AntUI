//go:build android

package gallery

import (
	"errors"
)

// ErrNoPermission is saving on a device old enough to need one.
var ErrNoPermission = errors.New("gallery: below Android 10 an app needs " +
	"WRITE_EXTERNAL_STORAGE to put anything in the device's library")
