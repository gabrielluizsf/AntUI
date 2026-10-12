//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// DetachLooper stops delivery. It must be called before the queue is
// destroyed, and from the looper's own thread.
func (q *InputQueue) DetachLooper() { C.AInputQueue_detachLooper(q.ptr) }
