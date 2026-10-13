//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

import "fmt"

func cameraError(what string, r C.camera_status_t) error {
	return fmt.Errorf("ndk: camera: %s failed (%d)", what, int(r))
}
