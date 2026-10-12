//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// Finish gives the event back, saying whether the app dealt with it. An
// unhandled key is passed on to the system, which is what makes the back
// button close the app when nothing else wanted it.
func (q *InputQueue) Finish(e *InputEvent, handled bool) {
	var h C.int
	if handled {
		h = 1
	}
	C.AInputQueue_finishEvent(q.ptr, e.ptr, h)
}
