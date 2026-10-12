//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// Pressure is 0 to 1 on a screen that measures it, and is 1 while down and 0
// while up on one that does not — so it cannot be used to tell whether the
// device reports pressure at all.
func (e *InputEvent) Pressure(index int) float32 {
	return float32(C.AMotionEvent_getPressure(e.ptr, C.size_t(index)))
}
