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
)

// StartPreview begins a stream of frames at the given size.
func (d *CameraDevice) StartPreview(width, height int) error {
	if d.running {
		return errors.New("ndk: the preview is already running")
	}
	if r := C.mw_preview_start(d.ptr, C.int32_t(width), C.int32_t(height),
		C.int32_t(FormatYUV420), &d.preview); r != C.ACAMERA_OK {
		C.mw_preview_stop(&d.preview)
		return cameraError("starting the preview", r)
	}
	d.running = true
	return nil
}
