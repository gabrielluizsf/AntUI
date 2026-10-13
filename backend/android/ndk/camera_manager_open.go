//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

import "unsafe"

// Open opens a camera. It needs the CAMERA permission, and fails rather than
// returning a black picture without it.
func (m *CameraManager) Open(id string) (*CameraDevice, error) {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	var dev *C.ACameraDevice
	if r := C.mw_open(m.ptr, cid, &dev); r != C.ACAMERA_OK {
		return nil, cameraError("opening camera "+id, r)
	}
	return &CameraDevice{ptr: dev}, nil
}
