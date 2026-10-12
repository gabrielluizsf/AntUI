//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Sensors opens the sensor manager for a package.
//
// The package name matters: since API 26 the platform uses it to decide
// which sensors an app may see, and the older call that took no name is
// deprecated. Passing an empty name uses the old one deliberately; passing a
// name uses the new one where it exists and the old one where it does not.
func Sensors(packageName string) (*SensorManager, error) {
	var m *C.ASensorManager
	if packageName == "" {
		m = C.ASensorManager_getInstance()
	} else {
		name := C.CString(packageName)
		defer C.free(unsafe.Pointer(name))
		// Falls back to the nameless call on a platform too old for the
		// other, which is what the C helper does.
		m = C.mw_manager_for(name)
	}
	if m == nil {
		return nil, errors.New("ndk: this device has no sensor manager")
	}
	return &SensorManager{ptr: m}, nil
}
