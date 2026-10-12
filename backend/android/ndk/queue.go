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

// Queue makes a stream on the given looper. Like the input queue, it is
// thread-local: readings arrive only on the thread the looper belongs to.
func (m *SensorManager) Queue(l *Looper) (*SensorQueue, error) {
	q := C.ASensorManager_createEventQueue(m.ptr, l.ptr, C.int(SensorID), nil, nil)
	if q == nil {
		return nil, errors.New("ndk: cannot make a sensor queue")
	}
	return &SensorQueue{ptr: q, manager: m}, nil
}
