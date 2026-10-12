//go:build android

package audio

import (
	"errors"
)

// ErrNoMicrophone is recording without the permission having been granted.
var ErrNoMicrophone = errors.New("audio: this app has not been granted " +
	"RECORD_AUDIO; ask for it with antui/backend/android/permission first")
