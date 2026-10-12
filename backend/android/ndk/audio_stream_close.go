//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl

#include <stdlib.h>

#include "audio_cgo.h"
*/
import "C"

// Close stops the stream and gives it back.
func (s *AudioStream) Close() error {
	if s.ptr == nil {
		return nil
	}
	C.mw_audio_stop(s.ptr)
	r := C.mw_audio_close(s.ptr)
	s.ptr = nil
	if r != C.AAUDIO_OK {
		return audioError("closing", r)
	}
	return nil
}
