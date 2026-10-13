//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

import "errors"

// SetRate asks for readings no more often than every micros microseconds.
// It is a request: the platform delivers at that rate or slower, never
// faster, and a rate below the sensor's own minimum is clamped to it.
func (q *SensorQueue) SetRate(s Sensor, micros int) error {
	if C.ASensorEventQueue_setEventRate(q.ptr, s.ptr, C.int32_t(micros)) < 0 {
		return errors.New("ndk: cannot set the rate of " + s.Name())
	}
	return nil
}
