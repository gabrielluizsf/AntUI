//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// Name, Vendor, Type, Resolution and MinDelay describe one sensor.
func (s Sensor) Name() string { return C.GoString(C.ASensor_getName(s.ptr)) }
