//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

import "errors"

// Enable starts a sensor at its default rate.
func (q *SensorQueue) Enable(s Sensor) error {
	if C.ASensorEventQueue_enableSensor(q.ptr, s.ptr) < 0 {
		return errors.New("ndk: cannot enable " + s.Name())
	}
	return nil
}
