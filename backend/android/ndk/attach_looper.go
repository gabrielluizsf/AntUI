//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// AttachLooper starts delivery to l. The callback argument the C function
// takes is deliberately null: a callback would be C calling into Go on a
// thread at a moment neither side chose, and there is nothing to gain from
// it when the frame loop is going to ask anyway.
func (q *InputQueue) AttachLooper(l *Looper) {
	C.AInputQueue_attachLooper(q.ptr, l.ptr, C.int(InputID), nil, nil)
}
