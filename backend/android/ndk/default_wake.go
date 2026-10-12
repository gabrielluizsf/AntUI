//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// DefaultWake is Default, asking specifically for a sensor that can wake the
// device or specifically for one that cannot. A wake-up sensor keeps the
// processor alive to deliver its readings; a non-wake-up one is silent while
// the device sleeps and reports what it collected when it wakes.
func (m *SensorManager) DefaultWake(kind int, wake bool) Sensor {
	var w C.bool
	if wake {
		w = true
	}
	return Sensor{ptr: C.ASensorManager_getDefaultSensorEx(m.ptr, C.int(kind), w)}
}
