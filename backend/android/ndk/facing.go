//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

import "unsafe"

// Facing is which way a camera points: [CameraBack], [CameraFront] or
// [CameraExternal].
func (m *CameraManager) Facing(id string) int {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	return int(C.mw_camera_int(m.ptr, cid, C.ACAMERA_LENS_FACING, C.int(CameraBack)))
}
