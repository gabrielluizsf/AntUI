//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// Size is how much of the screen the touch covers, normalised to 0..1. It is
// how a thumb is told from a fingertip.
func (e *InputEvent) Size(index int) float32 {
	return float32(C.AMotionEvent_getSize(e.ptr, C.size_t(index)))
}
