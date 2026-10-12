//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

import (
	"unsafe"
)

// List is every sensor the device has.
func (m *SensorManager) List() []Sensor {
	var list C.ASensorList
	n := int(C.ASensorManager_getSensorList(m.ptr, &list))
	if n <= 0 || list == nil {
		return nil
	}
	raw := unsafe.Slice((**C.ASensor)(unsafe.Pointer(list)), n)
	out := make([]Sensor, n)
	for i := range raw {
		out[i] = Sensor{ptr: raw[i]}
	}
	return out
}
