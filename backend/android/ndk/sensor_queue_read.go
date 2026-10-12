//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// Read takes up to len(into) readings, and reports how many there were.
func (q *SensorQueue) Read(into []SensorEvent) int {
	if len(into) == 0 {
		return 0
	}
	buf := make([]C.ASensorEvent, len(into))
	n := int(C.ASensorEventQueue_getEvents(q.ptr, &buf[0], C.size_t(len(into))))
	if n <= 0 {
		return 0
	}
	for i := range n {
		e := &buf[i]
		into[i].Type = int(C.mw_sensor_type(e))
		into[i].ID = int(C.mw_sensor_id(e))
		into[i].At = int64(C.mw_sensor_time(e))
		for j := range into[i].V {
			into[i].V[j] = float32(C.mw_sensor_data(e, C.int(j)))
		}
	}
	return n
}
