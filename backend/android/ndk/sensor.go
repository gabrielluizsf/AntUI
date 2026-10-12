//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// Sensor is one physical or computed sensor.
type Sensor struct {
	ptr *C.ASensor
}
