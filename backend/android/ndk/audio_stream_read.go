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

// Read takes frames from the microphone and waits for them. data is filled
// interleaved, one float per channel per frame, and the count of frames read
// comes back — which may be fewer than asked for.
func (s *AudioStream) Read(data []float32, timeout time.Duration) (int, error) {
	if s.ptr == nil {
		return 0, errors.New("ndk: the audio stream is closed")
	}
	frames := len(data) / s.Channels
	if frames == 0 {
		return 0, nil
	}
	r := C.mw_audio_read(s.ptr, (*C.float)(unsafe.Pointer(&data[0])),
		C.int32_t(frames), C.int64_t(timeout.Nanoseconds()))
	if r < 0 {
		return 0, audioError("reading", C.aaudio_result_t(r))
	}
	return int(r), nil
}
