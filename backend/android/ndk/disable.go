//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

import (
	"errors"
)

// Disable stops one.
func (q *SensorQueue) Disable(s Sensor) error {
	if C.ASensorEventQueue_disableSensor(q.ptr, s.ptr) < 0 {
		return errors.New("ndk: cannot disable " + s.Name())
	}
	return nil
}
