//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// SensorManager is the platform's list of sensors and the maker of queues.
type SensorManager struct {
	ptr *C.ASensorManager
}
