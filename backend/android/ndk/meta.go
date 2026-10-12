//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// Meta is the modifier keys held, for a key or a motion.
func (e *InputEvent) Meta() Meta {
	if e.Kind() == EventKey {
		return Meta(C.AKeyEvent_getMetaState(e.ptr))
	}
	return Meta(C.AMotionEvent_getMetaState(e.ptr))
}
