//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl

#include <stdlib.h>

#include "audio_cgo.h"
*/
import "C"

// AudioStream is one stream of sound going out.
//
// It is written to rather than pulling: a goroutine calls [AudioStream.Write]
// and the call blocks until there is room. AAudio also has a callback mode,
// which is lower latency and runs on a real-time thread — a thread Go's
// collector may stop, which is exactly what a real-time thread must not
// have happen to it.
type AudioStream struct {
	ptr *C.AAudioStream
	// Rate, Channels and Burst are what the device settled on, which is not
	// always what was asked for.
	Rate     int
	Channels int
	// Burst is how many frames the device moves at a time. Writing in
	// multiples of it is what keeps a stream from stuttering; writing less
	// than one is a write that cannot be satisfied without waiting.
	Burst int
}
