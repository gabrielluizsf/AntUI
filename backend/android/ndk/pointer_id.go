//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// PointerID is the identifier that stays with one finger from the moment it
// touches to the moment it lifts.
func (e *InputEvent) PointerID(index int) int {
	return int(C.AMotionEvent_getPointerId(e.ptr, C.size_t(index)))
}
