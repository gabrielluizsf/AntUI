//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

// Close stops the preview and gives the camera back.
func (d *CameraDevice) Close() {
	if d.running {
		C.mw_preview_stop(&d.preview)
		d.running = false
	}
	if d.ptr != nil {
		C.mw_device_close(d.ptr)
		d.ptr = nil
	}
}
