//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// Resolution is the smallest change the sensor can report, in its own units.
func (s Sensor) Resolution() float32 { return float32(C.ASensor_getResolution(s.ptr)) }
