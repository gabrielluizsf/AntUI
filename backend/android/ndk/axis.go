//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// Axis reads one axis of one pointer. Everything a game controller reports
// comes through here rather than through X and Y.
func (e *InputEvent) Axis(axis, index int) float32 {
	return float32(C.AMotionEvent_getAxisValue(e.ptr, C.int32_t(axis), C.size_t(index)))
}
