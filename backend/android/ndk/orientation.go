//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

import (
	"unsafe"
)

// Orientation is how far the camera's own image is rotated from the device's
// natural orientation, in degrees clockwise: 90 on almost every phone.
//
// It is why photographs come out sideways. The sensor is mounted turned, the
// image arrives turned, and nothing rotates it — the app has to.
func (m *CameraManager) Orientation(id string) int {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	return int(C.mw_camera_int(m.ptr, cid, C.ACAMERA_SENSOR_ORIENTATION, 90))
}
