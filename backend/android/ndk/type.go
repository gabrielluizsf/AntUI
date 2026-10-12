//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// Type is what the sensor measures, as one of the ASENSOR_TYPE values.
func (s Sensor) Type() int { return int(C.ASensor_getType(s.ptr)) }
