//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// MotionAction is what happened, with the pointer index masked off.
func (e *InputEvent) MotionAction() int {
	return int(C.AMotionEvent_getAction(e.ptr)) & C.AMOTION_EVENT_ACTION_MASK
}
