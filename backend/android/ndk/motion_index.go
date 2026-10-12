//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// MotionIndex is *which* pointer a POINTER_DOWN or POINTER_UP is about, as
// an index into this event's pointers — **not** a pointer id.
//
// This is the single most common way to get multi-touch wrong. The index is
// a position in the list of pointers in this event and changes as fingers
// come and go; the id from [InputEvent.PointerID] is what stays with a
// finger for as long as it is down. Track ids; use indices only to read this
// event.
func (e *InputEvent) MotionIndex() int {
	a := int(C.AMotionEvent_getAction(e.ptr))
	return (a & C.AMOTION_EVENT_ACTION_POINTER_INDEX_MASK) >>
		C.AMOTION_EVENT_ACTION_POINTER_INDEX_SHIFT
}
