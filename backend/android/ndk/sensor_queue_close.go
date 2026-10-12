//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// Close destroys the queue. Every sensor on it stops.
func (q *SensorQueue) Close() {
	if q.ptr != nil {
		C.ASensorManager_destroyEventQueue(q.manager.ptr, q.ptr)
		q.ptr = nil
	}
}
