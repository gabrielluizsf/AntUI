//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Frame converts the newest frame into dst, which holds 0xAARRGGBB words
// with stride words between rows, turning it by rotate degrees clockwise on
// the way — 0, 90, 180 or 270. The width and height that come back are of
// the *rotated* image, so they are swapped for a quarter turn.
//
// It reports whether there was a frame: a preview produces them on its own
// schedule and asking more often than it produces them is normal. Frames
// that arrived and were not read are dropped, because what is wanted from a
// camera is always the newest.
func (d *CameraDevice) Frame(dst []uint32, stride, rotate int) (w, h int, ok bool, err error) {
	if !d.running || len(dst) == 0 {
		return 0, 0, false, nil
	}
	var cw, ch C.int32_t
	r := C.mw_preview_frame(&d.preview, (*C.uint32_t)(unsafe.Pointer(&dst[0])),
		C.int(stride), C.int(rotate), &cw, &ch)
	switch r {
	case 0:
		return 0, 0, false, nil
	case 1:
		return int(cw), int(ch), true, nil
	}
	return 0, 0, false, errors.New("ndk: cannot read the camera's image")
}
