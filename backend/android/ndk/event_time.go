//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// EventTime is when it happened, in nanoseconds on the same clock as
// SystemClock.uptimeMillis. It is what a fling's speed is measured against.
func (e *InputEvent) EventTime() int64 {
	return int64(C.AMotionEvent_getEventTime(e.ptr))
}
