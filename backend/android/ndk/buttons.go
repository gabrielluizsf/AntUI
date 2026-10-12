//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// Buttons is which mouse buttons are down.
func (e *InputEvent) Buttons() int32 {
	return int32(C.AMotionEvent_getButtonState(e.ptr))
}
