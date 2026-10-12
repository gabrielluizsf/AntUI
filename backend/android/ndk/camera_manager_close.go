//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

// Close gives the manager back.
func (m *CameraManager) Close() {
	if m.ptr != nil {
		C.mw_manager_delete(m.ptr)
		m.ptr = nil
	}
}
