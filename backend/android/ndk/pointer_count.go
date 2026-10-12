//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// PointerCount is how many pointers this event reports on.
func (e *InputEvent) PointerCount() int {
	return int(C.AMotionEvent_getPointerCount(e.ptr))
}
