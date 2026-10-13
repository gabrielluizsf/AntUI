//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

import "errors"

// Cameras opens the camera manager.
func Cameras() (*CameraManager, error) {
	if C.mw_camera_load() != 1 {
		return nil, ErrNoCamera2
	}
	m := C.mw_manager_create()
	if m == nil {
		return nil, errors.New("ndk: cannot open the camera manager")
	}
	return &CameraManager{ptr: m}, nil
}
