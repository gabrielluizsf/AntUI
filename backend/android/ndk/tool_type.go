//go:build android

package ndk

/*
#cgo LDFLAGS: -landroid

#include <android/input.h>
#include <android/keycodes.h>
*/
import "C"

// ToolType says what is touching: a finger, a stylus, a mouse, or the far
// end of a stylus being used as a rubber.
func (e *InputEvent) ToolType(index int) int {
	return int(C.AMotionEvent_getToolType(e.ptr, C.size_t(index)))
}
