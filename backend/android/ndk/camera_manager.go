//go:build android

package ndk

/*
#cgo LDFLAGS: -ldl -landroid

#include <stdlib.h>

#include "camera_cgo.h"
*/
import "C"

// CameraManager is the list of cameras.
type CameraManager struct{ ptr *C.ACameraManager }
