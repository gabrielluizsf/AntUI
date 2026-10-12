//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// Register starts a sensor with a rate and a batching latency, both in
// microseconds.
//
// Latency is what makes batching: the platform is allowed to hold readings
// for that long and deliver them together, which lets the processor sleep in
// between. For a step counter running all day that is the difference between
// a percent of the battery and ten. A latency of 0 means deliver each one as
// it happens.
func (q *SensorQueue) Register(s Sensor, rateMicros int, latencyMicros int64) error {
	if C.mw_register(q.ptr, s.ptr, C.int32_t(rateMicros), C.int64_t(latencyMicros)) < 0 {
		return ErrNoBatching
	}
	return nil
}
