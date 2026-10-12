//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// Default is the sensor of a kind that the platform considers the best one,
// or a nil Sensor when the device has none.
func (m *SensorManager) Default(kind int) Sensor {
	return Sensor{ptr: C.ASensorManager_getDefaultSensor(m.ptr, C.int(kind))}
}
