//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// X and Y are in the window's pixels, and are floats because a touchscreen
// reports between pixels and a mouse on a scaled display does too.
func (e *InputEvent) X(index int) float32 {
	return float32(C.AMotionEvent_getX(e.ptr, C.size_t(index)))
}
