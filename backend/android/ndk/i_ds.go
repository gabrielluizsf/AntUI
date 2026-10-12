//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

// IDs is every camera the app may open. They are opaque strings, not
// indices: "0" and "1" on most phones, and anything at all on a device with
// an external one.
func (m *CameraManager) IDs() ([]string, error) {
	var list *C.ACameraIdList
	if r := C.mw_id_list(m.ptr, &list); r != C.ACAMERA_OK {
		return nil, cameraError("listing the cameras", r)
	}
	defer C.mw_id_list_delete(list)
	n := int(C.mw_id_count(list))
	out := make([]string, n)
	for i := range n {
		out[i] = C.GoString(C.mw_id_at(list, C.int(i)))
	}
	return out, nil
}
