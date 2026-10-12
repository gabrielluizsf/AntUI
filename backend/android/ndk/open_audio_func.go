//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl

#include <stdlib.h>

#include "audio_cgo.h"
*/
import "C"

func openAudio(rate, channels, capacity int, input bool) (*AudioStream, error) {
	if C.mw_aaudio_load() != 1 {
		return nil, ErrNoAAudio
	}
	var in C.int32_t
	if input {
		in = 1
	}
	var (
		stream      *C.AAudioStream
		gotRate     C.int32_t
		gotChannels C.int32_t
		burst       C.int32_t
	)
	r := C.mw_audio_open(C.int32_t(rate), C.int32_t(channels), C.int32_t(capacity),
		in, &stream, &gotRate, &gotChannels, &burst)
	if r != C.AAUDIO_OK {
		return nil, audioError("opening the stream", r)
	}
	s := &AudioStream{
		ptr:      stream,
		Rate:     int(gotRate),
		Channels: int(gotChannels),
		Burst:    int(burst),
	}
	// Two bursts is the smallest buffer that does not underrun on the first
	// hiccup, and is what every AAudio example starts from.
	if s.Burst > 0 {
		C.mw_audio_set_buffer(stream, C.int32_t(s.Burst*2))
	}
	if r := C.mw_audio_start(stream); r != C.AAUDIO_OK {
		C.mw_audio_close(stream)
		return nil, audioError("starting the stream", r)
	}
	return s, nil
}
