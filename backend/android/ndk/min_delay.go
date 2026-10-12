//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include "sensor_cgo.h"
*/
import "C"

// MinDelay is the shortest interval between readings, in microseconds. It is
// 0 for a sensor that only reports when something changes — a step detector,
// a proximity sensor — and asking such a sensor for a rate does nothing.
func (s Sensor) MinDelay() int { return int(C.ASensor_getMinDelay(s.ptr)) }
