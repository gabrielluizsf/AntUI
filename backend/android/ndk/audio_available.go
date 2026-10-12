//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl

#include <stdlib.h>

#include "audio_cgo.h"
*/
import "C"

// AudioAvailable reports whether AAudio can be used at all. It opens the
// library the first time and remembers.
func AudioAvailable() bool { return C.mw_aaudio_load() == 1 }
