//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl

#include <stdlib.h>

#include "audio_cgo.h"
*/
import "C"

// BufferFrames is how much the device is holding, which is what latency is
// made of.
func (s *AudioStream) BufferFrames() int {
	if s.ptr == nil {
		return 0
	}
	return int(C.mw_audio_buffer(s.ptr))
}
