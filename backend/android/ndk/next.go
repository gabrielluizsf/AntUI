//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// Next takes the next event, or reports false when there are none. Every
// event it hands out must be given back to [InputQueue.Finish] or
// [InputQueue.PreDispatch] — the system has a fixed number of them and an
// app that keeps one stops receiving input entirely.
func (q *InputQueue) Next() (*InputEvent, bool) {
	var ev *C.AInputEvent
	if C.AInputQueue_getEvent(q.ptr, &ev) < 0 {
		return nil, false
	}
	return &InputEvent{ptr: ev}, true
}
