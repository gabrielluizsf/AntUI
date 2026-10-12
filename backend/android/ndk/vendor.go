//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// Vendor is who made the sensor, which is the only way to tell two of
// the same kind apart.
func (s Sensor) Vendor() string { return C.GoString(C.ASensor_getVendor(s.ptr)) }
