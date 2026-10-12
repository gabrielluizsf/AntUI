//go:build android

package ndk

/*
#cgo LDFLAGS: -lmediandk -landroid

#include "media_cgo.h"
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"
)

// Samples decodes sound into dst, which is filled interleaved, and reports
// how many floats were written.
func (d *Decoder) Samples(dst []float32) (n int, pts time.Duration, err error) {
	if d.closed || len(dst) == 0 {
		return 0, 0, nil
	}
	var wrote C.int
	switch C.mw_audio_pump(&d.media, (*C.float)(unsafe.Pointer(&dst[0])),
		C.int(len(dst)), &wrote) {
	case pumpError:
		return 0, 0, errors.New("ndk: the audio decoder failed")
	case pumpEnd:
		d.ended = true
	}
	return int(wrote), time.Duration(d.media.pts) * time.Microsecond, nil
}
