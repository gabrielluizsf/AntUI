//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

// CameraAvailable reports whether the camera API can be used at all.
func CameraAvailable() bool { return C.mw_camera_load() == 1 }
