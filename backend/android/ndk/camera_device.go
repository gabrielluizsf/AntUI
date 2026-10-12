//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

// CameraDevice is an open camera.
type CameraDevice struct {
	ptr     *C.ACameraDevice
	preview C.mw_preview
	running bool
}
