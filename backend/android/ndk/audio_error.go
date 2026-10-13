//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl

#include <stdlib.h>

#include "audio_cgo.h"
*/
import "C"

import "errors"

// audioError turns AAudio's numbered result into its own words, which are
// better than anything this package would invent.
func audioError(what string, r C.aaudio_result_t) error {
	text := C.GoString(C.mw_audio_text(r))
	if text == "" {
		return errors.New("ndk: audio: " + what + " failed")
	}
	return errors.New("ndk: audio: " + what + ": " + text)
}
