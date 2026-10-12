//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl

#include <stdlib.h>

#include "audio_cgo.h"
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"
)

// Write sends frames and waits for room. data is interleaved: one float per
// channel per frame.
func (s *AudioStream) Write(data []float32, timeout time.Duration) (int, error) {
	if s.ptr == nil {
		return 0, errors.New("ndk: the audio stream is closed")
	}
	frames := len(data) / s.Channels
	if frames == 0 {
		return 0, nil
	}
	r := C.mw_audio_write(s.ptr, (*C.float)(unsafe.Pointer(&data[0])),
		C.int32_t(frames), C.int64_t(timeout.Nanoseconds()))
	if r < 0 {
		return 0, audioError("writing", C.aaudio_result_t(r))
	}
	return int(r), nil
}
