//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// HasEvents reports whether anything is waiting, without taking it.
func (q *SensorQueue) HasEvents() bool {
	return C.ASensorEventQueue_hasEvents(q.ptr) > 0
}
