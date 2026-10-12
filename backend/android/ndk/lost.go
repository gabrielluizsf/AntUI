//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

// Lost reports whether the camera was taken away — unplugged, or claimed by
// another app, which the platform allows and does not ask about.
func (d *CameraDevice) Lost() bool { return C.mw_camera_lost() != 0 }
