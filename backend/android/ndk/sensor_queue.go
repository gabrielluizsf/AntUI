//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// SensorQueue is a stream of readings, delivered to a looper.
type SensorQueue struct {
	ptr     *C.ASensorEventQueue
	manager *SensorManager
}
