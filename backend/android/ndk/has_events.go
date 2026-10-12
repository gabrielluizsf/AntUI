//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// HasEvents reports whether anything is waiting.
func (q *InputQueue) HasEvents() bool { return C.AInputQueue_hasEvents(q.ptr) > 0 }
