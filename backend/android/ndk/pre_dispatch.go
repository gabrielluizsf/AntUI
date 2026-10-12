//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// PreDispatch offers the event to the system first, and reports whether the
// system took it. **Every event has to go through this before it is looked
// at**: it is how the on-screen keyboard sees a key at all, and an app that
// skips it appears to work until someone types.
//
// When it reports true the event has been consumed and must not be finished
// — the system will return it later if it did not want it after all.
func (q *InputQueue) PreDispatch(e *InputEvent) bool {
	return C.AInputQueue_preDispatchEvent(q.ptr, e.ptr) != 0
}
