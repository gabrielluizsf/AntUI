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

// Frame decodes until a picture is ready, and converts it into dst.
//
// It reports whether there is one: a decoder needs several compressed
// samples before it produces its first frame, so a call that gives nothing
// is normal and not the end. The end is [Decoder.Ended].
func (d *Decoder) Frame(dst []uint32, stride int) (w, h int, pts time.Duration, ok bool, err error) {
	if d.closed || len(dst) == 0 {
		return 0, 0, 0, false, nil
	}
	switch C.mw_video_pump(&d.media) {
	case pumpError:
		return 0, 0, 0, false, errors.New("ndk: the video decoder failed")
	case pumpEnd:
		d.ended = true
	}
	var cw, ch C.int32_t
	got := C.mw_media_frame(&d.media, (*C.uint32_t)(unsafe.Pointer(&dst[0])),
		C.int(stride), C.int(d.rotate), &cw, &ch)
	if got != 1 {
		return 0, 0, 0, false, nil
	}
	return int(cw), int(ch), time.Duration(d.media.pts) * time.Microsecond, true, nil
}
