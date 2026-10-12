//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// Y is where a pointer is, down the window, in pixels.
func (e *InputEvent) Y(index int) float32 {
	return float32(C.AMotionEvent_getY(e.ptr, C.size_t(index)))
}
